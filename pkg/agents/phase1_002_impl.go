// Auto-generated for phase1-002
// Task: Implement PM agent for proposal generation and requirements analysis

```go
// pkg/agents/pm/types.go
package pm

import "time"

type ProposalStatus string

const (
	StatusDraft      ProposalStatus = "draft"
	StatusValidating ProposalStatus = "validating"
	StatusVoting     ProposalStatus = "voting"
	StatusApproved   ProposalStatus = "approved"
	StatusRejected   ProposalStatus = "rejected"
	StatusExecuting  ProposalStatus = "executing"
	StatusCompleted  ProposalStatus = "completed"
)

type RequirementType string

const (
	ReqFunctional    RequirementType = "functional"
	ReqNonFunctional RequirementType = "non-functional"
	ReqSecurity      RequirementType = "security"
	ReqCompliance    RequirementType = "compliance"
	ReqPerformance   RequirementType = "performance"
)

type Priority string

const (
	PriorityP0 Priority = "P0"
	PriorityP1 Priority = "P1"
	PriorityP2 Priority = "P2"
)

type DependencyType string

const (
	DepInternal DependencyType = "internal"
	DepExternal DependencyType = "external"
	DepData     DependencyType = "data"
	DepAPI      DependencyType = "api"
)

type ImpactLevel string

const (
	ImpactHigh   ImpactLevel = "high"
	ImpactMedium ImpactLevel = "medium"
	ImpactLow    ImpactLevel = "low"
	ImpactNone   ImpactLevel = "none"
)

type TShirtSize string

const (
	SizeS  TShirtSize = "S"
	SizeM  TShirtSize = "M"
	SizeL  TShirtSize = "L"
	SizeXL TShirtSize = "XL"
)

type RiskLevel string

const (
	RiskCritical RiskLevel = "critical"
	RiskHigh     RiskLevel = "high"
	RiskMedium   RiskLevel = "medium"
	RiskLow      RiskLevel = "low"
)

type SecurityLevel string

const (
	SecCritical SecurityLevel = "critical"
	SecHigh     SecurityLevel = "high"
	SecMedium   SecurityLevel = "medium"
	SecLow      SecurityLevel = "low"
)

type Proposal struct {
	ID               string                 `json:"id"`
	Title            string                 `json:"title"`
	CreatedAt        time.Time              `json:"created_at"`
	Author           string                 `json:"author"`
	Status           ProposalStatus         `json:"status"`
	ProblemStatement string                 `json:"problem_statement"`
	Goals            []string               `json:"goals"`
	SuccessMetrics   []SuccessMetric        `json:"success_metrics"`
	TechnicalApproach TechnicalApproach     `json:"technical_approach"`
	Requirements      []Requirement         `json:"requirements"`
	Dependencies      []Dependency          `json:"dependencies"`
	EdgeCases         []EdgeCase            `json:"edge_cases"`
	BusinessImpact    BusinessImpact        `json:"business_impact"`
	EffortEstimate    EffortEstimate        `json:"effort_estimate"`
	RoadmapAlignment  RoadmapAlignment      `json:"roadmap_alignment"`
	RiskAssessment    RiskAssessment        `json:"risk_assessment"`
	SecurityImpact    SecurityImpact        `json:"security_impact"`
	ComplianceImpact  ComplianceImpact      `json:"compliance_impact"`
	RequiresHumanApproval bool              `json:"requires_human_approval"`
	ApprovalReason        string            `json:"approval_reason,omitempty"`
	ConsensusScore    float64               `json:"consensus_score,omitempty"`
	AgentVotes        map[string]Vote       `json:"agent_votes,omitempty"`
}

type SuccessMetric struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	Measurable bool   `json:"measurable"`
	Baseline   string `json:"baseline,omitempty"`
}

type TechnicalApproach struct {
	Overview               string        `json:"overview"`
	Components             []string      `json:"components"`
	Technologies           []string      `json:"technologies"`
	IntegrationPoints      []string      `json:"integration_points"`
	AlternativesConsidered []Alternative `json:"alternatives_considered"`
}

type Alternative struct {
	Approach  string `json:"approach"`
	Pros      string `json:"pros"`
	Cons      string `json:"cons"`
	NotChosen string `json:"not_chosen_reason"`
}

type Requirement struct {
	ID          string          `json:"id"`
	Type        RequirementType `json:"type"`
	Description string          `json:"description"`
	Priority    Priority        `json:"priority"`
	Testable    bool            `json:"testable"`
	TestMethod  string          `json:"test_method,omitempty"`
	DependsOn   []string        `json:"depends_on,omitempty"`
}

type Dependency struct {
	Type        DependencyType `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Blocking    bool           `json:"blocking"`
}

type EdgeCase struct {
	Scenario   string `json:"scenario"`
	Impact     string `json:"impact"`
	Mitigation string `json:"mitigation"`
}

type BusinessImpact struct {
	RevenueImpact     ImpactLevel `json:"revenue_impact"`
	ComplianceImpact  ImpactLevel `json:"compliance_impact"`
	PerformanceImpact ImpactLevel `json:"performance_impact"`
	CustomerImpact    ImpactLevel `json:"customer_impact"`
	Rationale         string      `json:"rationale"`
}

type EffortEstimate struct {
	Size          TShirtSize `json:"size"`
	StoryPoints   int        `json:"story_points"`
	Engineering   string     `json:"engineering"`
	Testing       string     `json:"testing"`
	Documentation string     `json:"documentation"`
	Total         string     `json:"total"`
}

