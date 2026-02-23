// Auto-generated tests for phase1-011

```go
// pkg/execution/validators/consensus_test.go
package validators

import (
	"context"
	"testing"

	"github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

func TestConsensusValidator_Validate(t *testing.T) {
	tests := []struct {
		name        string
		minScore    float64
		votingScore float64
		wantErr     bool
	}{
		{
			name:        "sufficient consensus",
			minScore:    0.80,
			votingScore: 0.85,
			wantErr:     false,
		},
		{
			name:        "exact minimum consensus",
			minScore:    0.80,
			votingScore: 0.80,
			wantErr:     false,
		},
		{
			name:        "insufficient consensus",
			minScore:    0.80,
			votingScore: 0.75,
			wantErr:     true,
		},
		{
			name:        "zero voting score",
			minScore:    0.80,
			votingScore: 0.00,
			wantErr:     true,
		},
		{
			name:        "perfect consensus",
			minScore:    0.80,
			votingScore: 1.00,
			wantErr:     false,
		},
		{
			name:        "default min score",
			minScore:    0.00,
			votingScore: 0.85,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewConsensusValidator(tt.minScore)
			proposal := &execution.Proposal{
				VotingScore: tt.votingScore,
			}

			err := v.Validate(context.Background(), proposal)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// pkg/execution/validators/schema_test.go
package validators

import (
	"context"
	"testing"
	"time"

	"github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

func TestSchemaValidator_Validate(t *testing.T) {
	validProposal := &execution.Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: execution.TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		Changes:     map[string]interface{}{"key": "value"},
		VotingScore: 0.85,
		Metadata: execution.ProposalMetadata{
			Author:    "test-agent",
			CreatedAt: time.Now(),
			RiskLevel: "low",
		},
	}

	tests := []struct {
		name     string
		proposal *execution.Proposal
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid proposal",
			proposal: validProposal,
			wantErr:  false,
		},
		{
			name: "missing ID",
			proposal: &execution.Proposal{
				OperatorType: "configmap-operator",
				Action:       "create",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: true,
			errMsg:  "proposal ID is required",
		},
		{
			name: "missing operator type",
			proposal: &execution.Proposal{
				ID:     "prop-123",
				Action: "create",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: true,
			errMsg:  "operator type is required",
		},
		{
			name: "missing action",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: true,
			errMsg:  "action is required",
		},
		{
			name: "invalid action",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "invalid",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: true,
			errMsg:  "invalid action",
		},
		{
			name: "missing target kind",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "create",
				Target: execution.TargetResource{
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: true,
			errMsg:  "target kind is required",
		},
		{
			name: "missing target name",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "create",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
				},
			},
			wantErr: true,
			errMsg:  "target name is required",
		},
		{
			name: "missing target namespace",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "create",
				Target: execution.TargetResource{
					Kind: "ConfigMap",
					Name: "test",
				},
			},
			wantErr: true,
			errMsg:  "target namespace is required",
		},
		{
			name: "update action",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "update",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: false,
		},
		{
			name: "delete action",
			proposal: &execution.Proposal{
				ID:           "prop-123",
				OperatorType: "configmap-operator",
				Action:       "delete",
				Target: execution.TargetResource{
					Kind:      "ConfigMap",
					Namespace: "default",
					Name:      "test",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewSchemaValidator()
			err := v.Validate(context.Background(), tt.proposal)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}

func TestIsValidAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   bool
	}{
		{"create is valid", "create", true},
		{"update is valid", "update", true},
		{"delete is valid", "delete", true},
		{"invalid action", "invalid", false},
		{"empty action", "", false},
		{"uppercase action", "CREATE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidAction(tt.action); got != tt.want {
				t.Errorf("isValidAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

// pkg/execution/validators/ratelimit_test.go
package validators

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter()

	// Test initial allows
	for i := 0; i < 10; i++ {
		allowed, err := rl.Allow("test-operator")
		if !allowed || err != nil {
			t.Errorf("Allow() iteration %d: allowed = %v, err = %v, want allowed = true, err = nil", i, allowed, err)
		}
	}

	// Test rate limit exceeded
	allowed, err := rl.Allow("test-operator")
	if allowed || err == nil {
		t.Errorf("Allow() after exhaustion: allowed = %v, err = %v, want allowed = false, err != nil", allowed, err)
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl := NewRateLimiter()

	// Exhaust tokens
	for i := 0; i < 10; i++ {
		rl.Allow("test-operator")
	}

	// Verify exhausted
	allowed, _ := rl.Allow("test-operator")
	if allowed {
		t.Error("Allow() should be false after exhaustion")
	}

	// Reset
	err := rl.Reset("test-operator")
	if err != nil {
		t.Errorf("Reset() error = %v, want nil", err)
	}

	// Should allow again
	allowed, err = rl.Allow("test-operator")
	if !allowed || err != nil {
		t.Errorf("Allow() after reset: allowed = %v, err = %v, want allowed = true, err = nil", allowed, err)
	}
}

func TestRateLimiter_Reset_NonExistent(t *testing.T) {
	rl := NewRateLimiter()

	err := rl.Reset("non-existent")
	if err == nil {
		t.Error("Reset() for non-existent operator should return error")
	}
}

func TestRateLimiter_MultipleOperators(t *testing.T) {
	rl := NewRateLimiter()

	// Exhaust operator1
	for i := 0; i < 10; i++ {
		rl.Allow("operator1")
	}

	// operator1 should be blocked
	allowed, _ := rl.Allow("operator1")
	if allowed {
		t.Error("operator1 should be rate limited")
	}

	// operator2 should still be allowed
	allowed, err := rl.Allow("operator2")
	if !allowed || err != nil {
		t.Errorf("operator2 should be allowed: allowed = %v, err = %v", allowed, err)
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	bucket := &TokenBucket{
		tokens:     5,
		maxTokens:  10,
		refillRate: 100 * time.Millisecond,
		lastRefill: time.Now().Add(-200 * time.Millisecond),
	}

	bucket.refill()

	if bucket.tokens < 7 {
		t.Errorf("tokens after refill = %d, want >= 7", bucket.tokens)
	}
}

func TestTokenBucket_Consume(t *testing.T) {
	bucket := &TokenBucket{
		tokens:     3,
		maxTokens:  10,
		refillRate: 60 * time.Second,
		lastRefill: time.Now(),
	}

	// First consume should succeed
	allowed, err := bucket.consume()
	if !allowed || err != nil {
		t.Errorf("consume() = %v, %v, want true, nil", allowed, err)
	}

	if bucket.tokens != 2 {
		t.Errorf("tokens after consume = %d, want 2", bucket.tokens)
	}

	// Exhaust tokens
	bucket.consume()
	bucket.consume()

	// Should fail now
	allowed, err = bucket.consume()
	if allowed || err == nil {
		t.Errorf("consume() after exhaustion = %v, %v, want false, error", allowed, err)
	}
}

// pkg/execution/audit/logger_test.go
package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

func TestAuditLogger_LogExecution(t *testing.T) {
	tmpDir := t.TempDir()
	logger := NewAuditLogger(tmpDir)

	exec := &execution.ExecutionResult{
		ExecutionID: "exec-123",
		ProposalID:  "prop-456",
		Status:      execution.StatusSuccess,
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	err := logger.LogExecution(exec)
	if err != nil {
		t.Errorf("LogExecution() error = %v", err)
	}

	// Verify file exists
	date := time.Now().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, "executions", date+".jsonl")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file not created: %s", logFile)
	}
}

func TestAuditLogger_Query(t *testing.T) {
	tmpDir := t.TempDir()
	logger := NewAuditLogger(tmpDir)

	// Log multiple executions
	exec1 := &execution.ExecutionResult{
		ExecutionID: "exec-1",
		ProposalID:  "prop-1",
		Status:      execution.StatusSuccess,
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
	}
	exec2 := &execution.ExecutionResult{
		ExecutionID: "exec-2",
		ProposalID:  "prop-2",
		Status:      execution.StatusFailed,
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	logger.LogExecution(exec1)
	logger.LogExecution(exec2)

	tests := []struct {
		name       string
		filter     execution.AuditFilter
		wantCount  int
		wantStatus execution.ExecutionStatusType
	}{
		{
			name:      "no filter",
			filter:    execution.AuditFilter{},
			wantCount: 2,
		},
		{
			name: "filter by proposal ID",
			filter: execution.AuditFilter{
				ProposalID: "prop-1",
			},
			wantCount: 1,
		},
		{
			name: "filter by status",
			filter: execution.AuditFilter{
				Status: execution.StatusFailed,
			},
			wantCount:  1,
			wantStatus: execution.StatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := logger.Query(tt.filter)
			if err != nil {
				t.Errorf("Query() error = %v", err)
				return
			}

			if len(results) != tt.wantCount {
				t.Errorf("Query() returned %d results, want %d", len(results), tt.wantCount)
			}

			if tt.wantStatus != "" && len(results) > 0 {
				if results[0].Status != tt.wantStatus {
					t.Errorf("Query() result status = %v, want %v", results[0].Status, tt.wantStatus)
				}
			}
		})
	}
}

func TestAuditLogger_Query_DateRange(t *testing.T) {
	tmpDir := t.TempDir()
	logger := NewAuditLogger(tmpDir)

	now := time.Now()
	past := now.Add(-48 * time.Hour)

	exec := &execution.ExecutionResult{
		ExecutionID: "exec-1",
		ProposalID:  "prop-1",
		Status:      execution.StatusSuccess,
		StartedAt:   now,
		CompletedAt: now,
	}

	logger.LogExecution(exec)

	tests := []struct {
		name      string
		filter    execution.AuditFilter
		wantCount int
	}{
		{
			name: "within date range",
			filter: execution.AuditFilter{
				StartDate: past,
				EndDate:   now.Add(1 * time.Hour),
			},
			wantCount: 1,
		},
		{
			name: "outside date range",
			filter: execution.AuditFilter{
				StartDate: now.Add(1 * time.Hour),
				EndDate:   now.Add(2 * time.Hour),
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := logger.Query(tt.filter)
			if err != nil {
				t.Errorf("Query() error = %v", err)
				return
			}

			if len(results) != tt.wantCount {
				t.Errorf("Query() returned %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

// pkg/execution/engine_test.go
package execution

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestExecutionEngine_Execute_Success(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		Changes: map[string]interface{}{
			"key1": "value1",
		},
		VotingScore: 0.85,
		Metadata: ProposalMetadata{
			Author:    "test-agent",
			CreatedAt: time.Now(),
			RiskLevel: "low",
		},
	}

	result, err := engine.Execute(context.Background(), proposal, ExecutionOptions{
		SkipHealthCheck: true,
	})

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Status != StatusSuccess {
		t.Errorf("Execute() status = %v, want %v", result.Status, StatusSuccess)
	}

	if result.ProposalID != proposal.ID {
		t.Errorf("Execute() proposalID = %v, want %v", result.ProposalID, proposal.ID)
	}

	if len(result.AppliedChanges) == 0 {
		t.Error("Execute() should have applied changes")
	}
}

func TestExecutionEngine_Execute_ValidationFailed(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		VotingScore: 0.70, // Below threshold
	}

	result, err := engine.Execute(context.Background(), proposal, ExecutionOptions{})

	if err == nil {
		t.Fatal("Execute() should return error for low consensus")
	}

	if result.Status != StatusRejected {
		t.Errorf("Execute() status = %v, want %v", result.Status, StatusRejected)
	}
}

func TestExecutionEngine_Execute_MissingFields(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "", // Missing ID
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		VotingScore: 0.85,
	}

	result, err := engine.Execute(context.Background(), proposal, ExecutionOptions{})

	if err == nil {
		t.Fatal("Execute() should return error for missing ID")
	}

	if result.Status != StatusRejected {
		t.Errorf("Execute() status = %v, want %v", result.Status, StatusRejected)
	}
}

func TestExecutionEngine_DryRun(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "update",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		VotingScore: 0.85,
	}

	result, err := engine.DryRun(context.Background(), proposal)

	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}

	if !result.WouldApply {
		t.Error("DryRun() WouldApply should be true")
	}

	if result.PredictedDiff == "" {
		t.Error("DryRun() should have predicted diff")
	}

	if len(result.ValidationErrors) > 0 {
		t.Errorf("DryRun() should have no validation errors, got %v", result.ValidationErrors)
	}
}

func TestExecutionEngine_DryRun_ValidationErrors(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		VotingScore: 0.70,
	}

	result, err := engine.DryRun(context.Background(), proposal)

	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}

	if result.WouldApply {
		t.Error("DryRun() WouldApply should be false")
	}

	if len(result.ValidationErrors) == 0 {
		t.Error("DryRun() should have validation errors")
	}
}

func TestExecutionEngine_GetStatus(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	// Start an execution
	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		Changes:     map[string]interface{}{"key": "value"},
		VotingScore: 0.85,
	}

	result, _ := engine.Execute(context.Background(), proposal, ExecutionOptions{
		SkipHealthCheck: true,
	})

	// Get status
	status, err := engine.GetStatus(result.ExecutionID)
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	if status.ExecutionID != result.ExecutionID {
		t.Errorf("GetStatus() executionID = %v, want %v", status.ExecutionID, result.ExecutionID)
	}

	if status.Status != StatusSuccess {
		t.Errorf("GetStatus() status = %v, want %v", status.Status, StatusSuccess)
	}
}

func TestExecutionEngine_GetStatus_NotFound(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	_, err := engine.GetStatus("non-existent")
	if err == nil {
		t.Error("GetStatus() should return error for non-existent execution")
	}
}

func TestExecutionEngine_Execute_UpdateConfigMap(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()

	// Pre-create ConfigMap
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"key1": "old-value",
		},
	}
	k8sClient.CoreV1().ConfigMaps("default").Create(context.Background(), cm, metav1.CreateOptions{})

	tmpDir := t.TempDir()
	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "update",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		Changes: map[string]interface{}{
			"key1": "new-value",
			"key2": "added-value",
		},
		VotingScore: 0.85,
	}

	result, err := engine.Execute(context.Background(), proposal, ExecutionOptions{
		SkipHealthCheck: true,
	})

	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.Status != StatusSuccess {
		t.Errorf("Execute() status = %v, want %v", result.Status, StatusSuccess)
	}

	// Verify ConfigMap was updated
	updatedCM, _ := k8sClient.CoreV1().ConfigMaps("default").Get(context.Background(), "test-config", metav1.GetOptions{})
	if updatedCM.Data["key1"] != "new-value" {
		t.Errorf("ConfigMap key1 = %v, want new-value", updatedCM.Data["key1"])
	}
	if updatedCM.Data["key2"] != "added-value" {
		t.Errorf("ConfigMap key2 = %v, want added-value", updatedCM.Data["key2"])
	}
}

func TestExecutionEngine_Execute_RateLimitExceeded(t *testing.T) {
	k8sClient := fake.NewSimpleClientset()
	tmpDir := t.TempDir()

	engine := NewExecutionEngine(EngineOptions{
		K8sClient:    k8sClient,
		AuditPath:    tmpDir + "/audit",
		MinConsensus: 0.80,
		SnapshotDir:  tmpDir + "/snapshots",
	})

	proposal := &Proposal{
		ID:           "prop-123",
		OperatorType: "configmap-operator",
		Action:       "create",
		Target: TargetResource{
			Kind:      "ConfigMap",
			Namespace: "default",
			Name:      "test-config",
		},
		Changes:     map[string]interface{}{"key": "value"},
		VotingScore: 0.85,
	}

	// Exhaust rate limit
	for i := 0; i < 10; i++ {
		engine.Execute(context.Background(), proposal, ExecutionOptions{SkipHealthCheck: true})
	}

	// Next execution should fail
	result, err := engine.Execute(context.Background(), proposal, ExecutionOptions{})

	if err == nil {
		t.Fatal("Execute() should fail due to rate limit")
	}

	if result.Status != StatusRejected {
		t.Errorf("Execute() status = %v, want %v", result.Status, StatusRejected)
	}
}

func TestAllHealthChecksPassed(t *testing.T) {
	tests := []struct {
		name    string
		results []HealthCheckResult
		want    bool
	}{
		{
			name:    "empty results",
			results: []HealthCheckResult{},
			want:    true,
		},
		{
			name: "all passed",
			results: []HealthCheckResult{
				{CheckName: "check1", Passed: true},
				{CheckName: "check2", Passed: true},
			},
			want: true,
		},
		{
			name: "one failed",
			results: []HealthCheckResult{
				{CheckName: "check1", Passed: true},
				{CheckName: "check2", Passed: false},
			},
			want: false,
		},
		{
			name: "all failed",
			results: []HealthCheckResult{
				{CheckName: "check1", Passed: false},
				{CheckName: "check2", Passed: false},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allHealthChecksPassed(tt.results); got != tt.want {
				t.Errorf("allHealthChecksPassed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateExecutionID(t *testing.T) {
	id1 := generateExecutionID()
	id2 := generateExecutionID()

	if id1 == id2 {
		t.Error("generateExecutionID() should generate unique IDs")
	}

	if len(id1) == 0 {
		t.Error("generateExecutionID() should not return empty string")
	}

	if id1[:5] != "exec-" {
		t.Errorf("generateExecutionID() should start with 'exec-', got %s", id1)
	}
}
```