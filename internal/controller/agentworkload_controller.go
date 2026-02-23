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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	agenticv1alpha1 "github.com/shreyansh/agentic-operator/api/v1alpha1"
	"github.com/shreyansh/agentic-operator/pkg/mcp"
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
// 4. Updating status with proposed actions
// 5. Requeue after 30 seconds
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

	// Step 4: Update status with proposed action
	now := metav1.Now()
	action := agenticv1alpha1.Action{
		Name:        proposal["action"].(string),
		Description: proposal["description"].(string),
		Confidence:  proposal["confidence"].(string),
		Timestamp:   &now,
		Approved:    boolPtr(false),
	}

	workload.Status.ProposedActions = append(workload.Status.ProposedActions, action)
	workload.Status.Phase = "Running"
	workload.Status.LastReconcileTime = &now
	workload.Status.ReadyAgents = int32(len(workload.Spec.Agents))

	if err := r.Status().Update(ctx, &workload); err != nil {
		log.Error(err, "failed to update workload status")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	log.Info("Updated workload status with proposed action")

	// Step 5: Requeue after 30 seconds for next reconciliation
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
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
