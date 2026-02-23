// Auto-generated for phase1-001
// Task: Initialize Kubebuilder project with AgenticProposal CRD

```go
// api/v1alpha1/groupversion_info.go
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "agentic.clawdlinux.org", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
```

```go
// api/v1alpha1/agenticproposal_types.go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgenticProposalSpec defines the desired state of AgenticProposal
type AgenticProposalSpec struct {
	// Content is the proposal text (markdown/plain text)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:MaxLength=32768
	Content string `json:"content"`

	// SubmittedBy identifies the agent that created the proposal
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	SubmittedBy string `json:"submittedBy"`

	// ProposedAt is the submission timestamp
	// +kubebuilder:validation:Required
	ProposedAt metav1.Time `json:"proposedAt"`

	// Tags for categorization (e.g., "security", "performance")
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Priority affects processing order (1=low, 5=critical)
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:default=3
	// +optional
	Priority int `json:"priority,omitempty"`

	// RequiredApprovers lists agents that MUST approve (subset of voters)
	// +optional
	RequiredApprovers []string `json:"requiredApprovers,omitempty"`
}

// ProposalState represents the proposal lifecycle state
// +kubebuilder:validation:Enum=Pending;UnderReview;Approved;Rejected;Executing;Completed;Failed
type ProposalState string

const (
	StatePending     ProposalState = "Pending"
	StateUnderReview ProposalState = "UnderReview"
	StateApproved    ProposalState = "Approved"
	StateRejected    ProposalState = "Rejected"
	StateExecuting   ProposalState = "Executing"
	StateCompleted   ProposalState = "Completed"
	StateFailed      ProposalState = "Failed"
)

// VoteDecision represents an agent's vote
// +kubebuilder:validation:Enum=Approve;Reject;Abstain
type VoteDecision string

const (
	VoteApprove VoteDecision = "Approve"
	VoteReject  VoteDecision = "Reject"
	VoteAbstain VoteDecision = "Abstain"
)

// AgentVote contains vote details from a single agent
type AgentVote struct {
	// Agent name
	Agent string `json:"agent"`

	// Vote decision
	Decision VoteDecision `json:"decision"`

	// Reasoning for the vote
	// +optional
	Reasoning string `json:"reasoning,omitempty"`

	// Timestamp of vote
	VotedAt metav1.Time `json:"votedAt"`

	// Score (0-100) if quantitative feedback
	// +optional
	Score *int `json:"score,omitempty"`
}

// ExecutionResult captures the outcome of executed proposals
type ExecutionResult struct {
	// Success indicates if execution completed successfully
	Success bool `json:"success"`

	// Message provides human-readable execution summary
	// +optional
	Message string `json:"message,omitempty"`

	// ExecutedAt timestamp
	ExecutedAt metav1.Time `json:"executedAt"`

	// Artifacts generated (e.g., commit SHA, PR URL)
	// +optional
	Artifacts map[string]string `json:"artifacts,omitempty"`
}

// AgenticProposalStatus defines the observed state of AgenticProposal
type AgenticProposalStatus struct {
	// State is the current lifecycle state
	// +kubebuilder:default=Pending
	State ProposalState `json:"state"`

	// Votes is a map of agent name -> vote decision
	// Deprecated: Use VoteDetails for richer vote information
	// +optional
	Votes map[string]string `json:"votes,omitempty"`

	// VoteDetails contains structured vote information
	// +optional
	VoteDetails []AgentVote `json:"voteDetails,omitempty"`

	// ApprovalScore is the consensus percentage (0-100)
	// +optional
	ApprovalScore *float64 `json:"approvalScore,omitempty"`

	// ApprovedAt timestamp (if approved)
	// +optional
	ApprovedAt *metav1.Time `json:"approvedAt,omitempty"`

	// RejectedAt timestamp (if rejected)
	// +optional
	RejectedAt *metav1.Time `json:"rejectedAt,omitempty"`

	// ExecutionResult captures execution outcome
	// +optional
	ExecutionResult *ExecutionResult `json:"executionResult,omitempty"`

	// Conditions track state transitions
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration reflects the generation last processed
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=proposal;proposals
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Score",type=number,JSONPath=`.status.approvalScore`
// +kubebuilder:printcolumn:name="Submitter",type=string,JSONPath=`.spec.submittedBy`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// AgenticProposal is the Schema for the agenticproposals API
type AgenticProposal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgenticProposalSpec   `json:"spec,omitempty"`
	Status AgenticProposalStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AgenticProposalList contains a list of AgenticProposal
type AgenticProposalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgenticProposal `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgenticProposal{}, &AgenticProposalList{})
}
```

```go
// internal/controller/agenticproposal_controller.go
package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	agenticv1alpha1 "github.com/shreyanshjain7174/clawdlinux-operator/api/v1alpha1"
)