type RoadmapAlignment struct {
	Phase     string   `json:"phase"`
	Milestone string   `json:"milestone"`
	OKRs      []string `json:"okrs"`
	BlockedBy []string `json:"blocked_by,omitempty"`
}

type RiskAssessment struct {
	OverallRisk RiskLevel `json:"overall_risk"`
	Risks       []Risk    `json:"risks"`
	Mitigation  string    `json:"mitigation_strategy"`
}

type Risk struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Probability string `json:"probability"`
	Impact      string `json:"impact"`
	Mitigation  string `json:"mitigation"`
}

type SecurityImpact struct {
	Level          SecurityLevel `json:"level"`
	Considerations []string      `json:"considerations"`
	RequiresAudit  bool          `json:"requires_audit"`
}

type ComplianceImpact struct {
	AffectedStandards []string `json:"affected_standards"`
	RequiresReview    bool     `json:"requires_review"`
	Notes             string   `json:"notes"`
}

type Vote struct {
	Agent     string    `json:"agent"`
	Score     int       `json:"score"`
	Decision  string    `json:"decision"`
	Feedback  string    `json:"feedback"`
	Timestamp time.Time `json:"timestamp"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkspaceContext struct {
	Architecture  string
	Roadmap       string
	Memory        string
	RecentCommits []Commit
	ExistingPRDs  []*Proposal
}

type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
}

type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
	Score    int
}

// pkg/agents/pm/interfaces.go
package pm

import "context"

type PMAgent interface {
	GenerateProposal(ctx context.Context, task Task) (*Proposal, error)
	RefineProposal(ctx context.Context, proposal *Proposal, feedback []Vote) (*Proposal, error)
	ValidateProposal(ctx context.Context, proposal *Proposal) (ValidationResult, error)
	SubmitForVoting(ctx context.Context, proposal *Proposal) error
}

type ContextProvider interface {
	LoadArchitecture(ctx context.Context) (string, error)
	LoadRoadmap(ctx context.Context) (string, error)
	LoadMemory(ctx context.Context) (string, error)
	GetRecentCommits(ctx context.Context, limit int) ([]Commit, error)
	GetExistingProposals(ctx context.Context) ([]*Proposal, error)
	CheckDuplicates(ctx context.Context, title string) ([]string, error)
}

type ProposalGenerator interface {
	Generate(ctx context.Context, task Task, context WorkspaceContext) (*Proposal, error)
}

type RequirementsAnalyzer interface {
	AnalyzeRequirements(ctx context.Context, proposal *Proposal) ([]Requirement, error)
	BuildDependencyGraph(ctx context.Context, requirements []Requirement) ([]Dependency, error)
	EnumerateEdgeCases(ctx context.Context, proposal *Proposal) ([]EdgeCase, error)
}

type BusinessAlignmentEngine interface {
	CalculateImpact(ctx context.Context, proposal *Proposal) (BusinessImpact, error)
	EstimateEffort(ctx context.Context, proposal *Proposal) (EffortEstimate, error)
	MapToRoadmap(ctx context.Context, proposal *Proposal, roadmap string) (RoadmapAlignment, error)
}

type ProposalValidator interface {
	ValidateCompleteness(proposal *Proposal) ValidationResult
	ValidateClarity(proposal *Proposal) ValidationResult
	ValidateFeasibility(proposal *Proposal) ValidationResult
}

type GitIntegration interface {
	CommitProposal(ctx context.Context, proposal *Proposal) error
	UpdateProposal(ctx context.Context, proposal *Proposal) error
}

type MissionControlClient interface {
	GetPendingTasks(ctx context.Context) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status string) error
	LinkProposalToTask(ctx context.Context, taskID string, proposalID string) error
}

type ConsensusClient interface {
	SubmitProposal(ctx context.Context, proposal *Proposal) error
	GetVotingResults(ctx context.Context) ([]VotingResult, error)
}

type VotingResult struct {
	ProposalID string
	Approved   bool
	Score      float64
	Feedback   []Vote
}

type LLMClient interface {
	Complete(ctx context.Context, prompt string, model string) (string, error)
}

type GitClient interface {
	Add(path string) error
	Commit(message string) error
	Push() error
}

// pkg/agents/pm/context_provider.go
package pm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type contextProvider struct {
	workspaceRoot string
	gitClient     GitClient
}

func NewContextProvider(workspaceRoot string, gitClient GitClient) ContextProvider {
	return &contextProvider{
		workspaceRoot: workspaceRoot,
		gitClient:     gitClient,
	}
}

func (cp *contextProvider) LoadArchitecture(ctx context.Context) (string, error) {
	path := filepath.Join(cp.workspaceRoot, "clawdlinux", "ARCHITECTURE.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load architecture: %w", err)
	}
	return string(data), nil
}

func (cp *contextProvider) LoadRoadmap(ctx context.Context) (string, error) {
	path := filepath.Join(cp.workspaceRoot, "clawdlinux", "ROADMAP.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load roadmap: %w", err)
	}
	return string(data), nil
}

func (cp *contextProvider) LoadMemory(ctx context.Context) (string, error) {
	path := filepath.Join(cp.workspaceRoot, "MEMORY.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load memory: %w", err)
	}
	return string(data), nil
}

func (cp *contextProvider) GetRecentCommits(ctx context.Context, limit int) ([]Commit, error) {
	// TODO: Implement git log parsing
	return []Commit{}, nil
}

func (cp *contextProvider) GetExistingProposals(ctx context.Context) ([]*Proposal, error) {
	proposalsDir := filepath.Join(cp.workspaceRoot, "clawdlinux", "docs", "proposals")
	
	if _, err := os.Stat(proposalsDir); os.IsNotExist(err) {
		return []*Proposal{}, nil
	}
	
	files, err := filepath.Glob(filepath.Join(proposalsDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob proposals: %w", err)
	}
	
	proposals := make([]*Proposal, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		var proposal Proposal
		if err := json.Unmarshal(data, &proposal); err != nil {
			continue
		}
		
		proposals = append(proposals, &proposal)
	}
	
	return proposals, nil
}

func (cp *contextProvider) CheckDuplicates(ctx context.Context, title string) ([]string, error) {
	proposals, err := cp.GetExistingProposals(ctx)
	if err != nil {
		return nil, err
	}
	
	duplicates := []string{}
	titleLower := strings.ToLower(title)
	
	for _, p := range proposals {
		if strings.Contains(strings.ToLower(p.Title), titleLower) {
			duplicates = append(duplicates, p.ID)
		}
	}
	
	return duplicates, nil
}

// pkg/agents/pm/generator.go
package pm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type llmProposalGenerator struct {
	llmClient LLMClient
	model     string
}

func NewProposalGenerator(llmClient LLMClient, model string) ProposalGenerator {
	return &llmProposalGenerator{
		llmClient: llmClient,
		model:     model,
	}
}

func (g *llmProposalGenerator) Generate(ctx context.Context, task Task, wsContext WorkspaceContext) (*Proposal, error) {
	prompt := g.buildPrompt(task, wsContext)
	
	response, err := g.llmClient.Complete(ctx, prompt, g.model)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}
	
	proposal, err := g.parseProposal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proposal: %w", err)
	}
	
	proposal.ID = generateProposalID()
	proposal.CreatedAt = time.Now()
	proposal.Author = "pm-agent"
	proposal.Status = StatusDraft
	
	return proposal, nil
}

func (g *llmProposalGenerator) buildPrompt(task Task, ctx WorkspaceContext) string {
	commitsStr := g.formatCommits(ctx.RecentCommits)
	
	return fmt.Sprintf(`You are a Product Manager AI agent. Generate a structured proposal for:

TASK: %s
DESCRIPTION: %s
PRIORITY: %s

CONTEXT:
--- ARCHITECTURE ---
%s

--- ROADMAP ---
%s

--- MEMORY ---
%s

--- RECENT COMMITS ---
%s

REQUIREMENTS:
1. Problem statement (clear, specific)
2. Success metrics (measurable, no vague language)
3. Technical approach (components, integration points)
4. Requirements (functional, testable)
5. Dependencies (internal/external)
6. Edge cases (error states, boundaries)
7. Business impact (revenue/compliance/performance)
8. Effort estimate (T-shirt size: S/M/L/XL)
9. Risk assessment (with mitigations)
10. Security/compliance impact

OUTPUT FORMAT: JSON matching Proposal schema
BANNED WORDS: "should", "might", "could", "maybe", "probably"
REQUIRED: All fields must be filled with concrete values

Generate proposal:`, 
		task.Title, 
		task.Description,
		task.Priority,
		ctx.Architecture, 
		ctx.Roadmap,
		ctx.Memory,
		commitsStr)
}

func (g *llmProposalGenerator) formatCommits(commits []Commit) string {
	if len(commits) == 0 {
		return "No recent commits"
	}
	
	var sb strings.Builder
	for _, c := range commits {
		sb.WriteString(fmt.Sprintf("- %s: %s (%s)\n", c.Hash[:7], c.Message, c.Author))
	}
	return sb.String()
}

func (g *llmProposalGenerator) parseProposal(response string) (*Proposal, error) {
	// Extract JSON from response (handle markdown code blocks)
	jsonStr := response
	if idx := strings.Index(response, "```json"); idx != -1 {
		jsonStr = response[idx+7:]
		if idx := strings.Index(jsonStr, "```"); idx != -1 {
			jsonStr = jsonStr[:idx]
		}
	} else if idx := strings.Index(response, "```"); idx != -1 {
		jsonStr = response[idx+3:]
		if idx := strings.Index(jsonStr, "```"); idx != -1 {
			jsonStr = jsonStr[:idx]
		}
	}
	
	jsonStr = strings.TrimSpace(jsonStr)
	
	var proposal Proposal
	if err := json.Unmarshal([]byte(jsonStr), &proposal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proposal JSON: %w", err)
	}
	
	return &proposal, nil
}

