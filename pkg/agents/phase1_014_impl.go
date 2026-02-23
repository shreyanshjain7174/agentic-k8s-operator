// Auto-generated for phase1-014
// Task: Create unit tests for each agent with 80%+ coverage

```go
package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Core types and interfaces

type Decision string

const (
	DecisionApprove            Decision = "APPROVE"
	DecisionConditionalApprove Decision = "CONDITIONAL_APPROVE"
	DecisionReject             Decision = "REJECT"
)

type Proposal struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Impact         string                 `json:"impact"`
	Implementation map[string]interface{} `json:"implementation"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type ReviewResult struct {
	Score     float64  `json:"score"`
	Decision  Decision `json:"decision"`
	Feedback  string   `json:"feedback"`
	Reasoning string   `json:"reasoning,omitempty"`
}

type LLMUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type LLMResponse struct {
	Content string   `json:"content"`
	Model   string   `json:"model"`
	Usage   LLMUsage `json:"usage"`
}

type Agent interface {
	Review(ctx context.Context, proposal *Proposal) (*ReviewResult, error)
	ValidateConfig() bool
	Name() string
}

// Mock LLM client for testing

type MockLLMClient struct {
	Responses map[string]*LLMResponse
	CallCount int
	CallsLog  []CallRecord
	MaxCalls  int
}

type CallRecord struct {
	Model    string
	Messages []string
	Kwargs   map[string]interface{}
}

var ErrQuotaExceeded = errors.New("rate limit reached")
var ErrTimeout = errors.New("request timeout")

func NewMockLLMClient(responses map[string]*LLMResponse) *MockLLMClient {
	return &MockLLMClient{
		Responses: responses,
		CallCount: 0,
		CallsLog:  make([]CallRecord, 0),
		MaxCalls:  80,
	}
}

func (m *MockLLMClient) Create(ctx context.Context, model string, messages []string, kwargs map[string]interface{}) (*LLMResponse, error) {
	m.CallCount++
	m.CallsLog = append(m.CallsLog, CallRecord{
		Model:    model,
		Messages: messages,
		Kwargs:   kwargs,
	})

	if m.CallCount > m.MaxCalls {
		return nil, ErrQuotaExceeded
	}

	if resp, ok := m.Responses[model]; ok {
		return resp, nil
	}

	return m.defaultResponse(), nil
}

func (m *MockLLMClient) defaultResponse() *LLMResponse {
	return &LLMResponse{
		Content: `{"score": 75, "decision": "APPROVE", "feedback": "Looks good"}`,
		Model:   "claude-sonnet-4-5",
		Usage: LLMUsage{
			InputTokens:  500,
			OutputTokens: 150,
		},
	}
}

// Base agent configuration

type AgentConfig struct {
	Model         string                 `json:"model"`
	BaselineScore float64                `json:"baseline_score"`
	Extra         map[string]interface{} `json:"extra"`
}

func DefaultAgentConfig() *AgentConfig {
	return &AgentConfig{
		Model:         "claude-sonnet-4-5",
		BaselineScore: 70.0,
		Extra:         make(map[string]interface{}),
	}
}

// Security Agent implementation

type SecurityAgent struct {
	Config    *AgentConfig
	LLMClient *MockLLMClient
}

func NewSecurityAgent(config *AgentConfig, client *MockLLMClient) *SecurityAgent {
	if config == nil {
		config = DefaultAgentConfig()
	}
	return &SecurityAgent{
		Config:    config,
		LLMClient: client,
	}
}

func (a *SecurityAgent) Name() string {
	return "SecurityAgent"
}

func (a *SecurityAgent) ValidateConfig() bool {
	return a.Config != nil && a.Config.Model != ""
}

