// Auto-generated tests for phase1-002

```go
// pkg/agents/pm/types_test.go
package pm

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProposalJSONMarshaling(t *testing.T) {
	now := time.Now()
	proposal := &Proposal{
		ID:               "PRD-001",
		Title:            "Test Proposal",
		CreatedAt:        now,
		Author:           "test-agent",
		Status:           StatusDraft,
		ProblemStatement: "Test problem",
		Goals:            []string{"Goal 1", "Goal 2"},
		ConsensusScore:   85.5,
	}

	data, err := json.Marshal(proposal)
	if err != nil {
		t.Fatalf("failed to marshal proposal: %v", err)
	}

	var decoded Proposal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal proposal: %v", err)
	}

	if decoded.ID != proposal.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, proposal.ID)
	}
	if decoded.ConsensusScore != proposal.ConsensusScore {
		t.Errorf("ConsensusScore mismatch: got %f, want %f", decoded.ConsensusScore, proposal.ConsensusScore)
	}
}

func TestValidationResultJSONMarshaling(t *testing.T) {
	result := ValidationResult{
		Valid:    false,
		Errors:   []string{"error1", "error2"},
		Warnings: []string{"warning1"},
		Score:    75,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ValidationResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Valid != result.Valid {
		t.Errorf("Valid mismatch")
	}
	if decoded.Score != result.Score {
		t.Errorf("Score mismatch")
	}
	if len(decoded.Errors) != 2 {
		t.Errorf("Errors length mismatch")
	}
}

// pkg/agents/pm/validator_test.go
package pm

import (
	"testing"
)

func TestValidateCompleteness(t *testing.T) {
	tests := []struct {
		name          string
		proposal      *Proposal
		wantValid     bool
		wantMinScore  int
		wantErrors    int
		wantWarnings  int
	}{
		{
			name: "complete proposal",
			proposal: &Proposal{
				ProblemStatement: "Clear problem statement",
				Goals:            []string{"Goal 1"},
				SuccessMetrics:   []SuccessMetric{{Name: "Metric1", Target: "100%", Measurable: true}},
				Requirements:     []Requirement{{ID: "REQ-001", Testable: true, TestMethod: "Unit test"}},
				TechnicalApproach: TechnicalApproach{Overview: "Technical overview"},
				EdgeCases:        []EdgeCase{{Scenario: "Edge case 1"}},
			},
			wantValid:    true,
			wantMinScore: 90,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:         "empty proposal",
			proposal:     &Proposal{},
			wantValid:    false,
			wantMinScore: 0,
			wantErrors:   4, // missing problem, metrics, requirements, approach
			wantWarnings: 2, // no goals, no edge cases
		},
		{
			name: "missing problem statement",
			proposal: &Proposal{
				Goals:          []string{"Goal 1"},
				SuccessMetrics: []SuccessMetric{{Name: "M1"}},
				Requirements:   []Requirement{{ID: "REQ-001"}},
				TechnicalApproach: TechnicalApproach{Overview: "Overview"},
			},
			wantValid:    false,
			wantMinScore: 50,
			wantErrors:   1,
		},
		{
			name: "untestable requirements",
			proposal: &Proposal{
				ProblemStatement: "Problem",
				SuccessMetrics:   []SuccessMetric{{Name: "M1"}},
				Requirements: []Requirement{
					{ID: "REQ-001", Testable: false},
					{ID: "REQ-002", Testable: true, TestMethod: ""},
				},
				TechnicalApproach: TechnicalApproach{Overview: "Overview"},
			},
			wantValid:    true,
			wantMinScore: 70,
			wantWarnings: 3, // 1 untestable + 1 missing test method + no goals + no edge cases
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewProposalValidator(nil)
			result := validator.ValidateCompleteness(tt.proposal)

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", result.Valid, tt.wantValid)
			}
			if result.Score < tt.wantMinScore {
				t.Errorf("Score = %d, want >= %d", result.Score, tt.wantMinScore)
			}
			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d. Errors: %v", len(result.Errors), tt.wantErrors, result.Errors)
			}
			if tt.wantWarnings > 0 && len(result.Warnings) < tt.wantWarnings {
				t.Errorf("Warnings count = %d, want >= %d", len(result.Warnings), tt.wantWarnings)
			}
		})
	}
}

func TestValidateClarity(t *testing.T) {
	tests := []struct {
		name         string
		proposal     *Proposal
		wantValid    bool
		wantMinScore int
	}{
		{
			name: "clear proposal",
			proposal: &Proposal{
				ProblemStatement: "This is a clear problem statement with enough detail to understand the context and requirements for the proposed solution",
				TechnicalApproach: TechnicalApproach{
					Overview: "This is a detailed technical approach that explains the architecture, components, integration points, and implementation strategy in sufficient depth",
				},
				SuccessMetrics: []SuccessMetric{
					{Name: "Metric1", Target: "95%", Measurable: true},
				},
			},
			wantValid:    true,
			wantMinScore: 100,
		},
		{
			name: "vague language",
			proposal: &Proposal{
				ProblemStatement: "We should probably implement this feature because it might help users",
				TechnicalApproach: TechnicalApproach{
					Overview: "We could use this approach, maybe with some optimizations that should improve performance",
				},
				SuccessMetrics: []SuccessMetric{{Name: "M1", Measurable: true}},
			},
			wantValid:    false,
			wantMinScore: 0,
		},
		{
			name: "unmeasurable metrics",
			proposal: &Proposal{
				ProblemStatement: "Clear problem statement that is long enough to pass validation requirements",
				TechnicalApproach: TechnicalApproach{
					Overview: "Detailed technical approach with sufficient length to meet the minimum requirements for validation purposes",
				},
				SuccessMetrics: []SuccessMetric{
					{Name: "User happiness", Measurable: false},
					{Name: "Code quality", Measurable: false, Target: ""},
				},
			},
			wantValid:    false,
			wantMinScore: 50,
		},
		{
			name: "brief statements",
			proposal: &Proposal{
				ProblemStatement:  "Short problem",
				TechnicalApproach: TechnicalApproach{Overview: "Brief approach"},
				SuccessMetrics:    []SuccessMetric{{Name: "M1", Measurable: true, Target: "100%"}},
			},
			wantValid:    true,
			wantMinScore: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewProposalValidator(nil)
			result := validator.ValidateClarity(tt.proposal)

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v. Errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
			if result.Score < tt.wantMinScore {
				t.Errorf("Score = %d, want >= %d", result.Score, tt.wantMinScore)
			}
		})
	}
}

func TestValidateFeasibility(t *testing.T) {
	tests := []struct {
		name         string
		proposal     *Proposal
		wantValid    bool
		wantMinScore int
	}{
		{
			name: "feasible proposal",
			proposal: &Proposal{
				Requirements:     []Requirement{{ID: "REQ-001"}, {ID: "REQ-002"}},
				EffortEstimate:   EffortEstimate{Size: SizeS},
				RiskAssessment:   RiskAssessment{OverallRisk: RiskLow},
				Dependencies:     []Dependency{{Name: "dep1", Blocking: false}},
			},
			wantValid:    true,
			wantMinScore: 90,
		},
		{
			name: "critical risk without details",
			proposal: &Proposal{
				RiskAssessment: RiskAssessment{OverallRisk: RiskCritical, Risks: []Risk{}},
			},
			wantValid:    false,
			wantMinScore: 50,
		},
		{
			name: "size mismatch",
			proposal: &Proposal{
				Requirements:   make([]Requirement, 25),
				EffortEstimate: EffortEstimate{Size: SizeS},
				RiskAssessment: RiskAssessment{OverallRisk: RiskLow},
			},
			wantValid:    true,
			wantMinScore: 60,
		},
		{
			name: "blocking dependency without description",
			proposal: &Proposal{
				Dependencies: []Dependency{
					{Name: "dep1", Blocking: true, Description: ""},
				},
				RiskAssessment: RiskAssessment{OverallRisk: RiskLow},
			},
			wantValid:    true,
			wantMinScore: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewProposalValidator(nil)
			result := validator.ValidateFeasibility(tt.proposal)

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v. Errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
			if result.Score < tt.wantMinScore {
				t.Errorf("Score = %d, want >= %d", result.Score, tt.wantMinScore)
			}
		})
	}
}

func TestValidatorCustomBannedWords(t *testing.T) {
	validator := NewProposalValidator([]string{"bad", "terrible"})
	
	proposal := &Proposal{
		ProblemStatement: "This is a bad idea that terrible people made",
		TechnicalApproach: TechnicalApproach{
			Overview: "We will implement this terrible solution using bad practices to deliver the feature requirements",
		},
		SuccessMetrics: []SuccessMetric{{Name: "M1", Measurable: true}},
	}

	result := validator.ValidateClarity(proposal)
	if result.Valid {
		t.Error("Expected invalid due to banned words")
	}
	if len(result.Errors) < 2 {
		t.Errorf("Expected at least 2 errors for banned words, got %d", len(result.Errors))
	}
}

// pkg/agents/pm/generator_test.go
package pm

import (
	"context"
	"testing"
	"time"
)

type mockLLMClient struct {
	response string
	err      error
}

func (m *mockLLMClient) Complete(ctx context.Context, prompt string, model string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func TestGenerateProposal(t *testing.T) {
	validJSON := `{
		"title": "Test Proposal",
		"problem_statement": "Test problem",
		"goals": ["Goal 1"],
		"success_metrics": [{"name": "M1", "target": "100%", "measurable": true}],
		"technical_approach": {
			"overview": "Test approach",
			"components": ["C1"],
			"technologies": ["Go"],
			"integration_points": ["I1"],
			"alternatives_considered": []
		}
	}`

	tests := []struct {
		name        string
		llmResponse string
		llmError    error
		wantError   bool
	}{
		{
			name:        "valid json response",
			llmResponse: validJSON,
			wantError:   false,
		},
		{
			name:        "json with markdown",
			llmResponse: "```json\n" + validJSON + "\n```",
			wantError:   false,
		},
		{
			name:        "json with generic markdown",
			llmResponse: "```\n" + validJSON + "\n```",
			wantError:   false,
		},
		{
			name:        "llm error",
			llmError:    context.DeadlineExceeded,
			wantError:   true,
		},
		{
			name:        "invalid json",
			llmResponse: "{invalid json",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := &mockLLMClient{response: tt.llmResponse, err: tt.llmError}
			generator := NewProposalGenerator(llm, "test-model")

			task := Task{
				ID:          "task-1",
				Title:       "Test Task",
				Description: "Test Description",
				Priority:    PriorityP1,
			}

			wsContext := WorkspaceContext{
				Architecture: "Test arch",
				Roadmap:      "Test roadmap",
			}

			proposal, err := generator.Generate(context.Background(), task, wsContext)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if proposal.ID == "" {
				t.Error("Proposal ID is empty")
			}
			if proposal.Author != "pm-agent" {
				t.Errorf("Author = %s, want pm-agent", proposal.Author)
			}
			if proposal.Status != StatusDraft {
				t.Errorf("Status = %s, want %s", proposal.Status, StatusDraft)
			}
			if proposal.CreatedAt.IsZero() {
				t.Error("CreatedAt is zero")
			}
		})
	}
}

func TestFormatCommits(t *testing.T) {
	llm := &mockLLMClient{}
	generator := NewProposalGenerator(llm, "test-model").(*llmProposalGenerator)

	tests := []struct {
		name    string
		commits []Commit
		want    string
	}{
		{
			name:    "empty commits",
			commits: []Commit{},
			want:    "No recent commits",
		},
		{
			name: "single commit",
			commits: []Commit{
				{Hash: "abcdef1234567890", Author: "test-user", Message: "Test commit"},
			},
			want: "- abcdef1: Test commit (test-user)\n",
		},
		{
			name: "multiple commits",
			commits: []Commit{
				{Hash: "1111111111111111", Author: "user1", Message: "Commit 1"},
				{Hash: "2222222222222222", Author: "user2", Message: "Commit 2"},
			},
			want: "- 1111111: Commit 1 (user1)\n- 2222222: Commit 2 (user2)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.formatCommits(tt.commits)
			if got != tt.want {
				t.Errorf("formatCommits() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateProposalID(t *testing.T) {
	id1 := generateProposalID()
	time.Sleep(time.Millisecond)
	id2 := generateProposalID()

	if id1 == id2 {
		t.Error("Expected unique IDs")
	}
	if id1[:4] != "PRD-" {
		t.Errorf("ID should start with PRD-, got %s", id1)
	}
}

// pkg/agents/pm/requirements_analyzer_test.go
package pm

import (
	"context"
	"testing"
)

func TestAnalyzeRequirements(t *testing.T) {
	tests := []struct {
		name        string
		proposal    *Proposal
		llmResponse string
		wantError   bool
		wantCount   int
	}{
		{
			name: "existing requirements enhanced",
			proposal: &Proposal{
				Requirements: []Requirement{
					{Description: "Req 1", Testable: false},
				},
			},
			wantError: false,
			wantCount: 1,
		},
		{
			name: "generate new requirements",
			proposal: &Proposal{
				ProblemStatement: "Test problem",
				TechnicalApproach: TechnicalApproach{Overview: "Test approach"},
			},
			llmResponse: `[
				{
					"id": "REQ-001",
					"type": "functional",
					"description": "Test requirement",
					"priority": "P0",
					"testable": true,
					"test_method": "Unit test"
				}
			]`,
			wantError: false,
			wantCount: 1,
		},
		{
			name: "invalid json response",
			proposal: &Proposal{
				ProblemStatement: "Test problem",
			},
			llmResponse: "{invalid",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := &mockLLMClient{response: tt.llmResponse}
			analyzer := NewRequirementsAnalyzer(llm, "test-model")

			reqs, err := analyzer.AnalyzeRequirements(context.Background(), tt.proposal)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(reqs) != tt.wantCount {
				t.Errorf("Requirements count = %d, want %d", len(reqs), tt.wantCount)
			}

			// Check enhancement worked
			for _, req := range reqs {
				if req.ID == "" {
					t.Error("Requirement ID is empty")
				}
				if !req.Testable {
					t.Error("Requirement should be testable")
				}
			}
		})
	}
}

func TestBuildDependencyGraph(t *testing.T) {
	tests := []struct {
		name         string
		requirements []Requirement
		wantCount    int
	}{
		{
			name:         "no dependencies",
			requirements: []Requirement{{ID: "REQ-001"}},
			wantCount:    0,
		},
		{
			name: "single dependency",
			requirements: []Requirement{
				{ID: "REQ-001", DependsOn: []string{"REQ-002"}, Priority: PriorityP0},
			},
			wantCount: 1,
		},
		{
			name: "multiple dependencies",
			requirements: []Requirement{
				{ID: "REQ-001", DependsOn: []string{"REQ-002", "REQ-003"}},
				{ID: "REQ-004", DependsOn: []string{"REQ-001"}},
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewRequirementsAnalyzer(&mockLLMClient{}, "test-model")
			deps, err := analyzer.BuildDependencyGraph(context.Background(), tt.requirements)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(deps) != tt.wantCount {
				t.Errorf("Dependencies count = %d, want %d", len(deps), tt.wantCount)
			}

			// Verify blocking flag for P0
			for _, dep := range deps {
				if dep.Type != DepInternal {
					t.Error("Dependency should be internal")
				}
			}
		})
	}
}

func TestEnumerateEdgeCases(t *testing.T) {
	llmResponse := `[
		{
			"scenario": "Network timeout",
			"impact": "Request fails",
			"mitigation": "Retry logic"
		}
	]`

	llm := &mockLLMClient{response: llmResponse}
	analyzer := NewRequirementsAnalyzer(llm, "test-model")

	proposal := &Proposal{
		Title: "Test",
		Requirements: []Requirement{
			{ID: "REQ-001", Description: "Test requirement"},
		},
	}

	edgeCases, err := analyzer.EnumerateEdgeCases(context.Background(), proposal)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(edgeCases) != 1 {
		t.Errorf("EdgeCases count = %d, want 1", len(edgeCases))
	}

	if edgeCases[0].Scenario != "Network timeout" {
		t.Errorf("Scenario = %s, want 'Network timeout'", edgeCases[0].Scenario)
	}
}

// pkg/agents/pm/business_alignment_test.go
package pm

import (
	"context"
	"testing"
)

func TestEstimateEffort(t *testing.T) {
	tests := []struct {
		name        string
		numReqs     int
		numDeps     int
		wantSize    TShirtSize
		wantPoints  int
	}{
		{
			name:       "small task",
			numReqs:    2,
			numDeps:    0,
			wantSize:   SizeS,
			wantPoints: 1,
		},
		{
			name:       "medium task",
			numReqs:    5,
			numDeps:    1,
			wantSize:   SizeM,
			wantPoints: 3,
		},
		{
			name:       "large task",
			numReqs:    12,
			numDeps:    4,
			wantSize:   SizeL,
			wantPoints: 8,
		},
		{
			name:       "xl task",
			numReqs:    20,
			numDeps:    8,
			wantSize:   SizeXL,
			wantPoints: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewBusinessAlignmentEngine(&mockLLMClient{}, "test-model")

			proposal := &Proposal{
				Requirements: make([]Requirement, tt.numReqs),
				Dependencies: make([]Dependency, tt.numDeps),
			}

			estimate, err := engine.EstimateEffort(context.Background(), proposal)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if estimate.Size != tt.wantSize {
				t.Errorf("Size = %s, want %s", estimate.Size, tt.wantSize)
			}
			if estimate.StoryPoints != tt.wantPoints {
				t.Errorf("StoryPoints = %d, want %d", estimate.StoryPoints, tt.wantPoints)
			}
			if estimate.Total == "" {
				t.Error("Total estimate is empty")
			}
		})
	}
}

func TestCalculateImpact(t *testing.T) {
	llmResponse := `{
		"revenue_impact": "high",
		"compliance_impact": "medium",
		"performance_impact": "low",
		"customer_impact": "high",
		"rationale": "Test rationale"
	}`

	llm := &mockLLMClient{response: llmResponse}
	engine := NewBusinessAlignmentEngine(llm, "test-model")

	proposal := &Proposal{
		Title:            "Test Proposal",
		ProblemStatement: "Test problem",
		TechnicalApproach: TechnicalApproach{Overview: "Test approach"},
	}

	impact, err := engine.CalculateImpact(context.Background(), proposal)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if impact.RevenueImpact != ImpactHigh {
		t.Errorf("RevenueImpact = %s, want %s", impact.RevenueImpact, ImpactHigh)
	}
	if impact.Rationale == "" {
		t.Error("Rationale is empty")
	}
}

func TestMapToRoadmap(t *testing.T) {
	llmResponse := `{
		"phase": "Phase 1",
		"milestone": "MVP",
		"okrs": ["OKR-1", "OKR-2"],
		"blocked_by": []
	}`

	llm := &mockLLMClient{response: llmResponse}
	engine := NewBusinessAlignmentEngine(llm, "test-model")

	proposal := &Proposal{
		Title:            "Test",
		ProblemStatement: "Problem",
	}

	alignment, err := engine.MapToRoadmap(context.Background(), proposal, "Test roadmap")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if alignment.Phase != "Phase 1" {
		t.Errorf("Phase = %s, want Phase 1", alignment.Phase)
	}
	if len(alignment.OKRs) != 2 {
		t.Errorf("OKRs count = %d, want 2", len(alignment.OKRs))
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     string
	}{
		{
			name:     "plain json",
			response: `{"key": "value"}`,
			want:     `{"key": "value"}`,
		},
		{
			name:     "json with markdown",
			response: "```json\n{\"key\": \"value\"}\n```",
			want:     `{"key": "value"}`,
		},
		{
			name:     "json with text",
			response: "Here is the result: {\"key\": \"value\"} and more text",
			want:     `{"key": "value"}`,
		},
		{
			name:     "json with whitespace",
			response: "\n\n  {\"key\": \"value\"}  \n",
			want:     `{"key": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractJSON(tt.response)
			if got != tt.want {
				t.Errorf("extractJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

// pkg/agents/pm/git_integration_test.go
package pm

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type mockGitClient struct {
	addCalled    bool
	commitCalled bool
	addError     error
	commitError  error
}

func (m *mockGitClient) Add(path string) error {
	m.addCalled = true
	return m.addError
}

func (m *mockGitClient) Commit(message string) error {
	m.commitCalled = true
	return m.commitError
}

func (m *mockGitClient) Push() error {
	return nil
}

func TestCommitProposal(t *testing.T) {
	tmpDir := t.TempDir()

	proposal := &Proposal{
		ID:               "PRD-123",
		Title:            "Test Proposal",
		ProblemStatement: "Test problem",
		Status:           StatusDraft,
		EffortEstimate:   EffortEstimate{Size: SizeM},
	}

	git := &mockGitClient{}
	integration := NewGitIntegration(tmpDir, git)

	err := integration.CommitProposal(context.Background(), proposal)
	if err != nil {
		t.Fatalf("CommitProposal failed: %v", err)
	}

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "docs", "proposals", "PRD-123.json")
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read proposal file: %v", err)
	}

	var loaded Proposal
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal proposal: %v", err)
	}

	if loaded.ID != proposal.ID {
		t.Errorf("Loaded ID = %s, want %s", loaded.ID, proposal.ID)
	}

	// Verify git operations were called
	if !git.addCalled {
		t.Error("Git add was not called")
	}
	if !git.commitCalled {
		t.Error("Git commit was not called")
	}
}

func TestUpdateProposal(t *testing.T) {
	tmpDir := t.TempDir()

	proposal := &Proposal{
		ID:     "PRD-456",
		Title:  "Updated Proposal",
		Status: StatusApproved,
	}

	git := &mockGitClient{}
	integration := NewGitIntegration(tmpDir, git)

	err := integration.UpdateProposal(context.Background(), proposal)
	if err != nil {
		t.Fatalf("UpdateProposal failed: %v", err)
	}

	if !git.addCalled || !git.commitCalled {
		t.Error("Git operations were not called")
	}
}

// pkg/agents/pm/context_provider_test.go
package pm

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	clawdlinuxDir := filepath.Join(tmpDir, "clawdlinux")
	os.MkdirAll(clawdlinuxDir, 0755)

	archContent := "# Architecture\nTest architecture content"
	os.WriteFile(filepath.Join(clawdlinuxDir, "ARCHITECTURE.md"), []byte(archContent), 0644)

	roadmapContent := "# Roadmap\nTest roadmap content"
	os.WriteFile(filepath.Join(clawdlinuxDir, "ROADMAP.md"), []byte(roadmapContent), 0644)

	memoryContent := "# Memory\nTest memory content"
	os.WriteFile(filepath.Join(tmpDir, "MEMORY.md"), []byte(memoryContent), 0644)

	provider := NewContextProvider(tmpDir, &mockGitClient{})

	// Test LoadArchitecture
	arch, err := provider.LoadArchitecture(context.Background())
	if err != nil {
		t.Errorf("LoadArchitecture failed: %v", err)
	}
	if arch != archContent {
		t.Errorf("Architecture content mismatch")
	}

	// Test LoadRoadmap
	roadmap, err := provider.LoadRoadmap(context.Background())
	if err != nil {
		t.Errorf("LoadRoadmap failed: %v", err)
	}
	if roadmap != roadmapContent {
		t.Errorf("Roadmap content mismatch")
	}

	// Test LoadMemory
	memory, err := provider.LoadMemory(context.Background())
	if err != nil {
		t.Errorf("LoadMemory failed: %v", err)
	}
	if memory != memoryContent {
		t.Errorf("Memory content mismatch")
	}
}

func TestGetExistingProposals(t *testing.T) {
	tmpDir := t.TempDir()
	proposalsDir := filepath.Join(tmpDir, "clawdlinux", "docs", "proposals")
	os.MkdirAll(proposalsDir, 0755)

	// Create test proposals
	proposals := []*Proposal{
		{ID: "PRD-001", Title: "Proposal 1"},
		{ID: "PRD-002", Title: "Proposal 2"},
	}

	for _, p := range proposals {
		data, _ := json.Marshal(p)
		os.WriteFile(filepath.Join(proposalsDir, p.ID+".json"), data, 0644)
	}

	provider := NewContextProvider(tmpDir, &mockGitClient{})
	loaded, err := provider.GetExistingProposals(context.Background())
	if err != nil {
		t.Fatalf("GetExistingProposals failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("Loaded %d proposals, want 2", len(loaded))
	}
}

func TestCheckDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	proposalsDir := filepath.Join(tmpDir, "clawdlinux", "docs", "proposals")
	os.MkdirAll(proposalsDir, 0755)

	proposals := []*Proposal{
		{ID: "PRD-001", Title: "Implement Feature X"},
		{ID: "PRD-002", Title: "Feature Y Implementation"},
		{ID: "PRD-003", Title: "Fix Bug Z"},
	}

	for _, p := range proposals {
		data, _ := json.Marshal(p)
		os.WriteFile(filepath.Join(proposalsDir, p.ID+".json"), data, 0644)
	}

	provider := NewContextProvider(tmpDir, &mockGitClient{})

	tests := []struct {
		name      string
		title     string
		wantCount int
	}{
		{
			name:      "exact match",
			title:     "Implement Feature X",
			wantCount: 1,
		},
		{
			name:      "partial match",
			title:     "Feature",
			wantCount: 2,
		},
		{
			name:      "no match",
			title:     "Database Migration",
			wantCount: 0,
		},
		{
			name:      "case insensitive",
			title:     "FEATURE",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duplicates, err := provider.CheckDuplicates(context.Background(), tt.title)
			if err != nil {
				t.Fatalf("CheckDuplicates failed: %v", err)
			}
			if len(duplicates) != tt.wantCount {
				t.Errorf("Found %d duplicates, want %d", len(duplicates), tt.wantCount)
			}
		})
	}
}

// pkg/agents/pm/pm_agent_test.go
package pm

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockMissionControl struct {
	linkError error
}

func (m *mockMissionControl) GetPendingTasks(ctx context.Context) ([]Task, error) {
	return nil, nil
}

func (m *mockMissionControl) UpdateTaskStatus(ctx context.Context, taskID string, status string) error {
	return nil
}

func (m *mockMissionControl) LinkProposalToTask(ctx context.Context, taskID string, proposalID string) error {
	return m.linkError
}

type mockConsensusClient struct {
	submitError error
}

func (m *mockConsensusClient) SubmitProposal(ctx context.Context, proposal *Proposal) error {
	return m.submitError
}

func (m *mockConsensusClient) GetVotingResults(ctx context.Context) ([]VotingResult, error) {
	return nil, nil
}

func TestPMAgentGenerateProposal(t *testing.T) {
	tmpDir := t.TempDir()

	llmResponse := `{
		"title": "Test Proposal",
		"problem_statement": "This is a clear problem statement with enough detail to understand the requirements",
		"goals": ["Goal 1"],
		"success_metrics": [{"name": "M1", "target": "100%", "measurable": true}],
		"technical_approach": {
			"overview": "This is a detailed technical approach with sufficient information about the implementation strategy and architecture decisions",
			"components": ["C1"],
			"technologies": ["Go"],
			"integration_points": ["I1"],
			"alternatives_considered": []
		},
		"requirements": [{"id": "REQ-001", "type": "functional", "description": "Test", "priority": "P0", "testable": true, "test_method": "Unit test"}],
		"dependencies": [],
		"edge_cases": [{"scenario": "Test", "impact": "Low", "mitigation": "Handle it"}]
	}`

	config := PMAgentConfig{
		WorkspaceRoot:        tmpDir,
		LLMClient:            &mockLLMClient{response: llmResponse},
		Model:                "test-model",
		GitClient:            &mockGitClient{},
		ConsensusClient:      nil,
		MissionControlClient: nil,
	}

	agent := NewPMAgent(config)

	task := Task{
		ID:          "task-1",
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    PriorityP1,
		CreatedAt:   time.Now(),
	}

	proposal, err := agent.GenerateProposal(context.Background(), task)
	if err != nil {
		t.Fatalf("GenerateProposal failed: %v", err)
	}

	if proposal.ID == "" {
		t.Error("Proposal ID is empty")
	}
	if proposal.Status != StatusDraft {
		t.Errorf("Status = %s, want %s", proposal.Status, StatusDraft)
	}
	if len(proposal.Requirements) == 0 {
		t.Error("No requirements generated")
	}
}

func TestPMAgentValidateProposal(t *testing.T) {
	config := PMAgentConfig{
		WorkspaceRoot: t.TempDir(),
		LLMClient:     &mockLLMClient{},
		Model:         "test-model",
		GitClient:     &mockGitClient{},
	}

	agent := NewPMAgent(config)

	tests := []struct {
		name      string
		proposal  *Proposal
		wantValid bool
	}{
		{
			name: "valid proposal",
			proposal: &Proposal{
				ProblemStatement: "Clear problem statement with enough detail to pass validation requirements",
				Goals:            []string{"Goal 1"},
				SuccessMetrics:   []SuccessMetric{{Name: "M1", Target: "100%", Measurable: true}},
				Requirements:     []Requirement{{ID: "REQ-001", Testable: true, TestMethod: "Unit test"}},
				TechnicalApproach: TechnicalApproach{
					Overview: "Detailed technical approach with sufficient information to meet validation requirements for implementation strategy",
				},
				EdgeCases: []EdgeCase{{Scenario: "Test"}},
			},
			wantValid: true,
		},
		{
			name:      "invalid proposal",
			proposal:  &Proposal{},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.ValidateProposal(context.Background(), tt.proposal)
			if err != nil {
				t.Fatalf("ValidateProposal failed: %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v. Errors: %v", result.Valid, tt.wantValid, result.Errors)
			}
		})
	}
}

func TestFlagForHumanApproval(t *testing.T) {
	config := PMAgentConfig{
		WorkspaceRoot: t.TempDir(),
		LLMClient:     &mockLLMClient{},
		Model:         "test-model",
		GitClient:     &mockGitClient{},
	}

	agent := NewPMAgent(config).(*pmAgent)

	tests := []struct {
		name     string
		proposal *Proposal
		wantFlag bool
		wantReason string
	}{
		{
			name: "critical security",
			proposal: &Proposal{
				SecurityImpact: SecurityImpact{Level: SecCritical},
			},
			wantFlag:   true,
			wantReason: "High security impact",
		},
		{
			name: "high security",
			proposal: &Proposal{
				SecurityImpact: SecurityImpact{Level: SecHigh},
			},
			wantFlag:   true,
			wantReason: "High security impact",
		},
		{
			name: "critical risk",
			proposal: &Proposal{
				RiskAssessment: RiskAssessment{OverallRisk: RiskCritical},
			},
			wantFlag:   true,
			wantReason: "Critical risk level",
		},
		{
			name: "compliance impact",
			proposal: &Proposal{
				ComplianceImpact: ComplianceImpact{
					AffectedStandards: []string{"PCI-DSS", "HIPAA"},
				},
			},
			wantFlag:   true,
			wantReason: "Compliance impact",
		},
		{
			name: "no flags",
			proposal: &Proposal{
				SecurityImpact: SecurityImpact{Level: SecLow},
				RiskAssessment: RiskAssessment{OverallRisk: RiskLow},
			},
			wantFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent.flagForHumanApproval(tt.proposal)

			if tt.proposal.RequiresHumanApproval != tt.wantFlag {
				t.Errorf("RequiresHumanApproval = %v, want %v", tt.proposal.RequiresHumanApproval, tt.wantFlag)
			}
			if tt.wantFlag && tt.proposal.ApprovalReason != tt.wantReason {
				t.Errorf("ApprovalReason = %s, want %s", tt.proposal.ApprovalReason, tt.wantReason)
			}
		})
	}
}

func TestSubmitForVoting(t *testing.T) {
	tests := []struct {
		name          string
		consensus     ConsensusClient
		wantError     bool
	}{
		{
			name:      "success",
			consensus: &mockConsensusClient{},
			wantError: false,
		},
		{
			name:      "no consensus client",
			consensus: nil,
			wantError: true,
		},
		{
			name:      "submit error",
			consensus: &mockConsensusClient{submitError: errors.New("network error")},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := PMAgentConfig{
				WorkspaceRoot:   t.TempDir(),
				LLMClient:       &mockLLMClient{},
				Model:           "test-model",
				GitClient:       &mockGitClient{},
				ConsensusClient: tt.consensus,
			}

			agent := NewPMAgent(config)
			proposal := &Proposal{ID: "PRD-001"}

			err := agent.SubmitForVoting(context.Background(), proposal)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if proposal.Status != StatusVoting {
				t.Errorf("Status = %s, want %s", proposal.Status, StatusVoting)
			}
		})
	}
}
```