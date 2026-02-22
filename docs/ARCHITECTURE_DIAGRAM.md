# Architecture Diagrams

## High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AUTONOMOUS DEV SYSTEM                            │
│                                                                          │
│  ┌────────────────┐      ┌────────────────┐      ┌─────────────────┐   │
│  │   Task Queue   │      │  Orchestrator  │      │  GitHub Repo    │   │
│  │  (PostgreSQL)  │◄────►│   (main.py)    │─────►│  agentic-k8s-op │   │
│  └────────────────┘      └────────────────┘      └─────────────────┘   │
│         │                         │                                     │
│         │                         │                                     │
│         │                  ┌──────▼──────┐                              │
│         │                  │  LangGraph  │                              │
│         │                  │  Workflow   │                              │
│         │                  └──────┬──────┘                              │
│         │                         │                                     │
│         │          ┌──────────────┼──────────────┐                      │
│         │          │              │              │                      │
│    ┌────▼────┐ ┌───▼────┐ ┌──────▼─────┐ ┌─────▼─────┐                │
│    │PM Agent │ │Architect│ │Engineer    │ │QA Agent   │                │
│    │(Sonnet) │ │(Opus)★  │ │(Sonnet)    │ │(Sonnet)   │                │
│    └────┬────┘ └───┬────┘ └──────┬─────┘ └─────┬─────┘                │
│         │          │              │              │                      │
│         └──────────┼──────────────┼──────────────┘                      │
│                    │              │                                     │
│             ┌──────▼──────┐ ┌─────▼────────┐                           │
│             │Security     │ │GitManager    │                           │
│             │(Sonnet)     │ │(Python)      │                           │
│             └──────┬──────┘ └─────┬────────┘                           │
│                    │              │                                     │
│             ┌──────▼──────────────▼─────┐                              │
│             │  Consensus Voting          │                              │
│             │  (80% threshold)           │                              │
│             └────────────────────────────┘                              │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## LangGraph Agent Pipeline

```
                        Task from Queue
                              │
                              ▼
                    ┌─────────────────┐
                    │   PM Agent      │
                    │  (Sonnet 4.5)   │
                    │                 │
                    │ • Analyze task  │
                    │ • Write PRD     │
                    │ • Define goals  │
                    └────────┬────────┘
                             │
                    State: {task, prd}
                             │
                             ▼
                    ┌─────────────────┐
                    │ Architect Agent │
                    │  (Opus 4.5) ★   │
                    │                 │
                    │ • Design system │
                    │ • Components    │
                    │ • Score: 90/100 │
                    └────────┬────────┘
                             │
              State: {prd, architecture}
                             │
                             ▼
                    ┌─────────────────┐
                    │ Engineer Agent  │
                    │  (Sonnet 4.5)   │
                    │                 │
                    │ • Write Go code │
                    │ • Kubebuilder   │
                    │ • Score: 85/100 │
                    └────────┬────────┘
                             │
           State: {architecture, implementation}
                             │
                             ▼
                    ┌─────────────────┐
                    │   QA Agent      │
                    │  (Sonnet 4.5)   │
                    │                 │
                    │ • Review code   │
                    │ • Write tests   │
                    │ • Score: 80/100 │
                    └────────┬────────┘
                             │
              State: {implementation, tests}
                             │
                             ▼
                    ┌─────────────────┐
                    │ Security Agent  │
                    │  (Sonnet 4.5)   │
                    │                 │
                    │ • Risk scoring  │
                    │ • Baseline 70%  │
                    │ • Score: 75/100 │
                    └────────┬────────┘
                             │
      State: {all, consensus_votes: [90,85,80,75]}
                             │
                  Consensus: (90+85+80+75)/4 = 82.5%
                             │
                      [Threshold: 80%] → APPROVED ✅
                             │
                             ▼
                    ┌─────────────────┐
                    │  GitManager     │
                    │   (Python)      │
                    │                 │
                    │ • Create branch │
                    │ • Write files   │
                    │ • Commit & push │
                    │ • Create PR     │
                    └────────┬────────┘
                             │
                             ▼
                    GitHub PR Created
              (agentic-k8s-operator repo)
```