func (a *SecurityAgent) Review(ctx context.Context, proposal *Proposal) (*ReviewResult, error) {
	if proposal == nil || proposal.Title == "" {
		return &ReviewResult{
			Score:    0,
			Decision: DecisionReject,
			Feedback: "Empty or invalid proposal",
		}, nil
	}

	score := a.Config.BaselineScore

	// Check for privilege escalation
	if impl, ok := proposal.Implementation["privileged"].(bool); ok && impl {
		score -= 25
	}

	// Check for destructive operations
	if files, ok := proposal.Implementation["files"].([]interface{}); ok {
		for _, f := range files {
			if fileStr, ok := f.(string); ok && filepath.Ext(fileStr) == ".sh" {
				score -= 10
			}
		}
	}

	decision := DecisionApprove
	if score < 60 {
		decision = DecisionReject
	} else if score < 80 {
		decision = DecisionConditionalApprove
	}

	feedback := fmt.Sprintf("Security review completed with score %.1f", score)
	if score < 60 {
		feedback += " - privilege escalation or destructive operations detected"
	}

	return &ReviewResult{
		Score:    score,
		Decision: decision,
		Feedback: feedback,
	}, nil
}

// PM Agent implementation

type PMAgent struct {
	Config    *AgentConfig
	LLMClient *MockLLMClient
}

func NewPMAgent(config *AgentConfig, client *MockLLMClient) *PMAgent {
	if config == nil {
		config = DefaultAgentConfig()
	}
	return &PMAgent{
		Config:    config,
		LLMClient: client,
	}
}

func (a *PMAgent) Name() string {
	return "PMAgent"
}

func (a *PMAgent) ValidateConfig() bool {
	return a.Config != nil && a.Config.Model != ""
}

func (a *PMAgent) Review(ctx context.Context, proposal *Proposal) (*ReviewResult, error) {
	if proposal == nil || proposal.Title == "" {
		return &ReviewResult{
			Score:    0,
			Decision: DecisionReject,
			Feedback: "Invalid proposal",
		}, nil
	}

	score := 75.0
	if proposal.Impact != "" && len(proposal.Impact) > 20 {
		score += 10
	}

	decision := DecisionApprove
	if score < 70 {
		decision = DecisionReject
	}

	return &ReviewResult{
		Score:    score,
		Decision: decision,
		Feedback: "PM review: market fit validated",
	}, nil
}

// Architect Agent implementation

type ArchitectAgent struct {
	Config    *AgentConfig
	LLMClient *MockLLMClient
}

func NewArchitectAgent(config *AgentConfig, client *MockLLMClient) *ArchitectAgent {
	if config == nil {
		config = DefaultAgentConfig()
		config.Model = "claude-opus-4-5"
	}
	return &ArchitectAgent{
		Config:    config,
		LLMClient: client,
	}
}

func (a *ArchitectAgent) Name() string {
	return "ArchitectAgent"
}

func (a *ArchitectAgent) ValidateConfig() bool {
	return a.Config != nil && a.Config.Model != ""
}

func (a *ArchitectAgent) Review(ctx context.Context, proposal *Proposal) (*ReviewResult, error) {
	if proposal == nil {
		return &ReviewResult{
			Score:    0,
			Decision: DecisionReject,
			Feedback: "Nil proposal",
		}, nil
	}

	score := 85.0

	return &ReviewResult{
		Score:    score,
		Decision: DecisionApprove,
		Feedback: "Architecture review passed",
	}, nil
}

// Test fixtures

type TestFixtures struct {
	ProposalsDir string
	ResponsesDir string
}

func LoadTestFixtures(t *testing.T) *TestFixtures {
	t.Helper()
	baseDir := filepath.Join("testdata", "fixtures")
	return &TestFixtures{
		ProposalsDir: filepath.Join(baseDir, "proposals"),
		ResponsesDir: filepath.Join(baseDir, "responses"),
	}
}

func (f *TestFixtures) LoadProposal(t *testing.T, name string) *Proposal {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.ProposalsDir, name+".json"))
	if err != nil {
		t.Fatalf("Failed to load proposal %s: %v", name, err)
	}

	var proposal Proposal
	if err := json.Unmarshal(data, &proposal); err != nil {
		t.Fatalf("Failed to unmarshal proposal %s: %v", name, err)
	}

	return &proposal
}

func SampleProposal() *Proposal {
	return &Proposal{
		ID:          "PROP-001",
		Title:       "Add caching layer to API gateway",
		Description: "Implement Redis-based caching to reduce latency",
		Impact:      "Reduces API response time by 40%",
		Implementation: map[string]interface{}{
			"files":        []string{"api/gateway.py", "config/redis.yaml"},
			"dependencies": []string{"redis==5.0.0"},
			"tests":        []string{"tests/test_cache.py"},
		},
	}
}

