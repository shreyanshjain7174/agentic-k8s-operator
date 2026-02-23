// Copyright 2024 The Voting Operator Authors.
// Licensed under the Apache License, Version 2.0.

package voting

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockVoterRegistry for testing
type MockVoterRegistry struct {
	voters map[string]*votingv1alpha1.VoterRegistration
}

func NewMockVoterRegistry() *MockVoterRegistry {
	return &MockVoterRegistry{
		voters: make(map[string]*votingv1alpha1.VoterRegistration),
	}
}

func (m *MockVoterRegistry) GetVoter(ctx context.Context, namespace, name string) (*votingv1alpha1.VoterRegistration, error) {
	key := namespace + "/" + name
	if v, ok := m.voters[key]; ok {
		return v, nil
	}
	return nil, ErrVoterNotFound
}

func (m *MockVoterRegistry) AddVoter(namespace, name string, voter *votingv1alpha1.VoterRegistration) {
	key := namespace + "/" + name
	m.voters[key] = voter
}

var ErrVoterNotFound = fmt.Errorf("voter not found")

func TestEd25519Verifier_VerifyVote(t *testing.T) {
	// Generate test keypair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate keypair: %v", err)
	}

	registry := NewMockVoterRegistry()
	registry.AddVoter("default", "test-voter", &votingv1alpha1.VoterRegistration{
		Name:      "test-voter",
		PublicKey: base64.StdEncoding.EncodeToString(pubKey),
	})

	verifier := &Ed25519Verifier{VoterRegistry: registry}

	proposal := &votingv1alpha1.VoteProposal{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-proposal",
			Namespace: "default",
		},
	}

	timestamp := metav1.Now()
	vote := &votingv1alpha1.Vote{
		Voter:     "test-voter",
		Decision:  votingv1alpha1.VoteApprove,
		Timestamp: timestamp,
		Rationale: "Looks good",
	}

	// Create canonical message
	voteMsg := struct {
		ProposalID string `json:"proposalId"`
		Voter      string `json:"voter"`
		Decision   string `json:"decision"`
		Timestamp  string `json:"timestamp"`
	}{
		ProposalID: "default/test-proposal",
		Voter:      "test-voter",
		Decision:   string(votingv1alpha1.VoteApprove),
		Timestamp:  timestamp.Format(time.RFC3339),
	}

	msgBytes, _ := json.Marshal(voteMsg)
	signature := ed25519.Sign(privKey, msgBytes)
	vote.Signature = base64.StdEncoding.EncodeToString(signature)

	t.Run("valid signature", func(t *testing.T) {
		valid, err := verifier.VerifyVote(context.Background(), vote, proposal)
		if err != nil {
			t.Errorf("VerifyVote() error = %v", err)
		}
		if !valid {
			t.Errorf("VerifyVote() = false, want true")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		badVote := *vote
		badVote.Signature = base64.StdEncoding.EncodeToString([]byte("invalid signature data here!!!"))
		valid, err := verifier.VerifyVote(context.Background(), &badVote, proposal)
		if err != nil {
			// May error or return false
		}
		if valid {
			t.Errorf("VerifyVote() = true for invalid signature, want false")
		}
	})

	t.Run("unsigned vote - should allow", func(t *testing.T) {
		unsignedVote := *vote
		unsignedVote.Signature = ""
		valid, err := verifier.VerifyVote(context.Background(), &unsignedVote, proposal)
		if err != nil {
			t.Errorf("VerifyVote() error = %v", err)
		}
		if !valid {
			t.Errorf("VerifyVote() = false for unsigned vote, want true (allowed)")
		}
	})

	t.Run("voter not found", func(t *testing.T) {
		unknownVote := *vote
		unknownVote.Voter = "unknown-voter"
		_, err := verifier.VerifyVote(context.Background(), &unknownVote, proposal)
		if err == nil {
			t.Errorf("Expected error for unknown voter, got nil")
		}
	})

	t.Run("tampering detection - changed decision", func(t *testing.T) {
		tamperedVote := *vote
		tamperedVote.Decision = votingv1alpha1.VoteReject // Changed from APPROVE
		// Signature is still for APPROVE
		valid, err := verifier.VerifyVote(context.Background(), &tamperedVote, proposal)
		if valid {
			t.Errorf("VerifyVote() = true for tampered vote, want false")
		}
		if err != nil {
			// Error is acceptable
		}
	})
}

func TestDecodePublicKey(t *testing.T) {
	pubKey, _, _ := ed25519.GenerateKey(rand.Reader)

	t.Run("base64 encoded key", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString(pubKey)
		decoded, err := decodePublicKey(encoded)
		if err != nil {
			t.Errorf("decodePublicKey() error = %v", err)
		}
		if !decoded.Equal(pubKey) {
			t.Errorf("Decoded key doesn't match original")
		}
	})

	t.Run("invalid base64", func(t *testing.T) {
		_, err := decodePublicKey("not-valid-base64!!!")
		if err == nil {
			t.Errorf("Expected error for invalid base64, got nil")
		}
	})

	t.Run("wrong key length", func(t *testing.T) {
		shortKey := base64.StdEncoding.EncodeToString([]byte("tooshort"))
		_, err := decodePublicKey(shortKey)
		if err == nil {
			t.Errorf("Expected error for wrong key length, got nil")
		}
	})
}

func BenchmarkEd25519Verify(b *testing.B) {
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	message := []byte("benchmark message for signature verification")
	signature := ed25519.Sign(privKey, message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ed25519.Verify(pubKey, message, signature)
	}
}