func generateProposalID() string {
	return fmt.Sprintf("PRD-%d", time.Now().Unix())
}

// pkg/agents/pm/validator.go
package pm

import (
	"fmt"
	"strings"
)

type proposalValidator struct {
	bannedWords []string
}

func NewProposalValidator(bannedWords []string) ProposalValidator {
	if bannedWords == nil {
		bannedWords = []string{"should", "might", "could", "maybe", "probably"}
	}
	return &proposalValidator{
		bannedWords: bannedWords,
	}
}

func (v *proposalValidator) ValidateCompleteness(p *Proposal) ValidationResult {
	result := ValidationResult{Valid: true, Score: 100}
	
	if p.ProblemStatement == "" {
		result.Errors = append(result.Errors, "missing problem statement")
		result.Score -= 20
	}
	
	if len(p.Goals) == 0 {
		result.Warnings = append(result.Warnings, "no goals defined")
		result.Score -= 5
	}
	
	if len(p.SuccessMetrics) == 0 {
		result.Errors = append(result.Errors, "missing success metrics")
		result.Score -= 20
	}
	
	if len(p.Requirements) == 0 {
		result.Errors = append(result.Errors, "missing requirements")
		result.Score -= 20
	}
	
	if p.TechnicalApproach.Overview == "" {
		result.Errors = append(result.Errors, "missing technical approach")
		result.Score -= 15
	}
	
	if len(p.EdgeCases) == 0 {
		result.Warnings = append(result.Warnings, "no edge cases identified")
		result.Score -= 5
	}
	
	for i, req := range p.Requirements {
		if !req.Testable {
			result.Warnings = append(result.Warnings, 
				fmt.Sprintf("requirement %d (%s) not testable", i, req.ID))
			result.Score -= 5
		}
		if req.TestMethod == "" && req.Testable {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("requirement %d (%s) missing test method", i, req.ID))
			result.Score -= 3
		}
	}
	
	result.Valid = len(result.Errors) == 0
	if result.Score < 0 {
		result.Score = 0
	}
	
	return result
}

