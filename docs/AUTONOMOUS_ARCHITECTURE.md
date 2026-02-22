# Autonomous System Architecture

**Complete architecture of the autonomous development system (February 2026)**

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    AUTONOMOUS DEV SYSTEM                         │
│                   (DigitalOcean Droplet)                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
        ┌────────────────────────────────────────┐
        │      PostgreSQL Task Queue              │
        │  - tasks table (pending/in_progress)    │
        │  - Prioritized backlog (15 Phase 1)     │
        └────────────────────────────────────────┘
                              │
                              ▼
        ┌────────────────────────────────────────┐
        │   Main Orchestrator (main.py)          │
        │   - Polls queue every 60s               │
        │   - Systemd service: autonomous-dev     │
        └────────────────────────────────────────┘
                              │
                              ▼
        ┌────────────────────────────────────────┐
        │      LangGraph Workflow                 │
        │      (graph/workflow.py)                │
        └────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
        ┌──────────────────┐  ┌──────────────────┐
        │  Agent Pipeline   │  │  Git Manager     │
        │  (6 agents)       │  │  (PR creation)   │
        └──────────────────┘  └──────────────────┘
                    │                   │
                    ▼                   ▼
        ┌──────────────────┐  ┌──────────────────┐
        │ Consensus Voting  │  │ GitHub API       │
        │ (80% threshold)   │  │ (create PRs)     │
        └──────────────────┘  └──────────────────┘
                              │
                              ▼
        ┌────────────────────────────────────────┐
        │  agentic-k8s-operator (GitHub Repo)    │
        │  - PRs auto-created                     │
        │  - Branch protection enabled            │
        │  - Labels: autonomous, phase-1          │
        └────────────────────────────────────────┘
```

## Agent Pipeline (LangGraph Sequential Flow)

```
Task from Queue
      │
      ▼
┌──────────────┐
│  PM Agent    │  (Claude Sonnet 4.5)
│              │  - Analyze task
│              │  - Write PRD
│              │  - Define acceptance criteria
└──────┬───────┘
       │ State: {task, prd, consensus_votes[]}
       ▼
┌──────────────┐
│ Architect    │  (Claude Opus 4.5) ★ CRITICAL DECISIONS
│  Agent       │  - Design architecture
│              │  - Define components
│              │  - Risk assessment
│              │  - Vote: APPROVE (90/100)
└──────┬───────┘
       │ State: {prd, architecture, votes[]}
       ▼
┌──────────────┐
│ Engineer     │  (Claude Sonnet 4.5)
│  Agent       │  - Implement Go code
│              │  - Follow Kubebuilder patterns
│              │  - Error handling
│              │  - Vote: APPROVE (85/100)
└──────┬───────┘
       │ State: {architecture, implementation, votes[]}
       ▼
┌──────────────┐
│ QA Agent     │  (Claude Sonnet 4.5)
│              │  - Review code
│              │  - Write tests
│              │  - Edge cases
│              │  - Vote: APPROVE (80/100)
└──────┬───────┘
       │ State: {implementation, test_plan, votes[]}
       ▼
┌──────────────┐
│ Security     │  (Claude Sonnet 4.5)
│  Agent       │  - Risk-based scoring (baseline 70%)
│              │  - Check for new vulnerabilities
│              │  - Compliance review
│              │  - Vote: APPROVE (75/100)
└──────┬───────┘
       │ Consensus Score: Avg(90,85,80,75) = 82.5%
       │ Threshold: 80% → APPROVED ✅
       ▼
┌──────────────┐
│ GitManager   │  (Python, no LLM)
│              │  - Extract files from implementation
│              │  - Create feature branch
│              │  - Write files to repo
│              │  - Commit & push
│              │  - Create GitHub PR
└──────┬───────┘
       │
       ▼
GitHub PR Created
```

## Model Stack

| Agent       | Model              | Cost/1M tokens | When Used            |
|-------------|-------------------|----------------|---------------------|
| PM          | Sonnet 4.5        | $3.00          | Every task          |
| Architect   | **Opus 4.5** ★    | $15.00         | Every task (critical)|
| Engineer    | Sonnet 4.5        | $3.00          | Every task          |
| QA          | Sonnet 4.5        | $3.00          | Every task          |
| Security    | Sonnet 4.5        | $3.00          | Every task          |
| GitManager  | -                 | $0             | Code execution only  |

**Cost per task**: ~$0.50-1.00 (depending on task complexity)  
**Phase 1 estimate**: $10-15 (15 tasks)

## Consensus Scoring

```python
# Voting weights (equal weight for all agents)
votes = [
    {"agent": "PM", "vote": "APPROVE"},
    {"agent": "Architect", "score": 90},
    {"agent": "Engineer", "score": 85},
    {"agent": "QA", "score": 80},
    {"agent": "Security", "score": 75}
]

