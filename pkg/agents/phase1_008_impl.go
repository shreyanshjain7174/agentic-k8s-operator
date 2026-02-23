// Auto-generated for phase1-008
// Task: Implement consensus voting system (80% threshold for approval)

```go
// consensus/types.go
package consensus

import (
	"time"
)

// Vote represents a single agent's review decision
type Vote struct {
	VoteID     string                 `json:"vote_id"`
	ProposalID string                 `json:"proposal_id"`
	AgentID    string                 `json:"agent_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Decision   VoteDecision           `json:"decision"`
	Score      float64                `json:"score"`
	Criteria   map[string]string      `json:"criteria"`
	Feedback   string                 `json:"feedback"`
	Weight     float64                `json:"weight"`
}

// VoteDecision represents the outcome of a review
type VoteDecision string

const (
	VoteApprove            VoteDecision = "APPROVE"
	VoteConditionalApprove VoteDecision = "CONDITIONAL_APPROVE"
	VoteConditionalReject  VoteDecision = "CONDITIONAL_REJECT"
	VoteReject             VoteDecision = "REJECT"
)

// ToScore converts a decision to numeric score
func (d VoteDecision) ToScore() float64 {
	switch d {
	case VoteApprove:
		return 1.0
	case VoteConditionalApprove:
		return 0.75
	case VoteConditionalReject:
		return 0.25
	case VoteReject:
		return 0.0
	default:
		return 0.5
	}
}

// ConsensusResult aggregates all votes for a proposal
type ConsensusResult struct {
	ProposalID      string           `json:"proposal_id"`
	Timestamp       time.Time        `json:"timestamp"`
	TotalVotes      int              `json:"total_votes"`
	QuorumMet       bool             `json:"quorum_met"`
	RawScore        float64          `json:"raw_score"`
	WeightedScore   float64          `json:"weighted_score"`
	ApprovalPercent float64          `json:"approval_percent"`
	StandardDev     float64          `json:"standard_dev"`
	Unanimity       bool             `json:"unanimity"`
	Outcome         ConsensusOutcome `json:"outcome"`
	PassesThreshold bool             `json:"passes_threshold"`
	Votes           []Vote           `json:"votes"`
	AuditHash       string           `json:"audit_hash"`
}

// ConsensusOutcome represents the final decision
type ConsensusOutcome string

const (
	OutcomeApproved    ConsensusOutcome = "APPROVED"
	OutcomeConditional ConsensusOutcome = "CONDITIONAL"
	OutcomeRejected    ConsensusOutcome = "REJECTED"
	OutcomeNoQuorum    ConsensusOutcome = "NO_QUORUM"
)

// VotingConfig defines consensus parameters
type VotingConfig struct {
	ApprovalThreshold   float64                `json:"approval_threshold"`
	ConditionalMin      float64                `json:"conditional_min"`
	QuorumPercent       float64                `json:"quorum_percent"`
	AgentWeights        map[string]float64     `json:"agent_weights"`
	VotingDeadline      time.Duration          `json:"voting_deadline"`
	TiebreakerStrategy  TiebreakerStrategy     `json:"tiebreaker_strategy"`
	TiebreakerAgentID   string                 `json:"tiebreaker_agent_id"`
	EscalateConditional bool                   `json:"escalate_conditional"`
}

// TiebreakerStrategy defines how to handle edge cases
type TiebreakerStrategy string

const (
	TiebreakerWeighted        TiebreakerStrategy = "WEIGHTED"
	TiebreakerDesignatedVoter TiebreakerStrategy = "DESIGNATED"
	TiebreakerHumanReview     TiebreakerStrategy = "HUMAN_REVIEW"
	TiebreakerRejectByDefault TiebreakerStrategy = "REJECT"
)

// AuditEntry represents a single immutable log record
type AuditEntry struct {
	EntryID    string                 `json:"entry_id"`
	Timestamp  time.Time              `json:"timestamp"`
	EventType  string                 `json:"event_type"`
	ProposalID string                 `json:"proposal_id"`
	AgentID    string                 `json:"agent_id"`
	Data       map[string]interface{} `json:"data"`
	PrevHash   string                 `json:"prev_hash"`
	Hash       string                 `json:"hash"`
}
```

```go
// consensus/interfaces.go
package consensus