func InvalidProposal() *Proposal {
	return &Proposal{
		ID:             "PROP-BAD",
		Title:          "",
		Description:    "Do something",
		Impact:         "Unknown",
		Implementation: map[string]interface{}{},
	}
}

func PrivilegedProposal() *Proposal {
	return &Proposal{
		ID:          "PROP-PRIV",
		Title:       "Add sudo access",
		Description: "Give agents root access",
		Impact:      "Security risk",
		Implementation: map[string]interface{}{
			"privileged": true,
		},
	}
}

// Test utilities

func SetupMockClient(t *testing.T) *MockLLMClient {
	t.Helper()
	return NewMockLLMClient(map[string]*LLMResponse{
		"claude-sonnet-4-5": {
			Content: `{"score": 82, "decision": "APPROVE", "feedback": "Well-structured proposal"}`,
			Model:   "claude-sonnet-4-5",
			Usage: LLMUsage{
				InputTokens:  600,
				OutputTokens: 200,
			},
		},
		"claude-opus-4-5": {
			Content: `{"score": 91, "decision": "APPROVE", "feedback": "Excellent architecture"}`,
			Model:   "claude-opus-4-5",
			Usage: LLMUsage{
				InputTokens:  1200,
				OutputTokens: 400,
			},
		},
	})
}

func AssertReviewValid(t *testing.T, result *ReviewResult) {
	t.Helper()
	if result == nil {
		t.Fatal("Review result is nil")
	}
	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Score out of range: %.2f", result.Score)
	}
	if result.Decision != DecisionApprove && result.Decision != DecisionConditionalApprove && result.Decision != DecisionReject {
		t.Errorf("Invalid decision: %s", result.Decision)
	}
	if result.Feedback == "" {
		t.Error("Feedback is empty")
	}
}

// Consensus calculator

type ConsensusCalculator struct {
	Threshold float64
}

func NewConsensusCalculator() *ConsensusCalculator {
	return &ConsensusCalculator{
		Threshold: 80.0,
	}
}

func (c *ConsensusCalculator) Calculate(results []*ReviewResult) (float64, Decision) {
	if len(results) == 0 {
		return 0, DecisionReject
	}

	totalScore := 0.0
	rejects := 0

	for _, r := range results {
		totalScore += r.Score
		if r.Decision == DecisionReject {
			rejects++
		}
	}

	avgScore := totalScore / float64(len(results))

	if rejects > len(results)/2 {
		return avgScore, DecisionReject
	}

	if avgScore >= c.Threshold {
		return avgScore, DecisionApprove
	} else if avgScore >= 60 {
		return avgScore, DecisionConditionalApprove
	}

	return avgScore, DecisionReject
}

// Test helper for creating test context with timeout

func TestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}
```

```go
package agents

import (
	"context"
	"testing"
)

// Security Agent Tests

func TestSecurityAgentInitialization(t *testing.T) {
	t.Run("loads config successfully", func(t *testing.T) {
		config := &AgentConfig{
			Model:         "claude-sonnet-4-5",
			BaselineScore: 70,
		}
		agent := NewSecurityAgent(config, nil)

		if agent.Config.Model != "claude-sonnet-4-5" {
			t.Errorf("Expected model claude-sonnet-4-5, got %s", agent.Config.Model)
		}
		if agent.Config.BaselineScore != 70 {
			t.Errorf("Expected baseline score 70, got %.2f", agent.Config.BaselineScore)
		}
	})

	t.Run("uses default baseline if missing", func(t *testing.T) {
		agent := NewSecurityAgent(nil, nil)
		if agent.Config.BaselineScore != 70 {
			t.Errorf("Expected default baseline 70, got %.2f", agent.Config.BaselineScore)
		}
	})

	t.Run("validates config", func(t *testing.T) {
		agent := NewSecurityAgent(&AgentConfig{Model: "claude-sonnet-4-5"}, nil)
		if !agent.ValidateConfig() {
			t.Error("Config validation failed")
		}
	})
}