func (v *proposalValidator) ValidateClarity(p *Proposal) ValidationResult {
	result := ValidationResult{Valid: true, Score: 100}
	
	text := strings.ToLower(p.ProblemStatement + " " + p.TechnicalApproach.Overview)
	
	for _, word := range v.bannedWords {
		if strings.Contains(text, word) {
			result.Errors = append(result.Errors, 
				fmt.Sprintf("vague language detected: '%s'", word))
			result.Score -= 15
		}
	}
	
	for _, metric := range p.SuccessMetrics {
		if !metric.Measurable {
			result.Errors = append(result.Errors, 
				fmt.Sprintf("metric '%s' not measurable", metric.Name))
			result.Score -= 10
		}
		if metric.Target == "" {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("metric '%s' missing target", metric.Name))
			result.Score -= 5
		}
	}
	
	if len(p.ProblemStatement) < 50 {
		result.Warnings = append(result.Warnings, "problem statement too brief")
		result.Score -= 10
	}
	
	if p.TechnicalApproach.Overview != "" && len(p.TechnicalApproach.Overview) < 100 {
		result.Warnings = append(result.Warnings, "technical approach lacks detail")
		result.Score -= 10
	}
	
	result.Valid = len(result.Errors) == 0
	if result.Score < 0 {
		result.Score = 0
	}
	
	return result
}

func (v *proposalValidator) ValidateFeasibility(p *Proposal) ValidationResult {
	result := ValidationResult{Valid: true, Score: 100}
	
	if p.EffortEstimate.Size == "" {
		result.Warnings = append(result.Warnings, "missing effort estimate")
		result.Score -= 10
	}
	
	if len(p.Requirements) > 20 && p.EffortEstimate.Size == SizeS {
		result.Warnings = append(result.Warnings, 
			"effort estimate seems too small for number of requirements")
		result.Score -= 15
	}
	
	if p.RiskAssessment.OverallRisk == RiskCritical && len(p.RiskAssessment.Risks) == 0 {
		result.Errors = append(result.Errors, "critical risk level but no risks identified")
		result.Score -= 20
	}
	
	if len(p.Dependencies) > 0 {
		for _, dep := range p.Dependencies {
			if dep.Blocking && dep.Description == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("blocking dependency '%s' lacks description", dep.Name))
				result.Score -= 5
			}
		}
	}
	
	result.Valid = len(result.Errors) == 0
	if result.Score < 0 {
		result.Score = 0
	}
	
	return result
}

// pkg/agents/pm/requirements_analyzer.go
package pm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type requirementsAnalyzer struct {
	llmClient LLMClient
	model     string
}

func NewRequirementsAnalyzer(llmClient LLMClient, model string) RequirementsAnalyzer {
	return &requirementsAnalyzer{
		llmClient: llmClient,
		model:     model,
	}
}

func (ra *requirementsAnalyzer) AnalyzeRequirements(ctx context.Context, proposal *Proposal) ([]Requirement, error) {
	if len(proposal.Requirements) > 0 {
		return ra.enhanceRequirements(ctx, proposal.Requirements)
	}
	
	prompt := fmt.Sprintf(`Extract concrete, testable requirements from this proposal:

PROBLEM: %s
APPROACH: %s

For each requirement:
1. Must be specific and measurable
2. Must include test method
3. Identify dependencies
4. Classify type (functional/security/performance/compliance)
5. Assign priority (P0/P1/P2)

Output: JSON array of Requirements with schema:
{
  "id": "REQ-001",
  "type": "functional",
  "description": "...",
  "priority": "P0",
  "testable": true,
  "test_method": "...",
  "depends_on": []
}`, 
		proposal.ProblemStatement,
		proposal.TechnicalApproach.Overview)
	
	response, err := ra.llmClient.Complete(ctx, prompt, ra.model)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}
	
	requirements, err := ra.parseRequirements(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requirements: %w", err)
	}
	
	return requirements, nil
}

func (ra *requirementsAnalyzer) enhanceRequirements(ctx context.Context, reqs []Requirement) ([]Requirement, error) {
	for i := range reqs {
		if reqs[i].ID == "" {
			reqs[i].ID = fmt.Sprintf("REQ-%03d", i+1)
		}
		if !reqs[i].Testable {
			reqs[i].Testable = true
			if reqs[i].TestMethod == "" {
				reqs[i].TestMethod = "manual verification"
			}
		}
	}
	return reqs, nil
}