import (
	"context"
)

// VotingManager orchestrates the voting process
type VotingManager interface {
	CollectVotes(ctx context.Context, proposalID string, agents []string) ([]Vote, error)
	SubmitVote(ctx context.Context, vote Vote) error
	GetVotes(ctx context.Context, proposalID string) ([]Vote, error)
	CheckQuorum(ctx context.Context, proposalID string) (bool, error)
}

// ScoreAggregator calculates consensus metrics
type ScoreAggregator interface {
	CalculateConsensus(votes []Vote, config VotingConfig) (ConsensusResult, error)
	CalculateWeightedScore(votes []Vote, weights map[string]float64) float64
	CalculateStandardDeviation(votes []Vote) float64
	CheckUnanimity(votes []Vote) bool
}

// DecisionEngine applies business logic to vote outcomes
type DecisionEngine interface {
	DetermineOutcome(consensusResult ConsensusResult, config VotingConfig) ConsensusOutcome
	ApplyTiebreaker(consensusResult ConsensusResult, config VotingConfig) (ConsensusOutcome, error)
	ShouldEscalate(consensusResult ConsensusResult, config VotingConfig) bool
}

// AuditLogger maintains immutable vote history
type AuditLogger interface {
	LogVote(ctx context.Context, vote Vote) error
	LogConsensusResult(ctx context.Context, result ConsensusResult) error
	GetAuditTrail(ctx context.Context, proposalID string) ([]AuditEntry, error)
	VerifyIntegrity(ctx context.Context, proposalID string) (bool, error)
}

// Storage backend for votes and audit logs
type Storage interface {
	SaveVote(ctx context.Context, vote Vote) error
	GetVotes(ctx context.Context, proposalID string) ([]Vote, error)
	AppendAuditEntry(ctx context.Context, entry AuditEntry) error
	GetAuditEntries(ctx context.Context, proposalID string) ([]AuditEntry, error)
	GetLastHash(ctx context.Context) (string, error)
	SaveLastHash(ctx context.Context, hash string) error
}
```

```go
// consensus/aggregator.go
package consensus

import (
	"fmt"
	"math"
	"time"
)

// SimpleScoreAggregator implements ScoreAggregator interface
type SimpleScoreAggregator struct{}

// NewSimpleScoreAggregator creates a new score aggregator
func NewSimpleScoreAggregator() *SimpleScoreAggregator {
	return &SimpleScoreAggregator{}
}

// CalculateConsensus computes final score and outcome
func (a *SimpleScoreAggregator) CalculateConsensus(votes []Vote, config VotingConfig) (ConsensusResult, error) {
	if len(votes) == 0 {
		return ConsensusResult{}, fmt.Errorf("no votes provided")
	}

	rawScore := 0.0
	for _, vote := range votes {
		rawScore += vote.Score
	}
	rawScore /= float64(len(votes))

	weightedScore := a.CalculateWeightedScore(votes, config.AgentWeights)
	stdDev := a.CalculateStandardDeviation(votes)
	unanimous := a.CheckUnanimity(votes)

	result := ConsensusResult{
		ProposalID:      votes[0].ProposalID,
		Timestamp:       time.Now(),
		TotalVotes:      len(votes),
		RawScore:        rawScore,
		WeightedScore:   weightedScore,
		ApprovalPercent: weightedScore,
		StandardDev:     stdDev,
		Unanimity:       unanimous,
		Votes:           votes,
	}

	return result, nil
}

// CalculateWeightedScore applies agent weights
func (a *SimpleScoreAggregator) CalculateWeightedScore(votes []Vote, weights map[string]float64) float64 {
	if len(votes) == 0 {
		return 0.0
	}

	weightedScore := 0.0
	totalWeight := 0.0

	for _, vote := range votes {
		weight := weights[vote.AgentID]
		if weight == 0 {
			weight = 1.0
		}
		weightedScore += vote.Score * weight
		totalWeight += weight
	}

	return weightedScore / totalWeight
}

