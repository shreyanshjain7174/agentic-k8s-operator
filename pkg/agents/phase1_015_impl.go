// Auto-generated for phase1-015
// Task: Build integration tests for end-to-end proposal lifecycle

```go
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AgentDecision represents the decision made by an agent
type AgentDecision string

const (
	AgentDecisionApprove            AgentDecision = "APPROVE"
	AgentDecisionConditionalApprove AgentDecision = "CONDITIONAL_APPROVE"
	AgentDecisionReject             AgentDecision = "REJECT"
)

// ProposalStatus represents the current status of a proposal
type ProposalStatus string

const (
	ProposalStatusDraft     ProposalStatus = "DRAFT"
	ProposalStatusReview    ProposalStatus = "REVIEW"
	ProposalStatusApproved  ProposalStatus = "APPROVED"
	ProposalStatusRejected  ProposalStatus = "REJECTED"
	ProposalStatusExecuting ProposalStatus = "EXECUTING"
	ProposalStatusCompleted ProposalStatus = "COMPLETED"
	ProposalStatusFailed    ProposalStatus = "FAILED"
)

// AgentReview contains the review details from a single agent
type AgentReview struct {
	Agent      string        `json:"agent"`
	Score      int           `json:"score"`
	Decision   AgentDecision `json:"decision"`
	Reason     string        `json:"reason"`
	Timestamp  time.Time     `json:"timestamp"`
	LatencyMS  int           `json:"latency_ms"`
}

// ProposalContext contains all metadata for a proposal
type ProposalContext struct {
	ProposalID        string          `json:"proposal_id"`
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	Type              string          `json:"type"`
	Reviews           []AgentReview   `json:"reviews"`
	AvgScore          float64         `json:"avg_score"`
	ApprovalRate      float64         `json:"approval_rate"`
	ConsensusDecision AgentDecision   `json:"consensus_decision"`
	Status            ProposalStatus  `json:"status"`
	FilesChanged      []string        `json:"files_changed"`
	MCTaskID          *string         `json:"mc_task_id,omitempty"`
}

// QuotaState tracks API quota usage
type QuotaState struct {
	TotalCalls   int                `json:"total_calls"`
	ModelCalls   map[string]int     `json:"model_calls"`
	WindowStart  time.Time          `json:"window_start"`
	WindowEndIn  time.Duration      `json:"window_end_in"`
}

// IAgentCaller defines the interface for agent execution
type IAgentCaller interface {
	CallAgent(ctx context.Context, sessionID string, message string, timeout int) (string, int, error)
}

// IGitManager defines the interface for Git operations
type IGitManager interface {
	CreateBranch(ctx context.Context, branchName string) error
	CommitChanges(ctx context.Context, message string, files []string) error
	PushBranch(ctx context.Context, branchName string) error
	Rollback(ctx context.Context, commitHash string) error
}

// IMissionControlClient defines the interface for Mission Control API
type IMissionControlClient interface {
	CreateTask(ctx context.Context, title string, description string) (string, error)
	UpdateTask(ctx context.Context, taskID string, status string, notes string) error
}

// IQuotaManager defines the interface for quota tracking
type IQuotaManager interface {
	CheckQuota(ctx context.Context, model string) (bool, error)
	RecordCall(ctx context.Context, model string) error
	GetState(ctx context.Context) (QuotaState, error)
}

// MockAgentCaller implements IAgentCaller for testing
type MockAgentCaller struct {
	mu        sync.Mutex
	responses map[string]string
	calls     []AgentCall
	latency   time.Duration
}

type AgentCall struct {
	SessionID string
	Message   string
	Timeout   int
	Timestamp time.Time
}

func NewMockAgentCaller() *MockAgentCaller {
	return &MockAgentCaller{
		responses: make(map[string]string),
		calls:     make([]AgentCall, 0),
		latency:   100 * time.Millisecond,
	}
}

func (m *MockAgentCaller) SetResponse(agent string, response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[agent] = response
}

func (m *MockAgentCaller) SetLatency(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.latency = latency
}

func (m *MockAgentCaller) CallAgent(ctx context.Context, sessionID string, message string, timeout int) (string, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.calls = append(m.calls, AgentCall{
		SessionID: sessionID,
		Message:   message,
		Timeout:   timeout,
		Timestamp: time.Now(),
	})
	
	time.Sleep(m.latency)
	
	response, ok := m.responses[sessionID]
	if !ok {
		response = fmt.Sprintf(`{"score": 85, "decision": "APPROVE", "reason": "Mock approval"}`)
	}
	
	return response, int(m.latency.Milliseconds()), nil
}

func (m *MockAgentCaller) GetCalls() []AgentCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]AgentCall, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// MockGitManager implements IGitManager for testing
type MockGitManager struct {
	mu       sync.Mutex
	branches map[string]bool
	commits  []GitCommit
	pushes   []string
	failNext bool
}

type GitCommit struct {
	Branch  string
	Message string
	Files   []string
	Hash    string
	Time    time.Time
}

func NewMockGitManager() *MockGitManager {
	return &MockGitManager{
		branches: make(map[string]bool),
		commits:  make([]GitCommit, 0),
		pushes:   make([]string, 0),
	}
}

func (m *MockGitManager) CreateBranch(ctx context.Context, branchName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return fmt.Errorf("git: failed to create branch %s", branchName)
	}
	
	m.branches[branchName] = true
	return nil
}

func (m *MockGitManager) CommitChanges(ctx context.Context, message string, files []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return fmt.Errorf("git: failed to commit changes")
	}
	
	commit := GitCommit{
		Message: message,
		Files:   files,
		Hash:    fmt.Sprintf("abc%d", len(m.commits)),
		Time:    time.Now(),
	}
	m.commits = append(m.commits, commit)
	return nil
}

func (m *MockGitManager) PushBranch(ctx context.Context, branchName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return fmt.Errorf("git: failed to push branch %s", branchName)
	}
	
	m.pushes = append(m.pushes, branchName)
	return nil
}

func (m *MockGitManager) Rollback(ctx context.Context, commitHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return fmt.Errorf("git: failed to rollback to %s", commitHash)
	}
	
	return nil
}

func (m *MockGitManager) SetFailNext(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failNext = fail
}

func (m *MockGitManager) GetCommits() []GitCommit {
	m.mu.Lock()
	defer m.mu.Unlock()
	commits := make([]GitCommit, len(m.commits))
	copy(commits, m.commits)
	return commits
}

// MockMissionControlClient implements IMissionControlClient for testing
type MockMissionControlClient struct {
	mu       sync.Mutex
	tasks    map[string]MCTask
	failNext bool
}

type MCTask struct {
	ID          string
	Title       string
	Description string
	Status      string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewMockMissionControlClient() *MockMissionControlClient {
	return &MockMissionControlClient{
		tasks: make(map[string]MCTask),
	}
}

func (m *MockMissionControlClient) CreateTask(ctx context.Context, title string, description string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return "", fmt.Errorf("mission control: failed to create task")
	}
	
	taskID := fmt.Sprintf("task-%d", len(m.tasks)+1)
	m.tasks[taskID] = MCTask{
		ID:          taskID,
		Title:       title,
		Description: description,
		Status:      "open",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	return taskID, nil
}

func (m *MockMissionControlClient) UpdateTask(ctx context.Context, taskID string, status string, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return fmt.Errorf("mission control: failed to update task %s", taskID)
	}
	
	task, ok := m.tasks[taskID]
	if !ok {
		return fmt.Errorf("mission control: task %s not found", taskID)
	}
	
	task.Status = status
	task.Notes = notes
	task.UpdatedAt = time.Now()
	m.tasks[taskID] = task
	
	return nil
}

func (m *MockMissionControlClient) SetFailNext(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failNext = fail
}

func (m *MockMissionControlClient) GetTask(taskID string) (MCTask, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task, ok := m.tasks[taskID]
	return task, ok
}

// MockQuotaManager implements IQuotaManager for testing
type MockQuotaManager struct {
	mu          sync.Mutex
	totalCalls  int
	modelCalls  map[string]int
	windowStart time.Time
	quotaLimit  int
	failNext    bool
}

func NewMockQuotaManager(quotaLimit int) *MockQuotaManager {
	return &MockQuotaManager{
		modelCalls:  make(map[string]int),
		windowStart: time.Now(),
		quotaLimit:  quotaLimit,
	}
}

func (m *MockQuotaManager) CheckQuota(ctx context.Context, model string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.failNext {
		m.failNext = false
		return false, fmt.Errorf("quota: check failed")
	}
	
	return m.totalCalls < m.quotaLimit, nil
}

func (m *MockQuotaManager) RecordCall(ctx context.Context, model string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.totalCalls++
	m.modelCalls[model]++
	return nil
}

func (m *MockQuotaManager) GetState(ctx context.Context) (QuotaState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	windowEnd := m.windowStart.Add(5 * time.Hour)
	return QuotaState{
		TotalCalls:  m.totalCalls,
		ModelCalls:  m.modelCalls,
		WindowStart: m.windowStart,
		WindowEndIn: time.Until(windowEnd),
	}, nil
}

func (m *MockQuotaManager) SetFailNext(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failNext = fail
}

func (m *MockQuotaManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalCalls = 0
	m.modelCalls = make(map[string]int)
	m.windowStart = time.Now()
}

// SwarmCoordinator orchestrates the proposal lifecycle
type SwarmCoordinator struct {
	agentCaller   IAgentCaller
	gitManager    IGitManager
	mcClient      IMissionControlClient
	quotaManager  IQuotaManager
	
	mu        sync.Mutex
	proposals map[string]*ProposalContext
}

func NewSwarmCoordinator(
	agentCaller IAgentCaller,
	gitManager IGitManager,
	mcClient IMissionControlClient,
	quotaManager IQuotaManager,
) *SwarmCoordinator {
	return &SwarmCoordinator{
		agentCaller:  agentCaller,
		gitManager:   gitManager,
		mcClient:     mcClient,
		quotaManager: quotaManager,
		proposals:    make(map[string]*ProposalContext),
	}
}

func (s *SwarmCoordinator) CreateProposal(ctx context.Context, title string, description string, proposalType string) (*ProposalContext, error) {
	proposalID := fmt.Sprintf("proposal-%d", time.Now().Unix())
	
	proposal := &ProposalContext{
		ProposalID:  proposalID,
		Title:       title,
		Description: description,
		Type:        proposalType,
		Reviews:     make([]AgentReview, 0),
		Status:      ProposalStatusDraft,
	}
	
	s.mu.Lock()
	s.proposals[proposalID] = proposal
	s.mu.Unlock()
	
	return proposal, nil
}

func (s *SwarmCoordinator) ReviewProposal(ctx context.Context, proposalID string) error {
	s.mu.Lock()
	proposal, ok := s.proposals[proposalID]
	s.mu.Unlock()
	
	if !ok {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	
	proposal.Status = ProposalStatusReview
	
	agents := []string{"architect", "engineer", "qa", "security"}
	reviews := make([]AgentReview, 0, len(agents))
	
	for _, agent := range agents {
		hasQuota, err := s.quotaManager.CheckQuota(ctx, "sonnet")
		if err != nil {
			return fmt.Errorf("quota check failed: %w", err)
		}
		if !hasQuota {
			return fmt.Errorf("quota exhausted")
		}
		
		start := time.Now()
		response, latency, err := s.agentCaller.CallAgent(ctx, agent, proposal.Description, 30)
		if err != nil {
			return fmt.Errorf("agent %s failed: %w", agent, err)
		}
		
		if err := s.quotaManager.RecordCall(ctx, "sonnet"); err != nil {
			return fmt.Errorf("quota record failed: %w", err)
		}
		
		var reviewData struct {
			Score    int    `json:"score"`
			Decision string `json:"decision"`
			Reason   string `json:"reason"`
		}
		
		if err := json.Unmarshal([]byte(response), &reviewData); err != nil {
			return fmt.Errorf("failed to parse agent response: %w", err)
		}
		
		review := AgentReview{
			Agent:     agent,
			Score:     reviewData.Score,
			Decision:  AgentDecision(reviewData.Decision),
			Reason:    reviewData.Reason,
			Timestamp: start,
			LatencyMS: latency,
		}
		
		reviews = append(reviews, review)
	}
	
	proposal.Reviews = reviews
	
	totalScore := 0
	approvals := 0
	for _, review := range reviews {
		totalScore += review.Score
		if review.Decision == AgentDecisionApprove || review.Decision == AgentDecisionConditionalApprove {
			approvals++
		}
	}
	
	proposal.AvgScore = float64(totalScore) / float64(len(reviews))
	proposal.ApprovalRate = float64(approvals) / float64(len(reviews)) * 100
	
	if proposal.ApprovalRate >= 80 {
		proposal.ConsensusDecision = AgentDecisionApprove
		proposal.Status = ProposalStatusApproved
	} else {
		proposal.ConsensusDecision = AgentDecisionReject
		proposal.Status = ProposalStatusRejected
	}
	
	return nil
}

func (s *SwarmCoordinator) ExecuteProposal(ctx context.Context, proposalID string) error {
	s.mu.Lock()
	proposal, ok := s.proposals[proposalID]
	s.mu.Unlock()
	
	if !ok {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	
	if proposal.Status != ProposalStatusApproved {
		return fmt.Errorf("proposal not approved: status=%s", proposal.Status)
	}
	
	proposal.Status = ProposalStatusExecuting
	
	branchName := fmt.Sprintf("proposal-%s", proposalID)
	if err := s.gitManager.CreateBranch(ctx, branchName); err != nil {
		proposal.Status = ProposalStatusFailed
		return fmt.Errorf("failed to create branch: %w", err)
	}
	
	files := []string{"main.go", "README.md"}
	proposal.FilesChanged = files
	
	if err := s.gitManager.CommitChanges(ctx, proposal.Title, files); err != nil {
		proposal.Status = ProposalStatusFailed
		return fmt.Errorf("failed to commit: %w", err)
	}
	
	if err := s.gitManager.PushBranch(ctx, branchName); err != nil {
		proposal.Status = ProposalStatusFailed
		return fmt.Errorf("failed to push: %w", err)
	}
	
	taskID, err := s.mcClient.CreateTask(ctx, proposal.Title, proposal.Description)
	if err != nil {
		proposal.Status = ProposalStatusFailed
		return fmt.Errorf("failed to create mission control task: %w", err)
	}
	
	proposal.MCTaskID = &taskID
	
	if err := s.mcClient.UpdateTask(ctx, taskID, "completed", "Proposal executed successfully"); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	
	proposal.Status = ProposalStatusCompleted
	return nil
}

func (s *SwarmCoordinator) GetProposal(proposalID string) (*ProposalContext, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	proposal, ok := s.proposals[proposalID]
	return proposal, ok
}
```