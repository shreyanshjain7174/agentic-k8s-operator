// Auto-generated for phase1-011
// Task: Implement proposal execution engine for applying operator changes to cluster

```go
// pkg/execution/types.go
package execution

import (
    "context"
    "time"
)

// ExecutionEngine orchestrates the full execution lifecycle
type ExecutionEngine interface {
    Execute(ctx context.Context, proposal *Proposal, opts ExecutionOptions) (*ExecutionResult, error)
    DryRun(ctx context.Context, proposal *Proposal) (*DryRunResult, error)
    Rollback(ctx context.Context, executionID string) error
    GetStatus(executionID string) (*ExecutionStatus, error)
}

// Proposal represents an operator change request
type Proposal struct {
    ID           string                 `json:"id"`
    OperatorType string                 `json:"operator_type"`
    Action       string                 `json:"action"`
    Target       TargetResource         `json:"target"`
    Changes      map[string]interface{} `json:"changes"`
    VotingScore  float64                `json:"voting_score"`
    Metadata     ProposalMetadata       `json:"metadata"`
}

// TargetResource identifies the Kubernetes resource
type TargetResource struct {
    Kind      string `json:"kind"`
    Namespace string `json:"namespace"`
    Name      string `json:"name"`
}

// ProposalMetadata contains proposal metadata
type ProposalMetadata struct {
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    RiskLevel string    `json:"risk_level"`
}

// ExecutionOptions configures execution behavior
type ExecutionOptions struct {
    DryRun          bool
    SkipHealthCheck bool
    RollbackTimeout time.Duration
}

// ExecutionResult contains execution outcome
type ExecutionResult struct {
    ExecutionID        string              `json:"execution_id"`
    ProposalID         string              `json:"proposal_id"`
    Status             ExecutionStatusType `json:"status"`
    AppliedChanges     []Change            `json:"applied_changes"`
    RollbackID         string              `json:"rollback_id,omitempty"`
    Error              string              `json:"error,omitempty"`
    StartedAt          time.Time           `json:"started_at"`
    CompletedAt        time.Time           `json:"completed_at"`
    HealthCheckResults []HealthCheckResult `json:"health_checks"`
}

// ExecutionStatusType represents execution state
type ExecutionStatusType string

const (
    StatusPending    ExecutionStatusType = "pending"
    StatusValidating ExecutionStatusType = "validating"
    StatusExecuting  ExecutionStatusType = "executing"
    StatusVerifying  ExecutionStatusType = "verifying"
    StatusSuccess    ExecutionStatusType = "success"
    StatusFailed     ExecutionStatusType = "failed"
    StatusRolledBack ExecutionStatusType = "rolled_back"
    StatusRejected   ExecutionStatusType = "rejected"
)

// ExecutionStatus provides current execution status
type ExecutionStatus struct {
    ExecutionID string              `json:"execution_id"`
    Status      ExecutionStatusType `json:"status"`
    UpdatedAt   time.Time           `json:"updated_at"`
}

// Validator validates proposals before execution
type Validator interface {
    Validate(ctx context.Context, proposal *Proposal) error
}

// RateLimiter controls execution rate
type RateLimiter interface {
    Allow(operatorType string) (bool, error)
    Reset(operatorType string) error
}

// StateSnapshot manages resource snapshots
type StateSnapshot interface {
    Create(ctx context.Context, proposal *Proposal) (*Snapshot, error)
    Restore(ctx context.Context, snapshotID string) error
    Delete(ctx context.Context, snapshotID string) error
    Get(snapshotID string) (*Snapshot, error)
}

// Snapshot contains resource state
type Snapshot struct {
    ID             string             `json:"id"`
    ProposalID     string             `json:"proposal_id"`
    CreatedAt      time.Time          `json:"created_at"`
    Resources      []ResourceSnapshot `json:"resources"`
    EtcdBackupPath string             `json:"etcd_backup_path,omitempty"`
}

// ResourceSnapshot captures resource state
type ResourceSnapshot struct {
    Kind      string                 `json:"kind"`
    Namespace string                 `json:"namespace"`
    Name      string                 `json:"name"`
    State     map[string]interface{} `json:"state"`
    Version   string                 `json:"version"`
}

// OperatorExecutor executes operator changes
type OperatorExecutor interface {
    Execute(ctx context.Context, proposal *Proposal, snapshot *Snapshot) ([]Change, error)
    SupportedOperators() []string
}

// Change represents an applied modification
type Change struct {
    ResourceKind string    `json:"resource_kind"`
    ResourceName string    `json:"resource_name"`
    Namespace    string    `json:"namespace"`
    Action       string    `json:"action"`
    Diff         string    `json:"diff"`
    AppliedAt    time.Time `json:"applied_at"`
}

// HealthChecker verifies resource health
type HealthChecker interface {
    Check(ctx context.Context, proposal *Proposal, changes []Change) ([]HealthCheckResult, error)
}

// HealthCheckResult contains check outcome
type HealthCheckResult struct {
    CheckName string    `json:"check_name"`
    Passed    bool      `json:"passed"`
    Message   string    `json:"message"`
    CheckedAt time.Time `json:"checked_at"`
}

// RollbackManager handles rollbacks
type RollbackManager interface {
    Rollback(ctx context.Context, executionID string) error
    AutoRollback(ctx context.Context, execution *ExecutionResult) error
}

// AuditLogger logs execution events
type AuditLogger interface {
    LogExecution(execution *ExecutionResult) error
    Query(filter AuditFilter) ([]*ExecutionResult, error)
}

// AuditFilter filters audit logs
type AuditFilter struct {
    ProposalID   string
    OperatorType string
    Status       ExecutionStatusType
    StartDate    time.Time
    EndDate      time.Time
}

// DryRunResult contains dry-run outcome
type DryRunResult struct {
    ProposalID       string   `json:"proposal_id"`
    WouldApply       bool     `json:"would_apply"`
    PredictedDiff    string   `json:"predicted_diff"`
    ValidationErrors []string `json:"validation_errors,omitempty"`
}

// EngineOptions configures the execution engine
type EngineOptions struct {
    K8sClient      interface{}
    EtcdClient     interface{}
    AuditPath      string
    MinConsensus   float64
    SnapshotDir    string
}
```

