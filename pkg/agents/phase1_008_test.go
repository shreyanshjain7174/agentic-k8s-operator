// Auto-generated tests for phase1-008

```go
// consensus/types_test.go
package consensus

import (
	"testing"
)

func TestVoteDecision_ToScore(t *testing.T) {
	tests := []struct {
		name     string
		decision VoteDecision
		want     float64
	}{
		{"Approve", VoteApprove, 1.0},
		{"ConditionalApprove", VoteConditionalApprove, 0.75},
		{"ConditionalReject", VoteConditionalReject, 0.25},
		{"Reject", VoteReject, 0.0},
		{"Unknown", VoteDecision("UNKNOWN"), 0.5},
		{"Empty", VoteDecision(""), 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.decision.ToScore()
			if got != tt.want {
				t.Errorf("ToScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

```go
// consensus/aggregator_test.go
package consensus

import (
	"math"
	"testing"
	"time"
)

func TestSimpleScoreAggregator_CalculateConsensus(t *testing.T) {
	agg := NewSimpleScoreAggregator()
	now := time.Now()

	tests := []struct {
		name    string
		votes   []Vote
		config  VotingConfig
		wantErr bool
	}{
		{
			name: "single vote",
			votes: []Vote{
				{ProposalID: "p1", AgentID: "a1", Score: 0.8, Decision: VoteApprove, Timestamp: now},
			},
			config: VotingConfig{AgentWeights: map[string]float64{"a1": 1.0}},
		},
		{
			name: "multiple votes equal weight",
			votes: []Vote{
				{ProposalID: "p1", AgentID: "a1", Score: 1.0, Decision: VoteApprove, Timestamp: now},
				{ProposalID: "p1", AgentID: "a2", Score: 0.75, Decision: VoteConditionalApprove, Timestamp: now},
				{ProposalID: "p1", AgentID: "a3", Score: 0.5, Decision: VoteConditionalReject, Timestamp: now},
			},
			config: VotingConfig{AgentWeights: map[string]float64{}},
		},
		{
			name:    "no votes",
			votes:   []Vote{},
			config:  VotingConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agg.CalculateConsensus(tt.votes, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateConsensus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.ProposalID != tt.votes[0].ProposalID {
					t.Errorf("ProposalID = %v, want %v", result.ProposalID, tt.votes[0].ProposalID)
				}
				if result.TotalVotes != len(tt.votes) {
					t.Errorf("TotalVotes = %v, want %v", result.TotalVotes, len(tt.votes))
				}
			}
		})
	}
}