consensus_score = sum(v['score'] for v in votes) / len(votes)
# Result: 82.5%

# Decision logic
if consensus_score >= 80:
    action = "CREATE_PR"
elif consensus_score >= 60:
    action = "CONDITIONAL_APPROVE"  # Flag for review
else:
    action = "REJECT"
```

## Security Agent Risk Scoring

**Baseline approach** (not keyword-based):
```python
score = 70  # Baseline (assume secure cluster)

# Deduct for NEW risks
if introduces_privilege_escalation:
    score -= 25
if destructive_operations_without_safeguards:
    score -= 30
if unvalidated_external_api:
    score -= 10
if missing_audit_logging:
    score -= 10

# Add for enhancements
if encryption_at_rest_transit:
    score += 10
if advanced_audit_logging:
    score += 10

# Clamp to 0-100
score = max(0, min(100, score))
```

## Cronjob Architecture (Future)

**Note**: Currently NO cronjobs configured. Planned for monitoring:

```bash
# Planned cronjobs (not yet deployed)

# 1. Task health check (every 30 min)
*/30 * * * * /root/autonomous-dev-system/venv/bin/python3 /root/autonomous-dev-system/tools/check_queue.py

# 2. Daily progress report (9 PM IST)
0 15 * * * /root/autonomous-dev-system/venv/bin/python3 /root/autonomous-dev-system/tools/daily_report.py

# 3. Weekly summary (Sundays 10 AM)
0 4 * * 0 /root/autonomous-dev-system/venv/bin/python3 /root/autonomous-dev-system/tools/weekly_summary.py
```

**Current monitoring**: Manual via `journalctl -u autonomous-dev -f`

## State Management

### PostgreSQL Schema

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_id VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    pr_url TEXT,
    consensus_score FLOAT,
    result JSONB
);

CREATE INDEX idx_status ON tasks(status);
CREATE INDEX idx_priority ON tasks(priority DESC);
```

### Task Lifecycle

```
pending → in_progress → completed
                    ↘ failed
```

## Deployment Architecture

```
Local Mac (Goodra)
      │
      │ 1. Package autonomous-dev-system
      │ 2. Upload via SCP
      ▼
DigitalOcean Droplet (167.71.228.21)
      │
      ├─ PostgreSQL (localhost:5432)
      │   └─ autonomous_dev database
      │
      ├─ /root/autonomous-dev-system/
      │   ├─ venv/ (Python 3.12 virtualenv)
      │   ├─ agents/ (6 agent implementations)
      │   ├─ graph/ (LangGraph workflow)
      │   ├─ db/ (PostgreSQL integration)
      │   ├─ tools/ (utilities)
      │   └─ main.py (orchestrator)
      │
      ├─ /root/agentic-k8s-operator/ (cloned from GitHub)
      │   └─ Target repo for PRs
      │
      └─ Systemd Service: autonomous-dev.service
           ├─ Runs: main.py
           ├─ Auto-restart on failure
           └─ Logs: journalctl -u autonomous-dev
```

## GitHub Integration

### Branch Strategy

```
main (protected)
  ├─ feature/task-phase1-001  (PR #1)
  ├─ feature/task-phase1-002  (PR #2)
  ├─ feature/task-phase1-003  (PR #3)
  └─ ... (15 Phase 1 PRs)
```

### PR Metadata

Each autonomous PR includes:
- **Title**: `[Autonomous] <Task description>`
- **Body**: 
  - Consensus score (e.g., 82.5%)
  - Full PRD
  - Agent votes breakdown
- **Labels**:
  - `autonomous` (all PRs)
  - `phase-1` (current phase)
  - `consensus-approved` (score ≥80%)
  - `security-review` (if security score <80%)

### Branch Protection

- **main** requires:
  - 1 approving review
  - No force pushes
  - No deletions

