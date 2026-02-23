# Agentic Kubernetes Operator

**Production-grade Kubernetes operator for orchestrating AI agent swarms using LangGraph, Argo Workflows, and Kubebuilder.**

## Vision

Enable safe, durable execution of containerized AI agent workloads on Kubernetes with automatic checkpointing, human-in-the-loop review, and native DAG orchestration.

## Technology Stack (PoC)

| Component | Choice | Why |
|-----------|--------|-----|
| **Agent Framework** | LangGraph v1.0 | Durable execution with automatic checkpointing (MIT license) |
| **Operator** | Kubebuilder v4.11 | Production-grade Go scaffolding (kubernetes-sigs) |
| **Workflow Orchestration** | Argo Workflows v4.0 | Native DAG templates, suspend/resume, artifact passing (CNCF Graduated) |
| **Browser Automation** | Browserless + Playwright | Centralized browser pool with HPA autoscaling |
| **LLM Access** | LiteLLM | Unified OpenAI-compatible endpoint (100+ providers) |
| **Observability** | Langfuse + OpenTelemetry | MIT-licensed, self-hosted, complete trace visibility |
| **Infrastructure** | k3s on DigitalOcean | Single 8 GiB droplet ($48/month), 4.2 months on $200 budget |

## Architecture Overview

```
┌─────────────── k3s Cluster ($48/mo) ──────────────┐
│                                                    │
│  Operator Namespace                Agent Namespace│
│  ┌──────────────┐              ┌──────────────────┤
│  │AgentWorkload │──────────────▶│Argo Workflow     │
│  │  Operator    │              │ (DAG)            │
│  │(Kubebuilder) │              ├─ Scraper Jobs   │
│  └──────────────┘              ├─ Analyzer Jobs  │
│                                ├─ Suspend Node   │
│  JWT License ✓                 └─ Report Job     │
│  OpenMeter ✓                                     │
│                                                  │
│  Browser Pool (Browserless/Steel)                │
│  LiteLLM Proxy → OpenAI/Anthropic/Local vLLM    │
│  Langfuse + OTel Collector                       │
│                                                  │
└─────────────────────────────────────────────────┘
```

## PoC Development Plan

### Phase 1: Foundation
- [ ] Clean repository structure
- [ ] Kubebuilder operator scaffolding
- [ ] AgentWorkload CRD definition
- [ ] Basic RBAC and webhook setup

### Phase 2: Agent Execution
- [ ] LangGraph agent pod templates
- [ ] Argo Workflow integration
- [ ] Basic reconciliation loop

### Phase 3: Observability
- [ ] Langfuse deployment
- [ ] OpenTelemetry instrumentation
- [ ] Structured logging

### Phase 4: Browser & LLM
- [ ] Browserless centralized pool
- [ ] LiteLLM proxy integration
- [ ] Vision model support

### Phase 5: Production
- [ ] JWT license enforcement
- [ ] OpenMeter usage metering
- [ ] Helm 4 umbrella charts
- [ ] Network policies (Cilium)
- [ ] Pod security standards

## Key Architectural Decisions

### 1. Durable Execution (LangGraph Checkpointing)
LangGraph automatically checkpoints at every graph node. When a Kubernetes pod gets OOMKilled or preempted mid-execution, the workflow resumes from exactly where it stopped — no state loss.

### 2. Centralized Browser Pool
One Browserless/Steel deployment with HPA autoscaling, not browser sidecars per pod. More resource-efficient, enables session reuse, scales independently.

### 3. API-Mode LLM Inference
GPT-4o Mini at $0.15/$0.60 per million tokens (~$0.001 per analysis). Only cost-effective to run local vLLM above 100K requests/month.

### 4. Single Droplet MVP
k3s on 8 GiB DigitalOcean droplet provides 6.4 GiB usable RAM after k3s overhead — sufficient for operator, agents, Argo, browser pool, PostgreSQL, and Langfuse.

### 5. Secure Skills Only
All tooling via verified OpenClaw skills from clawhub.ai and skills.sh — no custom untrusted code.

## Quick Start (When Ready)

```bash
# 1. Initialize Kubebuilder operator
kubebuilder init --domain ai.example.com --repo github.com/org/agentic-operator

# 2. Create AgentWorkload API
kubebuilder create api --group agents --version v1alpha1 --kind AgentWorkload

# 3. Set up Argo Workflows
helm repo add argo https://argoproj.github.io/argo-helm
helm install argo-workflows argo/argo-workflows -n argo

# 4. Deploy Browserless
helm install browserless -n browsers

# 5. Deploy LiteLLM
helm install litellm ...
```

## Security

- **Network Isolation:** Cilium DNS-based FQDN egress policies per agent
- **Pod Security:** Restricted PSS profile + custom seccomp for browser containers
- **RBAC:** Least-privilege ServiceAccounts, operator gets API access, agents get none
- **Secrets:** External Secrets Operator syncs API keys from Vault/AWS Secrets Manager

## Documentation

- `/docs/diagrams/` - Architecture visualizations
- `ARCHITECTURE.md` - Detailed system design (coming)
- `DEVELOPMENT.md` - Development guide (coming)

## License

MIT (operator code) + Apache 2.0 (dependencies)

## Status

🟡 **PoC Phase 1: Foundation** - Repository cleaned, ready for development