```go
// pkg/execution/engine.go
package execution

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/google/uuid"
    "k8s.io/client-go/kubernetes"
)

// executionEngine implements ExecutionEngine
type executionEngine struct {
    validators  []Validator
    rateLimiter RateLimiter
    snapshot    StateSnapshot
    executor    OperatorExecutor
    health      HealthChecker
    rollback    RollbackManager
    audit       AuditLogger
    k8sClient   *kubernetes.Clientset

    statusMu sync.RWMutex
    statuses map[string]*ExecutionStatus
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(opts EngineOptions) ExecutionEngine {
    k8sClient := opts.K8sClient.(*kubernetes.Clientset)
    
    snapshotMgr := NewSnapshotManager(SnapshotOptions{
        K8sClient:   k8sClient,
        SnapshotDir: opts.SnapshotDir,
    })
    
    auditLogger := NewAuditLogger(opts.AuditPath)
    
    return &executionEngine{
        validators: []Validator{
            NewConsensusValidator(opts.MinConsensus),
            NewSchemaValidator(),
        },
        rateLimiter: NewRateLimiter(),
        snapshot:    snapshotMgr,
        executor:    NewOperatorExecutor(k8sClient),
        health:      NewHealthChecker(k8sClient),
        rollback:    NewRollbackManager(snapshotMgr, k8sClient, auditLogger),
        audit:       auditLogger,
        k8sClient:   k8sClient,
        statuses:    make(map[string]*ExecutionStatus),
    }
}

// Execute runs the full execution pipeline
func (e *executionEngine) Execute(ctx context.Context, proposal *Proposal, opts ExecutionOptions) (*ExecutionResult, error) {
    result := &ExecutionResult{
        ExecutionID: generateExecutionID(),
        ProposalID:  proposal.ID,
        Status:      StatusPending,
        StartedAt:   time.Now(),
    }

    e.updateStatus(result.ExecutionID, StatusPending)

    // Phase 1: Validation
    result.Status = StatusValidating
    e.updateStatus(result.ExecutionID, StatusValidating)
    
    if err := e.runValidation(ctx, proposal); err != nil {
        result.Status = StatusRejected
        result.Error = err.Error()
        result.CompletedAt = time.Now()
        e.updateStatus(result.ExecutionID, StatusRejected)
        e.audit.LogExecution(result)
        return result, err
    }

    // Phase 2: Rate limiting
    if allowed, err := e.rateLimiter.Allow(proposal.OperatorType); !allowed {
        result.Status = StatusRejected
        result.Error = fmt.Sprintf("rate limit exceeded: %v", err)
        result.CompletedAt = time.Now()
        e.updateStatus(result.ExecutionID, StatusRejected)
        e.audit.LogExecution(result)
        return result, err
    }

    // Phase 3: Create snapshot
    snapshot, err := e.snapshot.Create(ctx, proposal)
    if err != nil {
        result.Status = StatusFailed
        result.Error = fmt.Sprintf("snapshot failed: %v", err)
        result.CompletedAt = time.Now()
        e.updateStatus(result.ExecutionID, StatusFailed)
        e.audit.LogExecution(result)
        return result, err
    }
    result.RollbackID = snapshot.ID

    // Phase 4: Execute
    result.Status = StatusExecuting
    e.updateStatus(result.ExecutionID, StatusExecuting)
    
    changes, err := e.executor.Execute(ctx, proposal, snapshot)
    if err != nil {
        result.Status = StatusFailed
        result.Error = err.Error()
        result.CompletedAt = time.Now()
        e.updateStatus(result.ExecutionID, StatusFailed)
        e.rollback.AutoRollback(ctx, result)
        e.audit.LogExecution(result)
        return result, err
    }
    result.AppliedChanges = changes

    // Phase 5: Health checks
    result.Status = StatusVerifying
    e.updateStatus(result.ExecutionID, StatusVerifying)
    
    if !opts.SkipHealthCheck {
        timeout := opts.RollbackTimeout
        if timeout == 0 {
            timeout = 5 * time.Minute
        }

        healthResults, err := e.runHealthChecks(ctx, proposal, changes, timeout)
        result.HealthCheckResults = healthResults

        if err != nil || !allHealthChecksPassed(healthResults) {
            result.Status = StatusFailed
            result.Error = "health checks failed"
            result.CompletedAt = time.Now()
            e.updateStatus(result.ExecutionID, StatusFailed)
            e.rollback.AutoRollback(ctx, result)
            e.audit.LogExecution(result)
            return result, fmt.Errorf("health checks failed, rolled back")
        }
    }

    // Success
    result.Status = StatusSuccess
    result.CompletedAt = time.Now()
    e.updateStatus(result.ExecutionID, StatusSuccess)
    e.audit.LogExecution(result)

    return result, nil
}

// DryRun simulates execution without applying changes
func (e *executionEngine) DryRun(ctx context.Context, proposal *Proposal) (*DryRunResult, error) {
    result := &DryRunResult{
        ProposalID:       proposal.ID,
        WouldApply:       true,
        ValidationErrors: []string{},
    }

    // Run validation
    for _, validator := range e.validators {
        if err := validator.Validate(ctx, proposal); err != nil {
            result.WouldApply = false
            result.ValidationErrors = append(result.ValidationErrors, err.Error())
        }
    }

    // Check rate limit
    if allowed, err := e.rateLimiter.Allow(proposal.OperatorType); !allowed {
        result.WouldApply = false
        result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("rate limit: %v", err))
    }

    // Generate predicted diff (simplified)
    result.PredictedDiff = fmt.Sprintf("Would %s %s/%s in namespace %s",
        proposal.Action, proposal.Target.Kind, proposal.Target.Name, proposal.Target.Namespace)

    return result, nil
}

// Rollback reverses an execution
func (e *executionEngine) Rollback(ctx context.Context, executionID string) error {
    return e.rollback.Rollback(ctx, executionID)
}

// GetStatus retrieves execution status
func (e *executionEngine) GetStatus(executionID string) (*ExecutionStatus, error) {
    e.statusMu.RLock()
    defer e.statusMu.RUnlock()

    status, ok := e.statuses[executionID]
    if !ok {
        return nil, fmt.Errorf("execution not found: %s", executionID)
    }

    return status, nil
}

// runValidation runs all validators
func (e *executionEngine) runValidation(ctx context.Context, proposal *Proposal) error {
    for _, validator := range e.validators {
        if err := validator.Validate(ctx, proposal); err != nil {
            return err
        }
    }
    return nil
}

// runHealthChecks monitors health with timeout
func (e *executionEngine) runHealthChecks(ctx context.Context, proposal *Proposal, changes []Change, timeout time.Duration) ([]HealthCheckResult, error) {
    deadline := time.Now().Add(timeout)
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for time.Now().Before(deadline) {
        results, err := e.health.Check(ctx, proposal, changes)
        if err == nil && allHealthChecksPassed(results) {
            return results, nil
        }

        select {
        case <-ticker.C:
            continue
        case <-ctx.Done():
            return results, ctx.Err()
        }
    }

    return e.health.Check(ctx, proposal, changes)
}

// updateStatus updates execution status
func (e *executionEngine) updateStatus(executionID string, status ExecutionStatusType) {
    e.statusMu.Lock()
    defer e.statusMu.Unlock()

    e.statuses[executionID] = &ExecutionStatus{
        ExecutionID: executionID,
        Status:      status,
        UpdatedAt:   time.Now(),
    }
}

// allHealthChecksPassed checks if all health checks passed
func allHealthChecksPassed(results []HealthCheckResult) bool {
    for _, result := range results {
        if !result.Passed {
            return false
        }
    }
    return true
}

// generateExecutionID generates a unique execution ID
func generateExecutionID() string {
    return "exec-" + uuid.New().String()[:8]
}
```

