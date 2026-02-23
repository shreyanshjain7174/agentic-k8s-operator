// Auto-generated for phase1-010
// Task: Create GitHub integration for automated PR creation from approved proposals

```go
// pkg/models/proposal.go
package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type AgentRole string

const (
	RolePM            AgentRole = "pm"
	RoleArchitect     AgentRole = "architect"
	RoleEngineer      AgentRole = "engineer"
	RoleQA            AgentRole = "qa"
	RoleSecurity      AgentRole = "security"
	RoleDocumentation AgentRole = "documentation"
)

type VoteDecision string

const (
	DecisionApprove            VoteDecision = "APPROVE"
	DecisionConditionalApprove VoteDecision = "CONDITIONAL_APPROVE"
	DecisionReject             VoteDecision = "REJECT"
)

type AgentVote struct {
	AgentRole AgentRole    `json:"agent_role"`
	Decision  VoteDecision `json:"decision"`
	Score     float64      `json:"score"` // 0.0 - 100.0
	Reasoning string       `json:"reasoning"`
	Timestamp time.Time    `json:"timestamp"`
}

type ProposalMetadata struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LoopNumber  int       `json:"loop_number"`
}

type ApprovedProposal struct {
	Metadata              ProposalMetadata  `json:"metadata"`
	Votes                 []AgentVote       `json:"votes"`
	ConsensusScore        float64           `json:"consensus_score"`
	ImplementationFiles   map[string]string `json:"implementation_files"`
	ImplementationSummary string            `json:"implementation_summary"`
}

type PRMetadata struct {
	ProposalID string     `json:"proposal_id"`
	BranchName string     `json:"branch_name"`
	PRNumber   *int       `json:"pr_number,omitempty"`
	PRURL      *string    `json:"pr_url,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	Error      *string    `json:"error,omitempty"`
}

// GetBranchName generates branch name: swarm/proposal-123-title-slug
func (p *ApprovedProposal) GetBranchName() string {
	slug := strings.ToLower(p.Metadata.Title)
	if len(slug) > 40 {
		slug = slug[:40]
	}

	// Replace non-alphanumeric chars with hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	proposalIDShort := p.Metadata.ID
	if len(proposalIDShort) > 8 {
		proposalIDShort = proposalIDShort[:8]
	}

	return fmt.Sprintf("swarm/proposal-%s-%s", proposalIDShort, slug)
}

// Title returns the title-cased agent role name
func (r AgentRole) Title() string {
	return strings.Title(string(r))
}
```

```go
// pkg/github/manager.go
package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v58/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Manager handles GitHub API interactions
type Manager struct {
	client *github.Client
	owner  string
	repo   string
	logger *logrus.Logger
}

// NewManager creates a new GitHub manager
func NewManager(token, owner, repo string, logger *logrus.Logger) *Manager {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Manager{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
		logger: logger,
	}
}

// CreateBranch creates a new branch from base branch
func (m *Manager) CreateBranch(ctx context.Context, branchName, baseBranch string) error {
	// Get base branch reference
	baseRef, _, err := m.client.Git.GetRef(ctx, m.owner, m.repo, fmt.Sprintf("heads/%s", baseBranch))
	if err != nil {
		m.logger.WithError(err).Errorf("Failed to get base branch %s", baseBranch)
		return fmt.Errorf("get base ref: %w", err)
	}

	// Create new branch
	newRef := &github.Reference{
		Ref: github.String(fmt.Sprintf("refs/heads/%s", branchName)),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}

	_, _, err = m.client.Git.CreateRef(ctx, m.owner, m.repo, newRef)
	if err != nil {
		m.logger.WithError(err).Errorf("Failed to create branch %s", branchName)
		return fmt.Errorf("create branch: %w", err)
	}

	m.logger.Infof("Created branch: %s", branchName)
	return nil
}

