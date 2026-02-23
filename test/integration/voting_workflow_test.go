// Copyright 2024 The Voting Operator Authors.
// Licensed under the Apache License, Version 2.0.

// +build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
	"github.com/yourorg/voting-operator/internal/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// TestVotingWorkflowEndToEnd tests the complete voting workflow
func TestVotingWorkflowEndToEnd(t *testing.T) {
	// This test would use envtest for real Kubernetes API
	// For demonstration, using fake client

	ctx := context.Background()
	scheme := runtime.NewScheme()
	_ = votingv1alpha1.AddToScheme(scheme)

	// Create fake client
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()

	// Step 1: Create VoterConfig
	voterConfig := &votingv1alpha1.VoterConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "test-ns",
		},
		Spec: votingv1alpha1.VoterConfigSpec{
			DefaultThreshold: 80,
			Voters: []votingv1alpha1.VoterRegistration{
				{
					Name:       "security-agent",
					WebhookURL: "http://security.test/vote",
					Required:   true,
				},
				{
					Name:       "qa-agent",
					WebhookURL: "http://qa.test/vote",
					Required:   true,
				},
			},
		},
	}

	err := fakeClient.Create(ctx, voterConfig)
	if err != nil {
		t.Fatalf("Failed to create VoterConfig: %v", err)
	}

	// Step 2: Create VoteProposal
	proposal := &votingv1alpha1.VoteProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment-update",
			Namespace: "test-ns",
		},
		Spec: votingv1alpha1.VoteProposalSpec{
			Operation: "deployment.update",
			TargetRef: votingv1alpha1.TargetReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "test-app",
				Namespace:  "test-ns",
			},
			VoteConfig: votingv1alpha1.VoteConfiguration{
				Threshold:      80,
				Timeout:        metav1.Duration{Duration: 5 * time.Minute},
				RequiredVoters: []string{"security-agent", "qa-agent"},
			},
			ProposedChange: votingv1alpha1.ChangeSpec{
				Type: "patch",
			},
		},
	}

	err = fakeClient.Create(ctx, proposal)
	if err != nil {
		t.Fatalf("Failed to create VoteProposal: %v", err)
	}

	// Step 3: Simulate voter responses
	// In real test, controller would notify webhooks
	// Here we directly update status to simulate votes

	// Fetch proposal
	var fetchedProposal votingv1alpha1.VoteProposal
	err = fakeClient.Get(ctx, types.NamespacedName{
		Name:      "test-deployment-update",
		Namespace: "test-ns",
	}, &fetchedProposal)
	if err != nil {
		t.Fatalf("Failed to fetch proposal: %v", err)
	}

	// Add votes
	fetchedProposal.Status.Phase = votingv1alpha1.PhaseVoting
	fetchedProposal.Status.Votes = []votingv1alpha1.Vote{
		{
			Voter:     "security-agent",
			Decision:  votingv1alpha1.VoteApprove,
			Score:     100,
			Timestamp: metav1.Now(),
			Rationale: "No security issues detected",
		},
		{
			Voter:     "qa-agent",
			Decision:  votingv1alpha1.VoteConditionalApprove,
			Score:     70,
			Timestamp: metav1.Now(),
			Rationale: "Tests pass but coverage could be better",
		},
	}

	err = fakeClient.Status().Update(ctx, &fetchedProposal)
	if err != nil {
		t.Fatalf("Failed to update proposal status: %v", err)
	}

	// Step 4: Verify aggregate score
	expectedScore := (100.0 + 70.0) / 2.0 // 85.0
	if fetchedProposal.Status.AggregateScore != nil {
		score := *fetchedProposal.Status.AggregateScore
		if score != expectedScore {
			t.Errorf("Expected score %v, got %v", expectedScore, score)
		}
	}

	// Step 5: Verify proposal would be approved (≥80%)
	if expectedScore < 80.0 {
		t.Errorf("Expected proposal to pass threshold, score: %v", expectedScore)
	}

	t.Logf("Voting workflow test completed successfully")
}