## Data Flow

```
┌──────────────┐
│ load_tasks.py│
│              │
│ • Phase 1    │
│ • 15 tasks   │
└──────┬───────┘
       │
       ▼
┌────────────────────────────────┐
│     PostgreSQL Database         │
│                                 │
│  tasks table:                   │
│  ┌──────────────────────────┐  │
│  │ id | task_id | status    │  │
│  ├──────────────────────────┤  │
│  │ 1  | p1-001  | pending   │  │
│  │ 2  | p1-002  | pending   │  │
│  │ 3  | p1-003  | in_prog   │  │
│  │ 4  | p1-004  | completed │  │
│  └──────────────────────────┘  │
│                                 │
│  • priority DESC                │
│  • created_at ASC               │
└────────┬───────────────────────┘
         │
         │ Poll (60s interval)
         ▼
┌──────────────────────┐
│  main.py             │
│  (Orchestrator)      │
│                      │
│  while True:         │
│    task = get_next() │
│    process(task)     │
│    sleep(60)         │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ graph/workflow.py    │
│                      │
│ workflow.invoke({    │
│   task_id,           │
│   task,              │
│   repo_path          │
│ })                   │
└──────┬───────────────┘
       │
       │ Sequential execution
       ▼
┌─────────────────────────────┐
│  6 Agents Execute            │
│  (LangGraph state machine)   │
│                              │
│  State flows:                │
│  {} → {prd} → {arch} →       │
│  {impl} → {tests} →          │
│  {security} → {consensus}    │
└──────┬──────────────────────┘
       │
       │ If consensus ≥ 80%
       ▼
┌──────────────────────┐
│  agents/git_manager.py│
│                      │
│  • Parse files       │
│  • git checkout -b   │
│  • git add/commit    │
│  • git push origin   │
│  • gh pr create      │
└──────┬───────────────┘
       │
       ▼
┌────────────────────────────┐
│   GitHub Repository         │
│   agentic-k8s-operator     │
│                            │
│   Pull Requests:           │
│   ┌──────────────────────┐ │
│   │ #1 [Autonomous]      │ │
│   │ Consensus: 85%       │ │
│   │ Labels: autonomous   │ │
│   └──────────────────────┘ │
│   ┌──────────────────────┐ │
│   │ #2 [Autonomous]      │ │
│   │ Consensus: 82%       │ │
│   └──────────────────────┘ │
└────────────────────────────┘
```

## Deployment Infrastructure

```
┌──────────────────────────────────────────────┐
│         Local Mac (Development)               │
│                                               │
│  ~/.openclaw/workspace/autonomous-dev-system/ │
│    ├─ agents/                                 │
│    ├─ graph/                                  │
│    ├─ db/                                     │
│    └─ deploy.sh                               │
└─────────────────┬────────────────────────────┘
                  │
                  │ SCP upload
                  │
                  ▼
┌──────────────────────────────────────────────────────┐
│     DigitalOcean Droplet (167.71.228.21)             │
│                                                      │
│  ┌────────────────────────────────────────────┐     │
│  │  /root/autonomous-dev-system/              │     │
│  │    ├─ venv/ (Python 3.12)                  │     │
│  │    ├─ agents/                              │     │
│  │    ├─ graph/                               │     │
│  │    ├─ db/                                  │     │
│  │    ├─ main.py                              │     │
│  │    └─ .env                                 │     │
│  └────────────────────────────────────────────┘     │
│                                                      │
│  ┌────────────────────────────────────────────┐     │
│  │  PostgreSQL (localhost:5432)               │     │
│  │    └─ autonomous_dev database              │     │
│  └────────────────────────────────────────────┘     │
│                                                      │
│  ┌────────────────────────────────────────────┐     │
│  │  /root/agentic-k8s-operator/               │     │
│  │    └─ Git clone of target repo             │     │
│  └────────────────────────────────────────────┘     │
│                                                      │
│  ┌────────────────────────────────────────────┐     │
│  │  Systemd Service                           │     │
│  │  /etc/systemd/system/autonomous-dev.service│     │
│  │    ExecStart: /root/.../main.py            │     │
│  │    Restart: always                         │     │
│  └────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────┘
```