func (ra *requirementsAnalyzer) parseRequirements(response string) ([]Requirement, error) {
	jsonStr := response
	if idx := strings.Index(response, "```json"); idx != -1 {
		jsonStr = response[idx+7:]
		if idx := strings.Index(jsonStr, "```"); idx != -1 {
			jsonStr = jsonStr[:idx]
		}
	} else if idx := strings.Index(response, "["); idx != -1 {
		jsonStr = response[idx:]
		if idx := strings.LastIndex(jsonStr, "]"); idx != -1 {
			jsonStr = jsonStr[:idx+1]
		}
	}
	
	jsonStr = strings.TrimSpace(jsonStr)
	
	var requirements []Requirement
	if err := json.Unmarshal([]byte(jsonStr), &requirements); err != nil {
		return nil, fmt.Errorf("failed to unmarshal requirements: %w", err)
	}
	
	return requirements, nil
}

func (ra *requirementsAnalyzer) BuildDependencyGraph(ctx context.Context, requirements []Requirement) ([]Dependency, error) {
	deps := []Dependency{}
	
	for _, req := range requirements {
		if len(req.DependsOn) > 0 {
			for _, depID := range req.DependsOn {
				deps = append(deps, Dependency{
					Type:        DepInternal,
					Name:        depID,
					Description: fmt.Sprintf("Required by %s", req.ID),
					Blocking:    req.Priority == PriorityP0,
				})
			}
		}
	}
	
	return deps, nil
}

func (ra *requirementsAnalyzer) EnumerateEdgeCases(ctx context.Context, proposal *Proposal) ([]EdgeCase, error) {
	reqsStr := ra.formatRequirements(proposal.Requirements)
	
	prompt := fmt.Sprintf(`Enumerate edge cases for:

PROPOSAL: %s
REQUIREMENTS:
%s

Consider:
- Error states (network failures, timeouts)
- Boundary conditions (empty input, max limits)
- Concurrent operations
- Failure scenarios
- Race conditions

Output: JSON array of EdgeCases with schema:
{
  "scenario": "...",
  "impact": "...",
  "mitigation": "..."
}`,
		proposal.Title,
		reqsStr)
	
	response, err := ra.llmClient.Complete(ctx, prompt, ra.model)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}
	
	edgeCases, err := ra.parseEdgeCases(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse edge cases: %w", err)
	}
	
	return edgeCases, nil
}

func (ra *requirementsAnalyzer) formatRequirements(reqs []Requirement) string {
	var sb strings.Builder
	for _, req := range reqs {
		sb.WriteString(fmt.Sprintf("- [%s] %s (Priority: %s)\n", req.ID, req.Description, req.Priority))
	}
	return sb.String()
}

func (ra *requirementsAnalyzer) parseEdgeCases(response string) ([]EdgeCase, error) {
	jsonStr := response
	if idx := strings.Index(response, "```json"); idx != -1 {
		jsonStr = response[idx+7:]
		if idx := strings.Index(jsonStr, "```"); idx != -1 {
			jsonStr = jsonStr[:idx]
		}
	} else if idx := strings.Index(response, "["); idx != -1 {
		jsonStr = response[idx:]
		if idx := strings.LastIndex(jsonStr, "]"); idx != -1 {
			jsonStr = jsonStr[:idx+1]
		}
	}
	
	jsonStr = strings.TrimSpace(jsonStr)
	
	var edgeCases []EdgeCase
	if err := json.Unmarshal([]byte(jsonStr), &edgeCases); err != nil {
		return nil, fmt.Errorf("failed to unmarshal edge cases: %w", err)
	}
	
	return edgeCases, nil
}

// pkg/agents/pm/business_alignment.go
package pm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type businessAlignmentEngine struct {
	llmClient LLMClient
	model     string
}

func NewBusinessAlignmentEngine(llmClient LLMClient, model string) BusinessAlignmentEngine {
	return &businessAlignmentEngine{
		llmClient: llmClient,
		model:     model,
	}
}

func (bae *businessAlignmentEngine) CalculateImpact(ctx context.Context, proposal *Proposal) (BusinessImpact, error) {
	prompt := fmt.Sprintf(`Assess business impact for:

PROPOSAL: %s
PROBLEM: %s
APPROACH: %s

Rate impact (high/medium/low/none) for:
1. Revenue (does this enable monetization?)
2. Compliance (PCI-DSS, HIPAA, SOX, GDPR)
3. Performance (speed, scalability)
4. Customer (UX, reliability)

Provide rationale.

Output: JSON BusinessImpact with schema:
{
  "revenue_impact": "medium",
  "compliance_impact": "high",
  "performance_impact": "low",
  "customer_impact": "medium",
  "rationale": "..."
}`,
		proposal.Title,
		proposal.ProblemStatement,
		proposal.TechnicalApproach.Overview)
	
	response, err := bae.llmClient.Complete(ctx, prompt, bae.model)
	if err != nil {
		return BusinessImpact{}, fmt.Errorf("LLM call failed: %w", err)
	}
	
	impact, err := bae.parseBusinessImpact(response)
	if err != nil {
		return BusinessImpact{}, fmt.Errorf("failed to parse business impact: %w", err)
	}
	
	return impact, nil
}