// CalculateStandardDeviation measures vote dispersion
func (a *SimpleScoreAggregator) CalculateStandardDeviation(votes []Vote) float64 {
	if len(votes) == 0 {
		return 0.0
	}

	mean := 0.0
	for _, vote := range votes {
		mean += vote.Score
	}
	mean /= float64(len(votes))

	variance := 0.0
	for _, vote := range votes {
		diff := vote.Score - mean
		variance += diff * diff
	}
	variance /= float64(len(votes))

	return math.Sqrt(variance)
}

// CheckUnanimity determines if all votes agree
func (a *SimpleScoreAggregator) CheckUnanimity(votes []Vote) bool {
	if len(votes) == 0 {
		return false
	}

	firstDecision := votes[0].Decision
	for _, vote := range votes[1:] {
		if vote.Decision != firstDecision {
			return false
		}
	}

	return true
}
```

```go
// consensus/decision.go
package consensus

import (
	"fmt"
)

// ThresholdDecisionEngine implements DecisionEngine interface
type ThresholdDecisionEngine struct{}

// NewThresholdDecisionEngine creates a new decision engine
func NewThresholdDecisionEngine() *ThresholdDecisionEngine {
	return &ThresholdDecisionEngine{}
}

// DetermineOutcome applies threshold logic
func (e *ThresholdDecisionEngine) DetermineOutcome(result ConsensusResult, config VotingConfig) ConsensusOutcome {
	score := result.WeightedScore

	if score >= config.ApprovalThreshold {
		return OutcomeApproved
	} else if score >= config.ConditionalMin {
		return OutcomeConditional
	} else {
		return OutcomeRejected
	}
}

// ApplyTiebreaker handles edge cases
func (e *ThresholdDecisionEngine) ApplyTiebreaker(result ConsensusResult, config VotingConfig) (ConsensusOutcome, error) {
	if result.WeightedScore >= 75.0 && result.WeightedScore < config.ApprovalThreshold {
		switch config.TiebreakerStrategy {
		case TiebreakerWeighted:
			return e.DetermineOutcome(result, config), nil

		case TiebreakerDesignatedVoter:
			for _, vote := range result.Votes {
				if vote.AgentID == config.TiebreakerAgentID {
					if vote.Decision == VoteApprove {
						return OutcomeApproved, nil
					}
				}
			}
			return OutcomeConditional, nil

		case TiebreakerHumanReview:
			return OutcomeConditional, fmt.Errorf("requires human review")

		case TiebreakerRejectByDefault:
			return OutcomeRejected, nil
		}
	}

	return e.DetermineOutcome(result, config), nil
}

// ShouldEscalate determines if human review needed
func (e *ThresholdDecisionEngine) ShouldEscalate(result ConsensusResult, config VotingConfig) bool {
	if config.EscalateConditional && result.Outcome == OutcomeConditional {
		return true
	}

	if result.StandardDev > 20.0 {
		return true
	}

	if result.WeightedScore >= 75.0 && result.WeightedScore < config.ApprovalThreshold {
		return true
	}

	return false
}
```

```go
// consensus/voting.go
package consensus

import (
	"context"
	"fmt"
	"sync"
)

// DefaultVotingManager implements VotingManager interface
type DefaultVotingManager struct {
	storage Storage
	mu      sync.RWMutex
	votes   map[string][]Vote
}

// NewDefaultVotingManager creates a new voting manager
func NewDefaultVotingManager(storage Storage) *DefaultVotingManager {
	return &DefaultVotingManager{
		storage: storage,
		votes:   make(map[string][]Vote),
	}
}

// CollectVotes gathers votes from all agents for a proposal
func (m *DefaultVotingManager) CollectVotes(ctx context.Context, proposalID string, agents []string) ([]Vote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	votes, err := m.storage.GetVotes(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}

	return votes, nil
}

// SubmitVote records a single agent's vote
func (m *DefaultVotingManager) SubmitVote(ctx context.Context, vote Vote) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.storage.SaveVote(ctx, vote); err != nil {
		return fmt.Errorf("failed to save vote: %w", err)
	}

	m.votes[vote.ProposalID] = append(m.votes[vote.ProposalID], vote)

	return nil
}

