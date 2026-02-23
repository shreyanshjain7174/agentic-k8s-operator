/*
Copyright 2026 ClawdLinux.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AgenticProposalSpec defines the desired state of AgenticProposal
type AgenticProposalSpec struct {
	// Title is a short, human-readable summary of the proposal
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:MaxLength=200
	Title string `json:"title"`

	// Description is the full proposal content (markdown supported)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=50
	Description string `json:"description"`

	// ProposedBy identifies the agent that created this proposal
	// +kubebuilder:validation:Required
	ProposedBy ProposerInfo `json:"proposedBy"`

	// ConsensusThreshold is the minimum weighted score (0.0-1.0) required for approval
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0.5
	// +kubebuilder:validation:Maximum=1.0
	ConsensusThreshold float64 `json:"consensusThreshold"`

	// Voters is the list of agents authorized to vote on this proposal
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Voters []VoterConfig `json:"voters"`

	// ExecutionPlan defines how the approved proposal will be executed
	// +kubebuilder:validation:Required
	ExecutionPlan ExecutionPlanSpec `json:"executionPlan"`

	// Category is a classification tag (e.g., "feature", "bugfix", "refactor")
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="general"
	Category string `json:"category,omitempty"`

	// Tags are arbitrary labels for filtering/grouping
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty"`

	// RelatedProposals are references to dependent/parent proposals
	// +kubebuilder:validation:Optional
	RelatedProposals []string `json:"relatedProposals,omitempty"`
}

// ProposerInfo captures agent identity and timestamp
type ProposerInfo struct {
	// Agent is the unique identifier of the proposing agent (e.g., "architect", "pm")
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Agent string `json:"agent"`

	// Timestamp is when the proposal was created (auto-populated by webhook)
	// +kubebuilder:validation:Optional
	Timestamp *metav1.Time `json:"timestamp,omitempty"`
}

// VoterConfig defines an agent's voting weight
type VoterConfig struct {
	// Agent is the unique identifier of the voter
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Agent string `json:"agent"`

	// Weight is the multiplier for this agent's vote (default: 1.0)
	// +kubebuilder:validation:Minimum=0.1
	// +kubebuilder:validation:Maximum=5.0
	// +kubebuilder:default=1.0
	Weight float64 `json:"weight"`
}

// ExecutionPlanSpec defines the execution strategy
type ExecutionPlanSpec struct {
	// Type determines the execution engine to use
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=code-change;gitops-sync;webhook;script;manual
	Type string `json:"type"`

	// TargetFiles are the files affected by this proposal (for code-change type)
	// +kubebuilder:validation:Optional
	TargetFiles []string `json:"targetFiles,omitempty"`

	// EstimatedImpact is a risk assessment (low/medium/high)
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=low;medium;high
	EstimatedImpact string `json:"estimatedImpact,omitempty"`

	// RollbackPlan is human-readable instructions for reverting changes
	// +kubebuilder:validation:Optional
	RollbackPlan string `json:"rollbackPlan,omitempty"`

	// Parameters are execution-specific key-value pairs (JSON)
	// +kubebuilder:validation:Optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Parameters map[string]string `json:"parameters,omitempty"`
}

// AgenticProposalStatus defines the observed state of AgenticProposal
type AgenticProposalStatus struct {
	// Phase is the current state in the proposal lifecycle
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=Pending
	Phase ProposalPhase `json:"phase,omitempty"`

	// Votes is the accumulated list of agent decisions
	// +kubebuilder:validation:Optional
	Votes []Vote `json:"votes,omitempty"`

	// ConsensusScore is the weighted average of all votes (0-100)
	// +kubebuilder:validation:Optional
	ConsensusScore *float64 `json:"consensusScore,omitempty"`

	// ConsensusReached indicates if threshold was met
	// +kubebuilder:validation:Optional
	ConsensusReached bool `json:"consensusReached"`

	// ConsensusTimestamp is when consensus was achieved
	// +kubebuilder:validation:Optional
	ConsensusTimestamp *metav1.Time `json:"consensusTimestamp,omitempty"`

	// ExecutionStartTime marks when execution began
	// +kubebuilder:validation:Optional
	ExecutionStartTime *metav1.Time `json:"executionStartTime,omitempty"`

	// ExecutionEndTime marks when execution completed/failed
	// +kubebuilder:validation:Optional
	ExecutionEndTime *metav1.Time `json:"executionEndTime,omitempty"`

	// ExecutionResult captures the outcome of execution
	// +kubebuilder:validation:Optional
	ExecutionResult *ExecutionResult `json:"executionResult,omitempty"`

	// Conditions are K8s standard status conditions
	// +kubebuilder:validation:Optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// ProposalPhase represents the lifecycle state
// +kubebuilder:validation:Enum=Pending;InReview;Calculating;Approved;Rejected;Executing;Completed;Failed
type ProposalPhase string

const (
	PhasePending     ProposalPhase = "Pending"
	PhaseInReview    ProposalPhase = "InReview"
	PhaseCalculating ProposalPhase = "Calculating"
	PhaseApproved    ProposalPhase = "Approved"
	PhaseRejected    ProposalPhase = "Rejected"
	PhaseExecuting   ProposalPhase = "Executing"
	PhaseCompleted   ProposalPhase = "Completed"
	PhaseFailed      ProposalPhase = "Failed"
)

// Vote captures an agent's decision
type Vote struct {
	// Agent is the voter's identifier
	// +kubebuilder:validation:Required
	Agent string `json:"agent"`

	// Decision is APPROVE, CONDITIONAL_APPROVE, or REJECT
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=APPROVE;CONDITIONAL_APPROVE;REJECT
	Decision VoteDecision `json:"decision"`

	// Score is the numerical rating (0-100)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Score int `json:"score"`

	// Feedback is human-readable justification
	// +kubebuilder:validation:Optional
	Feedback string `json:"feedback,omitempty"`

	// Timestamp is when the vote was cast
	// +kubebuilder:validation:Required
	Timestamp metav1.Time `json:"timestamp"`
}

// VoteDecision enum
// +kubebuilder:validation:Enum=APPROVE;CONDITIONAL_APPROVE;REJECT
type VoteDecision string

const (
	VoteApprove            VoteDecision = "APPROVE"
	VoteConditionalApprove VoteDecision = "CONDITIONAL_APPROVE"
	VoteReject             VoteDecision = "REJECT"
)

// ExecutionResult captures execution outcome
type ExecutionResult struct {
	// Success indicates if execution completed without errors
	// +kubebuilder:validation:Required
	Success bool `json:"success"`

	// GitCommit is the resulting commit SHA (for code-change type)
	// +kubebuilder:validation:Optional
	GitCommit string `json:"gitCommit,omitempty"`

	// ArtifactsCreated lists files modified/created
	// +kubebuilder:validation:Optional
	ArtifactsCreated []Artifact `json:"artifactsCreated,omitempty"`

	// ErrorMessage captures failure reason
	// +kubebuilder:validation:Optional
	ErrorMessage string `json:"errorMessage,omitempty"`

	// Logs are execution logs (truncated to 10KB)
	// +kubebuilder:validation:Optional
	Logs string `json:"logs,omitempty"`
}

// Artifact represents a created/modified file
type Artifact struct {
	// File is the path relative to repo root
	// +kubebuilder:validation:Required
	File string `json:"file"`

	// LinesChanged is the diff size
	// +kubebuilder:validation:Optional
	LinesChanged int `json:"linesChanged,omitempty"`

	// Operation is "created", "modified", or "deleted"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=created;modified;deleted
	Operation string `json:"operation,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=proposal;proposals
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Consensus",type=string,JSONPath=`.status.consensusScore`,priority=1
// +kubebuilder:printcolumn:name="Threshold",type=string,JSONPath=`.spec.consensusThreshold`,priority=1
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