// AgenticProposalReconciler reconciles a AgenticProposal object
type AgenticProposalReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// ConsensusThreshold is the minimum approval percentage (e.g., 80.0)
	ConsensusThreshold float64

	// ExecutionEnabled controls whether approved proposals auto-execute
	ExecutionEnabled bool
}

// +kubebuilder:rbac:groups=agentic.clawdlinux.org,resources=agenticproposals,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=agentic.clawdlinux.org,resources=agenticproposals/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=agentic.clawdlinux.org,resources=agenticproposals/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile implements the reconciliation loop
func (r *AgenticProposalReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the AgenticProposal instance
	var proposal agenticv1alpha1.AgenticProposal
	if err := r.Get(ctx, req.NamespacedName, &proposal); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("AgenticProposal resource not found, ignoring")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get AgenticProposal")
		return ctrl.Result{}, err
	}

	logger.Info("Reconciling AgenticProposal",
		"name", proposal.Name,
		"namespace", proposal.Namespace,
		"state", proposal.Status.State,
		"submitter", proposal.Spec.SubmittedBy)

	// Initialize status if needed
	if proposal.Status.State == "" {
		proposal.Status.State = agenticv1alpha1.StatePending
		proposal.Status.ObservedGeneration = proposal.Generation
		if err := r.Status().Update(ctx, &proposal); err != nil {
			logger.Error(err, "Failed to initialize status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// State machine logic
	switch proposal.Status.State {
	case agenticv1alpha1.StatePending:
		return r.handlePending(ctx, &proposal)
	case agenticv1alpha1.StateUnderReview:
		return r.handleUnderReview(ctx, &proposal)
	case agenticv1alpha1.StateApproved:
		return r.handleApproved(ctx, &proposal)
	case agenticv1alpha1.StateRejected:
		return r.handleRejected(ctx, &proposal)
	case agenticv1alpha1.StateExecuting:
		return r.handleExecuting(ctx, &proposal)
	case agenticv1alpha1.StateCompleted, agenticv1alpha1.StateFailed:
		// Terminal states, no action needed
		return ctrl.Result{}, nil
	default:
		logger.Info("Unknown state, resetting to Pending")
		proposal.Status.State = agenticv1alpha1.StatePending
		return r.updateStatus(ctx, &proposal)
	}
}

// SetupWithManager sets up the controller with the Manager
func (r *AgenticProposalReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agenticv1alpha1.AgenticProposal{}).
		Complete(r)
}

// handlePending transitions proposal from Pending to UnderReview
func (r *AgenticProposalReconciler) handlePending(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Processing pending proposal", "name", proposal.Name)

	// Transition to UnderReview
	proposal.Status.State = agenticv1alpha1.StateUnderReview
	if proposal.Status.Votes == nil {
		proposal.Status.Votes = make(map[string]string)
	}
	if proposal.Status.VoteDetails == nil {
		proposal.Status.VoteDetails = []agenticv1alpha1.AgentVote{}
	}
	proposal.Status.ObservedGeneration = proposal.Generation

	// Add condition
	now := metav1.Now()
	condition := metav1.Condition{
		Type:               "UnderReview",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "TransitionedFromPending",
		Message:            "Proposal is now under review by agents",
	}
	proposal.Status.Conditions = append(proposal.Status.Conditions, condition)

	return r.updateStatus(ctx, proposal)
}

// handleUnderReview calculates consensus and transitions to Approved/Rejected
func (r *AgenticProposalReconciler) handleUnderReview(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Calculate approval score from votes
	score := r.calculateConsensusScore(proposal)
	proposal.Status.ApprovalScore = &score

	logger.Info("Consensus score calculated",
		"score", score,
		"threshold", r.ConsensusThreshold,
		"votes", len(proposal.Status.VoteDetails))

	// Check if we have enough votes (require at least 3 for consensus)
	minVotes := 3
	if len(proposal.Status.VoteDetails) < minVotes {
		logger.Info("Waiting for more votes",
			"current", len(proposal.Status.VoteDetails),
			"minimum", minVotes)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Check required approvers if specified
	if len(proposal.Spec.RequiredApprovers) > 0 {
		if !r.checkRequiredApprovers(proposal) {
			logger.Info("Required approvers have not all approved")
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
	}

	now := metav1.Now()

	// Check if consensus reached
	if score >= r.ConsensusThreshold {
		proposal.Status.State = agenticv1alpha1.StateApproved
		proposal.Status.ApprovedAt = &now

		condition := metav1.Condition{
			Type:               "Approved",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: now,
			Reason:             "ConsensusReached",
			Message:            "Proposal approved by consensus",
		}
		proposal.Status.Conditions = append(proposal.Status.Conditions, condition)

		logger.Info("Proposal approved", "score", score)
	} else {
		proposal.Status.State = agenticv1alpha1.StateRejected
		proposal.Status.RejectedAt = &now

		condition := metav1.Condition{
			Type:               "Rejected",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: now,
			Reason:             "ConsensusNotReached",
			Message:            "Proposal rejected - insufficient approval",
		}
		proposal.Status.Conditions = append(proposal.Status.Conditions, condition)

		logger.Info("Proposal rejected", "score", score)
	}

	return r.updateStatus(ctx, proposal)
}

// handleApproved handles approved proposals
func (r *AgenticProposalReconciler) handleApproved(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Proposal approved", "name", proposal.Name)

	// Phase 1: Just log, execution deferred to Phase 2
	if !r.ExecutionEnabled {
		logger.Info("Execution disabled, skipping")
		return ctrl.Result{}, nil
	}

	// Future: Trigger execution workflow
	proposal.Status.State = agenticv1alpha1.StateExecuting
	now := metav1.Now()
	condition := metav1.Condition{
		Type:               "Executing",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "ExecutionStarted",
		Message:            "Proposal execution started",
	}
	proposal.Status.Conditions = append(proposal.Status.Conditions, condition)

	return r.updateStatus(ctx, proposal)
}

// handleRejected handles rejected proposals (terminal state)
func (r *AgenticProposalReconciler) handleRejected(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Proposal rejected", "name", proposal.Name)
	// Terminal state, no action needed
	return ctrl.Result{}, nil
}

// handleExecuting handles proposals that are being executed
func (r *AgenticProposalReconciler) handleExecuting(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Proposal executing", "name", proposal.Name)

	// Phase 1: Stub - just mark as completed
	// Phase 2: Integrate with actual execution engine

	now := metav1.Now()
	proposal.Status.State = agenticv1alpha1.StateCompleted
	proposal.Status.ExecutionResult = &agenticv1alpha1.ExecutionResult{
		Success:    true,
		Message:    "Execution completed (stub)",
		ExecutedAt: now,
		Artifacts:  map[string]string{},
	}

	condition := metav1.Condition{
		Type:               "Completed",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: now,
		Reason:             "ExecutionCompleted",
		Message:            "Proposal execution completed successfully",
	}
	proposal.Status.Conditions = append(proposal.Status.Conditions, condition)

	return r.updateStatus(ctx, proposal)
}

// calculateConsensusScore computes approval percentage from votes
func (r *AgenticProposalReconciler) calculateConsensusScore(proposal *agenticv1alpha1.AgenticProposal) float64 {
	if len(proposal.Status.VoteDetails) == 0 {
		return 0.0
	}

	approvals := 0
	totalVotes := 0

	for _, vote := range proposal.Status.VoteDetails {
		// Count Approve votes, ignore Abstain
		if vote.Decision == agenticv1alpha1.VoteApprove {
			approvals++
			totalVotes++
		} else if vote.Decision == agenticv1alpha1.VoteReject {
			totalVotes++
		}
		// Abstain doesn't count toward total
	}

	if totalVotes == 0 {
		return 0.0
	}

	return (float64(approvals) / float64(totalVotes)) * 100.0
}

// checkRequiredApprovers verifies all required approvers have voted Approve
func (r *AgenticProposalReconciler) checkRequiredApprovers(proposal *agenticv1alpha1.AgenticProposal) bool {
	approvedBy := make(map[string]bool)
	for _, vote := range proposal.Status.VoteDetails {
		if vote.Decision == agenticv1alpha1.VoteApprove {
			approvedBy[vote.Agent] = true
		}
	}

	for _, required := range proposal.Spec.RequiredApprovers {
		if !approvedBy[required] {
			return false
		}
	}

	return true
}

// updateStatus updates the proposal status
func (r *AgenticProposalReconciler) updateStatus(ctx context.Context, proposal *agenticv1alpha1.AgenticProposal) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if err := r.Status().Update(ctx, proposal); err != nil {
		logger.Error(err, "Failed to update AgenticProposal status")
		return ctrl.Result{}, err
	}

	logger.Info("Status updated successfully",
		"state", proposal.Status.State,
		"score", proposal.Status.ApprovalScore)

	return ctrl.Result{}, nil
}
```

```go
// api/v1alpha1/agenticproposal_types_test.go
package v1alpha1_test

import (
	"testing"

	agenticv1alpha1 "github.com/shreyanshjain7174/clawdlinux-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAgenticProposalValidation(t *testing.T) {
	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
		Spec: agenticv1alpha1.AgenticProposalSpec{
			Content:     "Test proposal content with minimum 10 chars",
			SubmittedBy: "agent-pm",
			ProposedAt:  metav1.Now(),
		},
	}

	// Basic sanity check
	if proposal.Spec.Content == "" {
		t.Error("Content should not be empty")
	}

	if len(proposal.Spec.Content) < 10 {
		t.Error("Content should be at least 10 characters")
	}

	if proposal.Spec.SubmittedBy == "" {
		t.Error("SubmittedBy should not be empty")
	}
}

func TestProposalStateConstants(t *testing.T) {
	states := []agenticv1alpha1.ProposalState{
		agenticv1alpha1.StatePending,
		agenticv1alpha1.StateUnderReview,
		agenticv1alpha1.StateApproved,
		agenticv1alpha1.StateRejected,
		agenticv1alpha1.StateExecuting,
		agenticv1alpha1.StateCompleted,
		agenticv1alpha1.StateFailed,
	}

	for _, state := range states {
		if state == "" {
			t.Errorf("State constant should not be empty")
		}
	}
}

func TestVoteDecisionConstants(t *testing.T) {
	decisions := []agenticv1alpha1.VoteDecision{
		agenticv1alpha1.VoteApprove,
		agenticv1alpha1.VoteReject,
		agenticv1alpha1.VoteAbstain,
	}

	for _, decision := range decisions {
		if decision == "" {
			t.Errorf("VoteDecision constant should not be empty")
		}
	}
}

func TestAgentVoteStructure(t *testing.T) {
	score := 85
	vote := agenticv1alpha1.AgentVote{
		Agent:     "agent-architect",
		Decision:  agenticv1alpha1.VoteApprove,
		Reasoning: "Architecture looks sound",
		VotedAt:   metav1.Now(),
		Score:     &score,
	}

	if vote.Agent == "" {
		t.Error("Agent name should not be empty")
	}

	if vote.Decision != agenticv1alpha1.VoteApprove {
		t.Error("Decision should be VoteApprove")
	}

	if vote.Score == nil || *vote.Score != 85 {
		t.Error("Score should be 85")
	}
}

func TestExecutionResultStructure(t *testing.T) {
	result := agenticv1alpha1.ExecutionResult{
		Success:    true,
		Message:    "Execution completed successfully",
		ExecutedAt: metav1.Now(),
		Artifacts: map[string]string{
			"commitSHA": "abc123",
			"prURL":     "https://github.com/org/repo/pull/42",
		},
	}

	if !result.Success {
		t.Error("Success should be true")
	}

	if result.Message == "" {
		t.Error("Message should not be empty")
	}

	if len(result.Artifacts) != 2 {
		t.Error("Should have 2 artifacts")
	}
}
```

```go
// cmd/main.go
package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	agenticv1alpha1 "github.com/shreyanshjain7174/clawdlinux-operator/api/v1alpha1"
	"github.com/shreyanshjain7174/clawdlinux-operator/internal/controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(agenticv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var consensusThreshold float64
	var executionEnabled bool

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Float64Var(&consensusThreshold, "consensus-threshold", 80.0,
		"Minimum approval percentage required for consensus (0-100)")
	flag.BoolVar(&executionEnabled, "execution-enabled", false,
		"Enable automatic execution of approved proposals")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "agentic-operator.clawdlinux.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.AgenticProposalReconciler{
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		ConsensusThreshold: consensusThreshold,
		ExecutionEnabled:   executionEnabled,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AgenticProposal")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager",
		"consensusThreshold", consensusThreshold,
		"executionEnabled", executionEnabled)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
```

```go
// internal/controller/suite_test.go
package controller

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	agenticv1alpha1 "github.com/shreyanshjain7174/clawdlinux-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = agenticv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&AgenticProposalReconciler{
		Client:             k8sManager.GetClient(),
		Scheme:             k8sManager.GetScheme(),
		ConsensusThreshold: 80.0,
		ExecutionEnabled:   false,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("AgenticProposal Controller", func() {
	const timeout = time.Second * 10
	const interval = time.Millisecond * 250

	Context("When creating a new AgenticProposal", func() {
		It("Should transition from Pending to UnderReview", func() {
			// Test implementation would go here
			// This is a placeholder for the test structure
		})
	})
})
```