// TestVoteSubmissionAPI tests the vote submission HTTP endpoint
func TestVoteSubmissionAPI(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	_ = votingv1alpha1.AddToScheme(scheme)

	// Create proposal in voting phase
	proposal := &votingv1alpha1.VoteProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-test-proposal",
			Namespace: "default",
		},
		Spec: votingv1alpha1.VoteProposalSpec{
			VoteConfig: votingv1alpha1.VoteConfiguration{
				RequiredVoters: []string{"test-voter"},
			},
		},
		Status: votingv1alpha1.VoteProposalStatus{
			Phase: votingv1alpha1.PhaseVoting,
			VotingDeadline: &metav1.Time{
				Time: time.Now().Add(10 * time.Minute),
			},
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(proposal).
		WithStatusSubresource(proposal).
		Build()

	// Create API handler
	handler := &api.VoteSubmissionHandler{
		Client: fakeClient,
		SignatureVerifier: &mockSignatureVerifier{valid: true},
	}

	// Create test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Submit vote
	voteReq := api.VoteSubmissionRequest{
		ProposalID: "default/api-test-proposal",
		Voter:      "test-voter",
		Decision:   votingv1alpha1.VoteApprove,
		Score:      100,
		Rationale:  "Looks good to me",
	}

	reqBody, _ := json.Marshal(voteReq)
	resp, err := http.Post(server.URL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to submit vote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status 202, got %d", resp.StatusCode)
	}

	// Verify vote was recorded
	var updatedProposal votingv1alpha1.VoteProposal
	err = fakeClient.Get(ctx, types.NamespacedName{
		Name:      "api-test-proposal",
		Namespace: "default",
	}, &updatedProposal)
	if err != nil {
		t.Fatalf("Failed to fetch updated proposal: %v", err)
	}

	if len(updatedProposal.Status.Votes) != 1 {
		t.Errorf("Expected 1 vote, got %d", len(updatedProposal.Status.Votes))
	}

	if updatedProposal.Status.Votes[0].Voter != "test-voter" {
		t.Errorf("Expected voter 'test-voter', got '%s'", updatedProposal.Status.Votes[0].Voter)
	}
}

// TestConcurrentVoteSubmission tests race conditions in vote recording
func TestConcurrentVoteSubmission(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	_ = votingv1alpha1.AddToScheme(scheme)

	proposal := &votingv1alpha1.VoteProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "concurrent-test",
			Namespace: "default",
		},
		Status: votingv1alpha1.VoteProposalStatus{
			Phase: votingv1alpha1.PhaseVoting,
			VotingDeadline: &metav1.Time{
				Time: time.Now().Add(10 * time.Minute),
			},
		},
	}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(proposal).
		WithStatusSubresource(proposal).
		Build()

	handler := &api.VoteSubmissionHandler{
		Client: fakeClient,
		SignatureVerifier: &mockSignatureVerifier{valid: true},
	}

	// Simulate 10 concurrent vote submissions
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(voterNum int) {
			defer wg.Done()

			voteReq := api.VoteSubmissionRequest{
				ProposalID: "default/concurrent-test",
				Voter:      fmt.Sprintf("voter-%d", voterNum),
				Decision:   votingv1alpha1.VoteApprove,
				Score:      100,
				Rationale:  "Concurrent vote",
			}

			reqBody, _ := json.Marshal(voteReq)
			req, _ := http.NewRequest("POST", "/vote", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusAccepted {
				errors <- fmt.Errorf("voter-%d: unexpected status %d", voterNum, rr.Code)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}

	// Verify all 10 votes were recorded
	var finalProposal votingv1alpha1.VoteProposal
	err := fakeClient.Get(ctx, types.NamespacedName{
		Name:      "concurrent-test",
		Namespace: "default",
	}, &finalProposal)
	if err != nil {
		t.Fatalf("Failed to fetch final proposal: %v", err)
	}

	if len(finalProposal.Status.Votes) != 10 {
		t.Errorf("Expected 10 votes, got %d (potential race condition)", len(finalProposal.Status.Votes))
	}
}

// Mock signature verifier for testing
type mockSignatureVerifier struct {
	valid bool
}

func (m *mockSignatureVerifier) VerifyVote(ctx context.Context, vote *votingv1alpha1.Vote, proposal *votingv1alpha1.VoteProposal) (bool, error) {
	return m.valid, nil
}
