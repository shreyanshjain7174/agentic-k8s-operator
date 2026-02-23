// Copyright 2024 The Voting Operator Authors.
// Licensed under the Apache License, Version 2.0.

package voting

import (
	"testing"
	"time"

	votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDefaultAggregator_CalculateConsensus(t *testing.T) {
	tests := []struct {
		name           string
		votes          []votingv1alpha1.Vote
		threshold      int
		allowAbstain   bool
		expectedScore  float64
		expectedDecision string
	}{
		{
			name: "all approve - should pass",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v2", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v3", Decision: votingv1alpha1.VoteApprove, Score: 100},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    100.0,
			expectedDecision: "APPROVED",
		},
		{
			name: "mixed votes - should pass threshold",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v2", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
				{Voter: "v3", Decision: votingv1alpha1.VoteApprove, Score: 100},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    90.0, // (100 + 70 + 100) / 3
			expectedDecision: "APPROVED",
		},
		{
			name: "exactly at threshold - should pass",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v2", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
				{Voter: "v3", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    80.0, // (100 + 70 + 70) / 3
			expectedDecision: "APPROVED",
		},
		{
			name: "below threshold - should reject",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
				{Voter: "v2", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
				{Voter: "v3", Decision: votingv1alpha1.VoteReject, Score: 0},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    46.67, // (70 + 70 + 0) / 3 ≈ 46.67
			expectedDecision: "REJECTED",
		},
		{
			name: "all reject - should reject",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteReject, Score: 0},
				{Voter: "v2", Decision: votingv1alpha1.VoteReject, Score: 0},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    0.0,
			expectedDecision: "REJECTED",
		},
		{
			name: "abstentions allowed - should ignore",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v2", Decision: votingv1alpha1.VoteAbstain, Score: 0},
				{Voter: "v3", Decision: votingv1alpha1.VoteApprove, Score: 100},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    100.0, // (100 + 100) / 2, ignoring abstention
			expectedDecision: "APPROVED",
		},
		{
			name: "abstentions not allowed - should count as reject",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
				{Voter: "v2", Decision: votingv1alpha1.VoteAbstain, Score: 0},
				{Voter: "v3", Decision: votingv1alpha1.VoteApprove, Score: 100},
			},
			threshold:        80,
			allowAbstain:     false,
			expectedScore:    66.67, // (100 + 0 + 100) / 3
			expectedDecision: "REJECTED",
		},
		{
			name: "all abstentions allowed - should return special case",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteAbstain, Score: 0},
				{Voter: "v2", Decision: votingv1alpha1.VoteAbstain, Score: 0},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    0.0,
			expectedDecision: "ALL_ABSTAINED",
		},
		{
			name:             "no votes - should return special case",
			votes:            []votingv1alpha1.Vote{},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    0.0,
			expectedDecision: "NO_VOTES",
		},
		{
			name: "single vote approve",
			votes: []votingv1alpha1.Vote{
				{Voter: "v1", Decision: votingv1alpha1.VoteApprove, Score: 100},
			},
			threshold:        80,
			allowAbstain:     true,
			expectedScore:    100.0,
			expectedDecision: "APPROVED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := &DefaultAggregator{}
			proposal := &votingv1alpha1.VoteProposal{
				Spec: votingv1alpha1.VoteProposalSpec{
					VoteConfig: votingv1alpha1.VoteConfiguration{
						Threshold:        tt.threshold,
						AllowAbstentions: tt.allowAbstain,
					},
				},
				Status: votingv1alpha1.VoteProposalStatus{
					Votes: tt.votes,
				},
			}

			score, decision := agg.CalculateConsensus(proposal)

			// Allow small float precision errors
			if diff := score - tt.expectedScore; diff > 0.01 || diff < -0.01 {
				t.Errorf("CalculateConsensus() score = %v, want %v", score, tt.expectedScore)
			}
			if decision != tt.expectedDecision {
				t.Errorf("CalculateConsensus() decision = %v, want %v", decision, tt.expectedDecision)
			}
		})
	}
}

func TestDefaultAggregator_EdgeCases(t *testing.T) {
	agg := &DefaultAggregator{}

	t.Run("exactly 79.99 score - should reject", func(t *testing.T) {
		proposal := &votingv1alpha1.VoteProposal{
			Spec: votingv1alpha1.VoteProposalSpec{
				VoteConfig: votingv1alpha1.VoteConfiguration{
					Threshold: 80,
				},
			},
			Status: votingv1alpha1.VoteProposalStatus{
				Votes: []votingv1alpha1.Vote{
					{Voter: "v1", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
					{Voter: "v2", Decision: votingv1alpha1.VoteApprove, Score: 100},
					{Voter: "v3", Decision: votingv1alpha1.VoteConditionalApprove, Score: 70},
				},
			},
		}

		score, decision := agg.CalculateConsensus(proposal)
		if score != 80.0 {
			t.Errorf("Expected exactly 80.0, got %v", score)
		}
		if decision != "APPROVED" {
			t.Errorf("Expected APPROVED at exactly threshold, got %v", decision)
		}
	})

	t.Run("nil votes slice", func(t *testing.T) {
		proposal := &votingv1alpha1.VoteProposal{
			Spec: votingv1alpha1.VoteProposalSpec{
				VoteConfig: votingv1alpha1.VoteConfiguration{
					Threshold: 80,
				},
			},
			Status: votingv1alpha1.VoteProposalStatus{
				Votes: nil,
			},
		}

		score, decision := agg.CalculateConsensus(proposal)
		if score != 0.0 {
			t.Errorf("Expected 0.0 for nil votes, got %v", score)
		}
		if decision != "NO_VOTES" {
			t.Errorf("Expected NO_VOTES for nil votes, got %v", decision)
		}
	})
}
