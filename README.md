# Agentic Kubernetes Operator - PoC Development

**A production-grade Kubernetes operator for orchestrating durable AI agent swarms with LangGraph checkpointing, Argo Workflows DAG execution, and human-in-the-loop reviews.**

**Status:** 🟡 Phase 1 - Foundation (Repository cleaned, ready to scaffold)
**Updated:** 2026-02-23
**Budget:** $200 / 4+ months on k3s ($48/mo)

---

## What We're Building

### The Problem
AI agents on Kubernetes are ephemeral. Pod crashes mid-execution = lost state. No native orchestration for multi-step AI workflows. No built-in durable execution.

### The Solution
An operator that:
1. **Manages AgentWorkload resources** - Declare agent jobs as Kubernetes objects
2. **Executes with durable checkpointing** - Resume from exact point after pod crashes via LangGraph
3. **Orchestrates DAGs** - Multi-step workflows (scrape → analyze → review → report)
4. **Provides visibility** - Full trace tracing, cost tracking, human review gates
5. **Enforces safety** - RBAC isolation, network policies, license enforcement
6. **Works on $48/month** - Single k3s droplet, API-mode LLMs, no GPUs needed

### The Use Case
**Visual Market Analysis:** Scrape 50 competitor websites in parallel → analyze DOM structure + screenshots → generate strategic report with human review gate.

**Cost:** ~$0.001 per analysis (GPT-4o Mini) = $10.50 for 10,000 analyses/month

---

## How We're Building It (Methodology)

### The Approach
- **No custom autonomous dev systems** - Just verified OpenClaw skills
- **Incremental phases** - Build one component, test thoroughly, commit, move on
- **Brainstorm first** - Free models (Gemini/Kimi) + Opus for critical decisions
- **Secure skills only** - Everything from clawhub.ai / skills.sh
- **Real testing** - Each phase tested on actual k3s before moving forward

### The Loop
```
Brainstorm → Plan → Build → Test → Commit → Repeat
```

**Brainstorm:** Free models + Opus, document findings
**Plan:** Find secure skills, define success criteria
**Build:** One component, only trusted skills
**Test:** Real k3s cluster, validate functionality
**Commit:** Working code to GitHub
**Repeat:** Next phase

---

## Technology Stack (Rationale)

| Component | Choice | Why | Rationale |
|-----------|--------|-----|-----------|
| **Agent Framework** | LangGraph v1.0 | Durable checkpointing at every node | Kubernetes pods are ephemeral; LangGraph survives OOMKill/preemption |
| **Operator** | Kubebuilder v4.11 | Production Go scaffolding | kubernetes-sigs standard, generates secure RBAC patterns |
| **Workflows** | Argo Workflows v4.0 | Native DAG + suspend/resume | Human review gates, artifact passing, CNCF Graduated |
| **LLM Access** | LiteLLM | 100+ provider routing | Vendor-agnostic, automatic fallback, cost tracking |
| **Browser Pool** | Browserless + Playwright | Centralized + HPA | vs sidecars: 3x more efficient, shared cache, independent scaling |
| **Observability** | Langfuse + OpenTelemetry | MIT self-hosted | Full trace visibility, cost per analysis, MIT licensed (no vendor lock) |
| **Infrastructure** | k3s on DO 8GiB | Single droplet | $48/mo, 6.4 GiB usable, 4.2 months on $200 budget |

---

## Current Phase: Phase 1 - Foundation

### What We're Building Now
- AgentWorkload Kubernetes operator (Go, Kubebuilder v4.11)
- Basic CRD that accepts agent job specs
- RBAC framework (operator, agents, browser pool isolated)
- Webhook validation for AgentWorkload resources
- Build passes, deploys to k3s, accepts sample resources