```go
// pkg/execution/validators/consensus.go
package validators

import (
    "context"
    "fmt"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

// ConsensusValidator validates voting score
type ConsensusValidator struct {
    minScore float64
}

// NewConsensusValidator creates a new consensus validator
func NewConsensusValidator(minScore float64) *ConsensusValidator {
    if minScore == 0 {
        minScore = 0.80
    }
    return &ConsensusValidator{
        minScore: minScore,
    }
}

// Validate checks if proposal has sufficient consensus
func (v *ConsensusValidator) Validate(ctx context.Context, proposal *execution.Proposal) error {
    if proposal.VotingScore < v.minScore {
        return fmt.Errorf("insufficient consensus: %.2f < %.2f", proposal.VotingScore, v.minScore)
    }
    return nil
}
```

```go
// pkg/execution/validators/schema.go
package validators

import (
    "context"
    "fmt"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

// SchemaValidator validates proposal schema
type SchemaValidator struct{}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
    return &SchemaValidator{}
}

// Validate checks proposal schema validity
func (v *SchemaValidator) Validate(ctx context.Context, proposal *execution.Proposal) error {
    if proposal.ID == "" {
        return fmt.Errorf("proposal ID is required")
    }

    if proposal.OperatorType == "" {
        return fmt.Errorf("operator type is required")
    }

    if proposal.Action == "" {
        return fmt.Errorf("action is required")
    }

    if !isValidAction(proposal.Action) {
        return fmt.Errorf("invalid action: %s (must be create, update, or delete)", proposal.Action)
    }

    if proposal.Target.Kind == "" {
        return fmt.Errorf("target kind is required")
    }

    if proposal.Target.Name == "" {
        return fmt.Errorf("target name is required")
    }

    if proposal.Target.Namespace == "" {
        return fmt.Errorf("target namespace is required")
    }

    return nil
}

// isValidAction checks if action is valid
func isValidAction(action string) bool {
    validActions := map[string]bool{
        "create": true,
        "update": true,
        "delete": true,
    }
    return validActions[action]
}
```

```go
// pkg/execution/validators/ratelimit.go
package validators

import (
    "fmt"
    "sync"
    "time"
)

// TokenBucket implements rate limiting
type TokenBucket struct {
    tokens     int
    maxTokens  int
    refillRate time.Duration
    lastRefill time.Time
    mu         sync.Mutex
}

// RateLimiter manages rate limits per operator type
type RateLimiter struct {
    buckets map[string]*TokenBucket
    mu      sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        buckets: make(map[string]*TokenBucket),
    }
}

// Allow checks if operation is allowed
func (r *RateLimiter) Allow(operatorType string) (bool, error) {
    r.mu.Lock()
    bucket, exists := r.buckets[operatorType]
    if !exists {
        bucket = &TokenBucket{
            tokens:     10,
            maxTokens:  10,
            refillRate: 60 * time.Second,
            lastRefill: time.Now(),
        }
        r.buckets[operatorType] = bucket
    }
    r.mu.Unlock()

    return bucket.consume()
}

// Reset resets rate limit for operator type
func (r *RateLimiter) Reset(operatorType string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    bucket, exists := r.buckets[operatorType]
    if !exists {
        return fmt.Errorf("no rate limit bucket for operator: %s", operatorType)
    }

    bucket.mu.Lock()
    bucket.tokens = bucket.maxTokens
    bucket.lastRefill = time.Now()
    bucket.mu.Unlock()

    return nil
}

// consume attempts to consume a token
func (tb *TokenBucket) consume() (bool, error) {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    tb.refill()

    if tb.tokens > 0 {
        tb.tokens--
        return true, nil
    }

    return false, fmt.Errorf("rate limit exceeded, try again in %v", tb.refillRate)
}

// refill adds tokens based on elapsed time
func (tb *TokenBucket) refill() {
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill)

    tokensToAdd := int(elapsed / tb.refillRate)
    if tokensToAdd > 0 {
        tb.tokens = min(tb.tokens+tokensToAdd, tb.maxTokens)
        tb.lastRefill = now
    }
}

// min returns minimum of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

```go
// pkg/execution/snapshot/manager.go
package snapshot

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/kubernetes"
)

// SnapshotManager implements StateSnapshot
type SnapshotManager struct {
    k8sClient   *kubernetes.Clientset
    dynClient   dynamic.Interface
    snapshotDir string
    mu          sync.RWMutex
    snapshots   map[string]*execution.Snapshot
}