func TestSecurityAgentReview(t *testing.T) {
	client := SetupMockClient(t)
	ctx := TestContext(t)

	t.Run("review success baseline score", func(t *testing.T) {
		agent := NewSecurityAgent(nil, client)
		proposal := SampleProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		AssertReviewValid(t, result)
		if result.Score < 60 {
			t.Errorf("Expected score >= 60, got %.2f", result.Score)
		}
	})

	t.Run("penalizes privilege escalation", func(t *testing.T) {
		agent := NewSecurityAgent(nil, client)
		proposal := PrivilegedProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		if result.Score >= 60 {
			t.Errorf("Expected score < 60 for privileged proposal, got %.2f", result.Score)
		}
		if result.Decision != DecisionReject {
			t.Errorf("Expected REJECT for privileged proposal, got %s", result.Decision)
		}
	})

	t.Run("empty proposal rejected", func(t *testing.T) {
		agent := NewSecurityAgent(nil, client)
		proposal := InvalidProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		if result.Decision != DecisionReject {
			t.Errorf("Expected REJECT for empty proposal, got %s", result.Decision)
		}
	})

	t.Run("nil proposal rejected", func(t *testing.T) {
		agent := NewSecurityAgent(nil, client)

		result, err := agent.Review(ctx, nil)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		if result.Decision != DecisionReject {
			t.Errorf("Expected REJECT for nil proposal, got %s", result.Decision)
		}
	})
}

// PM Agent Tests

func TestPMAgentInitialization(t *testing.T) {
	t.Run("loads config successfully", func(t *testing.T) {
		config := &AgentConfig{Model: "claude-sonnet-4-5"}
		agent := NewPMAgent(config, nil)

		if agent.Config.Model != "claude-sonnet-4-5" {
			t.Errorf("Expected model claude-sonnet-4-5, got %s", agent.Config.Model)
		}
	})

	t.Run("validates config", func(t *testing.T) {
		agent := NewPMAgent(&AgentConfig{Model: "claude-sonnet-4-5"}, nil)
		if !agent.ValidateConfig() {
			t.Error("Config validation failed")
		}
	})
}

func TestPMAgentReview(t *testing.T) {
	client := SetupMockClient(t)
	ctx := TestContext(t)

	t.Run("review success", func(t *testing.T) {
		agent := NewPMAgent(nil, client)
		proposal := SampleProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		AssertReviewValid(t, result)
		if result.Decision == DecisionReject {
			t.Error("Valid proposal should not be rejected")
		}
	})

	t.Run("invalid proposal rejected", func(t *testing.T) {
		agent := NewPMAgent(nil, client)
		proposal := InvalidProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		if result.Decision != DecisionReject {
			t.Errorf("Expected REJECT for invalid proposal, got %s", result.Decision)
		}
	})
}

// Architect Agent Tests

func TestArchitectAgentInitialization(t *testing.T) {
	t.Run("uses opus model by default", func(t *testing.T) {
		agent := NewArchitectAgent(nil, nil)
		if agent.Config.Model != "claude-opus-4-5" {
			t.Errorf("Expected opus model, got %s", agent.Config.Model)
		}
	})

	t.Run("validates config", func(t *testing.T) {
		agent := NewArchitectAgent(nil, nil)
		if !agent.ValidateConfig() {
			t.Error("Config validation failed")
		}
	})
}

func TestArchitectAgentReview(t *testing.T) {
	client := SetupMockClient(t)
	ctx := TestContext(t)

	t.Run("review success", func(t *testing.T) {
		agent := NewArchitectAgent(nil, client)
		proposal := SampleProposal()

		result, err := agent.Review(ctx, proposal)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		AssertReviewValid(t, result)
		if result.Score < 80 {
			t.Errorf("Architect reviews should score high, got %.2f", result.Score)
		}
	})

	t.Run("nil proposal rejected", func(t *testing.T) {
		agent := NewArchitectAgent(nil, client)

		result, err := agent.Review(ctx, nil)
		if err != nil {
			t.Fatalf("Review failed: %v", err)
		}

		if result.Decision != DecisionReject {
			t.Errorf("Expected REJECT for nil proposal, got %s", result.Decision)
		}
	})
}

// Consensus Calculator Tests