func (bae *businessAlignmentEngine) parseBusinessImpact(response string) (BusinessImpact, error) {
	jsonStr := extractJSON(response)
	
	var impact BusinessImpact
	if err := json.Unmarshal([]byte(jsonStr), &impact); err != nil {
		return BusinessImpact{}, fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	return impact, nil
}

func (bae *businessAlignmentEngine) EstimateEffort(ctx context.Context, proposal *Proposal) (EffortEstimate, error) {
	numReqs := len(proposal.Requirements)
	numDeps := len(proposal.Dependencies)
	
	var size TShirtSize
	var storyPoints int
	var engineering, testing, documentation, total string
	
	if numReqs <= 3 && numDeps == 0 {
		size = SizeS
		storyPoints = 1
		engineering = "0.5-1 day"
		testing = "0.5 day"
		documentation = "0.5 day"
		total = "1-2 days"
	} else if numReqs <= 8 && numDeps <= 2 {
		size = SizeM
		storyPoints = 3
		engineering = "2-3 days"
		testing = "1 day"
		documentation = "0.5 day"
		total = "3.5-4.5 days"
	} else if numReqs <= 15 && numDeps <= 5 {
		size = SizeL
		storyPoints = 8
		engineering = "1-2 weeks"
		testing = "3-5 days"
		documentation = "1-2 days"
		total = "2-3 weeks"
	} else {
		size = SizeXL
		storyPoints = 13
		engineering = "2-4 weeks"
		testing = "1-2 weeks"
		documentation = "3-5 days"
		total = "4-6 weeks"
	}
	
	return EffortEstimate{
		Size:          size,
		StoryPoints:   storyPoints,
		Engineering:   engineering,
		Testing:       testing,
		Documentation: documentation,
		Total:         total,
	}, nil
}

func (bae *businessAlignmentEngine) MapToRoadmap(ctx context.Context, proposal *Proposal, roadmap string) (RoadmapAlignment, error) {
	prompt := fmt.Sprintf(`Map this proposal to the roadmap:

PROPOSAL: %s
PROBLEM: %s

ROADMAP:
%s

Determine:
1. Which phase (Phase 1/2/3/4)
2. Which milestone (MVP, Scale, Enterprise, SaaS)
3. Relevant OKRs
4. Any blockers

Output: JSON RoadmapAlignment with schema:
{
  "phase": "Phase 1",
  "milestone": "MVP",
  "okrs": ["..."],
  "blocked_by": []
}`,
		proposal.Title,
		proposal.ProblemStatement,
		roadmap)
	
	response, err := bae.llmClient.Complete(ctx, prompt, bae.model)
	if err != nil {
		return RoadmapAlignment{}, fmt.Errorf("LLM call failed: %w", err)
	}
	
	alignment, err := bae.parseRoadmapAlignment(response)
	if err != nil {
		return RoadmapAlignment{}, fmt.Errorf("failed to parse roadmap alignment: %w", err)
	}
	
	return alignment, nil
}

func (bae *businessAlignmentEngine) parseRoadmapAlignment(response string) (RoadmapAlignment, error) {
	jsonStr := extractJSON(response)
	
	var alignment RoadmapAlignment
	if err := json.Unmarshal([]byte(jsonStr), &alignment); err != nil {
		return RoadmapAlignment{}, fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	return alignment, nil
}

func extractJSON(response string) string {
	jsonStr := response
	if idx := strings.Index(response, "```json"); idx != -1 {
		jsonStr = response[idx+7:]
		if idx := strings.Index(jsonStr, "```"); idx != -1 {
			jsonStr = jsonStr[:idx]
		}
	} else if idx := strings.Index(response, "{"); idx != -1 {
		jsonStr = response[idx:]
		if idx := strings.LastIndex(jsonStr, "}"); idx != -1 {
			jsonStr = jsonStr[:idx+1]
		}
	}
	
	return strings.TrimSpace(jsonStr)
}

// pkg/agents/pm/git_integration.go
package pm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type gitIntegration struct {
	repoPath  string
	gitClient GitClient
}

func NewGitIntegration(repoPath string, gitClient GitClient) GitIntegration {
	return &gitIntegration{
		repoPath:  repoPath,
		gitClient: gitClient,
	}
}

func (gi *gitIntegration) CommitProposal(ctx context.Context, proposal *Proposal) error {
	data, err := json.MarshalIndent(proposal, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal proposal: %w", err)
	}
	
	proposalPath := filepath.Join(gi.repoPath, "docs", "proposals", fmt.Sprintf("%s.json", proposal.ID))
	
	if err := os.MkdirAll(filepath.Dir(proposalPath), 0755); err != nil {
		return fmt.Errorf("failed to create proposals dir: %w", err)
	}
	
	if err := os.WriteFile(proposalPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write proposal: %w", err)
	}
	
	if err := gi.gitClient.Add(proposalPath); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}
	
	commitMsg := fmt.Sprintf("[PM Agent] Add proposal: %s\n\n%s\nStatus: %s\nEffort: %s",
		proposal.Title,
		proposal.ProblemStatement,
		proposal.Status,
		proposal.EffortEstimate.Size)
	
	if err := gi.gitClient.Commit(commitMsg); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}
	
	return nil
}

func (gi *gitIntegration) UpdateProposal(ctx context.Context, proposal *Proposal) error {
	return gi.CommitProposal(ctx, proposal)
}

// pkg/agents/pm/pm_agent.go
package pm

import (
	"context"
	"fmt"
)

type pmAgent struct {
	contextProvider      ContextProvider
	generator            ProposalGenerator
	requirementsAnalyzer RequirementsAnalyzer
	businessEngine       BusinessAlignmentEngine
	validator            ProposalValidator
	gitIntegration       GitIntegration
	consensusClient      ConsensusClient
	missionControl       MissionControlClient
}

type PMAgentConfig struct {
	WorkspaceRoot        string
	LLMClient            LLMClient
	Model                string
	GitClient            GitClient
	ConsensusClient      ConsensusClient
	MissionControlClient MissionControlClient
}