// SnapshotOptions configures snapshot manager
type SnapshotOptions struct {
    K8sClient   *kubernetes.Clientset
    DynClient   dynamic.Interface
    SnapshotDir string
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager(opts SnapshotOptions) *SnapshotManager {
    if opts.SnapshotDir == "" {
        opts.SnapshotDir = "/tmp/clawdlinux/snapshots"
    }

    os.MkdirAll(opts.SnapshotDir, 0755)

    return &SnapshotManager{
        k8sClient:   opts.K8sClient,
        dynClient:   opts.DynClient,
        snapshotDir: opts.SnapshotDir,
        snapshots:   make(map[string]*execution.Snapshot),
    }
}

// Create creates a snapshot of current state
func (m *SnapshotManager) Create(ctx context.Context, proposal *execution.Proposal) (*execution.Snapshot, error) {
    snapshot := &execution.Snapshot{
        ID:         generateSnapshotID(),
        ProposalID: proposal.ID,
        CreatedAt:  time.Now(),
        Resources:  []execution.ResourceSnapshot{},
    }

    // Capture current resource state
    resourceSnapshot, err := m.captureResourceState(ctx, proposal.Target)
    if err != nil {
        // Resource might not exist (e.g., for create operations)
        if proposal.Action != "create" {
            return nil, fmt.Errorf("failed to capture resource state: %w", err)
        }
    } else {
        snapshot.Resources = append(snapshot.Resources, *resourceSnapshot)
    }

    // Save snapshot to disk
    if err := m.saveSnapshot(snapshot); err != nil {
        return nil, fmt.Errorf("failed to save snapshot: %w", err)
    }

    // Store in memory
    m.mu.Lock()
    m.snapshots[snapshot.ID] = snapshot
    m.mu.Unlock()

    return snapshot, nil
}

// Restore restores from a snapshot
func (m *SnapshotManager) Restore(ctx context.Context, snapshotID string) error {
    snapshot, err := m.Get(snapshotID)
    if err != nil {
        return err
    }

    for _, resourceSnapshot := range snapshot.Resources {
        if err := m.restoreResource(ctx, &resourceSnapshot); err != nil {
            return fmt.Errorf("failed to restore %s/%s: %w",
                resourceSnapshot.Kind, resourceSnapshot.Name, err)
        }
    }

    return nil
}

// Delete removes a snapshot
func (m *SnapshotManager) Delete(ctx context.Context, snapshotID string) error {
    m.mu.Lock()
    delete(m.snapshots, snapshotID)
    m.mu.Unlock()

    snapshotPath := filepath.Join(m.snapshotDir, snapshotID+".json")
    return os.Remove(snapshotPath)
}

// Get retrieves a snapshot
func (m *SnapshotManager) Get(snapshotID string) (*execution.Snapshot, error) {
    m.mu.RLock()
    snapshot, exists := m.snapshots[snapshotID]
    m.mu.RUnlock()

    if exists {
        return snapshot, nil
    }

    // Load from disk
    return m.loadSnapshot(snapshotID)
}

// captureResourceState captures current resource state
func (m *SnapshotManager) captureResourceState(ctx context.Context, target execution.TargetResource) (*execution.ResourceSnapshot, error) {
    gvr := m.kindToGVR(target.Kind)

    var resource *unstructured.Unstructured
    var err error

    if m.dynClient != nil {
        resource, err = m.dynClient.Resource(gvr).Namespace(target.Namespace).Get(ctx, target.Name, metav1.GetOptions{})
    } else {
        // Fallback to typed client for common resources
        return m.captureResourceStateTyped(ctx, target)
    }

    if err != nil {
        return nil, err
    }

    state := resource.Object
    version := resource.GetResourceVersion()

    return &execution.ResourceSnapshot{
        Kind:      target.Kind,
        Namespace: target.Namespace,
        Name:      target.Name,
        State:     state,
        Version:   version,
    }, nil
}

// captureResourceStateTyped captures state using typed client
func (m *SnapshotManager) captureResourceStateTyped(ctx context.Context, target execution.TargetResource) (*execution.ResourceSnapshot, error) {
    var state map[string]interface{}
    var version string

    switch target.Kind {
    case "ConfigMap":
        cm, err := m.k8sClient.CoreV1().ConfigMaps(target.Namespace).Get(ctx, target.Name, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
        data, _ := json.Marshal(cm)
        json.Unmarshal(data, &state)
        version = cm.ResourceVersion

    case "Deployment":
        deploy, err := m.k8sClient.AppsV1().Deployments(target.Namespace).Get(ctx, target.Name, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
        data, _ := json.Marshal(deploy)
        json.Unmarshal(data, &state)
        version = deploy.ResourceVersion

    default:
        return nil, fmt.Errorf("unsupported resource kind: %s", target.Kind)
    }

    return &execution.ResourceSnapshot{
        Kind:      target.Kind,
        Namespace: target.Namespace,
        Name:      target.Name,
        State:     state,
        Version:   version,
    }, nil
}

// restoreResource restores a resource from snapshot
func (m *SnapshotManager) restoreResource(ctx context.Context, snapshot *execution.ResourceSnapshot) error {
    gvr := m.kindToGVR(snapshot.Kind)

    obj := &unstructured.Unstructured{
        Object: snapshot.State,
    }

    if m.dynClient != nil {
        _, err := m.dynClient.Resource(gvr).Namespace(snapshot.Namespace).Update(ctx, obj, metav1.UpdateOptions{})
        return err
    }

    return fmt.Errorf("dynamic client not available for restore")
}

// saveSnapshot saves snapshot to disk
func (m *SnapshotManager) saveSnapshot(snapshot *execution.Snapshot) error {
    snapshotPath := filepath.Join(m.snapshotDir, snapshot.ID+".json")

    data, err := json.MarshalIndent(snapshot, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(snapshotPath, data, 0644)
}

// loadSnapshot loads snapshot from disk
func (m *SnapshotManager) loadSnapshot(snapshotID string) (*execution.Snapshot, error) {
    snapshotPath := filepath.Join(m.snapshotDir, snapshotID+".json")

    data, err := os.ReadFile(snapshotPath)
    if err != nil {
        return nil, err
    }

    var snapshot execution.Snapshot
    if err := json.Unmarshal(data, &snapshot); err != nil {
        return nil, err
    }

    m.mu.Lock()
    m.snapshots[snapshotID] = &snapshot
    m.mu.Unlock()

    return &snapshot, nil
}

// kindToGVR converts Kind to GroupVersionResource
func (m *SnapshotManager) kindToGVR(kind string) schema.GroupVersionResource {
    gvrMap := map[string]schema.GroupVersionResource{
        "ConfigMap":  {Group: "", Version: "v1", Resource: "configmaps"},
        "Deployment": {Group: "apps", Version: "v1", Resource: "deployments"},
        "Service":    {Group: "", Version: "v1", Resource: "services"},
    }

    if gvr, ok := gvrMap[kind]; ok {
        return gvr
    }

    return schema.GroupVersionResource{Resource: kind}
}

// generateSnapshotID generates unique snapshot ID
func generateSnapshotID() string {
    return "snap-" + uuid.New().String()[:8]
}
```

```go
// pkg/execution/operators/executor.go
package operators

import (
    "context"
    "fmt"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    "k8s.io/client-go/kubernetes"
)

// OperatorExecutor executes operator changes
type OperatorExecutor struct {
    k8sClient *kubernetes.Clientset
    operators map[string]Operator
}

// Operator handles specific resource types
type Operator interface {
    Execute(ctx context.Context, proposal *execution.Proposal, snapshot *execution.Snapshot) ([]execution.Change, error)
    Supports(kind string) bool
}

// NewOperatorExecutor creates a new operator executor
func NewOperatorExecutor(k8sClient *kubernetes.Clientset) *OperatorExecutor {
    executor := &OperatorExecutor{
        k8sClient: k8sClient,
        operators: make(map[string]Operator),
    }

    // Register operators
    executor.registerOperator(NewConfigMapOperator(k8sClient))
    executor.registerOperator(NewDeploymentOperator(k8sClient))

    return executor
}

// Execute executes a proposal
func (e *OperatorExecutor) Execute(ctx context.Context, proposal *execution.Proposal, snapshot *execution.Snapshot) ([]execution.Change, error) {
    operator := e.operators[proposal.Target.Kind]
    if operator == nil {
        return nil, fmt.Errorf("no operator registered for kind: %s", proposal.Target.Kind)
    }

    return operator.Execute(ctx, proposal, snapshot)
}

// SupportedOperators returns list of supported operators
func (e *OperatorExecutor) SupportedOperators() []string {
    kinds := make([]string, 0, len(e.operators))
    for kind := range e.operators {
        kinds = append(kinds, kind)
    }
    return kinds
}

// registerOperator registers an operator
func (e *OperatorExecutor) registerOperator(operator Operator) {
    for _, kind := range []string{"ConfigMap", "Deployment", "Service"} {
        if operator.Supports(kind) {
            e.operators[kind] = operator
        }
    }
}
```

```go
// pkg/execution/operators/configmap.go
package operators

import (
    "context"
    "fmt"
    "time"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// ConfigMapOperator handles ConfigMap operations
type ConfigMapOperator struct {
    k8sClient *kubernetes.Clientset
}

// NewConfigMapOperator creates a ConfigMap operator
func NewConfigMapOperator(k8sClient *kubernetes.Clientset) *ConfigMapOperator {
    return &ConfigMapOperator{
        k8sClient: k8sClient,
    }
}

// Execute executes ConfigMap change
func (o *ConfigMapOperator) Execute(ctx context.Context, proposal *execution.Proposal, snapshot *execution.Snapshot) ([]execution.Change, error) {
    changes := []execution.Change{}

    switch proposal.Action {
    case "create":
        change, err := o.createConfigMap(ctx, proposal)
        if err != nil {
            return nil, err
        }
        changes = append(changes, *change)

    case "update":
        change, err := o.updateConfigMap(ctx, proposal)
        if err != nil {
            return nil, err
        }
        changes = append(changes, *change)

    case "delete":
        change, err := o.deleteConfigMap(ctx, proposal)
        if err != nil {
            return nil, err
        }
        changes = append(changes, *change)

    default:
        return nil, fmt.Errorf("unsupported action: %s", proposal.Action)
    }

    return changes, nil
}

// Supports checks if operator supports kind
func (o *ConfigMapOperator) Supports(kind string) bool {
    return kind == "ConfigMap"
}

// createConfigMap creates a ConfigMap
func (o *ConfigMapOperator) createConfigMap(ctx context.Context, proposal *execution.Proposal) (*execution.Change, error) {
    data := make(map[string]string)
    for k, v := range proposal.Changes {
        if str, ok := v.(string); ok {
            data[k] = str
        }
    }

    cm := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      proposal.Target.Name,
            Namespace: proposal.Target.Namespace,
        },
        Data: data,
    }

    _, err := o.k8sClient.CoreV1().ConfigMaps(proposal.Target.Namespace).Create(ctx, cm, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    return &execution.Change{
        ResourceKind: "ConfigMap",
        ResourceName: proposal.Target.Name,
        Namespace:    proposal.Target.Namespace,
        Action:       "created",
        Diff:         fmt.Sprintf("Created ConfigMap with %d keys", len(data)),
        AppliedAt:    time.Now(),
    }, nil
}

// updateConfigMap updates a ConfigMap
func (o *ConfigMapOperator) updateConfigMap(ctx context.Context, proposal *execution.Proposal) (*execution.Change, error) {
    cm, err := o.k8sClient.CoreV1().ConfigMaps(proposal.Target.Namespace).Get(ctx, proposal.Target.Name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    oldData := make(map[string]string)
    for k, v := range cm.Data {
        oldData[k] = v
    }

    for k, v := range proposal.Changes {
        if str, ok := v.(string); ok {
            cm.Data[k] = str
        }
    }

    _, err = o.k8sClient.CoreV1().ConfigMaps(proposal.Target.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
    if err != nil {
        return nil, err
    }

    return &execution.Change{
        ResourceKind: "ConfigMap",
        ResourceName: proposal.Target.Name,
        Namespace:    proposal.Target.Namespace,
        Action:       "updated",
        Diff:         fmt.Sprintf("Updated %d keys", len(proposal.Changes)),
        AppliedAt:    time.Now(),
    }, nil
}

// deleteConfigMap deletes a ConfigMap
func (o *ConfigMapOperator) deleteConfigMap(ctx context.Context, proposal *execution.Proposal) (*execution.Change, error) {
    err := o.k8sClient.CoreV1().ConfigMaps(proposal.Target.Namespace).Delete(ctx, proposal.Target.Name, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return nil, err
    }

    return &execution.Change{
        ResourceKind: "ConfigMap",
        ResourceName: proposal.Target.Name,
        Namespace:    proposal.Target.Namespace,
        Action:       "deleted",
        Diff:         "Deleted ConfigMap",
        AppliedAt:    time.Now(),
    }, nil
}
```

```go
// pkg/execution/operators/deployment.go
package operators

import (
    "context"
    "fmt"
    "time"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// DeploymentOperator handles Deployment operations
type DeploymentOperator struct {
    k8sClient *kubernetes.Clientset
}

// NewDeploymentOperator creates a Deployment operator
func NewDeploymentOperator(k8sClient *kubernetes.Clientset) *DeploymentOperator {
    return &DeploymentOperator{
        k8sClient: k8sClient,
    }
}

// Execute executes Deployment change
func (o *DeploymentOperator) Execute(ctx context.Context, proposal *execution.Proposal, snapshot *execution.Snapshot) ([]execution.Change, error) {
    changes := []execution.Change{}

    switch proposal.Action {
    case "update":
        change, err := o.updateDeployment(ctx, proposal)
        if err != nil {
            return nil, err
        }
        changes = append(changes, *change)

    case "delete":
        change, err := o.deleteDeployment(ctx, proposal)
        if err != nil {
            return nil, err
        }
        changes = append(changes, *change)

    default:
        return nil, fmt.Errorf("unsupported action for Deployment: %s", proposal.Action)
    }

    return changes, nil
}

// Supports checks if operator supports kind
func (o *DeploymentOperator) Supports(kind string) bool {
    return kind == "Deployment"
}

// updateDeployment updates a Deployment
func (o *DeploymentOperator) updateDeployment(ctx context.Context, proposal *execution.Proposal) (*execution.Change, error) {
    deploy, err := o.k8sClient.AppsV1().Deployments(proposal.Target.Namespace).Get(ctx, proposal.Target.Name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    // Apply changes (simplified - only replicas for now)
    if replicas, ok := proposal.Changes["replicas"].(float64); ok {
        replicaCount := int32(replicas)
        deploy.Spec.Replicas = &replicaCount
    }

    _, err = o.k8sClient.AppsV1().Deployments(proposal.Target.Namespace).Update(ctx, deploy, metav1.UpdateOptions{})
    if err != nil {
        return nil, err
    }

    return &execution.Change{
        ResourceKind: "Deployment",
        ResourceName: proposal.Target.Name,
        Namespace:    proposal.Target.Namespace,
        Action:       "updated",
        Diff:         fmt.Sprintf("Updated deployment with %d changes", len(proposal.Changes)),
        AppliedAt:    time.Now(),
    }, nil
}

// deleteDeployment deletes a Deployment
func (o *DeploymentOperator) deleteDeployment(ctx context.Context, proposal *execution.Proposal) (*execution.Change, error) {
    err := o.k8sClient.AppsV1().Deployments(proposal.Target.Namespace).Delete(ctx, proposal.Target.Name, metav1.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        return nil, err
    }

    return &execution.Change{
        ResourceKind: "Deployment",
        ResourceName: proposal.Target.Name,
        Namespace:    proposal.Target.Namespace,
        Action:       "deleted",
        Diff:         "Deleted Deployment",
        AppliedAt:    time.Now(),
    }, nil
}
```

```go
// pkg/execution/health/checker.go
package health

import (
    "context"
    "fmt"
    "time"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// HealthChecker performs health checks
type HealthChecker struct {
    k8sClient *kubernetes.Clientset
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(k8sClient *kubernetes.Clientset) *HealthChecker {
    return &HealthChecker{
        k8sClient: k8sClient,
    }
}

// Check performs health checks on changes
func (h *HealthChecker) Check(ctx context.Context, proposal *execution.Proposal, changes []execution.Change) ([]execution.HealthCheckResult, error) {
    results := []execution.HealthCheckResult{}

    for _, change := range changes {
        switch change.ResourceKind {
        case "Deployment":
            result := h.checkDeployment(ctx, change)
            results = append(results, result)

        case "ConfigMap":
            result := h.checkConfigMap(ctx, change)
            results = append(results, result)

        case "Service":
            result := h.checkService(ctx, change)
            results = append(results, result)
        }
    }

    return results, nil
}

// checkDeployment checks deployment health
func (h *HealthChecker) checkDeployment(ctx context.Context, change execution.Change) execution.HealthCheckResult {
    if change.Action == "deleted" {
        return execution.HealthCheckResult{
            CheckName: "deployment_deleted",
            Passed:    true,
            Message:   "Deployment successfully deleted",
            CheckedAt: time.Now(),
        }
    }

    deployment, err := h.k8sClient.AppsV1().Deployments(change.Namespace).Get(ctx, change.ResourceName, metav1.GetOptions{})
    if err != nil {
        return execution.HealthCheckResult{
            CheckName: "deployment_exists",
            Passed:    false,
            Message:   err.Error(),
            CheckedAt: time.Now(),
        }
    }

    if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
        return execution.HealthCheckResult{
            CheckName: "deployment_ready",
            Passed:    true,
            Message:   fmt.Sprintf("%d/%d replicas ready", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas),
            CheckedAt: time.Now(),
        }
    }

    return execution.HealthCheckResult{
        CheckName: "deployment_ready",
        Passed:    false,
        Message:   fmt.Sprintf("only %d/%d replicas ready", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas),
        CheckedAt: time.Now(),
    }
}

// checkConfigMap checks ConfigMap health
func (h *HealthChecker) checkConfigMap(ctx context.Context, change execution.Change) execution.HealthCheckResult {
    if change.Action == "deleted" {
        return execution.HealthCheckResult{
            CheckName: "configmap_deleted",
            Passed:    true,
            Message:   "ConfigMap successfully deleted",
            CheckedAt: time.Now(),
        }
    }

    _, err := h.k8sClient.CoreV1().ConfigMaps(change.Namespace).Get(ctx, change.ResourceName, metav1.GetOptions{})
    if err != nil {
        return execution.HealthCheckResult{
            CheckName: "configmap_exists",
            Passed:    false,
            Message:   err.Error(),
            CheckedAt: time.Now(),
        }
    }

    return execution.HealthCheckResult{
        CheckName: "configmap_exists",
        Passed:    true,
        Message:   "ConfigMap exists and is accessible",
        CheckedAt: time.Now(),
    }
}

// checkService checks Service health
func (h *HealthChecker) checkService(ctx context.Context, change execution.Change) execution.HealthCheckResult {
    if change.Action == "deleted" {
        return execution.HealthCheckResult{
            CheckName: "service_deleted",
            Passed:    true,
            Message:   "Service successfully deleted",
            CheckedAt: time.Now(),
        }
    }

    svc, err := h.k8sClient.CoreV1().Services(change.Namespace).Get(ctx, change.ResourceName, metav1.GetOptions{})
    if err != nil {
        return execution.HealthCheckResult{
            CheckName: "service_exists",
            Passed:    false,
            Message:   err.Error(),
            CheckedAt: time.Now(),
        }
    }

    if svc.Spec.ClusterIP != "" {
        return execution.HealthCheckResult{
            CheckName: "service_ready",
            Passed:    true,
            Message:   fmt.Sprintf("Service has ClusterIP: %s", svc.Spec.ClusterIP),
            CheckedAt: time.Now(),
        }
    }

    return execution.HealthCheckResult{
        CheckName: "service_ready",
        Passed:    false,
        Message:   "Service has no ClusterIP assigned",
        CheckedAt: time.Now(),
    }
}
```

```go
// pkg/execution/rollback/manager.go
package rollback

import (
    "context"
    "fmt"
    "sync"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    "k8s.io/client-go/kubernetes"
)

// RollbackManager handles rollback operations
type RollbackManager struct {
    snapshotMgr execution.StateSnapshot
    k8sClient   *kubernetes.Clientset
    audit       execution.AuditLogger
    mu          sync.RWMutex
    executions  map[string]*execution.ExecutionResult
}

// NewRollbackManager creates a new rollback manager
func NewRollbackManager(snapshotMgr execution.StateSnapshot, k8sClient *kubernetes.Clientset, audit execution.AuditLogger) *RollbackManager {
    return &RollbackManager{
        snapshotMgr: snapshotMgr,
        k8sClient:   k8sClient,
        audit:       audit,
        executions:  make(map[string]*execution.ExecutionResult),
    }
}

// Rollback performs manual rollback
func (r *RollbackManager) Rollback(ctx context.Context, executionID string) error {
    r.mu.RLock()
    exec, exists := r.executions[executionID]
    r.mu.RUnlock()

    if !exists {
        return fmt.Errorf("execution not found: %s", executionID)
    }

    return r.AutoRollback(ctx, exec)
}

// AutoRollback automatically rolls back failed execution
func (r *RollbackManager) AutoRollback(ctx context.Context, exec *execution.ExecutionResult) error {
    if exec.RollbackID == "" {
        return fmt.Errorf("no rollback snapshot available for execution %s", exec.ExecutionID)
    }

    // Store execution for later rollback
    r.mu.Lock()
    r.executions[exec.ExecutionID] = exec
    r.mu.Unlock()

    // Restore snapshot
    if err := r.snapshotMgr.Restore(ctx, exec.RollbackID); err != nil {
        return fmt.Errorf("failed to restore snapshot: %w", err)
    }

    // Update execution status
    exec.Status = execution.StatusRolledBack
    r.audit.LogExecution(exec)

    return nil
}
```

```go
// pkg/execution/audit/logger.go
package audit

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
)

// AuditLogger logs execution events
type AuditLogger struct {
    auditDir string
    mu       sync.Mutex
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(auditDir string) *AuditLogger {
    if auditDir == "" {
        auditDir = "/tmp/clawdlinux/audit"
    }

    os.MkdirAll(filepath.Join(auditDir, "executions"), 0755)

    return &AuditLogger{
        auditDir: auditDir,
    }
}

// LogExecution logs an execution result
func (l *AuditLogger) LogExecution(exec *execution.ExecutionResult) error {
    l.mu.Lock()
    defer l.mu.Unlock()

    date := time.Now().Format("2006-01-02")
    filename := filepath.Join(l.auditDir, "executions", fmt.Sprintf("%s.jsonl", date))

    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(exec)
}

// Query retrieves executions matching filter
func (l *AuditLogger) Query(filter execution.AuditFilter) ([]*execution.ExecutionResult, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    results := []*execution.ExecutionResult{}

    // Read all JSONL files in date range
    files, err := filepath.Glob(filepath.Join(l.auditDir, "executions", "*.jsonl"))
    if err != nil {
        return nil, err
    }

    for _, filename := range files {
        fileResults, err := l.queryFile(filename, filter)
        if err != nil {
            continue
        }
        results = append(results, fileResults...)
    }

    return results, nil
}

// queryFile queries a single JSONL file
func (l *AuditLogger) queryFile(filename string, filter execution.AuditFilter) ([]*execution.ExecutionResult, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    results := []*execution.ExecutionResult{}
    decoder := json.NewDecoder(file)

    for decoder.More() {
        var exec execution.ExecutionResult
        if err := decoder.Decode(&exec); err != nil {
            continue
        }

        if l.matchesFilter(&exec, filter) {
            results = append(results, &exec)
        }
    }

    return results, nil
}

// matchesFilter checks if execution matches filter
func (l *AuditLogger) matchesFilter(exec *execution.ExecutionResult, filter execution.AuditFilter) bool {
    if filter.ProposalID != "" && exec.ProposalID != filter.ProposalID {
        return false
    }

    if filter.Status != "" && exec.Status != filter.Status {
        return false
    }

    if !filter.StartDate.IsZero() && exec.StartedAt.Before(filter.StartDate) {
        return false
    }

    if !filter.EndDate.IsZero() && exec.CompletedAt.After(filter.EndDate) {
        return false
    }

    return true
}
```

```go
// cmd/clawdctl/execution/execute.go
package execution

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    "github.com/spf13/cobra"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

// ExecuteCmd returns the execute command
func ExecuteCmd() *cobra.Command {
    var (
        proposalFile    string
        proposalID      string
        dryRun          bool
        skipHealthCheck bool
        timeout         time.Duration
    )

    cmd := &cobra.Command{
        Use:   "execute",
        Short: "Execute an approved proposal",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()

            // Load kubeconfig
            kubeconfig := os.Getenv("KUBECONFIG")
            if kubeconfig == "" {
                kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
            }

            config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
            if err != nil {
                return fmt.Errorf("failed to load kubeconfig: %w", err)
            }

            k8sClient, err := kubernetes.NewForConfig(config)
            if err != nil {
                return fmt.Errorf("failed to create k8s client: %w", err)
            }

            // Create execution engine
            engine := execution.NewExecutionEngine(execution.EngineOptions{
                K8sClient:    k8sClient,
                AuditPath:    "/tmp/clawdlinux/audit",
                MinConsensus: 0.80,
                SnapshotDir:  "/tmp/clawdlinux/snapshots",
            })

            // Load proposal
            var proposal execution.Proposal
            if proposalFile != "" {
                data, err := os.ReadFile(proposalFile)
                if err != nil {
                    return fmt.Errorf("failed to read proposal file: %w", err)
                }
                if err := json.Unmarshal(data, &proposal); err != nil {
                    return fmt.Errorf("failed to parse proposal: %w", err)
                }
            } else if proposalID != "" {
                // Load from proposal store (simplified)
                return fmt.Errorf("loading by ID not yet implemented, use --file")
            } else {
                return fmt.Errorf("either --file or --proposal-id required")
            }

            // Execute or dry-run
            if dryRun {
                result, err := engine.DryRun(ctx, &proposal)
                if err != nil {
                    return err
                }

                fmt.Printf("Dry-run result for proposal %s:\n", result.ProposalID)
                fmt.Printf("Would apply: %v\n", result.WouldApply)
                fmt.Printf("Predicted changes:\n%s\n", result.PredictedDiff)
                if len(result.ValidationErrors) > 0 {
                    fmt.Printf("Validation errors:\n")
                    for _, err := range result.ValidationErrors {
                        fmt.Printf("  - %s\n", err)
                    }
                }
                return nil
            }

            // Real execution
            result, err := engine.Execute(ctx, &proposal, execution.ExecutionOptions{
                DryRun:          false,
                SkipHealthCheck: skipHealthCheck,
                RollbackTimeout: timeout,
            })

            fmt.Printf("Execution ID: %s\n", result.ExecutionID)
            fmt.Printf("Status: %s\n", result.Status)
            fmt.Printf("Started: %s\n", result.StartedAt.Format(time.RFC3339))
            fmt.Printf("Completed: %s\n", result.CompletedAt.Format(time.RFC3339))

            if result.Error != "" {
                fmt.Printf("Error: %s\n", result.Error)
            }

            if len(result.AppliedChanges) > 0 {
                fmt.Printf("\nApplied changes:\n")
                for _, change := range result.AppliedChanges {
                    fmt.Printf("  - %s %s/%s: %s\n", change.Action, change.ResourceKind, change.ResourceName, change.Diff)
                }
            }

            if len(result.HealthCheckResults) > 0 {
                fmt.Printf("\nHealth checks:\n")
                for _, hc := range result.HealthCheckResults {
                    status := "✓"
                    if !hc.Passed {
                        status = "✗"
                    }
                    fmt.Printf("  %s %s: %s\n", status, hc.CheckName, hc.Message)
                }
            }

            return err
        },
    }

    cmd.Flags().StringVar(&proposalFile, "file", "", "Proposal JSON file")
    cmd.Flags().StringVar(&proposalID, "proposal-id", "", "Proposal ID")
    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate execution without applying changes")
    cmd.Flags().BoolVar(&skipHealthCheck, "skip-health-check", false, "Skip health checks")
    cmd.Flags().DurationVar(&timeout, "timeout", 5*time.Minute, "Health check timeout")

    return cmd
}
```

```go
// cmd/clawdctl/execution/rollback.go
package execution

import (
    "context"
    "fmt"
    "os"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    "github.com/spf13/cobra"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

// RollbackCmd returns the rollback command
func RollbackCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "rollback EXECUTION_ID",
        Short: "Rollback a failed execution",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            executionID := args[0]

            // Load kubeconfig
            kubeconfig := os.Getenv("KUBECONFIG")
            if kubeconfig == "" {
                kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
            }

            config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
            if err != nil {
                return fmt.Errorf("failed to load kubeconfig: %w", err)
            }

            k8sClient, err := kubernetes.NewForConfig(config)
            if err != nil {
                return fmt.Errorf("failed to create k8s client: %w", err)
            }

            // Create execution engine
            engine := execution.NewExecutionEngine(execution.EngineOptions{
                K8sClient:    k8sClient,
                AuditPath:    "/tmp/clawdlinux/audit",
                MinConsensus: 0.80,
                SnapshotDir:  "/tmp/clawdlinux/snapshots",
            })

            // Perform rollback
            if err := engine.Rollback(ctx, executionID); err != nil {
                return fmt.Errorf("rollback failed: %w", err)
            }

            fmt.Printf("Successfully rolled back execution %s\n", executionID)
            return nil
        },
    }

    return cmd
}
```

```go
// cmd/clawdctl/execution/status.go
package execution

import (
    "fmt"
    "os"

    "github.com/shreyanshjain7174/clawdlinux/pkg/execution"
    "github.com/spf13/cobra"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

// StatusCmd returns the status command
func StatusCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "status EXECUTION_ID",
        Short: "Get execution status",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            executionID := args[0]

            // Load kubeconfig
            kubeconfig := os.Getenv("KUBECONFIG")
            if kubeconfig == "" {
                kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
            }

            config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
            if err != nil {
                return fmt.Errorf("failed to load kubeconfig: %w", err)
            }

            k8sClient, err := kubernetes.NewForConfig(config)
            if err != nil {
                return fmt.Errorf("failed to create k8s client: %w", err)
            }

            // Create execution engine
            engine := execution.NewExecutionEngine(execution.EngineOptions{
                K8sClient:    k8sClient,
                AuditPath:    "/tmp/clawdlinux/audit",
                MinConsensus: 0.80,
                SnapshotDir:  "/tmp/clawdlinux/snapshots",
            })

            // Get status
            status, err := engine.GetStatus(executionID)
            if err != nil {
                return fmt.Errorf("failed to get status: %w", err)
            }

            fmt.Printf("Execution ID: %s\n", status.ExecutionID)
            fmt.Printf("Status: %s\n", status.Status)
            fmt.Printf("Updated: %s\n", status.UpdatedAt.Format(time.RFC3339))

            return nil
        },
    }

    return cmd
}
```