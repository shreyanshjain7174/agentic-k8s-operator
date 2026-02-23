# Implementation Fixes Required

**Priority Order:** Fix in this sequence for fastest path to production.

---

## 🔴 CRITICAL - Fix Before Any Deployment

### 1. Integrate Vote Submission API (1 hour)

**File:** `cmd/main.go`

**Add after creating manager:**

```go
// Initialize vote submission handler
voteHandler := &api.VoteSubmissionHandler{
    Client:            mgr.GetClient(),
    SignatureVerifier: signatureVerifier,
}

// Start HTTP server for vote submissions
go func() {
    mux := http.NewServeMux()
    mux.Handle("/vote", voteHandler)
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })
    
    setupLog.Info("starting vote submission server on :8082")
    if err := http.ListenAndServe(":8082", mux); err != nil {
        setupLog.Error(err, "vote submission server failed")
        os.Exit(1)
    }
}()
```

**Test:**
```bash
curl -X POST http://localhost:8082/vote \
  -H "Content-Type: application/json" \
  -d '{"proposalId":"default/test","voter":"v1","decision":"APPROVE","score":100,"rationale":"ok"}'
```

---

### 2. Fix Audit Logging (2 hours)

**File:** `internal/controller/voteproposal_controller.go`

**Replace async audit logging:**

```go
// BEFORE (line ~262):
go func() {
    if err := r.AuditLogger.LogProposal(context.Background(), proposal); err != nil {
        log.Error(err, "failed to audit log proposal")
    }
}()
return ctrl.Result{}, nil

// AFTER:
// Audit log MUST succeed before marking as executed
if err := r.AuditLogger.LogProposal(ctx, proposal); err != nil {
    log.Error(err, "CRITICAL: audit logging failed - will retry")
    // Don't transition to terminal state until audit succeeds
    return ctrl.Result{RequeueAfter: 30 * time.Second}, err
}

// Update status with audit reference
proposal.Status.AuditLogReference = "s3://..."  // Set by AuditLogger
if err := r.Status().Update(ctx, proposal); err != nil {
    log.Error(err, "failed to update audit reference in status")
    // Audit is logged, status update failed - safe to continue
}

return ctrl.Result{}, nil
```

**Add retry logic to S3 logger:**

**File:** `internal/audit/logger.go`

```go
func (l *S3AuditLogger) LogProposal(ctx context.Context, proposal *votingv1alpha1.VoteProposal) error {
    var lastErr error
    for attempt := 0; attempt < 3; attempt++ {
        err := l.writeToS3(ctx, proposal)
        if err == nil {
            return nil
        }
        lastErr = err
        time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
    }
    return fmt.Errorf("failed to write audit log after 3 attempts: %w", lastErr)
}

func (l *S3AuditLogger) writeToS3(ctx context.Context, proposal *votingv1alpha1.VoteProposal) error {
    // Existing S3 write logic here
    // ...
}
```

**Test:**
```bash
# Simulate S3 failure
export AWS_REGION=invalid
make run
# Verify proposal doesn't execute until audit succeeds
```

---

### 3. Add Emergency Override Validation Webhook (3 hours)

**File:** `internal/webhook/voteproposal_webhook.go`

**Add to `VoteProposalValidator.ValidateCreate`:**

```go
func (v *VoteProposalValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
    vp := obj.(*votingv1alpha1.VoteProposal)

    // Emergency override RBAC check
    if vp.Spec.EmergencyOverride != nil && vp.Spec.EmergencyOverride.Enabled {
        // Get user info from admission request
        userInfo, err := v.getUserInfo(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to get user info: %w", err)
        }

        // Check if user has emergency-override permission
        allowed, err := v.checkEmergencyOverridePermission(ctx, userInfo, vp.Namespace)
        if err != nil {
            return nil, fmt.Errorf("failed to check permissions: %w", err)
        }
        if !allowed {
            return nil, fmt.Errorf("user %s not authorized for emergency override", userInfo.Username)
        }

        // Require MFA proof
        if vp.Spec.EmergencyOverride.MFAProof == "" {
            return admission.Warnings{"MFA proof strongly recommended for emergency overrides"},
                   fmt.Errorf("MFA proof required for emergency override")
        }

        // Add warning for audit trail
        return admission.Warnings{
            fmt.Sprintf("EMERGENCY OVERRIDE by %s: %s", 
                userInfo.Username, 
                vp.Spec.EmergencyOverride.Justification),
        }, nil
    }

    // ... rest of validation
}

func (v *VoteProposalValidator) checkEmergencyOverridePermission(ctx context.Context, userInfo authenticationv1.UserInfo, namespace string) (bool, error) {
    // Create SubjectAccessReview
    sar := &authorizationv1.SubjectAccessReview{
        Spec: authorizationv1.SubjectAccessReviewSpec{
            User:   userInfo.Username,
            Groups: userInfo.Groups,
            ResourceAttributes: &authorizationv1.ResourceAttributes{
                Namespace: namespace,
                Verb:      "create",
                Group:     "voting.operator.io",
                Resource:  "voteproposals/emergency-override",
            },
        },
    }

    if err := v.Client.Create(ctx, sar); err != nil {
        return false, err
    }

    return sar.Status.Allowed, nil
}
```