// GetVotes retrieves all votes for a proposal
func (m *DefaultVotingManager) GetVotes(ctx context.Context, proposalID string) ([]Vote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	votes, err := m.storage.GetVotes(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}

	return votes, nil
}

// CheckQuorum verifies if enough agents voted
func (m *DefaultVotingManager) CheckQuorum(ctx context.Context, proposalID string) (bool, error) {
	votes, err := m.GetVotes(ctx, proposalID)
	if err != nil {
		return false, err
	}

	return len(votes) >= 4, nil
}
```

```go
// consensus/audit.go
package consensus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ChainedAuditLogger implements AuditLogger interface
type ChainedAuditLogger struct {
	storage  Storage
	lastHash string
}

// NewChainedAuditLogger creates a new audit logger
func NewChainedAuditLogger(storage Storage) *ChainedAuditLogger {
	lastHash, _ := storage.GetLastHash(context.Background())
	return &ChainedAuditLogger{
		storage:  storage,
		lastHash: lastHash,
	}
}

// LogVote records a vote to audit trail
func (l *ChainedAuditLogger) LogVote(ctx context.Context, vote Vote) error {
	entry := AuditEntry{
		EntryID:    uuid.New().String(),
		Timestamp:  time.Now(),
		EventType:  "VOTE_CAST",
		ProposalID: vote.ProposalID,
		AgentID:    vote.AgentID,
		Data: map[string]interface{}{
			"decision": vote.Decision,
			"score":    vote.Score,
			"feedback": vote.Feedback,
		},
		PrevHash: l.lastHash,
	}

	entry.Hash = l.computeHash(entry)
	l.lastHash = entry.Hash

	if err := l.storage.AppendAuditEntry(ctx, entry); err != nil {
		return fmt.Errorf("failed to append audit entry: %w", err)
	}

	if err := l.storage.SaveLastHash(ctx, entry.Hash); err != nil {
		return fmt.Errorf("failed to save last hash: %w", err)
	}

	return nil
}

// LogConsensusResult records final outcome
func (l *ChainedAuditLogger) LogConsensusResult(ctx context.Context, result ConsensusResult) error {
	entry := AuditEntry{
		EntryID:    uuid.New().String(),
		Timestamp:  time.Now(),
		EventType:  "CONSENSUS_REACHED",
		ProposalID: result.ProposalID,
		AgentID:    "",
		Data: map[string]interface{}{
			"outcome":         result.Outcome,
			"weighted_score":  result.WeightedScore,
			"passes_threshold": result.PassesThreshold,
		},
		PrevHash: l.lastHash,
	}

	entry.Hash = l.computeHash(entry)
	l.lastHash = entry.Hash

	if err := l.storage.AppendAuditEntry(ctx, entry); err != nil {
		return fmt.Errorf("failed to append audit entry: %w", err)
	}

	if err := l.storage.SaveLastHash(ctx, entry.Hash); err != nil {
		return fmt.Errorf("failed to save last hash: %w", err)
	}

	return nil
}

// GetAuditTrail retrieves full history for a proposal
func (l *ChainedAuditLogger) GetAuditTrail(ctx context.Context, proposalID string) ([]AuditEntry, error) {
	entries, err := l.storage.GetAuditEntries(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit entries: %w", err)
	}

	return entries, nil
}

// VerifyIntegrity checks cryptographic hashes
func (l *ChainedAuditLogger) VerifyIntegrity(ctx context.Context, proposalID string) (bool, error) {
	entries, err := l.GetAuditTrail(ctx, proposalID)
	if err != nil {
		return false, err
	}

	if len(entries) == 0 {
		return true, nil
	}

	for i := 1; i < len(entries); i++ {
		expectedHash := l.computeHash(entries[i])
		if entries[i].Hash != expectedHash {
			return false, fmt.Errorf("hash mismatch at entry %d", i)
		}
		if entries[i].PrevHash != entries[i-1].Hash {
			return false, fmt.Errorf("chain broken at entry %d", i)
		}
	}

	return true, nil
}