### Architecture for Phase 1
```
┌────────────────────────────────────────────┐
│ Operator Namespace                         │
├────────────────────────────────────────────┤
│ AgentWorkload CRD                          │
│ ├─ Spec: image, replicas, resources       │
│ ├─ Status: ready replicas, phase, events  │
│ └─ Validation: webhook checks              │
│                                            │
│ Operator Reconciler                       │
│ ├─ Create/Update Argo Workflows           │
│ ├─ Track Job completion                   │
│ └─ Update AgentWorkload status            │
│                                            │
│ RBAC + Webhook + Finalizers               │
└────────────────────────────────────────────┘
```

### Success Criteria for Phase 1
- [ ] Kubebuilder project scaffolded and builds without errors
- [ ] AgentWorkload CRD defined (v1alpha1) with proper spec/status fields
- [ ] Operator deploys to k3s without crashes
- [ ] Sample AgentWorkload resource creates and updates without errors
- [ ] Webhook validates invalid resources (rejects bad requests)
- [ ] RBAC ClusterRole grants only necessary permissions
- [ ] Finalizers prevent premature deletion
- [ ] Code compiles, lints, and follows Kubebuilder patterns
- [ ] Documented in code comments and README

### Testing for Phase 1
1. **Local:** `make test` passes all unit tests
2. **Integration:** Deploy to kind cluster, create sample resources
3. **Real k3s:** Deploy to actual k3s droplet, verify operator works
4. **Validation:** Webhook rejects invalid specs correctly

### Deliverables for Phase 1
- Kubebuilder project in GitHub repo
- AgentWorkload CRD YAML
- Operator source code (Go)
- Unit tests + integration tests
- Deployment guide (kind)

---

## Improvement Tracking

### How We Know We're On Track
1. **Code commits** - Working code pushed after each phase
2. **Tests passing** - 100% test pass rate before moving to next phase
3. **Real k3s validation** - Not just local testing
4. **Documentation** - Clear README + code comments
5. **Budget tracking** - Staying within $48/month infrastructure cost

### How We Improve (Feedback Loop)

#### After Each Phase
- [ ] **What worked?** Document successes (architecture decisions, patterns)
- [ ] **What was hard?** Identify pain points, challenges, gotchas
- [ ] **What would we do differently?** Plan improvements for next phase
- [ ] **Update roadmap** - Adjust timeline if needed
- [ ] **Update README** - Keep this document current

#### Red Flags (Stop and Rethink)
- ❌ Tests failing but pushing anyway
- ❌ Untested code in GitHub
- ❌ Budget exceeded before Phase 2
- ❌ Can't explain why a decision was made
- ❌ Diverging from the stack (using unauthorized tools)
- ❌ Not testing on real k3s

---

## Phase Breakdown

### Phase 1: Foundation (1-2 weeks)
**Goal:** Kubebuilder operator scaffolding works, deploys, handles AgentWorkload CRDs

**Decisions to Make:**
- AgentWorkload spec fields (image? replicas? resources? suspend?)
- Status tracking (phase enum: Pending/Creating/Running/Failed?)
- CRD versioning (v1alpha1 until stable)
- Reconciliation interval (30s polling or watch-based?)

**Skills Needed:**
- Go + Kubebuilder (secure skill)
- Kubernetes API design (secure skill)
- CRD validation patterns

**Success = operator deployable + webhook validates + sample resource works**

### Phase 2: Agent Execution (1-2 weeks)
**Goal:** Operator creates and manages Argo Workflows, agents execute in pods

**Decisions:**
- Argo Workflow template structure (DAG vs sequential?)
- Pod resource limits (500m CPU? 512Mi RAM?)
- Agent image specification (LangGraph base image?)
- Status propagation (how does Argo status → AgentWorkload status?)

**Success = operator creates Argo workflows, agents execute, status tracked**

### Phase 3: Observability (1 week)
**Goal:** Langfuse tracing, OpenTelemetry instrumentation, structured logging

### Phase 4: Browser & LLM (1-2 weeks)
**Goal:** Browserless pool integration, LiteLLM proxy routing

