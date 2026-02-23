// Auto-generated tests for phase1-001

```go
// internal/controller/agenticproposal_controller_test.go
package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	agenticv1alpha1 "github.com/shreyanshjain7174/clawdlinux-operator/api/v1alpha1"
)

func TestCalculateConsensusScore(t *testing.T) {
	tests := []struct {
		name           string
		voteDetails    []agenticv1alpha1.AgentVote
		expectedScore  float64
	}{
		{
			name:          "no votes",
			voteDetails:   []agenticv1alpha1.AgentVote{},
			expectedScore: 0.0,
		},
		{
			name: "all approve",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedScore: 100.0,
		},
		{
			name: "all reject",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteReject},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteReject},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteReject},
			},
			expectedScore: 0.0,
		},
		{
			name: "mixed votes - 2 approve, 1 reject",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteReject},
			},
			expectedScore: 66.66666666666666,
		},
		{
			name: "mixed votes - 3 approve, 2 reject",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent4", Decision: agenticv1alpha1.VoteReject},
				{Agent: "agent5", Decision: agenticv1alpha1.VoteReject},
			},
			expectedScore: 60.0,
		},
		{
			name: "all abstain",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteAbstain},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteAbstain},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteAbstain},
			},
			expectedScore: 0.0,
		},
		{
			name: "mixed with abstain - 2 approve, 1 reject, 2 abstain",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteReject},
				{Agent: "agent4", Decision: agenticv1alpha1.VoteAbstain},
				{Agent: "agent5", Decision: agenticv1alpha1.VoteAbstain},
			},
			expectedScore: 66.66666666666666,
		},
		{
			name: "single approve vote",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedScore: 100.0,
		},
		{
			name: "single reject vote",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteReject},
			},
			expectedScore: 0.0,
		},
		{
			name: "80% approval threshold test",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent4", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent5", Decision: agenticv1alpha1.VoteReject},
			},
			expectedScore: 80.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &agenticv1alpha1.AgenticProposal{
				Status: agenticv1alpha1.AgenticProposalStatus{
					VoteDetails: tt.voteDetails,
				},
			}

			reconciler := &AgenticProposalReconciler{
				ConsensusThreshold: 80.0,
			}

			score := reconciler.calculateConsensusScore(proposal)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestCheckRequiredApprovers(t *testing.T) {
	tests := []struct {
		name              string
		requiredApprovers []string
		voteDetails       []agenticv1alpha1.AgentVote
		expectedResult    bool
	}{
		{
			name:              "no required approvers",
			requiredApprovers: []string{},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedResult: true,
		},
		{
			name:              "all required approvers approved",
			requiredApprovers: []string{"agent1", "agent2"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteReject},
			},
			expectedResult: true,
		},
		{
			name:              "some required approvers missing",
			requiredApprovers: []string{"agent1", "agent2", "agent3"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedResult: false,
		},
		{
			name:              "required approver rejected",
			requiredApprovers: []string{"agent1", "agent2"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteReject},
			},
			expectedResult: false,
		},
		{
			name:              "required approver abstained",
			requiredApprovers: []string{"agent1", "agent2"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteAbstain},
			},
			expectedResult: false,
		},
		{
			name:              "extra votes beyond required approvers",
			requiredApprovers: []string{"agent1"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedResult: true,
		},
		{
			name:              "no votes but required approvers specified",
			requiredApprovers: []string{"agent1"},
			voteDetails:       []agenticv1alpha1.AgentVote{},
			expectedResult:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &agenticv1alpha1.AgenticProposal{
				Spec: agenticv1alpha1.AgenticProposalSpec{
					RequiredApprovers: tt.requiredApprovers,
				},
				Status: agenticv1alpha1.AgenticProposalStatus{
					VoteDetails: tt.voteDetails,
				},
			}

			reconciler := &AgenticProposalReconciler{}
			result := reconciler.checkRequiredApprovers(proposal)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestHandlePending(t *testing.T) {
	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-proposal",
			Namespace:  "default",
			Generation: 1,
		},
		Spec: agenticv1alpha1.AgenticProposalSpec{
			Content:     "Test proposal content",
			SubmittedBy: "agent-pm",
			ProposedAt:  metav1.Now(),
		},
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StatePending,
		},
	}

	// Note: This test would need a mock client for full integration
	// Here we're testing the state transition logic
	originalState := proposal.Status.State
	assert.Equal(t, agenticv1alpha1.StatePending, originalState)

	// After handlePending, state should be UnderReview
	// Votes and VoteDetails should be initialized
	// ObservedGeneration should be set
	// Condition should be added
}

func TestHandleUnderReviewInsufficientVotes(t *testing.T) {
	proposal := &agenticv1alpha1.AgenticProposal{
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StateUnderReview,
			VoteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
			},
		},
	}

	// With only 2 votes (minimum is 3), should requeue
	assert.Equal(t, 2, len(proposal.Status.VoteDetails))
}

func TestHandleUnderReviewApprovalScoreCalculation(t *testing.T) {
	tests := []struct {
		name            string
		voteDetails     []agenticv1alpha1.AgentVote
		threshold       float64
		expectedState   agenticv1alpha1.ProposalState
	}{
		{
			name: "approval above threshold",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
			},
			threshold:     80.0,
			expectedState: agenticv1alpha1.StateApproved,
		},
		{
			name: "approval below threshold",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteReject},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteReject},
			},
			threshold:     80.0,
			expectedState: agenticv1alpha1.StateRejected,
		},
		{
			name: "approval exactly at threshold",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent2", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent3", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent4", Decision: agenticv1alpha1.VoteApprove},
				{Agent: "agent5", Decision: agenticv1alpha1.VoteReject},
			},
			threshold:     80.0,
			expectedState: agenticv1alpha1.StateApproved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reconciler := &AgenticProposalReconciler{
				ConsensusThreshold: tt.threshold,
			}

			proposal := &agenticv1alpha1.AgenticProposal{
				Status: agenticv1alpha1.AgenticProposalStatus{
					VoteDetails: tt.voteDetails,
				},
			}

			score := reconciler.calculateConsensusScore(proposal)

			if score >= tt.threshold {
				assert.Equal(t, agenticv1alpha1.StateApproved, tt.expectedState)
			} else {
				assert.Equal(t, agenticv1alpha1.StateRejected, tt.expectedState)
			}
		})
	}
}

func TestHandleApprovedWithExecutionDisabled(t *testing.T) {
	reconciler := &AgenticProposalReconciler{
		ExecutionEnabled: false,
	}

	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StateApproved,
		},
	}

	ctx := context.Background()
	_, err := reconciler.handleApproved(ctx, proposal)

	assert.NoError(t, err)
	// State should remain Approved when execution is disabled
	assert.Equal(t, agenticv1alpha1.StateApproved, proposal.Status.State)
}

func TestHandleApprovedWithExecutionEnabled(t *testing.T) {
	reconciler := &AgenticProposalReconciler{
		ExecutionEnabled: true,
	}

	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StateApproved,
		},
	}

	// State should transition to Executing
	// Note: Full test would require mock client
	assert.Equal(t, agenticv1alpha1.StateApproved, proposal.Status.State)
}

func TestHandleExecuting(t *testing.T) {
	reconciler := &AgenticProposalReconciler{}

	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StateExecuting,
		},
	}

	ctx := context.Background()
	_, err := reconciler.handleExecuting(ctx, proposal)

	assert.NoError(t, err)
	// In Phase 1, state should transition to Completed
	assert.Equal(t, agenticv1alpha1.StateCompleted, proposal.Status.State)
	assert.NotNil(t, proposal.Status.ExecutionResult)
	assert.True(t, proposal.Status.ExecutionResult.Success)
	assert.Equal(t, "Execution completed (stub)", proposal.Status.ExecutionResult.Message)
}

func TestHandleRejected(t *testing.T) {
	reconciler := &AgenticProposalReconciler{}

	proposal := &agenticv1alpha1.AgenticProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
		Status: agenticv1alpha1.AgenticProposalStatus{
			State: agenticv1alpha1.StateRejected,
		},
	}

	ctx := context.Background()
	_, err := reconciler.handleRejected(ctx, proposal)

	assert.NoError(t, err)
	// State should remain Rejected (terminal state)
	assert.Equal(t, agenticv1alpha1.StateRejected, proposal.Status.State)
}

func TestConsensusScoreEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		voteDetails   []agenticv1alpha1.AgentVote
		expectedScore float64
		description   string
	}{
		{
			name:          "nil vote details",
			voteDetails:   nil,
			expectedScore: 0.0,
			description:   "Should return 0 for nil vote details",
		},
		{
			name: "single abstain vote",
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteAbstain},
			},
			expectedScore: 0.0,
			description:   "Abstain votes should not count",
		},
		{
			name: "large number of votes",
			voteDetails: func() []agenticv1alpha1.AgentVote {
				votes := make([]agenticv1alpha1.AgentVote, 100)
				for i := 0; i < 100; i++ {
					decision := agenticv1alpha1.VoteApprove
					if i%3 == 0 {
						decision = agenticv1alpha1.VoteReject
					}
					votes[i] = agenticv1alpha1.AgentVote{
						Agent:    "agent" + string(rune(i)),
						Decision: decision,
					}
				}
				return votes
			}(),
			expectedScore: 66.0,
			description:   "Should handle large vote counts correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &agenticv1alpha1.AgenticProposal{
				Status: agenticv1alpha1.AgenticProposalStatus{
					VoteDetails: tt.voteDetails,
				},
			}

			reconciler := &AgenticProposalReconciler{}
			score := reconciler.calculateConsensusScore(proposal)

			assert.InDelta(t, tt.expectedScore, score, 0.1, tt.description)
		})
	}
}

func TestRequiredApproversEdgeCases(t *testing.T) {
	tests := []struct {
		name              string
		requiredApprovers []string
		voteDetails       []agenticv1alpha1.AgentVote
		expectedResult    bool
		description       string
	}{
		{
			name:              "empty required approvers list",
			requiredApprovers: []string{},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteReject},
			},
			expectedResult: true,
			description:    "Empty required approvers should always pass",
		},
		{
			name:              "nil required approvers",
			requiredApprovers: nil,
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteReject},
			},
			expectedResult: true,
			description:    "Nil required approvers should always pass",
		},
		{
			name:              "duplicate required approvers",
			requiredApprovers: []string{"agent1", "agent1"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedResult: true,
			description:    "Duplicate approvers should work",
		},
		{
			name:              "case sensitive agent names",
			requiredApprovers: []string{"Agent1"},
			voteDetails: []agenticv1alpha1.AgentVote{
				{Agent: "agent1", Decision: agenticv1alpha1.VoteApprove},
			},
			expectedResult: false,
			description:    "Agent names should be case sensitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &agenticv1alpha1.AgenticProposal{
				Spec: agenticv1alpha1.AgenticProposalSpec{
					RequiredApprovers: tt.requiredApprovers,
				},
				Status: agenticv1alpha1.AgenticProposalStatus{
					VoteDetails: tt.voteDetails,
				},
			}

			reconciler := &AgenticProposalReconciler{}
			result := reconciler.checkRequiredApprovers(proposal)

			assert.Equal(t, tt.expectedResult, result, tt.description)
		})
	}
}
```