// CommitFiles commits multiple files to a branch
func (m *Manager) CommitFiles(ctx context.Context, branchName string, files map[string]string, commitMessage string) error {
	// Get branch reference
	ref, _, err := m.client.Git.GetRef(ctx, m.owner, m.repo, fmt.Sprintf("heads/%s", branchName))
	if err != nil {
		return fmt.Errorf("get branch ref: %w", err)
	}

	// Get base tree
	baseCommit, _, err := m.client.Git.GetCommit(ctx, m.owner, m.repo, *ref.Object.SHA)
	if err != nil {
		return fmt.Errorf("get base commit: %w", err)
	}

	// Create tree entries for each file
	var entries []*github.TreeEntry
	for filepath, content := range files {
		// Create blob for file content
		blob := &github.Blob{
			Content:  github.String(content),
			Encoding: github.String("utf-8"),
		}
		createdBlob, _, err := m.client.Git.CreateBlob(ctx, m.owner, m.repo, blob)
		if err != nil {
			return fmt.Errorf("create blob for %s: %w", filepath, err)
		}

		entries = append(entries, &github.TreeEntry{
			Path: github.String(filepath),
			Mode: github.String("100644"),
			Type: github.String("blob"),
			SHA:  createdBlob.SHA,
		})
	}

	// Create new tree
	tree, _, err := m.client.Git.CreateTree(ctx, m.owner, m.repo, *baseCommit.Tree.SHA, entries)
	if err != nil {
		return fmt.Errorf("create tree: %w", err)
	}

	// Create commit
	commit := &github.Commit{
		Message: github.String(commitMessage),
		Tree:    tree,
		Parents: []*github.Commit{baseCommit},
	}
	newCommit, _, err := m.client.Git.CreateCommit(ctx, m.owner, m.repo, commit, nil)
	if err != nil {
		return fmt.Errorf("create commit: %w", err)
	}

	// Update branch reference
	ref.Object.SHA = newCommit.SHA
	_, _, err = m.client.Git.UpdateRef(ctx, m.owner, m.repo, ref, false)
	if err != nil {
		return fmt.Errorf("update ref: %w", err)
	}

	m.logger.Infof("Committed %d files to %s", len(files), branchName)
	return nil
}

// CreatePullRequest creates a PR and returns (pr_number, pr_url, error)
func (m *Manager) CreatePullRequest(ctx context.Context, branchName, title, body, base string) (int, string, error) {
	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branchName),
		Base:  github.String(base),
		Body:  github.String(body),
	}

	createdPR, _, err := m.client.PullRequests.Create(ctx, m.owner, m.repo, pr)
	if err != nil {
		m.logger.WithError(err).Error("Failed to create PR")
		return 0, "", fmt.Errorf("create PR: %w", err)
	}

	m.logger.Infof("Created PR #%d: %s", *createdPR.Number, *createdPR.HTMLURL)
	return *createdPR.Number, *createdPR.HTMLURL, nil
}

// AddLabels adds labels to a PR
func (m *Manager) AddLabels(ctx context.Context, prNumber int, labels []string) error {
	_, _, err := m.client.Issues.AddLabelsToIssue(ctx, m.owner, m.repo, prNumber, labels)
	if err != nil {
		m.logger.WithError(err).Errorf("Failed to add labels to PR #%d", prNumber)
		return fmt.Errorf("add labels: %w", err)
	}

	m.logger.Infof("Added labels to PR #%d: %v", prNumber, labels)
	return nil
}
```

```go
// pkg/github/pr_generator.go
package github

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shreyanshjain7174/clawdlinux/pkg/models"
)

// Generator generates PR titles and bodies from approved proposals
type Generator struct{}

// NewGenerator creates a new PR generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateTitle generates PR title
func (g *Generator) GenerateTitle(proposal *models.ApprovedProposal) string {
	return fmt.Sprintf("[Swarm] %s", proposal.Metadata.Title)
}

// GenerateBody generates PR body with metadata and checklist
func (g *Generator) GenerateBody(proposal *models.ApprovedProposal) string {
	votesTable := g.generateVotesTable(proposal.Votes)
	checklist := g.generateChecklist()

	body := fmt.Sprintf(`## 🤖 Auto-Generated by ClawdLinux Autonomous Swarm

### Proposal Summary
%s

### Implementation
%s

### Consensus Scores
**Overall Approval: %.1f/100**

%s

### Review Checklist
%s

---
**Metadata:**
- Proposal ID: ` + "`%s`" + `
- Loop: #%d
- Created: %s
- Branch: ` + "`%s`" + `

[View Proposal Details](../proposals/%s.json)
`,
		proposal.Metadata.Description,
		proposal.ImplementationSummary,
		proposal.ConsensusScore,
		votesTable,
		checklist,
		proposal.Metadata.ID,
		proposal.Metadata.LoopNumber,
		proposal.Metadata.CreatedAt.Format("2006-01-02 15:04:05 UTC"),
		proposal.GetBranchName(),
		proposal.Metadata.ID,
	)

	return body
}