func NewPMAgent(cfg PMAgentConfig) PMAgent {
	contextProvider := NewContextProvider(cfg.WorkspaceRoot, cfg.GitClient)
	generator := NewProposalGenerator(cfg.LLMClient, cfg.Model)
	requirementsAnalyzer := NewRequirementsAnalyzer(cfg.LLMClient, cfg.Model)
	businessEngine := NewBusinessAlignmentEngine(cfg.LLMClient, cfg.Model)
	validator := NewProposalValidator(nil)
	gitIntegration := NewGitIntegration(cfg.WorkspaceRoot, cfg.GitClient)
	
	return &pmAgent{
		contextProvider:      contextProvider,
		generator:            generator,
		requirementsAnalyzer: requirementsAnalyzer,
		businessEngine:       businessEngine,
		validator:            validator,
		gitIntegration:       gitIntegration,
		consensusClient:      cfg.ConsensusClient,
		missionControl:       cfg.MissionControlClient,
	}
}

func (pm *pmAgent) GenerateProposal(ctx context.Context, task Task) (*Proposal, error) {
	wsContext, err := pm.loadContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load context: %w", err)
	}
	
	duplicates, err := pm.contextProvider.CheckDuplicates(ctx, task.Title)
	if err != nil {
		return nil, fmt.Errorf("duplicate check failed: %w", err)
	}
	if len(duplicates) > 0 {
		return nil, fmt.Errorf("duplicate proposals found: %v", duplicates)
	}
	
	proposal, err := pm.generator.Generate(ctx, task, wsContext)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}
	
	reqs, err := pm.requirementsAnalyzer.AnalyzeRequirements(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("requirements analysis failed: %w", err)
	}
	proposal.Requirements = reqs
	
	deps, err := pm.requirementsAnalyzer.BuildDependencyGraph(ctx, reqs)
	if err != nil {
		return nil, fmt.Errorf("dependency analysis failed: %w", err)
	}
	proposal.Dependencies = deps
	
	edgeCases, err := pm.requirementsAnalyzer.EnumerateEdgeCases(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("edge case analysis failed: %w", err)
	}
	proposal.EdgeCases = edgeCases
	
	impact, err := pm.businessEngine.CalculateImpact(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("impact calculation failed: %w", err)
	}
	proposal.BusinessImpact = impact
	
	effort, err := pm.businessEngine.EstimateEffort(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("effort estimation failed: %w", err)
	}
	proposal.EffortEstimate = effort
	
	alignment, err := pm.businessEngine.MapToRoadmap(ctx, proposal, wsContext.Roadmap)
	if err != nil {
		return nil, fmt.Errorf("roadmap mapping failed: %w", err)
	}
	proposal.RoadmapAlignment = alignment
	
	result, err := pm.ValidateProposal(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	if !result.Valid {
		return nil, fmt.Errorf("proposal validation failed: %v", result.Errors)
	}
	
	pm.flagForHumanApproval(proposal)
	
	if err := pm.gitIntegration.CommitProposal(ctx, proposal); err != nil {
		return nil, fmt.Errorf("git commit failed: %w", err)
	}
	
	if pm.missionControl != nil {
		if err := pm.missionControl.LinkProposalToTask(ctx, task.ID, proposal.ID); err != nil {
			return nil, fmt.Errorf("mission control update failed: %w", err)
		}
	}
	
	return proposal, nil
}

func (pm *pmAgent) RefineProposal(ctx context.Context, proposal *Proposal, feedback []Vote) (*Proposal, error) {
	// TODO: Implement refinement logic
	return proposal, nil
}

func (pm *pmAgent) ValidateProposal(ctx context.Context, proposal *Proposal) (ValidationResult, error) {
	completeness := pm.validator.ValidateCompleteness(proposal)
	clarity := pm.validator.ValidateClarity(proposal)
	feasibility := pm.validator.ValidateFeasibility(proposal)
	
	result := ValidationResult{
		Valid: completeness.Valid && clarity.Valid && feasibility.Valid,
		Score: (completeness.Score + clarity.Score + feasibility.Score) / 3,
	}
	
	result.Errors = append(result.Errors, completeness.Errors...)
	result.Errors = append(result.Errors, clarity.Errors...)
	result.Errors = append(result.Errors, feasibility.Errors...)
	
	result.Warnings = append(result.Warnings, completeness.Warnings...)
	result.Warnings = append(result.Warnings, clarity.Warnings...)
	result.Warnings = append(result.Warnings, feasibility.Warnings...)
	
	return result, nil
}

func (pm *pmAgent) SubmitForVoting(ctx context.Context, proposal *Proposal) error {
	proposal.Status = StatusVoting
	
	if pm.consensusClient == nil {
		return fmt.Errorf("consensus client not configured")
	}
	
	return pm.consensusClient.SubmitProposal(ctx, proposal)
}

func (pm *pmAgent) flagForHumanApproval(proposal *Proposal) {
	if proposal.SecurityImpact.Level == SecCritical || 
	   proposal.SecurityImpact.Level == SecHigh {
		proposal.RequiresHumanApproval = true
		proposal.ApprovalReason = "High security impact"
		return
	}
	
	if proposal.RiskAssessment.OverallRisk == RiskCritical {
		proposal.RequiresHumanApproval = true
		proposal.ApprovalReason = "Critical risk level"
		return
	}
	
	if len(proposal.ComplianceImpact.AffectedStandards) > 0 {
		proposal.RequiresHumanApproval = true
		proposal.ApprovalReason = "Compliance impact"
		return
	}
}