func TestSimpleScoreAggregator_CalculateWeightedScore(t *testing.T) {
	agg := NewSimpleScoreAggregator()
	now := time.Now()

	tests := []struct {
		name    string
		votes   []Vote
		weights map[string]float64
		want    float64
	}{
		{
			name:    "empty votes",
			votes:   []Vote{},
			weights: map[string]float64{},
			want:    0.0,
		},
		{
			name: "equal weights",
			votes: []Vote{
				{AgentID: "a1", Score: 1.0, Timestamp: now},
				{AgentID: "a2", Score: 0.5, Timestamp: now},
			},
			weights: map[string]float64{"a1": 1.0, "a2": 1.0},
			want:    0.75,
		},
		{
			name: "different weights",
			votes: []Vote{
				{AgentID: "a1", Score: 1.0, Timestamp: now},
				{AgentID: "a2", Score: 0.0, Timestamp: now},
			},
			weights: map[string]float64{"a1": 2.0, "a2": 1.0},
			want:    2.0 / 3.0,
		},
		{
			name: "default weight for missing agent",
			votes: []Vote{
				{AgentID: "a1", Score: 0.8, Timestamp: now},
			},
			weights: map[string]float64{},
			want:    0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agg.CalculateWeightedScore(tt.votes, tt.weights)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("CalculateWeightedScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleScoreAggregator_CalculateStandardDeviation(t *testing.T) {
	agg := NewSimpleScoreAggregator()
	now := time.Now()

	tests := []struct {
		name  string
		votes []Vote
		want  float64
	}{
		{
			name:  "empty votes",
			votes: []Vote{},
			want:  0.0,
		},
		{
			name: "single vote",
			votes: []Vote{
				{Score: 0.5, Timestamp: now},
			},
			want: 0.0,
		},
		{
			name: "identical scores",
			votes: []Vote{
				{Score: 0.8, Timestamp: now},
				{Score: 0.8, Timestamp: now},
				{Score: 0.8, Timestamp: now},
			},
			want: 0.0,
		},
		{
			name: "varied scores",
			votes: []Vote{
				{Score: 0.0, Timestamp: now},
				{Score: 1.0, Timestamp: now},
			},
			want: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agg.CalculateStandardDeviation(tt.votes)
			if math.Abs(got-tt.want) > 0.0001 {
				t.Errorf("CalculateStandardDeviation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleScoreAggregator_CheckUnanimity(t *testing.T) {
	agg := NewSimpleScoreAggregator()
	now := time.Now()

	tests := []struct {
		name  string
		votes []Vote
		want  bool
	}{
		{
			name:  "empty votes",
			votes: []Vote{},
			want:  false,
		},
		{
			name: "single vote",
			votes: []Vote{
				{Decision: VoteApprove, Timestamp: now},
			},
			want: true,
		},
		{
			name: "unanimous approve",
			votes: []Vote{
				{Decision: VoteApprove, Timestamp: now},
				{Decision: VoteApprove, Timestamp: now},
				{Decision: VoteApprove, Timestamp: now},
			},
			want: true,
		},
		{
			name: "mixed decisions",
			votes: []Vote{
				{Decision: VoteApprove, Timestamp: now},
				{Decision: VoteReject, Timestamp: now},
			},
			want: false,
		},
		{
			name: "unanimous reject",
			votes: []Vote{
				{Decision: VoteReject, Timestamp: now},
				{Decision: VoteReject, Timestamp: now},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agg.CheckUnanimity(tt.votes)
			if got != tt.want {
				t.Errorf("CheckUnanimity() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

```go
// consensus/decision_test.go
package consensus

import (
	"testing"
	"time"
)

func TestThresholdDecisionEngine_DetermineOutcome(t *testing.T) {
	engine := NewThresholdDecisionEngine()

	tests := []struct {
		name   string
		result ConsensusResult
		config VotingConfig
		want   ConsensusOutcome
	}{
		{
			name:   "approved above threshold",
			result: ConsensusResult{WeightedScore: 85.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeApproved,
		},
		{
			name:   "approved at threshold",
			result: ConsensusResult{WeightedScore: 80.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeApproved,
		},
		{
			name:   "conditional in range",
			result: ConsensusResult{WeightedScore: 70.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeConditional,
		},
		{
			name:   "conditional at min",
			result: ConsensusResult{WeightedScore: 60.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeConditional,
		},
		{
			name:   "rejected below min",
			result: ConsensusResult{WeightedScore: 50.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeRejected,
		},
		{
			name:   "rejected at zero",
			result: ConsensusResult{WeightedScore: 0.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.DetermineOutcome(tt.result, tt.config)
			if got != tt.want {
				t.Errorf("DetermineOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThresholdDecisionEngine_ApplyTiebreaker(t *testing.T) {
	engine := NewThresholdDecisionEngine()
	now := time.Now()

	tests := []struct {
		name    string
		result  ConsensusResult
		config  VotingConfig
		want    ConsensusOutcome
		wantErr bool
	}{
		{
			name:   "no tiebreaker needed - above threshold",
			result: ConsensusResult{WeightedScore: 85.0},
			config: VotingConfig{ApprovalThreshold: 80.0},
			want:   OutcomeApproved,
		},
		{
			name:   "no tiebreaker needed - below 75",
			result: ConsensusResult{WeightedScore: 70.0},
			config: VotingConfig{ApprovalThreshold: 80.0, ConditionalMin: 60.0},
			want:   OutcomeConditional,
		},
		{
			name:   "weighted tiebreaker",
			result: ConsensusResult{WeightedScore: 77.0},
			config: VotingConfig{
				ApprovalThreshold:  80.0,
				ConditionalMin:     60.0,
				TiebreakerStrategy: TiebreakerWeighted,
			},
			want: OutcomeConditional,
		},
		{
			name: "designated voter approves",
			result: ConsensusResult{
				WeightedScore: 76.0,
				Votes: []Vote{
					{AgentID: "leader", Decision: VoteApprove, Timestamp: now},
					{AgentID: "agent1", Decision: VoteConditionalApprove, Timestamp: now},
				},
			},
			config: VotingConfig{
				ApprovalThreshold:  80.0,
				TiebreakerStrategy: TiebreakerDesignatedVoter,
				TiebreakerAgentID:  "leader",
			},
			want: OutcomeApproved,
		},
		{
			name: "designated voter not approve",
			result: ConsensusResult{
				WeightedScore: 76.0,
				Votes: []Vote{
					{AgentID: "leader", Decision: VoteConditionalApprove, Timestamp: now},
					{AgentID: "agent1", Decision: VoteApprove, Timestamp: now},
				},
			},
			config: VotingConfig{
				ApprovalThreshold:  80.0,
				TiebreakerStrategy: TiebreakerDesignatedVoter,
				TiebreakerAgentID:  "leader",
			},
			want: OutcomeConditional,
		},
		{
			name:   "human review required",
			result: ConsensusResult{WeightedScore: 77.0},
			config: VotingConfig{
				ApprovalThreshold:  80.0,
				TiebreakerStrategy: TiebreakerHumanReview,
			},
			want:    OutcomeConditional,
			wantErr: true,
		},
		{
			name:   "reject by default",
			result: ConsensusResult{WeightedScore: 78.0},
			config: VotingConfig{
				ApprovalThreshold:  80.0,
				TiebreakerStrategy: TiebreakerRejectByDefault,
			},
			want: OutcomeRejected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.ApplyTiebreaker(tt.result, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyTiebreaker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ApplyTiebreaker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThresholdDecisionEngine_ShouldEscalate(t *testing.T) {
	engine := NewThresholdDecisionEngine()

	tests := []struct {
		name   string
		result ConsensusResult
		config VotingConfig
		want   bool
	}{
		{
			name:   "escalate conditional when configured",
			result: ConsensusResult{Outcome: OutcomeConditional},
			config: VotingConfig{EscalateConditional: true},
			want:   true,
		},
		{
			name:   "no escalate conditional when not configured",
			result: ConsensusResult{Outcome: OutcomeConditional},
			config: VotingConfig{EscalateConditional: false},
			want:   false,
		},
		{
			name:   "escalate high standard deviation",
			result: ConsensusResult{StandardDev: 25.0},
			config: VotingConfig{},
			want:   true,
		},
		{
			name:   "no escalate low standard deviation",
			result: ConsensusResult{StandardDev: 15.0},
			config: VotingConfig{},
			want:   false,
		},
		{
			name:   "escalate score in tiebreaker range",
			result: ConsensusResult{WeightedScore: 77.0},
			config: VotingConfig{ApprovalThreshold: 80.0},
			want:   true,
		},
		{
			name:   "no escalate score outside range",
			result: ConsensusResult{WeightedScore: 85.0},
			config: VotingConfig{ApprovalThreshold: 80.0},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.ShouldEscalate(tt.result, tt.config)
			if got != tt.want {
				t.Errorf("ShouldEscalate() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

```go
// consensus/mock_storage_test.go
package consensus

import (
	"context"
	"fmt"
	"sync"
)

type MockStorage struct {
	mu           sync.RWMutex
	votes        map[string][]Vote
	auditEntries map[string][]AuditEntry
	lastHash     string
	saveVoteErr  error
	getVotesErr  error
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		votes:        make(map[string][]Vote),
		auditEntries: make(map[string][]AuditEntry),
	}
}

func (m *MockStorage) SaveVote(ctx context.Context, vote Vote) error {
	if m.saveVoteErr != nil {
		return m.saveVoteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.votes[vote.ProposalID] = append(m.votes[vote.ProposalID], vote)
	return nil
}

func (m *MockStorage) GetVotes(ctx context.Context, proposalID string) ([]Vote, error) {
	if m.getVotesErr != nil {
		return nil, m.getVotesErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.votes[proposalID], nil
}

func (m *MockStorage) AppendAuditEntry(ctx context.Context, entry AuditEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auditEntries[entry.ProposalID] = append(m.auditEntries[entry.ProposalID], entry)
	return nil
}

func (m *MockStorage) GetAuditEntries(ctx context.Context, proposalID string) ([]AuditEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.auditEntries[proposalID], nil
}

func (m *MockStorage) GetLastHash(ctx context.Context) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastHash, nil
}

func (m *MockStorage) SaveLastHash(ctx context.Context, hash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastHash = hash
	return nil
}
```

```go
// consensus/voting_test.go
package consensus

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDefaultVotingManager_SubmitVote(t *testing.T) {
	storage := NewMockStorage()
	manager := NewDefaultVotingManager(storage)
	ctx := context.Background()

	vote := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  time.Now(),
		Decision:   VoteApprove,
		Score:      0.9,
	}

	err := manager.SubmitVote(ctx, vote)
	if err != nil {
		t.Errorf("SubmitVote() error = %v", err)
	}

	votes, err := manager.GetVotes(ctx, "p1")
	if err != nil {
		t.Errorf("GetVotes() error = %v", err)
	}
	if len(votes) != 1 {
		t.Errorf("GetVotes() returned %d votes, want 1", len(votes))
	}
}

func TestDefaultVotingManager_SubmitVoteError(t *testing.T) {
	storage := NewMockStorage()
	storage.saveVoteErr = fmt.Errorf("storage error")
	manager := NewDefaultVotingManager(storage)
	ctx := context.Background()

	vote := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  time.Now(),
	}

	err := manager.SubmitVote(ctx, vote)
	if err == nil {
		t.Error("SubmitVote() expected error, got nil")
	}
}

func TestDefaultVotingManager_CollectVotes(t *testing.T) {
	storage := NewMockStorage()
	manager := NewDefaultVotingManager(storage)
	ctx := context.Background()

	votes := []Vote{
		{VoteID: "v1", ProposalID: "p1", AgentID: "a1", Timestamp: time.Now()},
		{VoteID: "v2", ProposalID: "p1", AgentID: "a2", Timestamp: time.Now()},
	}

	for _, v := range votes {
		manager.SubmitVote(ctx, v)
	}

	collected, err := manager.CollectVotes(ctx, "p1", []string{"a1", "a2"})
	if err != nil {
		t.Errorf("CollectVotes() error = %v", err)
	}
	if len(collected) != 2 {
		t.Errorf("CollectVotes() returned %d votes, want 2", len(collected))
	}
}

func TestDefaultVotingManager_CheckQuorum(t *testing.T) {
	storage := NewMockStorage()
	manager := NewDefaultVotingManager(storage)
	ctx := context.Background()

	tests := []struct {
		name       string
		voteCount  int
		wantQuorum bool
	}{
		{"less than quorum", 3, false},
		{"exactly quorum", 4, true},
		{"more than quorum", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMockStorage()
			manager := NewDefaultVotingManager(storage)

			for i := 0; i < tt.voteCount; i++ {
				vote := Vote{
					VoteID:     fmt.Sprintf("v%d", i),
					ProposalID: "p1",
					AgentID:    fmt.Sprintf("a%d", i),
					Timestamp:  time.Now(),
				}
				manager.SubmitVote(ctx, vote)
			}

			quorum, err := manager.CheckQuorum(ctx, "p1")
			if err != nil {
				t.Errorf("CheckQuorum() error = %v", err)
			}
			if quorum != tt.wantQuorum {
				t.Errorf("CheckQuorum() = %v, want %v", quorum, tt.wantQuorum)
			}
		})
	}
}
```

```go
// consensus/audit_test.go
package consensus

import (
	"context"
	"testing"
	"time"
)

func TestChainedAuditLogger_LogVote(t *testing.T) {
	storage := NewMockStorage()
	logger := NewChainedAuditLogger(storage)
	ctx := context.Background()

	vote := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  time.Now(),
		Decision:   VoteApprove,
		Score:      0.9,
		Feedback:   "looks good",
	}

	err := logger.LogVote(ctx, vote)
	if err != nil {
		t.Errorf("LogVote() error = %v", err)
	}

	entries, err := storage.GetAuditEntries(ctx, "p1")
	if err != nil {
		t.Errorf("GetAuditEntries() error = %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].EventType != "VOTE_CAST" {
		t.Errorf("EventType = %v, want VOTE_CAST", entries[0].EventType)
	}
}

func TestChainedAuditLogger_LogConsensusResult(t *testing.T) {
	storage := NewMockStorage()
	logger := NewChainedAuditLogger(storage)
	ctx := context.Background()

	result := ConsensusResult{
		ProposalID:      "p1",
		WeightedScore:   0.85,
		Outcome:         OutcomeApproved,
		PassesThreshold: true,
	}

	err := logger.LogConsensusResult(ctx, result)
	if err != nil {
		t.Errorf("LogConsensusResult() error = %v", err)
	}

	entries, err := storage.GetAuditEntries(ctx, "p1")
	if err != nil {
		t.Errorf("GetAuditEntries() error = %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].EventType != "CONSENSUS_REACHED" {
		t.Errorf("EventType = %v, want CONSENSUS_REACHED", entries[0].EventType)
	}
}

func TestChainedAuditLogger_VerifyIntegrity(t *testing.T) {
	storage := NewMockStorage()
	logger := NewChainedAuditLogger(storage)
	ctx := context.Background()

	// Log multiple votes to create a chain
	votes := []Vote{
		{VoteID: "v1", ProposalID: "p1", AgentID: "a1", Timestamp: time.Now(), Decision: VoteApprove, Score: 0.9},
		{VoteID: "v2", ProposalID: "p1", AgentID: "a2", Timestamp: time.Now(), Decision: VoteApprove, Score: 0.85},
		{VoteID: "v3", ProposalID: "p1", AgentID: "a3", Timestamp: time.Now(), Decision: VoteConditionalApprove, Score: 0.75},
	}

	for _, vote := range votes {
		err := logger.LogVote(ctx, vote)
		if err != nil {
			t.Errorf("LogVote() error = %v", err)
		}
	}

	valid, err := logger.VerifyIntegrity(ctx, "p1")
	if err != nil {
		t.Errorf("VerifyIntegrity() error = %v", err)
	}
	if !valid {
		t.Error("VerifyIntegrity() = false, want true")
	}
}

func TestChainedAuditLogger_VerifyIntegrityEmpty(t *testing.T) {
	storage := NewMockStorage()
	logger := NewChainedAuditLogger(storage)
	ctx := context.Background()

	valid, err := logger.VerifyIntegrity(ctx, "nonexistent")
	if err != nil {
		t.Errorf("VerifyIntegrity() error = %v", err)
	}
	if !valid {
		t.Error("VerifyIntegrity() for empty chain should be true")
	}
}

func TestChainedAuditLogger_GetAuditTrail(t *testing.T) {
	storage := NewMockStorage()
	logger := NewChainedAuditLogger(storage)
	ctx := context.Background()

	vote := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  time.Now(),
		Decision:   VoteApprove,
		Score:      0.9,
	}

	logger.LogVote(ctx, vote)

	trail, err := logger.GetAuditTrail(ctx, "p1")
	if err != nil {
		t.Errorf("GetAuditTrail() error = %v", err)
	}
	if len(trail) != 1 {
		t.Errorf("GetAuditTrail() returned %d entries, want 1", len(trail))
	}
}
```

```go
// consensus/storage_test.go
package consensus

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSQLiteStorage_SaveAndGetVotes(t *testing.T) {
	dbPath := "test_votes.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	ctx := context.Background()
	now := time.Now()

	vote := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  now,
		Decision:   VoteApprove,
		Score:      0.9,
		Criteria:   map[string]string{"quality": "high"},
		Feedback:   "excellent work",
		Weight:     1.5,
	}

	err = storage.SaveVote(ctx, vote)
	if err != nil {
		t.Errorf("SaveVote() error = %v", err)
	}

	votes, err := storage.GetVotes(ctx, "p1")
	if err != nil {
		t.Errorf("GetVotes() error = %v", err)
	}
	if len(votes) != 1 {
		t.Fatalf("GetVotes() returned %d votes, want 1", len(votes))
	}

	got := votes[0]
	if got.VoteID != vote.VoteID {
		t.Errorf("VoteID = %v, want %v", got.VoteID, vote.VoteID)
	}
	if got.Decision != vote.Decision {
		t.Errorf("Decision = %v, want %v", got.Decision, vote.Decision)
	}
	if got.Score != vote.Score {
		t.Errorf("Score = %v, want %v", got.Score, vote.Score)
	}
}

func TestSQLiteStorage_SaveVoteUniqueConstraint(t *testing.T) {
	dbPath := "test_unique.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	ctx := context.Background()
	now := time.Now()

	vote1 := Vote{
		VoteID:     "v1",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  now,
		Decision:   VoteApprove,
		Score:      0.9,
	}

	vote2 := Vote{
		VoteID:     "v2",
		ProposalID: "p1",
		AgentID:    "a1",
		Timestamp:  now,
		Decision:   VoteReject,
		Score:      0.1,
	}

	storage.SaveVote(ctx, vote1)
	storage.SaveVote(ctx, vote2)

	votes, _ := storage.GetVotes(ctx, "p1")
	if len(votes) != 1 {
		t.Errorf("Expected 1 vote (replaced), got %d", len(votes))
	}
	if votes[0].VoteID != "v2" {
		t.Errorf("Expected vote v2 to replace v1")
	}
}

func TestSQLiteStorage_AuditEntries(t *testing.T) {
	dbPath := "test_audit.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	ctx := context.Background()

	entry := AuditEntry{
		EntryID:    "e1",
		Timestamp:  time.Now(),
		EventType:  "VOTE_CAST",
		ProposalID: "p1",
		AgentID:    "a1",
		Data:       map[string]interface{}{"score": 0.9},
		PrevHash:   "prev123",
		Hash:       "hash456",
	}

	err = storage.AppendAuditEntry(ctx, entry)
	if err != nil {
		t.Errorf("AppendAuditEntry() error = %v", err)
	}

	entries, err := storage.GetAuditEntries(ctx, "p1")
	if err != nil {
		t.Errorf("GetAuditEntries() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("GetAuditEntries() returned %d entries, want 1", len(entries))
	}

	got := entries[0]
	if got.EventType != entry.EventType {
		t.Errorf("EventType = %v, want %v", got.EventType, entry.EventType)
	}
	if got.Hash != entry.Hash {
		t.Errorf("Hash = %v, want %v", got.Hash, entry.Hash)
	}
}

func TestSQLiteStorage_LastHash(t *testing.T) {
	dbPath := "test_hash.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	ctx := context.Background()

	hash, err := storage.GetLastHash(ctx)
	if err != nil {
		t.Errorf("GetLastHash() error = %v", err)
	}
	if hash != "" {
		t.Errorf("GetLastHash() = %v, want empty string", hash)
	}

	err = storage.SaveLastHash(ctx, "newhash123")
	if err != nil {
		t.Errorf("SaveLastHash() error = %v", err)
	}

	hash, err = storage.GetLastHash(ctx)
	if err != nil {
		t.Errorf("GetLastHash() error = %v", err)
	}
	if hash != "newhash123" {
		t.Errorf("GetLastHash() = %v, want newhash123", hash)
	}
}
```

```go
// consensus/service_test.go
package consensus

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestConsensusService_SubmitVote(t *testing.T) {
	dbPath := "test_service.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 80.0,
		ConditionalMin:    60.0,
		AgentWeights: map[string]float64{
			"a1": 1.0,
			"a2": 1.5,
		},
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	err = service.SubmitVote(ctx, "p1", "a1", VoteApprove, 0.9, "good work")
	if err != nil {
		t.Errorf("SubmitVote() error = %v", err)
	}

	votes, _ := storage.GetVotes(ctx, "p1")
	if len(votes) != 1 {
		t.Errorf("Expected 1 vote, got %d", len(votes))
	}
}

func TestConsensusService_FinalizeVoting(t *testing.T) {
	dbPath := "test_finalize.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 80.0,
		ConditionalMin:    60.0,
		AgentWeights:      map[string]float64{},
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	service.SubmitVote(ctx, "p1", "a1", VoteApprove, 1.0, "approve")
	service.SubmitVote(ctx, "p1", "a2", VoteApprove, 0.9, "approve")
	service.SubmitVote(ctx, "p1", "a3", VoteConditionalApprove, 0.75, "conditional")
	service.SubmitVote(ctx, "p1", "a4", VoteApprove, 0.85, "approve")

	result, err := service.FinalizeVoting(ctx, "p1")
	if err != nil {
		t.Errorf("FinalizeVoting() error = %v", err)
	}
	if result == nil {
		t.Fatal("FinalizeVoting() returned nil result")
	}
	if result.ProposalID != "p1" {
		t.Errorf("ProposalID = %v, want p1", result.ProposalID)
	}
	if result.TotalVotes != 4 {
		t.Errorf("TotalVotes = %v, want 4", result.TotalVotes)
	}
	if result.Outcome != OutcomeApproved {
		t.Errorf("Outcome = %v, want %v", result.Outcome, OutcomeApproved)
	}
	if !result.PassesThreshold {
		t.Error("PassesThreshold should be true")
	}
}

func TestConsensusService_FinalizeVotingNoVotes(t *testing.T) {
	dbPath := "test_novotes.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 80.0,
		ConditionalMin:    60.0,
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	_, err = service.FinalizeVoting(ctx, "nonexistent")
	if err == nil {
		t.Error("FinalizeVoting() expected error for no votes, got nil")
	}
}

func TestConsensusService_VerifyAuditTrail(t *testing.T) {
	dbPath := "test_verify.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 80.0,
		ConditionalMin:    60.0,
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	service.SubmitVote(ctx, "p1", "a1", VoteApprove, 0.9, "approve")
	service.SubmitVote(ctx, "p1", "a2", VoteApprove, 0.85, "approve")
	service.FinalizeVoting(ctx, "p1")

	valid, err := service.VerifyAuditTrail(ctx, "p1")
	if err != nil {
		t.Errorf("VerifyAuditTrail() error = %v", err)
	}
	if !valid {
		t.Error("VerifyAuditTrail() = false, want true")
	}
}

func TestConsensusService_StartVotingSession(t *testing.T) {
	dbPath := "test_session.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{}
	service := NewConsensusService(storage, config)
	ctx := context.Background()

	err = service.StartVotingSession(ctx, "p1")
	if err != nil {
		t.Errorf("StartVotingSession() error = %v", err)
	}
}

func TestConsensusService_WeightedVoting(t *testing.T) {
	dbPath := "test_weighted.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 80.0,
		ConditionalMin:    60.0,
		AgentWeights: map[string]float64{
			"architect": 2.0,
			"engineer":  1.0,
			"qa":        1.0,
		},
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	service.SubmitVote(ctx, "p1", "architect", VoteApprove, 1.0, "approved")
	service.SubmitVote(ctx, "p1", "engineer", VoteReject, 0.0, "reject")
	service.SubmitVote(ctx, "p1", "qa", VoteReject, 0.0, "reject")

	result, err := service.FinalizeVoting(ctx, "p1")
	if err != nil {
		t.Errorf("FinalizeVoting() error = %v", err)
	}

	if result.WeightedScore < 0.4 || result.WeightedScore > 0.6 {
		t.Errorf("WeightedScore = %v, expected around 0.5 (2*1.0 / (2+1+1))", result.WeightedScore)
	}
}
```

```go
// consensus/integration_test.go
package consensus

import (
	"context"
	"os"
	"testing"
)

func TestIntegration_FullVotingWorkflow(t *testing.T) {
	dbPath := "test_integration.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold:   85.0,
		ConditionalMin:      70.0,
		QuorumPercent:       75.0,
		EscalateConditional: false,
		AgentWeights: map[string]float64{
			"pm":       1.0,
			"architect": 1.5,
			"engineer":  1.0,
			"qa":        1.0,
			"security":  1.2,
			"business":  0.8,
		},
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	// Start voting session
	err = service.StartVotingSession(ctx, "proposal-001")
	if err != nil {
		t.Fatalf("StartVotingSession() error = %v", err)
	}

	// Submit votes from all agents
	votes := []struct {
		agentID  string
		decision VoteDecision
		score    float64
		feedback string
	}{
		{"pm", VoteApprove, 0.95, "Aligns with roadmap"},
		{"architect", VoteApprove, 0.90, "Sound technical design"},
		{"engineer", VoteConditionalApprove, 0.80, "Implementation feasible with minor concerns"},
		{"qa", VoteApprove, 0.88, "Testable requirements"},
		{"security", VoteApprove, 0.92, "No security concerns"},
		{"business", VoteConditionalApprove, 0.75, "ROI acceptable"},
	}

	for _, v := range votes {
		err = service.SubmitVote(ctx, "proposal-001", v.agentID, v.decision, v.score, v.feedback)
		if err != nil {
			t.Errorf("SubmitVote(%s) error = %v", v.agentID, err)
		}
	}

	// Finalize voting
	result, err := service.FinalizeVoting(ctx, "proposal-001")
	if err != nil {
		t.Fatalf("FinalizeVoting() error = %v", err)
	}

	// Verify results
	if result.TotalVotes != 6 {
		t.Errorf("TotalVotes = %v, want 6", result.TotalVotes)
	}

	if result.Outcome != OutcomeApproved {
		t.Errorf("Outcome = %v, want %v (weighted score: %v)", result.Outcome, OutcomeApproved, result.WeightedScore)
	}

	if !result.PassesThreshold {
		t.Error("PassesThreshold should be true")
	}

	// Verify audit trail
	valid, err := service.VerifyAuditTrail(ctx, "proposal-001")
	if err != nil {
		t.Errorf("VerifyAuditTrail() error = %v", err)
	}
	if !valid {
		t.Error("VerifyAuditTrail() = false, want true")
	}

	// Verify audit entries exist
	auditEntries, err := storage.GetAuditEntries(ctx, "proposal-001")
	if err != nil {
		t.Errorf("GetAuditEntries() error = %v", err)
	}
	// Should have 6 vote entries + 1 consensus result entry
	if len(auditEntries) != 7 {
		t.Errorf("Expected 7 audit entries, got %d", len(auditEntries))
	}
}

func TestIntegration_RejectedProposal(t *testing.T) {
	dbPath := "test_rejected.db"
	defer os.Remove(dbPath)

	storage, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStorage() error = %v", err)
	}
	defer storage.Close()

	config := VotingConfig{
		ApprovalThreshold: 85.0,
		ConditionalMin:    70.0,
		AgentWeights:      map[string]float64{},
	}

	service := NewConsensusService(storage, config)
	ctx := context.Background()

	votes := []struct {
		agentID  string
		decision VoteDecision
		score    float64
	}{
		{"a1", VoteReject, 0.2},
		{"a2", VoteReject, 0.1},
		{"a3", VoteConditionalReject, 0.3},
		{"a4", VoteReject, 0.0},
	}

	for _, v := range votes {
		service.SubmitVote(ctx, "p-bad", v.agentID, v.decision, v.score, "rejected")
	}

	result, err := service.FinalizeVoting(ctx, "p-bad")
	if err != nil {
		t.Fatalf("FinalizeVoting() error = %v", err)
	}

	if result.Outcome != OutcomeRejected {
		t.Errorf("Outcome = %v, want %v", result.Outcome, OutcomeRejected)
	}

	if result.PassesThreshold {
		t.Error("PassesThreshold should be false")
	}
}
```