## Monitoring & Observability

### Logs

```bash
# Real-time logs
journalctl -u autonomous-dev -f

# Last 100 lines
journalctl -u autonomous-dev -n 100

# JSON output
journalctl -u autonomous-dev -o json
```

### Metrics (Manual)

```python
# Check queue stats
from db.database import TaskDatabase
db = TaskDatabase('postgresql://...')
stats = db.get_stats()

# Example output:
{
  'pending': {'count': 10, 'avg_score': None},
  'in_progress': {'count': 1, 'avg_score': None},
  'completed': {'count': 4, 'avg_score': 82.5}
}
```

## Security Posture

1. **Droplet**: 
   - Private IP only (no public exposure except SSH)
   - SSH key-based auth only (password disabled)
   - UFW firewall enabled

2. **GitHub**:
   - Personal access token (fine-grained, repo scope only)
   - Stored in `.env` (not committed)

3. **Anthropic API**:
   - OAuth token (rate-limited)
   - Stored in `.env` (not committed)

4. **PostgreSQL**:
   - Localhost only (no remote access)
   - Default postgres user (single-user system)

## Failure Modes & Recovery

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Agent crashes | Exception in main.py | Systemd auto-restart (10s delay) |
| Consensus <80% | Voting logic | Task marked `failed`, next task starts |
| GitHub rate limit | API error | Exponential backoff (not implemented yet) |
| PostgreSQL down | Connection error | Systemd restart (depends on postgresql.service) |
| Anthropic quota exhausted | API error | Pause for 5 hours (not implemented yet) |

## Configuration Files

```
/root/autonomous-dev-system/
├─ .env                 (secrets, not committed)
│   ├─ DATABASE_URL
│   ├─ GITHUB_TOKEN
│   ├─ ANTHROPIC_API_KEY
│   └─ REPO_PATH
│
├─ requirements.txt     (Python deps)
├─ main.py             (orchestrator entry point)
├─ load_tasks.py       (populate task queue)
└─ deploy.sh           (deployment automation)
```

## Phase 1 Task Backlog (15 Tasks)

1. **phase1-001**: Initialize Kubebuilder project (priority 100)
2. **phase1-002**: Implement PM agent (priority 90)
3. **phase1-003**: Implement Architect agent (priority 90)
4. **phase1-004**: Implement Engineer agent (priority 85)
5. **phase1-005**: Implement QA agent (priority 85)
6. **phase1-006**: Implement Security agent (priority 80)
7. **phase1-007**: PostgreSQL schema (priority 95)
8. **phase1-008**: Consensus voting (priority 95)
9. **phase1-009**: LangGraph workflow (priority 90)
10. **phase1-010**: GitHub integration (priority 85)
11. **phase1-011**: Execution engine (priority 80)
12. **phase1-012**: Monitoring dashboard (priority 75)
13. **phase1-013**: README & docs (priority 70)
14. **phase1-014**: Unit tests (priority 70)
15. **phase1-015**: Integration tests (priority 65)

## Timeline Projection

| Time | Event |
|------|-------|
| T+0 | Deploy autonomous-dev.service |
| T+5m | First task starts (Kubebuilder init) |
| T+15m | First PR created (consensus 85%+) |
| T+1h | 3-5 PRs created |
| T+4h | 8-10 PRs created |
| T+24h | 15 PRs created (Phase 1 complete) |
| T+48h | PRs reviewed & merged → Phase 2 starts |

## Cost Analysis

| Component | Cost/Month | Notes |
|-----------|-----------|-------|
| DigitalOcean Droplet | $14 | 8GB RAM, 160GB SSD |
| Claude API (Phase 1) | ~$10-15 | 15 tasks × $0.70 avg |
| PostgreSQL | $0 | Self-hosted |
| GitHub | $0 | Private repo (free tier) |
| **Total Phase 1** | **~$24-29** | One-time + monthly infra |

## Next Steps After Phase 1

1. **Review & merge PRs** (human oversight)
2. **Load Phase 2 tasks** (security & scale)
3. **Add monitoring cronjobs**
4. **Implement quota management** (pause on Anthropic limits)
5. **Add Slack/Telegram notifications** (PR created alerts)

---

**Architecture Status**: ✅ DEPLOYED  
**Phase**: 1 of 3 (MVP)  
**Last Updated**: 2026-02-22