### Phase 5: Production (2-3 weeks)
**Goal:** JWT licensing, usage metering, Helm charts, security hardening

---

## How to Use This README

**For Development:**
- Start of each phase: Read the Phase section, update success criteria, execute
- During development: Check against success criteria
- End of phase: Mark complete, document learnings, move to next phase

**For Tracking:**
- [ ] = Not started
- [x] = Completed
- Update status frequently
- Keep dates current

**For Improvement:**
- Every Friday: Review what worked/didn't work
- Update "Improvement Tracking" section
- Adjust roadmap for next phase

---

## Architecture Diagrams

Real architecture research in `/docs/diagrams/`:
- `system-overview.png` - Full system architecture
- `operator-controlplane.png` - Operator reconciliation loop
- `agent-dag-pipeline.png` - Agent execution DAG
- `infra-services-test-droplet.png` - Infrastructure on single droplet
- `packaging-distribution.png` - Helm charts + distribution

See `/docs/` for detailed architecture documentation.

---

## Key Architectural Principles

### 1. Durable Execution
LangGraph checkpoints at every graph node. When a pod crashes, resume from exact point — no lost state.

### 2. Cost Efficiency
API-mode LLMs (GPT-4o Mini $0.15/$0.60 per million tokens) vs GPU infrastructure. Only run local vLLM above 100K requests/month.

### 3. Single Droplet MVP
k3s on 8 GiB DigitalOcean = $48/month = 4.2 months on $200 budget. Scales to full cluster later.

### 4. Secure Patterns
- RBAC: operator gets API access, agents get none
- Network: Cilium DNS egress policies per job
- Secrets: External Secrets Operator syncs from Vault
- Validation: CEL rules in CRD + webhook

### 5. Human-in-the-Loop
Argo Workflows suspend nodes enable approval gates before final report generation.

---

## Contributing

**Requirements:**
- Use only verified skills from clawhub.ai / skills.sh
- Test on real k3s before committing
- Document architecture decisions in code comments
- Update this README when phase completes
- No untrusted custom code

**Process:**
1. Pick a phase goal from above
2. Brainstorm with free models + Opus
3. Find matching secure skills
4. Build one component
5. Test on k3s
6. Commit to GitHub
7. Update README with results

---

## Security

- **Network Isolation:** Cilium DNS-based FQDN egress (only whitelisted domains)
- **Pod Security:** Restricted PSS profile + custom seccomp for Chrome
- **RBAC:** Least-privilege ServiceAccounts
- **Secrets:** External Secrets Operator syncing
- **Audit:** Full OpenTelemetry trace logging

---

## Budget Tracking

- **Monthly:** $48 (k3s on 8 GiB DO droplet)
- **Total Budget:** $200
- **Runway:** 4.2 months
- **Optimization:** Per-second billing means destroy droplet when not developing
- **Actual Cost** (working 8h/day weekdays): ~$11/month = 18+ month runway

---

## FAQ

**Q: Why not use CrewAI instead of LangGraph?**
A: CrewAI is great for role-based agents but lacks durable execution. Kubernetes pods crash frequently; LangGraph checkpoints survive crashes.

**Q: Why Browserless instead of Playwright sidecars?**
A: Centralized pool is 3× more resource-efficient, enables session reuse, scales independently. Sidecars waste resources.

**Q: Why not run vLLM locally?**
A: Not cost-effective below 100K requests/month. GPT-4o Mini API is cheaper and simpler.

**Q: How long until production?**
A: 6-10 weeks (5 phases × 1-3 weeks each). Incremental, tested approach = higher quality.

---

## Status Tracking

Last Updated: 2026-02-23 18:08 IST
Phase: 1 - Foundation
Progress: Repository cleaned, ready to scaffold
Next: Brainstorm Phase 1 architecture with Opus + secure Kubebuilder skill
