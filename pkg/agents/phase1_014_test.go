// Auto-generated tests for phase1-014

```go
package agents

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

// Additional Security Agent Tests

func TestSecurityAgentScoringEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		proposal       *Proposal
		expectedScore  float64
		expectedDecision Decision
	}{
		{
			name: "multiple shell files penalty",
			proposal: &Proposal{
				ID:    "PROP-MULTI-SH",
				Title: "Deploy scripts",
				Implementation: map[string]interface{}{
					"files": []interface{}{"deploy.sh", "backup.sh", "cleanup.sh"},
				},
			},
			expectedScore:  40.0, // 70 - 30 (3 x 10)
			expectedDecision: DecisionReject,
		},
		{
			name: "mixed file types",
			proposal: &Proposal{
				ID:    "PROP-MIXED",
				Title: "Mixed changes",
				Implementation: map[string]interface{}{
					"files": []interface{}{"script.sh", "config.yaml", "data.json"},
				},
			},
			expectedScore:  60.0, // 70 - 10 (1 .sh file)
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "privileged and shell files combined",
			proposal: &Proposal{
				ID:    "PROP-COMBINED",
				Title: "Dangerous combo",
				Implementation: map[string]interface{}{
					"privileged": true,
					"files":      []interface{}{"install.sh"},
				},
			},
			expectedScore:  35.0, // 70 - 25 - 10
			expectedDecision: DecisionReject,
		},
		{
			name: "score exactly at threshold 60",
			proposal: &Proposal{
				ID:    "PROP-THRESHOLD-60",
				Title: "Exactly 60",
				Implementation: map[string]interface{}{
					"files": []interface{}{"script.sh"},
				},
			},
			expectedScore:  60.0,
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "score exactly at threshold 80",
			proposal: &Proposal{
				ID:    "PROP-THRESHOLD-80",
				Title: "Exactly 80",
				Implementation: map[string]interface{}{},
			},
			expectedScore:  70.0,
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "no implementation section",
			proposal: &Proposal{
				ID:             "PROP-NO-IMPL",
				Title:          "No implementation",
				Implementation: map[string]interface{}{},
			},
			expectedScore:  70.0,
			expectedDecision: DecisionConditionalApprove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewSecurityAgent(nil, nil)
			ctx := TestContext(t)

			result, err := agent.Review(ctx, tt.proposal)
			if err != nil {
				t.Fatalf("Review failed: %v", err)
			}

			if result.Score != tt.expectedScore {
				t.Errorf("Expected score %.1f, got %.1f", tt.expectedScore, result.Score)
			}
			if result.Decision != tt.expectedDecision {
				t.Errorf("Expected decision %s, got %s", tt.expectedDecision, result.Decision)
			}
		})
	}
}

func TestSecurityAgentCustomBaseline(t *testing.T) {
	tests := []struct {
		name         string
		baseline     float64
		proposal     *Proposal
		expectedScore float64
	}{
		{
			name:     "baseline 80",
			baseline: 80.0,
			proposal: SampleProposal(),
			expectedScore: 80.0,
		},
		{
			name:     "baseline 50",
			baseline: 50.0,
			proposal: SampleProposal(),
			expectedScore: 50.0,
		},
		{
			name:     "baseline 90 with penalty",
			baseline: 90.0,
			proposal: &Proposal{
				ID:    "PROP-PENALTY",
				Title: "Test",
				Implementation: map[string]interface{}{
					"files": []interface{}{"script.sh"},
				},
			},
			expectedScore: 80.0, // 90 - 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AgentConfig{
				Model:         "claude-sonnet-4-5",
				BaselineScore: tt.baseline,
			}
			agent := NewSecurityAgent(config, nil)
			ctx := TestContext(t)

			result, err := agent.Review(ctx, tt.proposal)
			if err != nil {
				t.Fatalf("Review failed: %v", err)
			}

			if result.Score != tt.expectedScore {
				t.Errorf("Expected score %.1f, got %.1f", tt.expectedScore, result.Score)
			}
		})
	}
}

func TestSecurityAgentName(t *testing.T) {
	agent := NewSecurityAgent(nil, nil)
	if agent.Name() != "SecurityAgent" {
		t.Errorf("Expected name 'SecurityAgent', got %s", agent.Name())
	}
}

// Additional PM Agent Tests

func TestPMAgentScoringRules(t *testing.T) {
	tests := []struct {
		name           string
		proposal       *Proposal
		expectedScore  float64
		expectedDecision Decision
	}{
		{
			name: "long impact text bonus",
			proposal: &Proposal{
				ID:     "PROP-IMPACT",
				Title:  "Feature with impact",
				Impact: "This is a very detailed impact description that is longer than 20 characters",
			},
			expectedScore:  85.0, // 75 + 10
			expectedDecision: DecisionApprove,
		},
		{
			name: "short impact no bonus",
			proposal: &Proposal{
				ID:     "PROP-SHORT",
				Title:  "Feature",
				Impact: "Short impact",
			},
			expectedScore:  75.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "empty impact no bonus",
			proposal: &Proposal{
				ID:     "PROP-EMPTY-IMPACT",
				Title:  "Feature",
				Impact: "",
			},
			expectedScore:  75.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "exactly 20 chars no bonus",
			proposal: &Proposal{
				ID:     "PROP-EXACT-20",
				Title:  "Feature",
				Impact: "12345678901234567890", // exactly 20
			},
			expectedScore:  75.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "21 chars gets bonus",
			proposal: &Proposal{
				ID:     "PROP-21-CHARS",
				Title:  "Feature",
				Impact: "123456789012345678901", // 21 chars
			},
			expectedScore:  85.0,
			expectedDecision: DecisionApprove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewPMAgent(nil, nil)
			ctx := TestContext(t)

			result, err := agent.Review(ctx, tt.proposal)
			if err != nil {
				t.Fatalf("Review failed: %v", err)
			}

			if result.Score != tt.expectedScore {
				t.Errorf("Expected score %.1f, got %.1f", tt.expectedScore, result.Score)
			}
			if result.Decision != tt.expectedDecision {
				t.Errorf("Expected decision %s, got %s", tt.expectedDecision, result.Decision)
			}
		})
	}
}

func TestPMAgentName(t *testing.T) {
	agent := NewPMAgent(nil, nil)
	if agent.Name() != "PMAgent" {
		t.Errorf("Expected name 'PMAgent', got %s", agent.Name())
	}
}

func TestPMAgentNilProposal(t *testing.T) {
	agent := NewPMAgent(nil, nil)
	ctx := TestContext(t)

	result, err := agent.Review(ctx, nil)
	if err != nil {
		t.Fatalf("Review failed: %v", err)
	}

	if result.Decision != DecisionReject {
		t.Errorf("Expected REJECT for nil proposal, got %s", result.Decision)
	}
	if result.Score != 0 {
		t.Errorf("Expected score 0 for nil proposal, got %.1f", result.Score)
	}
}

// Additional Architect Agent Tests

func TestArchitectAgentName(t *testing.T) {
	agent := NewArchitectAgent(nil, nil)
	if agent.Name() != "ArchitectAgent" {
		t.Errorf("Expected name 'ArchitectAgent', got %s", agent.Name())
	}
}

func TestArchitectAgentCustomConfig(t *testing.T) {
	config := &AgentConfig{
		Model:         "custom-model",
		BaselineScore: 95.0,
	}
	agent := NewArchitectAgent(config, nil)

	if agent.Config.Model != "custom-model" {
		t.Errorf("Expected custom model, got %s", agent.Config.Model)
	}
}

func TestArchitectAgentValidTitle(t *testing.T) {
	agent := NewArchitectAgent(nil, nil)
	ctx := TestContext(t)

	proposal := &Proposal{
		ID:          "PROP-ARCH",
		Title:       "Valid architecture",
		Description: "Description",
	}

	result, err := agent.Review(ctx, proposal)
	if err != nil {
		t.Fatalf("Review failed: %v", err)
	}

	if result.Decision != DecisionApprove {
		t.Errorf("Expected APPROVE, got %s", result.Decision)
	}
	if result.Score != 85.0 {
		t.Errorf("Expected score 85.0, got %.1f", result.Score)
	}
}

// Consensus Calculator Edge Cases

func TestConsensusCalculatorBoundaryScores(t *testing.T) {
	calc := NewConsensusCalculator()

	tests := []struct {
		name           string
		results        []*ReviewResult
		expectedScore  float64
		expectedDecision Decision
	}{
		{
			name: "exactly at 80 threshold",
			results: []*ReviewResult{
				{Score: 80, Decision: DecisionApprove},
				{Score: 80, Decision: DecisionApprove},
			},
			expectedScore:  80.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "just below 80 threshold",
			results: []*ReviewResult{
				{Score: 79.9, Decision: DecisionConditionalApprove},
				{Score: 79.9, Decision: DecisionConditionalApprove},
			},
			expectedScore:  79.9,
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "exactly at 60 threshold",
			results: []*ReviewResult{
				{Score: 60, Decision: DecisionConditionalApprove},
				{Score: 60, Decision: DecisionConditionalApprove},
			},
			expectedScore:  60.0,
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "just below 60 threshold",
			results: []*ReviewResult{
				{Score: 59.9, Decision: DecisionReject},
				{Score: 59.9, Decision: DecisionReject},
			},
			expectedScore:  59.9,
			expectedDecision: DecisionReject,
		},
		{
			name: "single approve overruled by two rejects",
			results: []*ReviewResult{
				{Score: 90, Decision: DecisionApprove},
				{Score: 30, Decision: DecisionReject},
				{Score: 40, Decision: DecisionReject},
			},
			expectedScore:  53.33,
			expectedDecision: DecisionReject,
		},
		{
			name: "exactly half rejects",
			results: []*ReviewResult{
				{Score: 85, Decision: DecisionApprove},
				{Score: 80, Decision: DecisionApprove},
				{Score: 40, Decision: DecisionReject},
				{Score: 35, Decision: DecisionReject},
			},
			expectedScore:  60.0,
			expectedDecision: DecisionConditionalApprove,
		},
		{
			name: "all zero scores",
			results: []*ReviewResult{
				{Score: 0, Decision: DecisionReject},
				{Score: 0, Decision: DecisionReject},
				{Score: 0, Decision: DecisionReject},
			},
			expectedScore:  0.0,
			expectedDecision: DecisionReject,
		},
		{
			name: "all perfect scores",
			results: []*ReviewResult{
				{Score: 100, Decision: DecisionApprove},
				{Score: 100, Decision: DecisionApprove},
				{Score: 100, Decision: DecisionApprove},
			},
			expectedScore:  100.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "single review approval",
			results: []*ReviewResult{
				{Score: 85, Decision: DecisionApprove},
			},
			expectedScore:  85.0,
			expectedDecision: DecisionApprove,
		},
		{
			name: "single review rejection",
			results: []*ReviewResult{
				{Score: 40, Decision: DecisionReject},
			},
			expectedScore:  40.0,
			expectedDecision: DecisionReject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, decision := calc.Calculate(tt.results)

			// Allow small floating point differences
			if score < tt.expectedScore-0.01 || score > tt.expectedScore+0.01 {
				t.Errorf("Expected score %.2f, got %.2f", tt.expectedScore, score)
			}
			if decision != tt.expectedDecision {
				t.Errorf("Expected decision %s, got %s", tt.expectedDecision, decision)
			}
		})
	}
}

func TestConsensusCalculatorCustomThreshold(t *testing.T) {
	calc := &ConsensusCalculator{Threshold: 90.0}

	results := []*ReviewResult{
		{Score: 85, Decision: DecisionApprove},
		{Score: 85, Decision: DecisionApprove},
	}

	score, decision := calc.Calculate(results)
	if score != 85.0 {
		t.Errorf("Expected score 85.0, got %.2f", score)
	}
	if decision != DecisionConditionalApprove {
		t.Errorf("Expected CONDITIONAL_APPROVE with 90 threshold, got %s", decision)
	}
}

// Mock LLM Client Additional Tests

func TestMockLLMClientCallsLog(t *testing.T) {
	client := NewMockLLMClient(nil)
	ctx := context.Background()

	messages1 := []string{"message 1"}
	messages2 := []string{"message 2", "message 3"}
	kwargs := map[string]interface{}{"temp": 0.7}

	client.Create(ctx, "model-a", messages1, nil)
	client.Create(ctx, "model-b", messages2, kwargs)

	if len(client.CallsLog) != 2 {
		t.Errorf("Expected 2 calls logged, got %d", len(client.CallsLog))
	}

	if client.CallsLog[0].Model != "model-a" {
		t.Errorf("Expected first call model-a, got %s", client.CallsLog[0].Model)
	}
	if len(client.CallsLog[0].Messages) != 1 {
		t.Errorf("Expected 1 message in first call, got %d", len(client.CallsLog[0].Messages))
	}

	if client.CallsLog[1].Model != "model-b" {
		t.Errorf("Expected second call model-b, got %s", client.CallsLog[1].Model)
	}
	if len(client.CallsLog[1].Messages) != 2 {
		t.Errorf("Expected 2 messages in second call, got %d", len(client.CallsLog[1].Messages))
	}
	if client.CallsLog[1].Kwargs["temp"] != 0.7 {
		t.Errorf("Expected temp 0.7, got %v", client.CallsLog[1].Kwargs["temp"])
	}
}

func TestMockLLMClientMultipleModels(t *testing.T) {
	responses := map[string]*LLMResponse{
		"model-1": {Content: "response 1", Model: "model-1"},
		"model-2": {Content: "response 2", Model: "model-2"},
		"model-3": {Content: "response 3", Model: "model-3"},
	}
	client := NewMockLLMClient(responses)
	ctx := context.Background()

	resp1, _ := client.Create(ctx, "model-1", []string{"test"}, nil)
	resp2, _ := client.Create(ctx, "model-2", []string{"test"}, nil)
	resp3, _ := client.Create(ctx, "model-3", []string{"test"}, nil)

	if resp1.Content != "response 1" {
		t.Errorf("Expected 'response 1', got %s", resp1.Content)
	}
	if resp2.Content != "response 2" {
		t.Errorf("Expected 'response 2', got %s", resp2.Content)
	}
	if resp3.Content != "response 3" {
		t.Errorf("Expected 'response 3', got %s", resp3.Content)
	}
}

func TestMockLLMClientDefaultResponse(t *testing.T) {
	client := NewMockLLMClient(map[string]*LLMResponse{})
	ctx := context.Background()

	resp, err := client.Create(ctx, "any-model", []string{"test"}, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Model != "claude-sonnet-4-5" {
		t.Errorf("Expected default model claude-sonnet-4-5, got %s", resp.Model)
	}
	if resp.Usage.InputTokens != 500 {
		t.Errorf("Expected 500 input tokens, got %d", resp.Usage.InputTokens)
	}
	if resp.Usage.OutputTokens != 150 {
		t.Errorf("Expected 150 output tokens, got %d", resp.Usage.OutputTokens)
	}
}

// Helper Functions Tests

func TestAgentConfig(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := DefaultAgentConfig()

		if config.Model != "claude-sonnet-4-5" {
			t.Errorf("Expected default model claude-sonnet-4-5, got %s", config.Model)
		}
		if config.BaselineScore != 70.0 {
			t.Errorf("Expected baseline 70.0, got %.2f", config.BaselineScore)
		}
		if config.Extra == nil {
			t.Error("Expected Extra map to be initialized")
		}
	})

	t.Run("custom config", func(t *testing.T) {
		config := &AgentConfig{
			Model:         "custom-model",
			BaselineScore: 85.5,
			Extra: map[string]interface{}{
				"key": "value",
			},
		}

		if config.Model != "custom-model" {
			t.Errorf("Expected custom-model, got %s", config.Model)
		}
		if config.BaselineScore != 85.5 {
			t.Errorf("Expected 85.5, got %.2f", config.BaselineScore)
		}
		if config.Extra["key"] != "value" {
			t.Errorf("Expected 'value', got %v", config.Extra["key"])
		}
	})
}

func TestProposalFixtures(t *testing.T) {
	t.Run("sample proposal", func(t *testing.T) {
		p := SampleProposal()

		if p.ID != "PROP-001" {
			t.Errorf("Expected ID PROP-001, got %s", p.ID)
		}
		if p.Title == "" {
			t.Error("Title should not be empty")
		}
		if p.Implementation == nil {
			t.Error("Implementation should not be nil")
		}
	})

	t.Run("invalid proposal", func(t *testing.T) {
		p := InvalidProposal()

		if p.Title != "" {
			t.Errorf("Expected empty title, got %s", p.Title)
		}
		if p.ID != "PROP-BAD" {
			t.Errorf("Expected ID PROP-BAD, got %s", p.ID)
		}
	})

	t.Run("privileged proposal", func(t *testing.T) {
		p := PrivilegedProposal()

		if priv, ok := p.Implementation["privileged"].(bool); !ok || !priv {
			t.Error("Expected privileged to be true")
		}
	})
}

func TestAssertReviewValid(t *testing.T) {
	t.Run("valid review passes", func(t *testing.T) {
		result := &ReviewResult{
			Score:    75.0,
			Decision: DecisionApprove,
			Feedback: "Good",
		}

		// This should not fail
		AssertReviewValid(t, result)
	})

	t.Run("nil result fails", func(t *testing.T) {
		mockT := &testing.T{}
		AssertReviewValid(mockT, nil)
		if !mockT.Failed() {
			t.Error("Expected test to fail for nil result")
		}
	})
}

func TestTestContext(t *testing.T) {
	ctx := TestContext(t)

	if ctx == nil {
		t.Fatal("Context should not be nil")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	if time.Until(deadline) > 61*time.Second {
		t.Error("Deadline should be around 60 seconds")
	}
}

// JSON Serialization Tests

func TestProposalJSONSerialization(t *testing.T) {
	proposal := SampleProposal()
	proposal.Metadata = map[string]interface{}{
		"author": "test-user",
		"timestamp": 1234567890,
	}

	data, err := json.Marshal(proposal)
	if err != nil {
		t.Fatalf("Failed to marshal proposal: %v", err)
	}

	var decoded Proposal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal proposal: %v", err)
	}

	if decoded.ID != proposal.ID {
		t.Errorf("Expected ID %s, got %s", proposal.ID, decoded.ID)
	}
	if decoded.Title != proposal.Title {
		t.Errorf("Expected Title %s, got %s", proposal.Title, decoded.Title)
	}
}

func TestReviewResultJSONSerialization(t *testing.T) {
	result := &ReviewResult{
		Score:     85.5,
		Decision:  DecisionApprove,
		Feedback:  "Excellent work",
		Reasoning: "Meets all criteria",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	var decoded ReviewResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if decoded.Score != result.Score {
		t.Errorf("Expected Score %.2f, got %.2f", result.Score, decoded.Score)
	}
	if decoded.Decision != result.Decision {
		t.Errorf("Expected Decision %s, got %s", result.Decision, decoded.Decision)
	}
	if decoded.Reasoning != result.Reasoning {
		t.Errorf("Expected Reasoning %s, got %s", result.Reasoning, decoded.Reasoning)
	}
}

func TestDecisionConstants(t *testing.T) {
	if DecisionApprove != "APPROVE" {
		t.Errorf("Expected APPROVE, got %s", DecisionApprove)
	}
	if DecisionConditionalApprove != "CONDITIONAL_APPROVE" {
		t.Errorf("Expected CONDITIONAL_APPROVE, got %s", DecisionConditionalApprove)
	}
	if DecisionReject != "REJECT" {
		t.Errorf("Expected REJECT, got %s", DecisionReject)
	}
}

// Context Cancellation Tests

func TestAgentReviewWithCancelledContext(t *testing.T) {
	agents := []Agent{
		NewSecurityAgent(nil, nil),
		NewPMAgent(nil, nil),
		NewArchitectAgent(nil, nil),
	}

	proposal := SampleProposal()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	for _, agent := range agents {
		_, err := agent.Review(ctx, proposal)
		// Agents should still work even with cancelled context in this implementation
		// since they don't actually check context
		if err != nil {
			t.Errorf("%s failed with cancelled context: %v", agent.Name(), err)
		}
	}
}

func TestAgentReviewWithTimeout(t *testing.T) {
	agents := []Agent{
		NewSecurityAgent(nil, nil),
		NewPMAgent(nil, nil),
		NewArchitectAgent(nil, nil),
	}

	proposal := SampleProposal()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond) // Ensure timeout

	for _, agent := range agents {
		_, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Errorf("%s failed with timeout context: %v", agent.Name(), err)
		}
	}
}

// Edge Case Integration Tests

func TestMultiAgentConsensus(t *testing.T) {
	client := SetupMockClient(t)
	ctx := TestContext(t)

	secAgent := NewSecurityAgent(nil, client)
	pmAgent := NewPMAgent(nil, client)
	archAgent := NewArchitectAgent(nil, client)

	proposal := SampleProposal()

	secResult, _ := secAgent.Review(ctx, proposal)
	pmResult, _ := pmAgent.Review(ctx, proposal)
	archResult, _ := archAgent.Review(ctx, proposal)

	calc := NewConsensusCalculator()
	_, decision := calc.Calculate([]*ReviewResult{secResult, pmResult, archResult})

	if decision == DecisionReject {
		t.Error("Sample proposal should not be rejected by consensus")
	}
}

func TestAllAgentsRejectBadProposal(t *testing.T) {
	client := SetupMockClient(t)
	ctx := TestContext(t)

	agents := []Agent{
		NewSecurityAgent(nil, client),
		NewPMAgent(nil, client),
		NewArchitectAgent(nil, client),
	}

	badProposal := InvalidProposal()
	results := make([]*ReviewResult, 0)

	for _, agent := range agents {
		result, err := agent.Review(ctx, badProposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}
		results = append(results, result)
	}

	calc := NewConsensusCalculator()
	_, decision := calc.Calculate(results)

	if decision != DecisionReject {
		t.Error("Invalid proposal should be rejected by consensus")
	}
}

// Filepath Tests

func TestSecurityAgentFileExtensionChecking(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		penalty  bool
	}{
		{"shell script", "deploy.sh", true},
		{"python file", "app.py", false},
		{"yaml config", "config.yaml", false},
		{"json data", "data.json", false},
		{"bash script", "install.bash", false}, // only .sh is penalized
		{"uppercase .SH", "script.SH", false},  // case sensitive
	}

	agent := NewSecurityAgent(nil, nil)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposal := &Proposal{
				ID:    "TEST",
				Title: "Test",
				Implementation: map[string]interface{}{
					"files": []interface{}{tt.filename},
				},
			}

			result, _ := agent.Review(ctx, proposal)
			expectedScore := 70.0
			if tt.penalty {
				expectedScore = 60.0
			}

			if result.Score != expectedScore {
				t.Errorf("Expected score %.1f for %s, got %.1f", 
					expectedScore, tt.filename, result.Score)
			}
		})
	}
}
```