// Copyright 2024 The Voting Operator Authors.
// Licensed under the Apache License, Version 2.0.

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
	"github.com/yourorg/voting-operator/internal/voting"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// VoteSubmissionHandler handles incoming vote submissions from voters
type VoteSubmissionHandler struct {
	Client            client.Client
	SignatureVerifier voting.SignatureVerifier
}

// VoteSubmissionRequest is the incoming vote from a voter
type VoteSubmissionRequest struct {
	ProposalID string                      `json:"proposalId"` // namespace/name
	Voter      string                      `json:"voter"`
	Decision   votingv1alpha1.VoteDecision `json:"decision"`
	Score      int                         `json:"score"`
	Rationale  string                      `json:"rationale"`
	Signature  string                      `json:"signature,omitempty"`
	Metadata   map[string]interface{}      `json:"metadata,omitempty"`
}

// ServeHTTP implements http.Handler for vote submissions
func (h *VoteSubmissionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.FromContext(ctx)

	// Only accept POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req VoteSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err, "failed to parse vote submission")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateSubmission(&req); err != nil {
		log.Error(err, "invalid vote submission")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse proposal ID (namespace/name)
	parts := strings.Split(req.ProposalID, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid proposalId format (expected: namespace/name)", http.StatusBadRequest)
		return
	}

	namespace, name := parts[0], parts[1]

	// Fetch proposal
	var proposal votingv1alpha1.VoteProposal
	if err := h.Client.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, &proposal); err != nil {
		log.Error(err, "proposal not found")
		http.Error(w, "Proposal not found", http.StatusNotFound)
		return
	}

	// Check if voting is still open
	if proposal.Status.Phase != votingv1alpha1.PhaseVoting {
		http.Error(w, fmt.Sprintf("Proposal not in voting phase (current: %s)", proposal.Status.Phase), http.StatusConflict)
		return
	}

	// Check deadline
	if proposal.Status.VotingDeadline != nil && time.Now().After(proposal.Status.VotingDeadline.Time) {
		http.Error(w, "Voting deadline has passed", http.StatusGone)
		return
	}

	// Check if voter already voted
	for _, existingVote := range proposal.Status.Votes {
		if existingVote.Voter == req.Voter {
			http.Error(w, "Voter has already submitted a vote", http.StatusConflict)
			return
		}
	}

	// Build vote object
	vote := votingv1alpha1.Vote{
		Voter:     req.Voter,
		Decision:  req.Decision,
		Score:     req.Score,
		Timestamp: metav1.Now(),
		Rationale: req.Rationale,
		Signature: req.Signature,
	}

	// Verify signature if provided
	if req.Signature != "" {
		valid, err := h.SignatureVerifier.VerifyVote(ctx, &vote, &proposal)
		if err != nil || !valid {
			log.Error(err, "signature verification failed", "voter", req.Voter)
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Add vote to proposal status using retry on conflict
	err := h.addVoteWithRetry(ctx, &proposal, vote)
	if err != nil {
		log.Error(err, "failed to add vote")
		http.Error(w, "Failed to record vote", http.StatusInternalServerError)
		return
	}

	log.Info("vote recorded successfully",
		"voter", req.Voter,
		"decision", req.Decision,
		"proposal", req.ProposalID)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "accepted",
		"message": "Vote recorded successfully",
	})
}

// addVoteWithRetry adds a vote with optimistic locking retry
func (h *VoteSubmissionHandler) addVoteWithRetry(ctx context.Context, proposal *votingv1alpha1.VoteProposal, vote votingv1alpha1.Vote) error {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		// Add vote to status
		proposal.Status.Votes = append(proposal.Status.Votes, vote)

		// Try to update
		err := h.Client.Status().Update(ctx, proposal)
		if err == nil {
			return nil // Success
		}

		// Check if conflict error
		if !strings.Contains(err.Error(), "conflict") {
			return err // Non-conflict error, fail
		}

		// Conflict - refetch and retry
		if err := h.Client.Get(ctx, types.NamespacedName{
			Namespace: proposal.Namespace,
			Name:      proposal.Name,
		}, proposal); err != nil {
			return err
		}

		// Check if vote was already added by another request
		for _, existingVote := range proposal.Status.Votes {
			if existingVote.Voter == vote.Voter {
				return nil // Vote already recorded, success
			}
		}

		time.Sleep(time.Duration(i*100) * time.Millisecond) // Exponential backoff
	}

	return fmt.Errorf("failed to add vote after %d retries", maxRetries)
}

// validateSubmission validates the vote submission request
func (h *VoteSubmissionHandler) validateSubmission(req *VoteSubmissionRequest) error {
	if req.ProposalID == "" {
		return fmt.Errorf("proposalId is required")
	}
	if req.Voter == "" {
		return fmt.Errorf("voter is required")
	}
	if req.Decision == "" {
		return fmt.Errorf("decision is required")
	}
	if req.Score < 0 || req.Score > 100 {
		return fmt.Errorf("score must be between 0 and 100")
	}
	if req.Rationale == "" {
		return fmt.Errorf("rationale is required")
	}

	// Validate decision is valid enum
	switch req.Decision {
	case votingv1alpha1.VoteApprove, votingv1alpha1.VoteConditionalApprove, votingv1alpha1.VoteReject, votingv1alpha1.VoteAbstain:
		// Valid
	default:
		return fmt.Errorf("invalid decision: %s", req.Decision)
	}

	return nil
}
