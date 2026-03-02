# Model Routing Observability with OpenTelemetry

This document describes how to observe and debug model routing decisions using OpenTelemetry tracing.

## Overview

Model routing in Phase 3 adds OpenTelemetry instrumentation to provide complete visibility into:
- Task classification decisions (validation/analysis/reasoning)
- Provider and model selection
- API call execution
- Token usage and costs
- Error propagation

## Tracing Architecture

### Span Hierarchy

Every model routing operation creates a hierarchical span structure:

```
model.routing (root span)
├── task.classification
│   └── [classification result: validation/analysis/reasoning]
├── provider.resolution
│   └── [provider lookup from config + secret resolution]
└── model.call
    └── [actual API invocation + token counting]
```

### Span Attributes

**Root Span (model.routing):**
- `workload.name` - Name of the AgentWorkload CR
- `workload.namespace` - Kubernetes namespace
- `task.category` - Classified task type
- `provider.name` - Selected provider
- `model.name` - Selected model
- `tokens.input` - Input tokens used
- `tokens.output` - Output tokens used

**Classification Span (task.classification):**
- `instruction.length` - Length of the task instructions
- `task.category` - Classification result

**Provider Resolution Span (provider.resolution):**
- `provider.name` - Provider identifier
- `provider.type` - Type (openai-compatible, workers-ai, etc.)
- `provider.endpoint` - API endpoint URL

**Model Call Span (model.call):**
- `provider.name` - Provider used
- `model.name` - Model called
- `tokens.input` - Input tokens consumed
- `tokens.output` - Output tokens generated
- `call.success` - True if API call succeeded

### Span Events

Spans record events for key decisions and errors:

**Validation Events:**
- `validation_failed` - ModelMapping or providers invalid

**Classification Events:**
- `task_classified` - Task type determined

**Routing Events:**
- `routing_failed` - Task mapping or provider lookup failed
- `provider_not_found` - Provider not in config
- `provider_init_failed` - Credential resolution failed
- `model_call_failed` - API call failed
- `routing_completed` - Routing succeeded

## Setup & Configuration

### Prerequisites

- OpenTelemetry compatible collector (optional, can export directly)
- Loki for log-based trace storage
- Grafana for visualization
- Prometheus for metrics

### Enable Tracing

Traces are **automatically enabled** when OpenTelemetry SDK is initialized. No configuration needed.

### Export to Loki

If you're already using the Loki stack in agentic-prod DOKS cluster:

```yaml
# In operator helm values
opentelemetry:
  enabled: true
  loki:
    enabled: true
    endpoint: "http://loki.logging:3100"
```

### Export to OpenTelemetry Collector

For more advanced setups with Tempo or Jaeger:

```yaml
opentelemetry:
  enabled: true
  exporter:
    type: otlp
    endpoint: "otel-collector.monitoring:4317"
    samplingRate: 0.1  # Sample 10% in production
```

## Viewing Traces

### In Grafana Loki UI

**Example query to see all routing traces:**
```logql
{job="agentic-operator"} | json | routing_decision != ""
```

**Filter by task category:**
```logql
{job="agentic-operator"} | json | task_category="analysis"
```

**Find failed routing attempts:**
```logql
{job="agentic-operator"} | json | span_name="routing_failed"
```

**See slow API calls (>5s):**
```logql
{job="agentic-operator", span_name="model.call"} | json | duration_ms > 5000
```

### In Grafana Traces UI

Create a new Trace panel pointing to your trace backend:

1. **Data Source:** Select your trace backend (Loki, Tempo, Jaeger)
2. **Trace Selector:** 
   - Service: `agentic-operator`
   - Operation: `model.routing`
3. **View Details:**
   - Expand spans to see attributes
   - Check events for decision points
   - Review durations for performance

### Example: Trace a Market Intelligence Workload

1. Open Grafana Loki
2. Query:
   ```logql
   {job="agentic-operator", workload="market-intelligence"} | json
   ```
3. Click a trace to view the span hierarchy
4. Inspect each span:
   - **task.classification:** See what category was assigned
   - **provider.resolution:** Confirm which provider was selected
   - **model.call:** Check token usage and API latency

## Debugging with Traces

### Problem: Tasks routed to wrong model

**Solution:** Check the task.classification span:
- Look at `task.category` attribute
- Inspect `instruction.length`
- Review classifier logic if classification wrong

