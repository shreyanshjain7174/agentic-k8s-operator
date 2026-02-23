/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	agenticv1alpha1 "github.com/shreyansh/agentic-operator/api/v1alpha1"
	"github.com/shreyansh/agentic-operator/pkg/mcp"
	"github.com/shreyansh/agentic-operator/pkg/opa"
)

// AgentWorkloadReconciler reconciles a AgentWorkload object
type AgentWorkloadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=agentic.ninerewards.io,resources=agentworkloads,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=agentic.ninerewards.io,resources=agentworkloads/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=agentic.ninerewards.io,resources=agentworkloads/finalizers,verbs=update

// Reconcile reconciles the AgentWorkload by:
// 1. Fetching the AgentWorkload CR
// 2. Calling MCP to get status
// 3. Proposing actions via MCP
// 4. Evaluating action safety using OPA
// 5. Executing approved actions or marking for approval
// 6. Updating status
// 7. Requeue after 30 seconds
func (r *AgentWorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Step 1: Fetch the AgentWorkload
	var workload agenticv1alpha1.AgentWorkload
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &workload); err != nil {
		log.Error(err, "unable to fetch AgentWorkload")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling AgentWorkload", "name", workload.Name, "workloadType", workload.Spec.WorkloadType)

	// Step 2: Connect to MCP server and fetch status
	mcpClient := mcp.NewMCPClient(workload.Spec.MCPServerEndpoint)

	status, err := mcpClient.CallTool("get_status", map[string]interface{}{})
	if err != nil {
		log.Error(err, "failed to get status from MCP server")
		workload.Status.Phase = "Failed"
		if err := r.Status().Update(ctx, &workload); err != nil {
			log.Error(err, "failed to update workload status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Got status from MCP", "status", status)

	// Extract cluster health from status (default to 75 if not provided)
	clusterHealth := 75.0
	if health, ok := status["cluster_health"].(float64); ok {
		clusterHealth = health
	}

	// Step 3: Call MCP to propose an action
	proposalParams := map[string]interface{}{
		"objective": workload.Spec.Objective,
		"status":    status,
	}

	proposal, err := mcpClient.CallTool("propose_action", proposalParams)
	if err != nil {
		log.Error(err, "failed to propose action from MCP server")
		workload.Status.Phase = "Failed"
		if err := r.Status().Update(ctx, &workload); err != nil {
			log.Error(err, "failed to update workload status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Proposed action from MCP", "proposal", proposal)

	// Step 4: Evaluate action safety using OPA
	now := metav1.Now()
	
	// Convert confidence string to float
	confidenceStr := proposal["confidence"].(string)
	confidence, err := strconv.ParseFloat(confidenceStr, 64)
	if err != nil {
		log.Error(err, "failed to parse confidence", "confidence", confidenceStr)
		confidence = 0.0
	}

	// Create OPA evaluator and evaluate action
	evaluator := opa.NewPolicyEvaluator()
	opaInput := &opa.EvaluationInput{
		ActionType:         proposal["action"].(string),
		Confidence:         confidence,
		ClusterHealthScore: clusterHealth,
		OPAPolicyMode:      *workload.Spec.OPAPolicy,
	}

	opaResult := evaluator.Evaluate(opaInput)

	log.Info("OPA evaluation result", "allowed", opaResult.Allowed, "confidence", opaResult.Confidence, "reasons", opaResult.Reasons)

	// Step 5: Handle action execution or approval pending
	action := agenticv1alpha1.Action{
		Name:        proposal["action"].(string),
		Description: proposal["description"].(string),
		Confidence:  confidenceStr,
		Timestamp:   &now,
	}

	if opaResult.Allowed {
		// Step 5a: Execute approved action via MCP
		log.Info("OPA approved action, executing", "action", action.Name)

		executeParams := map[string]interface{}{
			"action":     action.Name,
			"params":     proposal,
			"confidence": confidenceStr,
		}

		execution, err := mcpClient.CallTool("execute_action", executeParams)
		if err != nil {
			log.Error(err, "failed to execute action", "action", action.Name)
			workload.Status.Phase = "Failed"
			action.Approved = boolPtr(false)
			workload.Status.ProposedActions = append(workload.Status.ProposedActions, action)
		} else {
			log.Info("Action executed successfully", "action", action.Name, "result", execution)
			action.Approved = boolPtr(true)
			workload.Status.ExecutedActions = append(workload.Status.ExecutedActions, action)
			workload.Status.Phase = "Completed"
		}
	} else {
		// Step 5b: Mark for human approval
		log.Info("OPA denied action, requiring human approval", "action", action.Name, "reasons", opaResult.Reasons)
		action.Approved = boolPtr(false)
		workload.Status.ProposedActions = append(workload.Status.ProposedActions, action)
		workload.Status.Phase = "PendingApproval"
	}

	// Step 6: Update status
	if workload.Status.Phase == "" || (workload.Status.Phase != "Completed" && workload.Status.Phase != "Failed" && workload.Status.Phase != "PendingApproval") {
		workload.Status.Phase = "Running"
	}
	workload.Status.LastReconcileTime = &now
	workload.Status.ReadyAgents = int32(len(workload.Spec.Agents))

	if err := r.Status().Update(ctx, &workload); err != nil {
		log.Error(err, "failed to update workload status")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Updated workload status", "phase", workload.Status.Phase)

	// Step 7: Determine requeue interval
	var requeueInterval time.Duration
	if workload.Status.Phase == "PendingApproval" {
		// For pending approval, requeue less frequently (1 hour) to wait for human approval
		requeueInterval = 1 * time.Hour
	} else {
		// For running, requeue quickly
		requeueInterval = 30 * time.Second
	}

	return ctrl.Result{RequeueAfter: requeueInterval}, nil
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentWorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agenticv1alpha1.AgentWorkload{}).
		Named("agentworkload").
		Complete(r)
}