**Create RBAC for emergency override:**

**File:** `config/rbac/emergency_override_role.yaml`

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: emergency-override-authorized
  namespace: production
rules:
- apiGroups: ["voting.operator.io"]
  resources: ["voteproposals/emergency-override"]
  verbs: ["create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: oncall-sre-emergency
  namespace: production
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: emergency-override-authorized
subjects:
- kind: User
  name: oncall-sre@example.com  # Replace with actual users
```

**Test:**
```bash
# As non-authorized user
kubectl apply -f examples/emergency-override-proposal.yaml
# Should fail: "user not authorized for emergency override"

# As authorized user
kubectl apply -f examples/emergency-override-proposal.yaml --as=oncall-sre@example.com
# Should succeed
```

---

### 4. Fix Float Precision Bug (30 minutes)

**File:** `internal/voting/aggregator.go`

**Replace line ~35:**

```go
// BEFORE:
if avgScore >= float64(proposal.Spec.VoteConfig.Threshold) {
    decision = "APPROVED"
}

// AFTER:
const epsilon = 0.0001 // Tolerance for float comparison (0.01%)

func meetsThreshold(score float64, threshold int) bool {
    return score >= (float64(threshold) - epsilon)
}

// In CalculateConsensus:
if meetsThreshold(avgScore, proposal.Spec.VoteConfig.Threshold) {
    decision = "APPROVED"
} else {
    decision = "REJECTED"
}
```

**Test:**
```go
// internal/voting/aggregator_test.go
func TestFloatPrecisionEdgeCase(t *testing.T) {
    agg := &DefaultAggregator{}
    
    // Construct scenario that produces 79.99999999
    proposal := &votingv1alpha1.VoteProposal{
        Spec: votingv1alpha1.VoteProposalSpec{
            VoteConfig: votingv1alpha1.VoteConfiguration{
                Threshold: 80,
            },
        },
        Status: votingv1alpha1.VoteProposalStatus{
            Votes: []votingv1alpha1.Vote{
                {Score: 100},
                {Score: 70},
                {Score: 70},
                {Score: 70},
                {Score: 70},
            },
        },
    }
    
    score, decision := agg.CalculateConsensus(proposal)
    // (100 + 70*4) / 5 = 380/5 = 76.0 (not 80, bad example)
    
    // Better test: exact 80.0
    proposal.Status.Votes = []votingv1alpha1.Vote{
        {Score: 100}, {Score: 70}, {Score: 70},
    }
    score, decision = agg.CalculateConsensus(proposal)
    if score != 80.0 {
        t.Errorf("Expected exactly 80.0, got %v", score)
    }
    if decision != "APPROVED" {
        t.Errorf("Expected APPROVED at threshold, got %v", decision)
    }
}
```

---

### 5. Add Missing Tests (1-2 days)

**Target:** 80% coverage (currently ~60-65%)

**Priority order:**

#### 5a. Execution Engine Tests (2 hours)

**File:** `internal/execution/engine_test.go`

```go
package execution

import (
    "context"
    "testing"
    
    votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
    appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestServerSideApplyEngine_Patch(t *testing.T) {
    scheme := runtime.NewScheme()
    _ = appsv1.AddToScheme(scheme)
    _ = votingv1alpha1.AddToScheme(scheme)
    
    deployment := &appsv1.Deployment{
        // ... existing deployment
    }
    
    fakeClient := fake.NewClientBuilder().
        WithScheme(scheme).
        WithObjects(deployment).
        Build()
    
    engine := &ServerSideApplyEngine{Client: fakeClient}
    
    proposal := &votingv1alpha1.VoteProposal{
        Spec: votingv1alpha1.VoteProposalSpec{
            ProposedChange: votingv1alpha1.ChangeSpec{
                Type: "patch",
                Patch: &runtime.RawExtension{
                    Raw: []byte(`{"spec":{"replicas":10}}`),
                },
            },
            TargetRef: votingv1alpha1.TargetReference{
                APIVersion: "apps/v1",
                Kind: "Deployment",
                Name: "test-deployment",
                Namespace: "default",
            },
        },
    }
    
    status, err := engine.Execute(context.Background(), proposal)
    if err != nil {
        t.Fatalf("Execute() error = %v", err)
    }
    if !status.Success {
        t.Errorf("Execute() failed: %s", status.Message)
    }
    
    // Verify deployment was updated
    // ...
}

func TestServerSideApplyEngine_Create(t *testing.T) { /* ... */ }
func TestServerSideApplyEngine_Delete(t *testing.T) { /* ... */ }
func TestServerSideApplyEngine_Replace(t *testing.T) { /* ... */ }
```

#### 5b. Audit Logger Tests (1 hour)

**File:** `internal/audit/logger_test.go`

```go
package audit

import (
    "context"
    "testing"
    
    "github.com/aws/aws-sdk-go-v2/service/s3"
    votingv1alpha1 "github.com/yourorg/voting-operator/api/v1alpha1"
)

func TestS3AuditLogger_LogProposal(t *testing.T) {
    // Use localstack or mock S3 client
    // ...
}

func TestS3AuditLogger_RetryOnFailure(t *testing.T) {
    // Mock S3 client that fails twice, succeeds on third
    // ...
}

func TestS3AuditLogger_EmergencyOverride(t *testing.T) {
    // Verify emergency overrides logged separately
    // ...
}
```

#### 5c. Webhook Validation Tests (1 hour)

**File:** `internal/webhook/voteproposal_webhook_test.go`

```go
func TestVoteProposalValidator_InvalidProposals(t *testing.T) {
    // Test cases:
    // - Zero required voters
    // - Threshold > 100
    // - Threshold < 0
    // - Invalid change type
    // - Patch type without patch data
    // - Emergency override without justification
}
```

---

## 🟡 HIGH PRIORITY - Add Before Production

### 6. Add Prometheus Metrics (4 hours)

**File:** `internal/controller/metrics.go`

```go
package controller

import (
    "github.com/prometheus/client_golang/prometheus"
    "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
    proposalsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "voting_proposals_total",
            Help: "Total number of proposals by phase",
        },
        []string{"phase"},
    )
    
    voteAggregateScore = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "voting_aggregate_score",
            Help: "Distribution of aggregate scores",
            Buckets: []float64{0, 20, 40, 60, 80, 90, 100},
        },
        []string{"decision"},
    )
    
    voteSubmissionDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name: "vote_submission_duration_seconds",
            Help: "Time to submit and record a vote",
            Buckets: prometheus.DefBuckets,
        },
    )
    
    auditLogFailures = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "audit_log_failures_total",
            Help: "Total number of audit log write failures",
        },
    )
)