### Problem: API calls failing

**Solution:** Check the model.call span:
- Look for `model_call_failed` event
- Check error message in event attributes
- Verify provider credentials in secret

### Problem: High token usage

**Solution:** Review tokens across spans:
- Check `tokens.input` and `tokens.output` in root span
- Compare with similar tasks
- Identify expensive model selections

### Problem: Slow routing decisions

**Solution:** Check span durations:
- `task.classification` should be <10ms
- `provider.resolution` should be <100ms
- `model.call` depends on API (typically 500ms-5s)

## Metrics from Traces

Traces complement existing Prometheus metrics:

| Metric | Source | Use |
|--------|--------|-----|
| `agentic_model_routing_total` | Prometheus | Count by category/provider/model |
| `agentic_tokens_used_total` | Prometheus | Token consumption by provider |
| `agentic_estimated_cost_usd` | Prometheus | Estimated API costs |
| Trace spans | OpenTelemetry | Individual request details |

**Example:** Use metrics for dashboards, traces for debugging individual workloads.

## Performance Impact

OpenTelemetry tracing is **zero-cost when no samples are collected**:

- Sampler rate 0.0 (no sampling) = negligible overhead
- Sampler rate 0.01 (1% sampling) = <1% overhead
- Sampler rate 1.0 (all) = <5% overhead

**Recommendation:**
- Development: `samplingRate: 1.0` (see all traces)
- Production: `samplingRate: 0.01` (1% sampling for cost control)

## Real-World Example

**Workload Definition:**
```yaml
apiVersion: agentic.clawdlinux.org/v1alpha1
kind: AgentWorkload
metadata:
  name: market-analysis
spec:
  modelStrategy: cost-aware
  objective: "Analyze Q1 2026 tech sector performance"
  modelMapping:
    validation: openai/gpt-3.5-turbo
    analysis: openai/gpt-4
    reasoning: openai/gpt-4-turbo
```

**Expected Trace Flow:**

1. **model.routing started** (root span)
   - `workload.name: market-analysis`
   - `workload.namespace: default`

2. **task.classification**
   - Input length: 45 chars
   - Keyword detected: "analyze"
   - Result: `task.category: analysis`

3. **provider.resolution**
   - Provider: `openai`
   - Endpoint: `https://api.openai.com/v1`
   - Secret retrieved: ✓

4. **model.call**
   - Model: `gpt-4`
   - Input tokens: 87
   - Output tokens: 342
   - Duration: 2.3s
   - Success: true

5. **model.routing completed**
   - Total duration: 2.4s
   - Estimated cost: $0.015

**Visible in Loki:**
```
task="market-analysis" task_category="analysis" 
provider="openai" model="gpt-4" 
tokens_in=87 tokens_out=342 
cost_usd=0.015 duration_ms=2400
```

## Troubleshooting

### Traces not appearing in Loki

1. Check OpenTelemetry exporter is initialized
2. Verify Loki endpoint is reachable: `curl http://loki:3100/api/v1/labels`
3. Check operator logs for trace errors: `kubectl logs -f deployment/agentic-operator`

### Missing attributes in traces

1. Verify span was created (check span names in Loki)
2. Check if attributes are being set correctly
3. Verify no sampling is filtering out traces during test

### High trace volume

1. Reduce sampling rate: `samplingRate: 0.01`
2. Enable Loki retention: `retention: 7d`
3. Archive old traces to storage backend

## Integration with Existing Stack

Model routing traces integrate seamlessly with agentic-prod DOKS monitoring:

- **Prometheus:** Routing metrics already exported
- **Grafana:** Use Loki datasource to view traces
- **Loki:** Receives structured trace data
- **Velero:** Archives long-term trace data if needed

No changes needed to existing monitoring setup.

## Next Steps

1. Deploy with `helm install` using `helm-with-tracing.yaml`
2. Create a test AgentWorkload with cost-aware routing
3. View traces in Grafana Loki
4. Create custom dashboards for your use cases
5. Set up alerts based on trace patterns

## Further Reading

- [OpenTelemetry Documentation](https://opentelemetry.io/)
- [Loki Log Querying](https://grafana.com/docs/loki/latest/logql/)
- [Grafana Traces Panel](https://grafana.com/docs/grafana/latest/panels/visualizations/traces/)