func (pm *pmAgent) loadContext(ctx context.Context) (WorkspaceContext, error) {
	arch, err := pm.contextProvider.LoadArchitecture(ctx)
	if err != nil {
		arch = ""
	}
	
	roadmap, err := pm.contextProvider.LoadRoadmap(ctx)
	if err != nil {
		roadmap = ""
	}
	
	memory, err := pm.contextProvider.LoadMemory(ctx)
	if err != nil {
		memory = ""
	}
	
	commits, err := pm.contextProvider.GetRecentCommits(ctx, 20)
	if err != nil {
		commits = []Commit{}
	}
	
	prds, err := pm.contextProvider.GetExistingProposals(ctx)
	if err != nil {
		prds = []*Proposal{}
	}
	
	return WorkspaceContext{
		Architecture:  arch,
		Roadmap:       roadmap,
		Memory:        memory,
		RecentCommits: commits,
		ExistingPRDs:  prds,
	}, nil
}

// cmd/pm-agent/main.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shreyanshjain7174/clawdlinux/pkg/agents/pm"
)

func main() {
	taskTitle := flag.String("task", "", "Task title")
	taskDesc := flag.String("desc", "", "Task description")
	priority := flag.String("priority", "P1", "Priority (P0/P1/P2)")
	workspaceRoot := flag.String("workspace", os.Getenv("HOME")+"/.openclaw/workspace", "Workspace root")
	flag.Parse()
	
	if *taskTitle == "" {
		log.Fatal("--task required")
	}
	
	agent := setupPMAgent(*workspaceRoot)
	
	task := pm.Task{
		ID:          fmt.Sprintf("task-%d", time.Now().Unix()),
		Title:       *taskTitle,
		Description: *taskDesc,
		Priority:    pm.Priority(*priority),
		Source:      "cli",
		CreatedAt:   time.Now(),
	}
	
	fmt.Printf("Generating proposal for task: %s\n", task.Title)
	fmt.Println("---")
	
	proposal, err := agent.GenerateProposal(context.Background(), task)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}
	
	data, _ := json.MarshalIndent(proposal, "", "  ")
	fmt.Println(string(data))
	
	result, _ := agent.ValidateProposal(context.Background(), proposal)
	fmt.Printf("\n---\nValidation Score: %d/100\n", result.Score)
	
	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
	
	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warn := range result.Warnings {
			fmt.Printf("  - %s\n", warn)
		}
	}
	
	if result.Valid {
		fmt.Println("\n✅ Proposal is valid and ready for voting")
	} else {
		fmt.Println("\n❌ Proposal has validation errors")
	}
}

func setupPMAgent(workspaceRoot string) pm.PMAgent {
	llmClient := &mockLLMClient{}
	gitClient := &mockGitClient{}
	
	config := pm.PMAgentConfig{
		WorkspaceRoot:        workspaceRoot,
		LLMClient:            llmClient,
		Model:                "anthropic/claude-sonnet-4-6",
		GitClient:            gitClient,
		ConsensusClient:      nil,
		MissionControlClient: nil,
	}
	
	return pm.NewPMAgent(config)
}

type mockLLMClient struct{}

func (m *mockLLMClient) Complete(ctx context.Context, prompt string, model string) (string, error) {
	return `{
		"title": "Sample Proposal",
		"problem_statement": "Need to implement the requested feature",
		"goals": ["Deliver working implementation", "Maintain code quality"],
		"success_metrics": [
			{"name": "Implementation Complete", "target": "100%", "measurable": true}
		],
		"technical_approach": {
			"overview": "Implement using Go with standard libraries",
			"components": ["Core logic", "Tests"],
			"technologies": ["Go"],
			"integration_points": ["Existing system"],
			"alternatives_considered": []
		},
		"requirements": [
			{
				"id": "REQ-001",
				"type": "functional",
				"description": "Core functionality implemented",
				"priority": "P0",
				"testable": true,
				"test_method": "Unit tests"
			}
		],
		"dependencies": [],
		"edge_cases": [
			{
				"scenario": "Invalid input",
				"impact": "Error handling required",
				"mitigation": "Validation layer"
			}
		],
		"business_impact": {
			"revenue_impact": "low",
			"compliance_impact": "none",
			"performance_impact": "low",
			"customer_impact": "medium",
			"rationale": "Improves user experience"
		},
		"effort_estimate": {
			"size": "M",
			"story_points": 3,
			"engineering": "2-3 days",
			"testing": "1 day",
			"documentation": "0.5 day",
			"total": "3.5-4.5 days"
		},
		"roadmap_alignment": {
			"phase": "Phase 1",
			"milestone": "MVP",
			"okrs": ["Build core features"],
			"blocked_by": []
		},
		"risk_assessment": {
			"overall_risk": "low",
			"risks": [],
			"mitigation_strategy": "Standard testing and review"
		},
		"security_impact": {
			"level": "low",
			"considerations": [],
			"requires_audit": false
		},
		"compliance_impact": {
			"affected_standards": [],
			"requires_review": false,
			"notes": ""
		}
	}`, nil
}

type mockGitClient struct{}

func (m *mockGitClient) Add(path string) error {
	fmt.Printf("Git add: %s\n", path)
	return nil
}

func (m *mockGitClient) Commit(message string) error {
	fmt.Printf("Git commit: %s\n", message)
	return nil
}

func (m *mockGitClient) Push() error {
	fmt.Println("Git push")
	return nil
}
```