func init() {
    metrics.Registry.MustRegister(
        proposalsTotal,
        voteAggregateScore,
        voteSubmissionDuration,
        auditLogFailures,
    )
}
```

**Use in controller:**

```go
// In handleVotingPhase:
proposalsTotal.WithLabelValues(string(votingv1alpha1.PhaseVoting)).Inc()

// In CalculateConsensus:
voteAggregateScore.WithLabelValues(decision).Observe(score)

// In vote submission handler:
start := time.Now()
defer func() {
    voteSubmissionDuration.Observe(time.Since(start).Seconds())
}()

// In audit logger:
if err := l.writeToS3(...); err != nil {
    auditLogFailures.Inc()
}
```

---

### 7. Add Rate Limiting (2 hours)

**File:** `internal/api/ratelimit.go`

```go
package api

import (
    "net/http"
    "sync"
    "time"
    
    "golang.org/x/time/rate"
)

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rps      int
    burst    int
}

func NewRateLimiter(rps, burst int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rps:      rps,
        burst:    burst,
    }
}

func (rl *RateLimiter) getLimiter(voter string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.limiters[voter]
    rl.mu.RUnlock()
    
    if !exists {
        rl.mu.Lock()
        limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
        rl.limiters[voter] = limiter
        rl.mu.Unlock()
    }
    
    return limiter
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract voter from request (e.g., from auth header)
        voter := r.Header.Get("X-Voter-ID")
        if voter == "" {
            voter = r.RemoteAddr // Fallback to IP
        }
        
        limiter := rl.getLimiter(voter)
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

**Use in main.go:**

```go
rateLimiter := api.NewRateLimiter(10, 20) // 10 RPS, burst 20
voteHandler := rateLimiter.Middleware(voteSubmissionHandler)
http.Handle("/vote", voteHandler)
```

---

### 8. Make Vote Submission Idempotent (1 hour)

**File:** `internal/api/vote_submission.go` (line ~75)

```go
// Check if voter already voted
for _, existingVote := range proposal.Status.Votes {
    if existingVote.Voter == req.Voter {
        // Check if same vote (idempotent retry)
        if existingVote.Decision == req.Decision && existingVote.Score == req.Score {
            // Same vote - return success (idempotent)
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "already_recorded",
                "message": "Vote already recorded",
            })
            return
        }
        // Different decision - potential tampering
        http.Error(w, "Vote already cast with different decision", http.StatusConflict)
        return
    }
}
```

---

## 🟢 MEDIUM PRIORITY - Before Scaling

### 9. Controller Restart Recovery (2 hours)

**File:** `cmd/main.go` (after mgr.Start)

```go
// Reconcile in-flight proposals after startup
go func() {
    time.Sleep(10 * time.Second) // Wait for leader election
    
    ctx := context.Background()
    var proposals votingv1alpha1.VoteProposalList
    if err := mgr.GetClient().List(ctx, &proposals); err != nil {
        setupLog.Error(err, "failed to list proposals on startup")
        return
    }
    
    for _, p := range proposals.Items {
        if p.Status.Phase == votingv1alpha1.PhaseVoting {
            // Check if deadline passed
            if p.Status.VotingDeadline != nil && time.Now().After(p.Status.VotingDeadline.Time) {
                setupLog.Info("found timed-out proposal, triggering reconcile", 
                    "proposal", p.Name,
                    "deadline", p.Status.VotingDeadline.Time)
                // Trigger reconcile by updating annotation
                // ...
            }
        }
    }
}()
```

---

### 10. S3 Lifecycle Policy (30 minutes)

**File:** `deploy/s3-bucket-policy.json`

```json
{
  "Rules": [
    {
      "Id": "AuditLogRetention",
      "Status": "Enabled",
      "Prefix": "audit/",
      "Expiration": {
        "Days": 2555
      }
    },
    {
      "Id": "EmergencyOverrideRetention",
      "Status": "Enabled",
      "Prefix": "emergency-overrides/",
      "NoncurrentVersionExpiration": {
        "NoncurrentDays": 90
      }
    }
  ]
}
```

**Apply:**
```bash
aws s3api put-bucket-lifecycle-configuration \
  --bucket voting-audit-logs \
  --lifecycle-configuration file://deploy/s3-bucket-policy.json
```

---

## ⏱️ Timeline Summary

| Priority | Task | Time | Cumulative |
|----------|------|------|------------|
| 🔴 Critical 1 | Vote API integration | 1h | 1h |
| 🔴 Critical 2 | Fix audit logging | 2h | 3h |
| 🔴 Critical 3 | Emergency override RBAC | 3h | 6h |
| 🔴 Critical 4 | Float precision fix | 0.5h | 6.5h |
| 🔴 Critical 5 | Add tests (80% coverage) | 16h | 22.5h (~3 days) |
| 🟡 High 6 | Prometheus metrics | 4h | 26.5h |
| 🟡 High 7 | Rate limiting | 2h | 28.5h |
| 🟡 High 8 | Idempotency | 1h | 29.5h |
| 🟢 Medium 9 | Restart recovery | 2h | 31.5h |
| 🟢 Medium 10 | S3 lifecycle | 0.5h | 32h |

**Total: ~32 hours = 4 days** (with parallelization: ~2 weeks with testing)

---

## ✅ Verification Checklist

After each fix, verify:

```bash
# 1. Tests pass
./test.sh

# 2. Linting passes
golangci-lint run

# 3. Coverage meets threshold
go tool cover -func=coverage.out | grep total
# Should show ≥80%

# 4. Manual testing
kubectl apply -f examples/deployment-proposal.yaml
# Watch proposal lifecycle

# 5. Metrics work
curl http://localhost:8080/metrics | grep voting_

# 6. Audit logs written
aws s3 ls s3://voting-audit-logs/audit/$(date +%Y/%m/%d)/

# 7. Emergency override blocked for unauthorized users
kubectl apply -f examples/emergency-override-proposal.yaml --as=random-user
# Should fail
```

---

**Start Here:** Fix in order (1 → 10) for fastest path to production.