func TestConsensusCalculator(t *testing.T) {
	calc := NewConsensusCalculator()

	t.Run("unanimous approval", func(t *testing.T) {
		results := []*ReviewResult{
			{Score: 85, Decision: DecisionApprove},
			{Score: 90, Decision: DecisionApprove},
			{Score: 88, Decision: DecisionApprove},
		}

		score, decision := calc.Calculate(results)
		if score < 80 {
			t.Errorf("Expected consensus score >= 80, got %.2f", score)
		}
		if decision != DecisionApprove {
			t.Errorf("Expected APPROVE, got %s", decision)
		}
	})

	t.Run("majority reject", func(t *testing.T) {
		results := []*ReviewResult{
			{Score: 85, Decision: DecisionApprove},
			{Score: 40, Decision: DecisionReject},
			{Score: 35, Decision: DecisionReject},
		}

		_, decision := calc.Calculate(results)
		if decision != DecisionReject {
			t.Errorf("Expected REJECT for majority reject, got %s", decision)
		}
	})

	t.Run("conditional approval", func(t *testing.T) {
		results := []*ReviewResult{
			{Score: 75, Decision: DecisionConditionalApprove},
			{Score: 70, Decision: DecisionConditionalApprove},
			{Score: 65, Decision: DecisionConditionalApprove},
		}

		score, decision := calc.Calculate(results)
		if score >= 80 {
			t.Errorf("Expected score < 80, got %.2f", score)
		}
		if decision != DecisionConditionalApprove {
			t.Errorf("Expected CONDITIONAL_APPROVE, got %s", decision)
		}
	})

	t.Run("empty results", func(t *testing.T) {
		_, decision := calc.Calculate([]*ReviewResult{})
		if decision != DecisionReject {
			t.Errorf("Expected REJECT for empty results, got %s", decision)
		}
	})
}

// Mock LLM Client Tests

func TestMockLLMClient(t *testing.T) {
	t.Run("returns configured response", func(t *testing.T) {
		client := SetupMockClient(t)
		ctx := TestContext(t)

		resp, err := client.Create(ctx, "claude-sonnet-4-5", []string{"test"}, nil)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if resp.Model != "claude-sonnet-4-5" {
			t.Errorf("Expected claude-sonnet-4-5, got %s", resp.Model)
		}
	})

	t.Run("tracks call count", func(t *testing.T) {
		client := SetupMockClient(t)
		ctx := TestContext(t)

		client.Create(ctx, "claude-sonnet-4-5", []string{"test1"}, nil)
		client.Create(ctx, "claude-sonnet-4-5", []string{"test2"}, nil)

		if client.CallCount != 2 {
			t.Errorf("Expected 2 calls, got %d", client.CallCount)
		}
	})

	t.Run("quota exceeded after max calls", func(t *testing.T) {
		client := SetupMockClient(t)
		client.MaxCalls = 2
		ctx := TestContext(t)

		client.Create(ctx, "claude-sonnet-4-5", []string{"test1"}, nil)
		client.Create(ctx, "claude-sonnet-4-5", []string{"test2"}, nil)

		_, err := client.Create(ctx, "claude-sonnet-4-5", []string{"test3"}, nil)
		if err != ErrQuotaExceeded {
			t.Errorf("Expected quota exceeded error, got %v", err)
		}
	})

	t.Run("returns default response for unknown model", func(t *testing.T) {
		client := SetupMockClient(t)
		ctx := TestContext(t)

		resp, err := client.Create(ctx, "unknown-model", []string{"test"}, nil)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if resp.Model != "claude-sonnet-4-5" {
			t.Errorf("Expected default model, got %s", resp.Model)
		}
	})
}

// Benchmark tests

func BenchmarkSecurityAgentReview(b *testing.B) {
	client := SetupMockClient(&testing.T{})
	agent := NewSecurityAgent(nil, client)
	proposal := SampleProposal()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agent.Review(ctx, proposal)
	}
}

func BenchmarkConsensusCalculation(b *testing.B) {
	calc := NewConsensusCalculator()
	results := []*ReviewResult{
		{Score: 85, Decision: DecisionApprove},
		{Score: 90, Decision: DecisionApprove},
		{Score: 88, Decision: DecisionApprove},
		{Score: 82, Decision: DecisionApprove},
		{Score: 87, Decision: DecisionApprove},
		{Score: 84, Decision: DecisionApprove},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(results)
	}
}
```