// generateVotesTable generates markdown table of agent votes
func (g *Generator) generateVotesTable(votes []models.AgentVote) string {
	// Sort votes by score descending
	sortedVotes := make([]models.AgentVote, len(votes))
	copy(sortedVotes, votes)
	sort.Slice(sortedVotes, func(i, j int) bool {
		return sortedVotes[i].Score > sortedVotes[j].Score
	})

	rows := []string{
		"| Agent | Decision | Score | Reasoning |",
		"|-------|----------|-------|-----------|",
	}

	for _, vote := range sortedVotes {
		decisionEmoji := map[models.VoteDecision]string{
			models.DecisionApprove:            "✅",
			models.DecisionConditionalApprove: "⚠️",
			models.DecisionReject:             "❌",
		}
		emoji := decisionEmoji[vote.Decision]
		if emoji == "" {
			emoji = "❓"
		}

		reasoning := vote.Reasoning
		if len(reasoning) > 60 {
			reasoning = reasoning[:60] + "..."
		}

		row := fmt.Sprintf("| %s | %s %s | %.1f | %s |",
			vote.AgentRole.Title(),
			emoji,
			vote.Decision,
			vote.Score,
			reasoning,
		)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

// generateChecklist generates human reviewer checklist
func (g *Generator) generateChecklist() string {
	return `- [ ] Code quality meets standards
- [ ] Tests pass (if applicable)
- [ ] Documentation updated
- [ ] No security vulnerabilities introduced
- [ ] Consensus scores reviewed and acceptable`
}
```

```go
// pkg/github/retry.go
package github

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// WithRetry executes a function with exponential backoff retry logic
func WithRetry(ctx context.Context, config *RetryConfig, logger *logrus.Logger, funcName string, fn RetryableFunc) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt == config.MaxAttempts {
			logger.WithError(err).Errorf("%s failed after %d attempts", funcName, attempt)
			return fmt.Errorf("%s failed after %d attempts: %w", funcName, config.MaxAttempts, err)
		}

		// Calculate exponential backoff delay
		delay := time.Duration(float64(config.BaseDelay) * math.Pow(2, float64(attempt-1)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		logger.WithError(err).Warnf(
			"%s failed (attempt %d/%d), retrying in %s",
			funcName,
			attempt,
			config.MaxAttempts,
			delay,
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}
```

```go
// pkg/github/orchestrator.go
package github

import (
	"context"
	"fmt"
	"time"

	"github.com/shreyanshjain7174/clawdlinux/pkg/models"
	"github.com/sirupsen/logrus"
)

// Orchestrator orchestrates the complete PR creation workflow
type Orchestrator struct {
	manager   *Manager
	generator *Generator
	logger    *logrus.Logger
	config    *RetryConfig
}

// NewOrchestrator creates a new PR orchestrator
func NewOrchestrator(manager *Manager, logger *logrus.Logger) *Orchestrator {
	return &Orchestrator{
		manager:   manager,
		generator: NewGenerator(),
		logger:    logger,
		config:    DefaultRetryConfig(),
	}
}

// CreatePRFromProposal executes complete PR creation workflow with retry logic
func (o *Orchestrator) CreatePRFromProposal(ctx context.Context, proposal *models.ApprovedProposal) (*models.PRMetadata, error) {
	branchName := proposal.GetBranchName()
	metadata := &models.PRMetadata{
		ProposalID: proposal.Metadata.ID,
		BranchName: branchName,
		CreatedAt:  time.Now().UTC(),
	}

	err := WithRetry(ctx, o.config, o.logger, "CreatePRFromProposal", func(ctx context.Context) error {
		// Step 1: Create branch
		o.logger.Infof("Creating branch: %s", branchName)
		if err := o.manager.CreateBranch(ctx, branchName, "main"); err != nil {
			return fmt.Errorf("branch creation failed: %w", err)
		}

		// Step 2: Commit files
		o.logger.Infof("Committing %d files", len(proposal.ImplementationFiles))
		commitMsg := fmt.Sprintf("%s\n\nConsensus: %.1f%%", proposal.Metadata.Title, proposal.ConsensusScore)
		if err := o.manager.CommitFiles(ctx, branchName, proposal.ImplementationFiles, commitMsg); err != nil {
			return fmt.Errorf("file commit failed: %w", err)
		}

		// Step 3: Create PR
		title := o.generator.GenerateTitle(proposal)
		body := o.generator.GenerateBody(proposal)

		o.logger.Info("Creating pull request")
		prNumber, prURL, err := o.manager.CreatePullRequest(ctx, branchName, title, body, "main")
		if err != nil {
			return fmt.Errorf("PR creation failed: %w", err)
		}

		metadata.PRNumber = &prNumber
		metadata.PRURL = &prURL

		// Step 4: Add labels
		labels := []string{"swarm-generated", "autonomous"}
		if proposal.ConsensusScore >= 90 {
			labels = append(labels, "high-confidence")
		}
		if err := o.manager.AddLabels(ctx, prNumber, labels); err != nil {
			// Log but don't fail on label errors
			o.logger.WithError(err).Warn("Failed to add labels (non-fatal)")
		}

		o.logger.Infof("✅ PR created successfully: %s", prURL)
		return nil
	})

	if err != nil {
		errStr := err.Error()
		metadata.Error = &errStr
		o.logger.WithError(err).Error("❌ PR creation failed")
		return metadata, err
	}

	return metadata, nil
}
```

```go
// pkg/github/integration.go
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/shreyanshjain7174/clawdlinux/pkg/models"
	"github.com/sirupsen/logrus"
)

// Config holds GitHub integration configuration
type Config struct {
	Token     string
	RepoOwner string
	RepoName  string
}

// LoadConfigFromEnv loads GitHub configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	owner := os.Getenv("GITHUB_REPO_OWNER")
	if owner == "" {
		owner = "shreyanshjain7174"
	}

	repo := os.Getenv("GITHUB_REPO_NAME")
	if repo == "" {
		repo = "clawdlinux"
	}

	return &Config{
		Token:     token,
		RepoOwner: owner,
		RepoName:  repo,
	}, nil
}

// PRCreationLogger logs PR creation events in structured JSON format
type PRCreationLogger struct {
	logger *logrus.Logger
}

// NewPRCreationLogger creates a new PR creation logger
func NewPRCreationLogger(logger *logrus.Logger) *PRCreationLogger {
	return &PRCreationLogger{logger: logger}
}

// LogPRCreation logs a PR creation event
func (l *PRCreationLogger) LogPRCreation(proposalID string, metadata *models.PRMetadata, success bool) {
	logEntry := map[string]interface{}{
		"event":       "pr_created",
		"proposal_id": proposalID,
		"pr_number":   metadata.PRNumber,
		"pr_url":      metadata.PRURL,
		"success":     success,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	if metadata.Error != nil {
		logEntry["error"] = *metadata.Error
	}

	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		l.logger.WithError(err).Error("Failed to marshal log entry")
		return
	}

	l.logger.Info(string(jsonBytes))
}

// EvaluateConsensusAndCreatePR evaluates consensus and triggers PR creation if approved
func EvaluateConsensusAndCreatePR(ctx context.Context, proposal *models.ApprovedProposal, orchestrator *Orchestrator, logger *logrus.Logger) error {
	prLogger := NewPRCreationLogger(logger)

	if proposal.ConsensusScore >= 80.0 {
		logger.Infof("Proposal %s approved with %.1f%% consensus, creating PR...", proposal.Metadata.ID, proposal.ConsensusScore)

		metadata, err := orchestrator.CreatePRFromProposal(ctx, proposal)
		if err != nil {
			prLogger.LogPRCreation(proposal.Metadata.ID, metadata, false)
			logger.WithError(err).Error("❌ Failed to auto-create PR")
			logger.Info("💡 Fallback: Manual PR creation required")
			return err
		}

		prLogger.LogPRCreation(proposal.Metadata.ID, metadata, true)
		logger.Infof("✅ Auto-created PR: %s", *metadata.PRURL)
		return nil
	}

	logger.Infof("Proposal %s rejected with %.1f%% consensus (threshold: 80%%)", proposal.Metadata.ID, proposal.ConsensusScore)
	return nil
}
```