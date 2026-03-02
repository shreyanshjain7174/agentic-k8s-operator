# Phase 8: Auto-Scaling Based on SLA

**Objective:** Dynamically scale LLM providers and downgrade models when tenant SLA is at risk.

**Timeline:** 1-2 hours | **Deployment:** agentic-prod cluster

---

## Architecture

```
┌─────────────────────────────────────────┐
│ SLA Monitor (Phase 7)                   │
├─────────────────────────────────────────┤
│ Tracks: success_rate, breach_count      │
│ Updates every workload completion       │
└────────────┬────────────────────────────┘
             │
             ↓
┌─────────────────────────────────────────┐
│ Auto-Scaler (Phase 8 NEW)               │
├─────────────────────────────────────────┤
│ Decision Logic:                         │
│  IF success_rate < 80% (warning)        │
│    → Scale LLM provider +1 replica      │
│    → Trigger circuit breaker            │
│  IF success_rate < 50% (critical)       │
│    → Scale LLM provider +2 replicas     │
│    → Downgrade model to cheaper option  │
│    → Alert ops team                     │
│  IF success_rate > 95% (healthy)        │
│    → Scale LLM provider -1 replica      │
└────────────┬────────────────────────────┘
             │
             ↓
┌─────────────────────────────────────────┐
│ Kubernetes Resource Adapter             │
├─────────────────────────────────────────┤
│ Updates Deployment replicas             │
│ Updates CRD ModelMapping (model switch) │
│ Records scaling events                  │
└─────────────────────────────────────────┘
```

---

## Implementation

### Step 1: Scaler Types

**File:** `pkg/autoscaling/types.go`

```go
type ScalingPolicy struct {
	WarningThreshold   float64 // 80% success rate
	CriticalThreshold  float64 // 50% success rate
	HealthyThreshold   float64 // 95% success rate
	MaxReplicas        int     // Cap scaling at N replicas
	MinReplicas        int     // Never scale below N
	CooldownSeconds    int     // Wait N seconds between scaling actions
	ModelDowngradeMap  map[string]string // gpt-4 → gpt-3.5, etc.
}

type ScalingEvent struct {
	TenantName      string
	TriggerType     string // warning, critical, healthy, cooldown
	Action          string // scale_up, scale_down, model_downgrade, noop
	PreviousReplicas int
	NewReplicas     int
	PreviousModel   string
	NewModel        string
	Timestamp       time.Time
	Reason          string
}
```

### Step 2: Scaler Logic

**File:** `pkg/autoscaling/scaler.go`

```go
type AutoScaler struct {
	slaMonitor *multitenancy.SLAMonitor
	policy     *ScalingPolicy
	kubeClient client.Client
	lastScaling map[string]time.Time // Track cooldown per tenant
	history    []ScalingEvent
}

func (s *AutoScaler) EvaluateAndScale(tenantName string) (*ScalingEvent, error) {
	// 1. Get SLA status
	status, err := s.slaMonitor.GetStatus(tenantName)
	if err != nil {
		return nil, err
	}
	
	// 2. Check cooldown
	if s.isInCooldown(tenantName) {
		return &ScalingEvent{
			TenantName:  tenantName,
			TriggerType: "cooldown",
			Action:      "noop",
		}, nil
	}
	
	// 3. Decide action
	if status.SuccessRatePercent < s.policy.CriticalThreshold {
		return s.scaleCritical(tenantName, status)
	}
	if status.SuccessRatePercent < s.policy.WarningThreshold {
		return s.scaleWarning(tenantName, status)
	}
	if status.SuccessRatePercent > s.policy.HealthyThreshold {
		return s.scaleDown(tenantName, status)
	}
	
	return &ScalingEvent{
		TenantName:  tenantName,
		TriggerType: "normal",
		Action:      "noop",
	}, nil
}

func (s *AutoScaler) scaleCritical(tenantName string, status *multitenancy.SLAStatus) (*ScalingEvent, error) {
	// Scale up replicas + downgrade model
	event := &ScalingEvent{
		TenantName:  tenantName,
		TriggerType: "critical",
		Action:      "scale_up_and_model_downgrade",
	}
	// Update Kubernetes deployment
	// Update CRD ModelMapping
	// Record event
	return event, nil
}
```

### Step 3: Kubernetes Integration

```go
func (s *AutoScaler) UpdateDeploymentReplicas(tenantName string, newReplicas int) error {
	// Find deployment for tenant
	// kubectl scale deployment <provider>-<tenant> --replicas=N
}

func (s *AutoScaler) UpdateModelMapping(tenantName string, newModel string) error {
	// Update AgentWorkload CRD ModelMapping
	// Apply cheaper model for cost optimization
}
```

### Step 4: Monitoring & Alerts

```go
func (s *AutoScaler) GetScalingHistory(tenantName string, limit int) []ScalingEvent {
	// Return recent scaling events for audit trail
}

func (s *AutoScaler) AlertOps(event *ScalingEvent) error {
	// Send to Slack/PagerDuty when critical threshold hit
}
```

---

## Tests

**File:** `pkg/autoscaling/scaler_test.go`

```go
func TestAutoScaler_CriticalThreshold(t *testing.T) {
	// 40% success rate → should scale critical
}

func TestAutoScaler_WarningThreshold(t *testing.T) {
	// 75% success rate → should scale warning (+1 replica)
}

func TestAutoScaler_HealthyThreshold(t *testing.T) {
	// 98% success rate → should scale down (-1 replica)
}

func TestAutoScaler_Cooldown(t *testing.T) {
	// Scale up, then try again immediately → should noop due to cooldown
}

func TestAutoScaler_ModelDowngrade(t *testing.T) {
	// Critical threshold → should downgrade gpt-4 to gpt-3.5
}
```

---

## Deployment

```bash
# 1. Create scaling policy ConfigMap
kubectl create configmap scaling-policy \
  -n agentic-system \
  --from-literal=warning-threshold=0.80 \
  --from-literal=critical-threshold=0.50 \
  --from-literal=healthy-threshold=0.95 \
  --from-literal=cooldown-seconds=300

# 2. Update controller with AutoScaler
# (see controller integration below)

# 3. Monitor scaling events
kubectl logs -n agentic-system -l app=agentic-operator | grep "scaling"
```

---

## Success Criteria

✅ Evaluate SLA status every reconciliation  
✅ Scale LLM provider replicas based on thresholds  
✅ Downgrade models when SLA critical  
✅ Cooldown prevents scaling thrashing  
✅ Scaling history tracked (audit trail)  
✅ Alerts sent on critical breach  
✅ All tests passing  

---

## Integration Points

1. **SLA Monitor** (Phase 7) — Provides success rate data
2. **Controller** — Calls AutoScaler after SLA recording
3. **Kubernetes API** — Updates Deployment + CRD
4. **Observability** — Tracks scaling events as metrics

---

**Status:** Ready to implement 🚀