## Security Agent Risk Scoring

```
                  Base Score: 70
                  (Assume secure cluster)
                        │
        ┌───────────────┼───────────────┐
        │                               │
   DEDUCTIONS                       BONUSES
        │                               │
        ▼                               ▼
┌──────────────────┐           ┌──────────────────┐
│ Privilege esc    │           │ Encryption       │
│   -25            │           │   +10            │
├──────────────────┤           ├──────────────────┤
│ Destructive ops  │           │ Audit logging    │
│   -30            │           │   +10            │
├──────────────────┤           ├──────────────────┤
│ External API     │           │ Secret mgmt      │
│   -10            │           │   +5             │
├──────────────────┤           └──────────────────┘
│ No audit log     │
│   -10            │
├──────────────────┤
│ Error exposure   │
│   -15            │
└──────────────────┘
        │
        ▼
  Final Score = 70 - deductions + bonuses
        │
        ▼
┌──────────────────────────────────┐
│  Decision Logic:                 │
│                                  │
│  ≥80  → APPROVE                  │
│  60-79 → CONDITIONAL_APPROVE     │
│  <60  → REJECT                   │
└──────────────────────────────────┘
```

## Task State Machine

```
     ┌─────────┐
     │ pending │
     └────┬────┘
          │
          │ get_next_task()
          │
          ▼
   ┌──────────────┐
   │ in_progress  │
   └──────┬───────┘
          │
          │ Workflow execution
          │
     ┌────┴────┐
     │         │
     ▼         ▼
┌─────────┐  ┌────────┐
│completed│  │ failed │
└─────────┘  └────────┘
     │            │
     │            │
     ▼            ▼
  (Archive)   (Archive)
```

## Model Cost Breakdown (Per Task)

```
Task Execution
     │
     ├─ PM Agent (Sonnet 4.5)
     │    └─ ~2K tokens in, 1K tokens out
     │    └─ Cost: ~$0.09
     │
     ├─ Architect Agent (Opus 4.5) ★
     │    └─ ~3K tokens in, 2K tokens out
     │    └─ Cost: ~$0.30
     │
     ├─ Engineer Agent (Sonnet 4.5)
     │    └─ ~4K tokens in, 3K tokens out
     │    └─ Cost: ~$0.21
     │
     ├─ QA Agent (Sonnet 4.5)
     │    └─ ~5K tokens in, 2K tokens out
     │    └─ Cost: ~$0.21
     │
     ├─ Security Agent (Sonnet 4.5)
     │    └─ ~5K tokens in, 1K tokens out
     │    └─ Cost: ~$0.18
     │
     └─ GitManager (No LLM)
          └─ Cost: $0
─────────────────────────────────
Total per task: ~$0.99

Phase 1 (15 tasks): ~$15
```

## Monitoring Dashboard (Planned)

```
┌─────────────────────────────────────────┐
│   Autonomous Dev Dashboard              │
├─────────────────────────────────────────┤
│                                         │
│  Queue Status:                          │
│    Pending:      10 tasks               │
│    In Progress:  1 task                 │
│    Completed:    4 tasks (80% avg)      │
│    Failed:       0 tasks                │
│                                         │
│  Current Task:                          │
│    ID:     phase1-005                   │
│    Agent:  QA Agent                     │
│    Score:  78% (in progress)            │
│                                         │
│  Recent PRs:                            │
│    #4 - Consensus: 85% ✅               │
│    #3 - Consensus: 82% ✅               │
│    #2 - Consensus: 90% ✅               │
│    #1 - Consensus: 78% ⚠️                │
│                                         │
│  Cost Tracking:                         │
│    Today:   $3.50                       │
│    Week:    $12.00                      │
│    Month:   $28.00                      │
│                                         │
└─────────────────────────────────────────┘
```

---

**Generated**: 2026-02-22  
**Version**: 1.0  
**Status**: Production-ready