func (l *ChainedAuditLogger) computeHash(entry AuditEntry) string {
	dataJSON, _ := json.Marshal(entry.Data)
	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		entry.EntryID,
		entry.Timestamp.Format(time.RFC3339),
		entry.EventType,
		entry.ProposalID,
		string(dataJSON),
		entry.PrevHash,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
```

```go
// consensus/storage.go
package consensus

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage implements Storage interface
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage backend
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS votes (
		vote_id TEXT PRIMARY KEY,
		proposal_id TEXT NOT NULL,
		agent_id TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		decision TEXT NOT NULL,
		score REAL NOT NULL,
		criteria TEXT,
		feedback TEXT,
		weight REAL DEFAULT 1.0,
		UNIQUE(proposal_id, agent_id)
	);

	CREATE TABLE IF NOT EXISTS audit_log (
		entry_id TEXT PRIMARY KEY,
		timestamp DATETIME NOT NULL,
		event_type TEXT NOT NULL,
		proposal_id TEXT NOT NULL,
		agent_id TEXT,
		data TEXT,
		prev_hash TEXT,
		hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS metadata (
		key TEXT PRIMARY KEY,
		value TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_votes_proposal ON votes(proposal_id);
	CREATE INDEX IF NOT EXISTS idx_audit_proposal ON audit_log(proposal_id);
	`

	_, err := db.Exec(schema)
	return err
}

// SaveVote stores a vote
func (s *SQLiteStorage) SaveVote(ctx context.Context, vote Vote) error {
	criteriaJSON, _ := json.Marshal(vote.Criteria)

	query := `
		INSERT OR REPLACE INTO votes (vote_id, proposal_id, agent_id, timestamp, decision, score, criteria, feedback, weight)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		vote.VoteID,
		vote.ProposalID,
		vote.AgentID,
		vote.Timestamp,
		vote.Decision,
		vote.Score,
		string(criteriaJSON),
		vote.Feedback,
		vote.Weight,
	)

	return err
}

// GetVotes retrieves all votes for a proposal
func (s *SQLiteStorage) GetVotes(ctx context.Context, proposalID string) ([]Vote, error) {
	query := `
		SELECT vote_id, proposal_id, agent_id, timestamp, decision, score, criteria, feedback, weight
		FROM votes
		WHERE proposal_id = ?
		ORDER BY timestamp
	`

	rows, err := s.db.QueryContext(ctx, query, proposalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []Vote
	for rows.Next() {
		var vote Vote
		var criteriaJSON string

		err := rows.Scan(
			&vote.VoteID,
			&vote.ProposalID,
			&vote.AgentID,
			&vote.Timestamp,
			&vote.Decision,
			&vote.Score,
			&criteriaJSON,
			&vote.Feedback,
			&vote.Weight,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(criteriaJSON), &vote.Criteria)
		votes = append(votes, vote)
	}

	return votes, nil
}

// AppendAuditEntry adds an audit log entry
func (s *SQLiteStorage) AppendAuditEntry(ctx context.Context, entry AuditEntry) error {
	dataJSON, _ := json.Marshal(entry.Data)

	query := `
		INSERT INTO audit_log (entry_id, timestamp, event_type, proposal_id, agent_id, data, prev_hash, hash)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		entry.EntryID,
		entry.Timestamp,
		entry.EventType,
		entry.ProposalID,
		entry.AgentID,
		string(dataJSON),
		entry.PrevHash,
		entry.Hash,
	)

	return err
}

// GetAuditEntries retrieves audit trail for a proposal
func (s *SQLiteStorage) GetAuditEntries(ctx context.Context, proposalID string) ([]AuditEntry, error) {
	query := `
		SELECT entry_id, timestamp, event_type, proposal_id, agent_id, data, prev_hash, hash
		FROM audit_log
		WHERE proposal_id = ?
		ORDER BY timestamp
	`

	rows, err := s.db.QueryContext(ctx, query, proposalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var entry AuditEntry
		var dataJSON string

		err := rows.Scan(
			&entry.EntryID,
			&entry.Timestamp,
			&entry.EventType,
			&entry.ProposalID,
			&entry.AgentID,
			&dataJSON,
			&entry.PrevHash,
			&entry.Hash,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(dataJSON), &entry.Data)
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetLastHash retrieves the last audit entry hash
func (s *SQLiteStorage) GetLastHash(ctx context.Context) (string, error) {
	query := `SELECT value FROM metadata WHERE key = 'last_hash'`

	var hash string
	err := s.db.QueryRowContext(ctx, query).Scan(&hash)
	if err == sql.ErrNoRows {
		return "", nil
	}

	return hash, err
}

// SaveLastHash stores the last audit entry hash
func (s *SQLiteStorage) SaveLastHash(ctx context.Context, hash string) error {
	query := `INSERT OR REPLACE INTO metadata (key, value) VALUES ('last_hash', ?)`
	_, err := s.db.ExecContext(ctx, query, hash)
	return err
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
```

```go
// consensus/service.go
package consensus

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ConsensusService orchestrates the complete voting workflow
type ConsensusService struct {
	votingManager VotingManager
	aggregator    ScoreAggregator
	decisionEngine DecisionEngine
	auditLogger   AuditLogger
	config        VotingConfig
}

// NewConsensusService creates a new consensus service
func NewConsensusService(storage Storage, config VotingConfig) *ConsensusService {
	return &ConsensusService{
		votingManager:  NewDefaultVotingManager(storage),
		aggregator:     NewSimpleScoreAggregator(),
		decisionEngine: NewThresholdDecisionEngine(),
		auditLogger:    NewChainedAuditLogger(storage),
		config:         config,
	}
}

// StartVotingSession initiates a new voting round
func (s *ConsensusService) StartVotingSession(ctx context.Context, proposalID string) error {
	// Voting session started - could add session tracking here
	return nil
}

// SubmitVote records an agent's vote
func (s *ConsensusService) SubmitVote(ctx context.Context, proposalID, agentID string, decision VoteDecision, score float64, feedback string) error {
	vote := Vote{
		VoteID:     uuid.New().String(),
		ProposalID: proposalID,
		AgentID:    agentID,
		Timestamp:  time.Now(),
		Decision:   decision,
		Score:      score,
		Feedback:   feedback,
		Weight:     s.config.AgentWeights[agentID],
	}

	if vote.Weight == 0 {
		vote.Weight = 1.0
	}

	if err := s.votingManager.SubmitVote(ctx, vote); err != nil {
		return fmt.Errorf("failed to submit vote: %w", err)
	}

	if err := s.auditLogger.LogVote(ctx, vote); err != nil {
		return fmt.Errorf("failed to log vote: %w", err)
	}

	return nil
}

// FinalizeVoting calculates consensus and determines outcome
func (s *ConsensusService) FinalizeVoting(ctx context.Context, proposalID string) (*ConsensusResult, error) {
	votes, err := s.votingManager.GetVotes(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes: %w", err)
	}

	if len(votes) == 0 {
		return nil, fmt.Errorf("no votes found for proposal %s", proposalID)
	}

	result, err := s.aggregator.CalculateConsensus(votes, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate consensus: %w", err)
	}

	outcome := s.decisionEngine.DetermineOutcome(result, s.config)
	result.Outcome = outcome
	result.PassesThreshold = outcome == OutcomeApproved

	if s.decisionEngine.ShouldEscalate(result, s.config) {
		fmt.Printf("ESCALATION NEEDED: Proposal %s requires human review\n", proposalID)
	}

	if err := s.auditLogger.LogConsensusResult(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to log consensus result: %w", err)
	}

	return &result, nil
}

// GetConsensusResult retrieves the voting result
func (s *ConsensusService) GetConsensusResult(ctx context.Context, proposalID string) (*ConsensusResult, error) {
	return s.FinalizeVoting(ctx, proposalID)
}

// VerifyAuditTrail checks integrity of voting records
func (s *ConsensusService) VerifyAuditTrail(ctx context.Context, proposalID string) (bool, error) {
	return s.auditLogger.VerifyIntegrity(ctx, proposalID)
}
```