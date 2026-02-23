// Auto-generated tests for phase1-015

```go
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSwarmCoordinator_CreateProposal(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		description  string
		proposalType string
	}{
		{
			name:         "valid proposal",
			title:        "Add new feature",
			description:  "Implement user authentication",
			proposalType: "feature",
		},
		{
			name:         "empty title",
			title:        "",
			description:  "Some description",
			proposalType: "bugfix",
		},
		{
			name:         "long description",
			title:        "Complex feature",
			description:  strings.Repeat("x", 10000),
			proposalType: "refactor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coordinator := NewSwarmCoordinator(
				NewMockAgentCaller(),
				NewMockGitManager(),
				NewMockMissionControlClient(),
				NewMockQuotaManager(100),
			)

			ctx := context.Background()
			proposal, err := coordinator.CreateProposal(ctx, tt.title, tt.description, tt.proposalType)

			if err != nil {
				t.Fatalf("CreateProposal() error = %v", err)
			}

			if proposal == nil {
				t.Fatal("CreateProposal() returned nil proposal")
			}

			if proposal.Title != tt.title {
				t.Errorf("Title = %v, want %v", proposal.Title, tt.title)
			}

			if proposal.Description != tt.description {
				t.Errorf("Description = %v, want %v", proposal.Description, tt.description)
			}

			if proposal.Type != tt.proposalType {
				t.Errorf("Type = %v, want %v", proposal.Type, tt.proposalType)
			}

			if proposal.Status != ProposalStatusDraft {
				t.Errorf("Status = %v, want %v", proposal.Status, ProposalStatusDraft)
			}

			if len(proposal.Reviews) != 0 {
				t.Errorf("Reviews length = %v, want 0", len(proposal.Reviews))
			}

			if !strings.HasPrefix(proposal.ProposalID, "proposal-") {
				t.Errorf("ProposalID = %v, should start with 'proposal-'", proposal.ProposalID)
			}

			// Verify proposal is stored
			stored, ok := coordinator.GetProposal(proposal.ProposalID)
			if !ok {
				t.Error("Proposal not stored in coordinator")
			}
			if stored.ProposalID != proposal.ProposalID {
				t.Errorf("Stored proposal ID = %v, want %v", stored.ProposalID, proposal.ProposalID)
			}
		})
	}
}

func TestSwarmCoordinator_ReviewProposal(t *testing.T) {
	tests := []struct {
		name              string
		agentResponses    map[string]string
		quotaLimit        int
		expectedStatus    ProposalStatus
		expectedDecision  AgentDecision
		expectedApproval  float64
		expectedAvgScore  float64
		expectError       bool
		errorContains     string
	}{
		{
			name: "all approve - proposal approved",
			agentResponses: map[string]string{
				"architect": `{"score": 90, "decision": "APPROVE", "reason": "Good design"}`,
				"engineer":  `{"score": 85, "decision": "APPROVE", "reason": "Clean code"}`,
				"qa":        `{"score": 88, "decision": "APPROVE", "reason": "Well tested"}`,
				"security":  `{"score": 92, "decision": "APPROVE", "reason": "Secure"}`,
			},
			quotaLimit:       100,
			expectedStatus:   ProposalStatusApproved,
			expectedDecision: AgentDecisionApprove,
			expectedApproval: 100.0,
			expectedAvgScore: 88.75,
			expectError:      false,
		},
		{
			name: "mixed approvals - proposal approved (80%)",
			agentResponses: map[string]string{
				"architect": `{"score": 85, "decision": "APPROVE", "reason": "Good"}`,
				"engineer":  `{"score": 80, "decision": "CONDITIONAL_APPROVE", "reason": "OK"}`,
				"qa":        `{"score": 75, "decision": "CONDITIONAL_APPROVE", "reason": "Acceptable"}`,
				"security":  `{"score": 60, "decision": "REJECT", "reason": "Concerns"}`,
			},
			quotaLimit:       100,
			expectedStatus:   ProposalStatusApproved,
			expectedDecision: AgentDecisionApprove,
			expectedApproval: 75.0,
			expectedAvgScore: 75.0,
			expectError:      false,
		},
		{
			name: "below threshold - proposal rejected",
			agentResponses: map[string]string{
				"architect": `{"score": 70, "decision": "CONDITIONAL_APPROVE", "reason": "Needs work"}`,
				"engineer":  `{"score": 60, "decision": "REJECT", "reason": "Poor quality"}`,
				"qa":        `{"score": 65, "decision": "REJECT", "reason": "Not tested"}`,
				"security":  `{"score": 50, "decision": "REJECT", "reason": "Insecure"}`,
			},
			quotaLimit:       100,
			expectedStatus:   ProposalStatusRejected,
			expectedDecision: AgentDecisionReject,
			expectedApproval: 25.0,
			expectedAvgScore: 61.25,
			expectError:      false,
		},
		{
			name: "exactly 80% approval - approved",
			agentResponses: map[string]string{
				"architect": `{"score": 85, "decision": "APPROVE", "reason": "Good"}`,
				"engineer":  `{"score": 80, "decision": "APPROVE", "reason": "Good"}`,
				"qa":        `{"score": 75, "decision": "CONDITIONAL_APPROVE", "reason": "OK"}`,
				"security":  `{"score": 60, "decision": "REJECT", "reason": "Concerns"}`,
			},
			quotaLimit:       100,
			expectedStatus:   ProposalStatusApproved,
			expectedDecision: AgentDecisionApprove,
			expectedApproval: 75.0,
			expectedAvgScore: 75.0,
			expectError:      false,
		},
		{
			name:           "quota exhausted",
			agentResponses: map[string]string{},
			quotaLimit:     0,
			expectError:    true,
			errorContains:  "quota exhausted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agentCaller := NewMockAgentCaller()
			for agent, response := range tt.agentResponses {
				agentCaller.SetResponse(agent, response)
			}
			agentCaller.SetLatency(1 * time.Millisecond)

			coordinator := NewSwarmCoordinator(
				agentCaller,
				NewMockGitManager(),
				NewMockMissionControlClient(),
				NewMockQuotaManager(tt.quotaLimit),
			)

			ctx := context.Background()
			proposal, _ := coordinator.CreateProposal(ctx, "Test", "Description", "feature")

			err := coordinator.ReviewProposal(ctx, proposal.ProposalID)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error = %v, want to contain %v", err.Error(), tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("ReviewProposal() error = %v", err)
			}

			if proposal.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", proposal.Status, tt.expectedStatus)
			}

			if proposal.ConsensusDecision != tt.expectedDecision {
				t.Errorf("ConsensusDecision = %v, want %v", proposal.ConsensusDecision, tt.expectedDecision)
			}

			if proposal.ApprovalRate != tt.expectedApproval {
				t.Errorf("ApprovalRate = %v, want %v", proposal.ApprovalRate, tt.expectedApproval)
			}

			if proposal.AvgScore != tt.expectedAvgScore {
				t.Errorf("AvgScore = %v, want %v", proposal.AvgScore, tt.expectedAvgScore)
			}

			if len(proposal.Reviews) != 4 {
				t.Errorf("Reviews length = %v, want 4", len(proposal.Reviews))
			}

			// Verify all agents were called
			calls := agentCaller.GetCalls()
			if len(calls) != 4 {
				t.Errorf("Agent calls = %v, want 4", len(calls))
			}

			expectedAgents := map[string]bool{
				"architect": false,
				"engineer":  false,
				"qa":        false,
				"security":  false,
			}
			for _, call := range calls {
				if _, exists := expectedAgents[call.SessionID]; exists {
					expectedAgents[call.SessionID] = true
				}
			}
			for agent, called := range expectedAgents {
				if !called {
					t.Errorf("Agent %s was not called", agent)
				}
			}
		})
	}
}

func TestSwarmCoordinator_ReviewProposal_NotFound(t *testing.T) {
	coordinator := NewSwarmCoordinator(
		NewMockAgentCaller(),
		NewMockGitManager(),
		NewMockMissionControlClient(),
		NewMockQuotaManager(100),
	)

	ctx := context.Background()
	err := coordinator.ReviewProposal(ctx, "non-existent-id")

	if err == nil {
		t.Fatal("Expected error for non-existent proposal")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error = %v, want 'not found'", err.Error())
	}
}

func TestSwarmCoordinator_ReviewProposal_InvalidJSON(t *testing.T) {
	agentCaller := NewMockAgentCaller()
	agentCaller.SetResponse("architect", "invalid json")
	agentCaller.SetLatency(1 * time.Millisecond)

	coordinator := NewSwarmCoordinator(
		agentCaller,
		NewMockGitManager(),
		NewMockMissionControlClient(),
		NewMockQuotaManager(100),
	)

	ctx := context.Background()
	proposal, _ := coordinator.CreateProposal(ctx, "Test", "Description", "feature")

	err := coordinator.ReviewProposal(ctx, proposal.ProposalID)

	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("Error = %v, want to contain 'parse'", err.Error())
	}
}

func TestSwarmCoordinator_ExecuteProposal(t *testing.T) {
	tests := []struct {
		name           string
		setupProposal  func(*SwarmCoordinator, context.Context) string
		setupMocks     func(*MockGitManager, *MockMissionControlClient)
		expectError    bool
		errorContains  string
		expectedStatus ProposalStatus
		verifyFunc     func(*testing.T, *SwarmCoordinator, string, *MockGitManager, *MockMissionControlClient)
	}{
		{
			name: "successful execution",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				proposal.Status = ProposalStatusApproved
				return proposal.ProposalID
			},
			setupMocks: func(gm *MockGitManager, mc *MockMissionControlClient) {},
			expectError: false,
			expectedStatus: ProposalStatusCompleted,
			verifyFunc: func(t *testing.T, sc *SwarmCoordinator, id string, gm *MockGitManager, mc *MockMissionControlClient) {
				proposal, _ := sc.GetProposal(id)
				
				if len(proposal.FilesChanged) != 2 {
					t.Errorf("FilesChanged length = %v, want 2", len(proposal.FilesChanged))
				}

				commits := gm.GetCommits()
				if len(commits) != 1 {
					t.Errorf("Commits = %v, want 1", len(commits))
				}

				if proposal.MCTaskID == nil {
					t.Error("MCTaskID is nil")
				} else {
					task, ok := mc.GetTask(*proposal.MCTaskID)
					if !ok {
						t.Error("Task not found in Mission Control")
					}
					if task.Status != "completed" {
						t.Errorf("Task status = %v, want 'completed'", task.Status)
					}
				}
			},
		},
		{
			name: "proposal not found",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				return "non-existent-id"
			},
			setupMocks:    func(gm *MockGitManager, mc *MockMissionControlClient) {},
			expectError:   true,
			errorContains: "not found",
		},
		{
			name: "proposal not approved - draft",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				return proposal.ProposalID
			},
			setupMocks:    func(gm *MockGitManager, mc *MockMissionControlClient) {},
			expectError:   true,
			errorContains: "not approved",
		},
		{
			name: "proposal not approved - rejected",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				proposal.Status = ProposalStatusRejected
				return proposal.ProposalID
			},
			setupMocks:    func(gm *MockGitManager, mc *MockMissionControlClient) {},
			expectError:   true,
			errorContains: "not approved",
		},
		{
			name: "git branch creation fails",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				proposal.Status = ProposalStatusApproved
				return proposal.ProposalID
			},
			setupMocks: func(gm *MockGitManager, mc *MockMissionControlClient) {
				gm.SetFailNext(true)
			},
			expectError:    true,
			errorContains:  "failed to create branch",
			expectedStatus: ProposalStatusFailed,
		},
		{
			name: "git commit fails",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				proposal.Status = ProposalStatusApproved
				return proposal.ProposalID
			},
			setupMocks: func(gm *MockGitManager, mc *MockMissionControlClient) {
				gm.branches["proposal-test"] = true
			},
			expectError:    true,
			errorContains:  "failed to commit",
			expectedStatus: ProposalStatusFailed,
			verifyFunc: func(t *testing.T, sc *SwarmCoordinator, id string, gm *MockGitManager, mc *MockMissionControlClient) {
				gm.SetFailNext(true)
			},
		},
		{
			name: "mission control task creation fails",
			setupProposal: func(sc *SwarmCoordinator, ctx context.Context) string {
				proposal, _ := sc.CreateProposal(ctx, "Feature", "Description", "feature")
				proposal.Status = ProposalStatusApproved
				return proposal.ProposalID
			},
			setupMocks: func(gm *MockGitManager, mc *MockMissionControlClient) {
				mc.SetFailNext(true)
			},
			expectError:    true,
			errorContains:  "failed to create mission control task",
			expectedStatus: ProposalStatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitManager := NewMockGitManager()
			mcClient := NewMockMissionControlClient()
			
			coordinator := NewSwarmCoordinator(
				NewMockAgentCaller(),
				gitManager,
				mcClient,
				NewMockQuotaManager(100),
			)

			ctx := context.Background()
			proposalID := tt.setupProposal(coordinator, ctx)
			tt.setupMocks(gitManager, mcClient)

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, coordinator, proposalID, gitManager, mcClient)
			}

			err := coordinator.ExecuteProposal(ctx, proposalID)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error = %v, want to contain %v", err.Error(), tt.errorContains)
				}
				if tt.expectedStatus != "" {
					proposal, _ := coordinator.GetProposal(proposalID)
					if proposal != nil && proposal.Status != tt.expectedStatus {
						t.Errorf("Status = %v, want %v", proposal.Status, tt.expectedStatus)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("ExecuteProposal() error = %v", err)
			}

			if tt.expectedStatus != "" {
				proposal, _ := coordinator.GetProposal(proposalID)
				if proposal.Status != tt.expectedStatus {
					t.Errorf("Status = %v, want %v", proposal.Status, tt.expectedStatus)
				}
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, coordinator, proposalID, gitManager, mcClient)
			}
		})
	}
}

func TestSwarmCoordinator_GetProposal(t *testing.T) {
	coordinator := NewSwarmCoordinator(
		NewMockAgentCaller(),
		NewMockGitManager(),
		NewMockMissionControlClient(),
		NewMockQuotaManager(100),
	)

	ctx := context.Background()
	created, _ := coordinator.CreateProposal(ctx, "Test", "Description", "feature")

	t.Run("existing proposal", func(t *testing.T) {
		proposal, ok := coordinator.GetProposal(created.ProposalID)
		if !ok {
			t.Fatal("GetProposal() returned false for existing proposal")
		}
		if proposal.ProposalID != created.ProposalID {
			t.Errorf("ProposalID = %v, want %v", proposal.ProposalID, created.ProposalID)
		}
	})

	t.Run("non-existent proposal", func(t *testing.T) {
		_, ok := coordinator.GetProposal("non-existent")
		if ok {
			t.Error("GetProposal() returned true for non-existent proposal")
		}
	})
}

func TestSwarmCoordinator_ConcurrentOperations(t *testing.T) {
	coordinator := NewSwarmCoordinator(
		NewMockAgentCaller(),
		NewMockGitManager(),
		NewMockMissionControlClient(),
		NewMockQuotaManager(1000),
	)

	ctx := context.Background()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent CreateProposal
	t.Run("concurrent create", func(t *testing.T) {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := coordinator.CreateProposal(ctx, fmt.Sprintf("Title %d", id), "Description", "feature")
				if err != nil {
					t.Errorf("CreateProposal failed: %v", err)
				}
			}(i)
		}
		wg.Wait()
	})

	// Verify all proposals were created
	coordinator.mu.Lock()
	count := len(coordinator.proposals)
	coordinator.mu.Unlock()

	if count != numGoroutines {
		t.Errorf("Created %d proposals, want %d", count, numGoroutines)
	}
}

func TestMockAgentCaller(t *testing.T) {
	mock := NewMockAgentCaller()
	mock.SetLatency(50 * time.Millisecond)
	
	response := `{"score": 95, "decision": "APPROVE", "reason": "Excellent"}`
	mock.SetResponse("test-agent", response)

	ctx := context.Background()
	start := time.Now()
	result, latency, err := mock.CallAgent(ctx, "test-agent", "test message", 30)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("CallAgent error = %v", err)
	}

	if result != response {
		t.Errorf("Response = %v, want %v", result, response)
	}

	if latency < 50 {
		t.Errorf("Latency = %v, want >= 50", latency)
	}

	if elapsed < 50*time.Millisecond {
		t.Errorf("Elapsed time = %v, want >= 50ms", elapsed)
	}

	calls := mock.GetCalls()
	if len(calls) != 1 {
		t.Errorf("Calls length = %v, want 1", len(calls))
	}

	if calls[0].SessionID != "test-agent" {
		t.Errorf("SessionID = %v, want 'test-agent'", calls[0].SessionID)
	}

	if calls[0].Message != "test message" {
		t.Errorf("Message = %v, want 'test message'", calls[0].Message)
	}
}

func TestMockGitManager(t *testing.T) {
	mock := NewMockGitManager()
	ctx := context.Background()

	t.Run("create branch", func(t *testing.T) {
		err := mock.CreateBranch(ctx, "feature-branch")
		if err != nil {
			t.Errorf("CreateBranch error = %v", err)
		}

		mock.mu.Lock()
		exists := mock.branches["feature-branch"]
		mock.mu.Unlock()

		if !exists {
			t.Error("Branch not created")
		}
	})

	t.Run("commit changes", func(t *testing.T) {
		files := []string{"file1.go", "file2.go"}
		err := mock.CommitChanges(ctx, "Test commit", files)
		if err != nil {
			t.Errorf("CommitChanges error = %v", err)
		}

		commits := mock.GetCommits()
		if len(commits) != 1 {
			t.Errorf("Commits length = %v, want 1", len(commits))
		}

		if commits[0].Message != "Test commit" {
			t.Errorf("Commit message = %v, want 'Test commit'", commits[0].Message)
		}
	})

	t.Run("fail next operation", func(t *testing.T) {
		mock.SetFailNext(true)
		err := mock.CreateBranch(ctx, "fail-branch")
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})
}

func TestMockMissionControlClient(t *testing.T) {
	mock := NewMockMissionControlClient()
	ctx := context.Background()

	t.Run("create task", func(t *testing.T) {
		taskID, err := mock.CreateTask(ctx, "Test Task", "Description")
		if err != nil {
			t.Errorf("CreateTask error = %v", err)
		}

		task, ok := mock.GetTask(taskID)
		if !ok {
			t.Error("Task not found")
		}

		if task.Title != "Test Task" {
			t.Errorf("Task title = %v, want 'Test Task'", task.Title)
		}

		if task.Status != "open" {
			t.Errorf("Task status = %v, want 'open'", task.Status)
		}
	})

	t.Run("update task", func(t *testing.T) {
		taskID, _ := mock.CreateTask(ctx, "Task", "Desc")
		err := mock.UpdateTask(ctx, taskID, "completed", "Done")
		if err != nil {
			t.Errorf("UpdateTask error = %v", err)
		}

		task, _ := mock.GetTask(taskID)
		if task.Status != "completed" {
			t.Errorf("Task status = %v, want 'completed'", task.Status)
		}

		if task.Notes != "Done" {
			t.Errorf("Task notes = %v, want 'Done'", task.Notes)
		}
	})

	t.Run("fail next operation", func(t *testing.T) {
		mock.SetFailNext(true)
		_, err := mock.CreateTask(ctx, "Fail", "Desc")
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})
}

func TestMockQuotaManager(t *testing.T) {
	mock := NewMockQuotaManager(10)
	ctx := context.Background()

	t.Run("check quota available", func(t *testing.T) {
		available, err := mock.CheckQuota(ctx, "sonnet")
		if err != nil {
			t.Errorf("CheckQuota error = %v", err)
		}
		if !available {
			t.Error("Expected quota to be available")
		}
	})

	t.Run("record calls", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			err := mock.RecordCall(ctx, "sonnet")
			if err != nil {
				t.Errorf("RecordCall error = %v", err)
			}
		}

		state, err := mock.GetState(ctx)
		if err != nil {
			t.Errorf("GetState error = %v", err)
		}

		if state.TotalCalls != 5 {
			t.Errorf("TotalCalls = %v, want 5", state.TotalCalls)
		}

		if state.ModelCalls["sonnet"] != 5 {
			t.Errorf("ModelCalls[sonnet] = %v, want 5", state.ModelCalls["sonnet"])
		}
	})

	t.Run("quota exhausted", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			mock.RecordCall(ctx, "opus")
		}

		available, err := mock.CheckQuota(ctx, "opus")
		if err != nil {
			t.Errorf("CheckQuota error = %v", err)
		}
		if available {
			t.Error("Expected quota to be exhausted")
		}
	})

	t.Run("reset quota", func(t *testing.T) {
		mock.Reset()
		state, _ := mock.GetState(ctx)
		if state.TotalCalls != 0 {
			t.Errorf("TotalCalls after reset = %v, want 0", state.TotalCalls)
		}
	})
}

func TestAgentDecisionConstants(t *testing.T) {
	tests := []struct {
		decision AgentDecision
		expected string
	}{
		{AgentDecisionApprove, "APPROVE"},
		{AgentDecisionConditionalApprove, "CONDITIONAL_APPROVE"},
		{AgentDecisionReject, "REJECT"},
	}

	for _, tt := range tests {
		if string(tt.decision) != tt.expected {
			t.Errorf("AgentDecision = %v, want %v", tt.decision, tt.expected)
		}
	}
}

func TestProposalStatusConstants(t *testing.T) {
	tests := []struct {
		status   ProposalStatus
		expected string
	}{
		{ProposalStatusDraft, "DRAFT"},
		{ProposalStatusReview, "REVIEW"},
		{ProposalStatusApproved, "APPROVED"},
		{ProposalStatusRejected, "REJECTED"},
		{ProposalStatusExecuting, "EXECUTING"},
		{ProposalStatusCompleted, "COMPLETED"},
		{ProposalStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("ProposalStatus = %v, want %v", tt.status, tt.expected)
		}
	}
}

func TestProposalContext_JSONSerialization(t *testing.T) {
	taskID := "task-123"
	proposal := &ProposalContext{
		ProposalID:        "proposal-1",
		Title:             "Test Proposal",
		Description:       "Description",
		Type:              "feature",
		Reviews:           []AgentReview{},
		AvgScore:          85.5,
		ApprovalRate:      100.0,
		ConsensusDecision: AgentDecisionApprove,
		Status:            ProposalStatusApproved,
		FilesChanged:      []string{"file1.go"},
		MCTaskID:          &taskID,
	}

	data, err := json.Marshal(proposal)
	if err != nil {
		t.Fatalf("JSON Marshal error = %v", err)
	}

	var decoded ProposalContext
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("JSON Unmarshal error = %v", err)
	}

	if decoded.ProposalID != proposal.ProposalID {
		t.Errorf("ProposalID = %v, want %v", decoded.ProposalID, proposal.ProposalID)
	}

	if *decoded.MCTaskID != *proposal.MCTaskID {
		t.Errorf("MCTaskID = %v, want %v", *decoded.MCTaskID, *proposal.MCTaskID)
	}
}
```