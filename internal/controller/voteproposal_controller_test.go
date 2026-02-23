// Copyright 2024 The Voting Operator Authors.
// Licensed under the Apache License, Version 2.0.

package controller

import (
	"context"
	"testing"
	"time"

	votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
	"github.com/yourorg/voting-operator/internal/audit"
	"github.com/yourorg/voting-operator/internal/execution"
	"github.com/yourorg/voting-operator/internal/voting"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// Mock implementations for testing

type MockVoterNotifier struct {
	NotifyCalled bool
	NotifyError  error
}

func (m *MockVoterNotifier) NotifyVoters(ctx context.Context, proposal *votingv1alpha1.VoteProposal, voters []string) error {
	m.NotifyCalled = true
	return m.NotifyError
}

type MockVoteAggregator struct {
	Score    float64
	Decision string
}

func (m *MockVoteAggregator) CalculateConsensus(proposal *votingv1alpha1.VoteProposal) (float64, string) {
	return m.Score, m.Decision
}

type MockExecutionEngine struct {
	ExecuteCalled bool
	ExecuteError  error
	ExecuteStatus *votingv1alpha1.ExecutionStatus
}

func (m *MockExecutionEngine) Execute(ctx context.Context, proposal *votingv1alpha1.VoteProposal) (*votingv1alpha1.ExecutionStatus, error) {
	m.ExecuteCalled = true
	if m.ExecuteError != nil {
		return nil, m.ExecuteError
	}
	if m.ExecuteStatus != nil {
		return m.ExecuteStatus, nil
	}
	return &votingv1alpha1.ExecutionStatus{
		Success: true,
		Message: "Executed successfully",
	}, nil
}

type MockAuditLogger struct {
	LogCalled          bool
	OverrideLogCalled  bool
	LogError           error
}

func (m *MockAuditLogger) LogProposal(ctx context.Context, proposal *votingv1alpha1.VoteProposal) error {
	m.LogCalled = true
	return m.LogError
}

func (m *MockAuditLogger) LogEmergencyOverride(ctx context.Context, proposal *votingv1alpha1.VoteProposal) error {
	m.OverrideLogCalled = true
	return m.LogError
}

type MockSignatureVerifier struct {
	VerifyResult bool
	VerifyError  error
}

func (m *MockSignatureVerifier) VerifyVote(ctx context.Context, vote *votingv1alpha1.Vote, proposal *votingv1alpha1.VoteProposal) (bool, error) {
	return m.VerifyResult, m.VerifyError
}

// Test suite
var _ = Describe("VoteProposalController", func() {
	var (
		scheme *runtime.Scheme
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		scheme = runtime.NewScheme()
		_ = votingv1alpha1.AddToScheme(scheme)
	})

	Context("when reconciling a new proposal", func() {
		It("should transition to Voting phase and notify voters", func() {
			proposal := &votingv1alpha1.VoteProposal{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-proposal",
					Namespace: "default",
				},
				Spec: votingv1alpha1.VoteProposalSpec{
					Operation: "deployment.update",
					TargetRef: votingv1alpha1.TargetReference{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "test-deployment",
						Namespace:  "default",
					},
					VoteConfig: votingv1alpha1.VoteConfiguration{
						Threshold:      80,
						Timeout:        metav1.Duration{Duration: 5 * time.Minute},
						RequiredVoters: []string{"voter1", "voter2"},
					},
					ProposedChange: votingv1alpha1.ChangeSpec{
						Type: "patch",
					},
				},
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(proposal).
				WithStatusSubresource(proposal).
				Build()

			mockNotifier := &MockVoterNotifier{}
			mockAggregator := &MockVoteAggregator{Score: 0, Decision: "NO_VOTES"}
			mockExecutor := &MockExecutionEngine{}
			mockAuditor := &MockAuditLogger{}
			mockVerifier := &MockSignatureVerifier{VerifyResult: true}

			reconciler := &VoteProposalReconciler{
				Client:            fakeClient,
				Scheme:            scheme,
				VoterNotifier:     mockNotifier,
				VoteAggregator:    mockAggregator,
				ExecutionEngine:   mockExecutor,
				AuditLogger:       mockAuditor,
				SignatureVerifier: mockVerifier,
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-proposal",
					Namespace: "default",
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(mockNotifier.NotifyCalled).To(BeTrue())

			// Check status was updated
			var updated votingv1alpha1.VoteProposal
			err = fakeClient.Get(ctx, types.NamespacedName{
				Name:      "test-proposal",
				Namespace: "default",
			}, &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Status.Phase).To(Equal(votingv1alpha1.PhaseVoting))
			Expect(updated.Status.VotingDeadline).NotTo(BeNil())
		})
	})

	Context("when proposal reaches threshold", func() {
		It("should transition to Approved phase", func() {
			proposal := &votingv1alpha1.VoteProposal{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-proposal",
					Namespace: "default",
				},
				Spec: votingv1alpha1.VoteProposalSpec{
					VoteConfig: votingv1alpha1.VoteConfiguration{
						Threshold:      80,
						RequiredVoters: []string{"voter1"},
					},
				},
				Status: votingv1alpha1.VoteProposalStatus{
					Phase: votingv1alpha1.PhaseVoting,
					Votes: []votingv1alpha1.Vote{
						{
							Voter:    "voter1",
							Decision: votingv1alpha1.VoteApprove,
							Score:    100,
						},
					},
				},
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(proposal).
				WithStatusSubresource(proposal).
				Build()

			mockAggregator := &MockVoteAggregator{Score: 100, Decision: "APPROVED"}
			mockVerifier := &MockSignatureVerifier{VerifyResult: true}

			reconciler := &VoteProposalReconciler{
				Client:            fakeClient,
				Scheme:            scheme,
				VoteAggregator:    mockAggregator,
				SignatureVerifier: mockVerifier,
				VoterNotifier:     &MockVoterNotifier{},
				ExecutionEngine:   &MockExecutionEngine{},
				AuditLogger:       &MockAuditLogger{},
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-proposal",
					Namespace: "default",
				},
			})

			Expect(err).NotTo(HaveOccurred())

			var updated votingv1alpha1.VoteProposal
			err = fakeClient.Get(ctx, types.NamespacedName{
				Name:      "test-proposal",
				Namespace: "default",
			}, &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Status.Phase).To(Equal(votingv1alpha1.PhaseApproved))
			Expect(*updated.Status.AggregateScore).To(Equal(100.0))
		})
	})

	Context("when proposal is approved", func() {
		It("should execute the change", func() {
			proposal := &votingv1alpha1.VoteProposal{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-proposal",
					Namespace: "default",
				},
				Spec: votingv1alpha1.VoteProposalSpec{
					Operation: "deployment.update",
					TargetRef: votingv1alpha1.TargetReference{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "test-deployment",
					},
					ProposedChange: votingv1alpha1.ChangeSpec{
						Type: "patch",
					},
				},
				Status: votingv1alpha1.VoteProposalStatus{
					Phase:          votingv1alpha1.PhaseApproved,
					AggregateScore: ptr(85.0),
				},
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(proposal).
				WithStatusSubresource(proposal).
				Build()

			mockExecutor := &MockExecutionEngine{
				ExecuteStatus: &votingv1alpha1.ExecutionStatus{
					Success: true,
					Message: "Executed",
				},
			}
			mockAuditor := &MockAuditLogger{}

			reconciler := &VoteProposalReconciler{
				Client:          fakeClient,
				Scheme:          scheme,
				ExecutionEngine: mockExecutor,
				AuditLogger:     mockAuditor,
				VoterNotifier:   &MockVoterNotifier{},
				VoteAggregator:  &MockVoteAggregator{},
				SignatureVerifier: &MockSignatureVerifier{},
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-proposal",
					Namespace: "default",
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(mockExecutor.ExecuteCalled).To(BeTrue())

			var updated votingv1alpha1.VoteProposal
			err = fakeClient.Get(ctx, types.NamespacedName{
				Name:      "test-proposal",
				Namespace: "default",
			}, &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Status.Phase).To(Equal(votingv1alpha1.PhaseExecuted))
			Expect(updated.Status.ExecutionStatus.Success).To(BeTrue())
		})
	})

	Context("emergency override", func() {
		It("should execute immediately with audit log", func() {
			proposal := &votingv1alpha1.VoteProposal{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "emergency-proposal",
					Namespace: "default",
				},
				Spec: votingv1alpha1.VoteProposalSpec{
					EmergencyOverride: &votingv1alpha1.EmergencyOverrideSpec{
						Enabled:       true,
						Justification: "Production outage - immediate fix required",
						Approver:      "oncall-sre",
					},
				},
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(proposal).
				WithStatusSubresource(proposal).
				Build()

			mockAuditor := &MockAuditLogger{}

			reconciler := &VoteProposalReconciler{
				Client:      fakeClient,
				Scheme:      scheme,
				AuditLogger: mockAuditor,
				VoterNotifier: &MockVoterNotifier{},
				VoteAggregator: &MockVoteAggregator{},
				ExecutionEngine: &MockExecutionEngine{},
				SignatureVerifier: &MockSignatureVerifier{},
			}

			_, err := reconciler.Reconcile(ctx, ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "emergency-proposal",
					Namespace: "default",
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(mockAuditor.OverrideLogCalled).To(BeTrue())

			var updated votingv1alpha1.VoteProposal
			err = fakeClient.Get(ctx, types.NamespacedName{
				Name:      "emergency-proposal",
				Namespace: "default",
			}, &updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Status.Decision).To(Equal("EMERGENCY_OVERRIDE"))
		})
	})
})

// Unit test for hasQuorum
func TestHasQuorum(t *testing.T) {
	reconciler := &VoteProposalReconciler{}

	tests := []struct {
		name            string
		requiredVoters  []string
		votes           []votingv1alpha1.Vote
		expectedQuorum  bool
	}{
		{
			name:           "all required voters voted",
			requiredVoters: []string{"v1", "v2"},
			votes: []votingv1alpha1.Vote{
				{Voter: "v1"},
				{Voter: "v2"},
			},
			expectedQuorum: true,
		},
		{
			name:           "missing one required voter",
			requiredVoters: []string{"v1", "v2", "v3"},
			votes: []votingv1alpha1.Vote{
				{Voter: "v1"},
				{Voter: "v2"},
			},
			expectedQuorum: false,
		},
		{
			name:           "extra optional voters",
			requiredVoters: []string{"v1"},
			votes: []votingv1alpha1.Vote{
				{Voter: "v1"},
				{Voter: "optional-v2"},
				{Voter: "optional-v3"},
			},
			expectedQuorum: true,
		},
		{
			name:           "no votes",
			requiredVoters: []string{"v1"},
			votes:          []votingv1alpha1.Vote{},
			expectedQuorum: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &votingv1alpha1.VoteProposal{
				Spec: votingv1alpha1.VoteProposalSpec{
					VoteConfig: votingv1alpha1.VoteConfiguration{
						RequiredVoters: tt.requiredVoters,
					},
				},
				Status: votingv1alpha1.VoteProposalStatus{
					Votes: tt.votes,
				},
			}

			result := reconciler.hasQuorum(proposal)
			if result != tt.expectedQuorum {
				t.Errorf("hasQuorum() = %v, want %v", result, tt.expectedQuorum)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestMain(m *testing.M) {
	// Setup test suite
	RegisterFailHandler(Fail)
	RunSpecs(m, "Controller Suite")
}
