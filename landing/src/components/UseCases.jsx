import { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence, useInView } from "framer-motion";
import { ArrowRight, TrendingDown, TrendingUp, Clock, DollarSign } from "lucide-react";

const YAML_COMPETITIVE = `apiVersion: agentic.clawdlinux.org/v1alpha1
kind: AgentWorkload
metadata:
  name: competitive-intel-agent
  namespace: agents
spec:
  model: claude-3-5-sonnet
  mcpServers:
    - name: browserless
      url: http://browserless:3000
    - name: minio
      url: http://minio:9000
  task: |
    Research top 5 competitors.
    Generate PDF report with insights.
  schedule: "0 9 * * MON"
  outputBucket: reports/competitive`;

const YAML_REMEDIATION = `apiVersion: agentic.clawdlinux.org/v1alpha1
kind: AgentWorkload
metadata:
  name: k8s-remediation-agent
  namespace: operators
spec:
  model: claude-3-5-sonnet
  triggers:
    - type: PodCrashLoop
    - type: OOMKilled
    - type: PodPending
  mcpServers:
    - name: kubectl-mcp
      url: http://kubectl-mcp:8080
  remediation:
    maxRetries: 3
    dryRun: false
  alertChannel: "#ops-alerts"`;

const YAML_SWARM = `apiVersion: agentic.clawdlinux.org/v1alpha1
kind: AgentWorkload
metadata:
  name: research-swarm
  namespace: agents
spec:
  model: claude-3-5-sonnet
  orchestration:
    type: langgraph
    checkpointer: postgresql
    parallel: 5
  agents:
    - role: market-analyst
    - role: sentiment-analyst
    - role: data-aggregator
    - role: report-writer
    - role: fact-checker
  outputFormat: synthesized-report
  timeout: "120s"`;

const tabs = [
  {
    id: "competitive",
    label: "Competitive Intelligence",
    problem:
      "Analysts spending 8 hours on competitor research. $320 per report, delayed insights, and stale data by the time it reaches leadership.",
    solution:
      "Deploy an AgentWorkload, get a full intelligence report in under 4 minutes. Automated, repeatable, and 99.9% cheaper.",
    metrics: {
      before: { time: "8 hr", cost: "$320" },
      after: { time: "4 min", cost: "<$0.01" },
      savings: "$320 saved per run",
      savingsColor: "#f59e0b",
    },
    yaml: YAML_COMPETITIVE,
  },
  {
    id: "remediation",
    label: "Autonomous K8s Remediation",
    problem:
      "Pod failures triggering 22-minute mean time to recovery. Each incident costs $28K in downtime with engineers pulled from feature work.",
    solution:
      "The operator detects issues, launches a remediation agent, and fixes in 47 seconds. Engineers sleep through the night.",
    metrics: {
      before: { time: "22 min", cost: "MTTR" },
      after: { time: "47 sec", cost: "MTTR" },
      savings: "$28K downtime protected",
      savingsColor: "#00d4aa",
    },
    yaml: YAML_REMEDIATION,
  },
  {
    id: "swarm",
    label: "Multi-Agent Research Swarm",
    problem:
      "6 hours of analyst time for deep market research. $400 per session with inconsistent quality and no audit trail.",
    solution:
      "Orchestrate 5 parallel agents with LangGraph. Synthesized, cited report in 90 seconds with full execution logs.",
    metrics: {
      before: { time: "6 hr", cost: "$400" },
      after: { time: "90 sec", cost: "$0.05" },
      savings: "4+ hours saved per run",
      savingsColor: "#6366f1",
    },
    yaml: YAML_SWARM,
  },
];

function MetricCard({ label, time, cost, highlight, highlightColor }) {
  return (
    <div
      className="flex-1 rounded-xl p-5"
      style={{
        background: highlight
          ? `rgba(${highlightColor}, 0.08)`
          : "rgba(255,255,255,0.03)",
        border: highlight
          ? `1px solid rgba(${highlightColor}, 0.25)`
          : "1px solid rgba(255,255,255,0.06)",
      }}
    >
      <div
        className="text-xs font-semibold uppercase tracking-widest mb-3"
        style={{
          fontFamily: "'DM Sans', sans-serif",
          color: highlight ? `rgb(${highlightColor})` : "#94a3b8",
        }}
      >
        {label}
      </div>
      <div
        className="text-2xl font-bold mb-1"
        style={{
          fontFamily: "'Syne', sans-serif",
          color: highlight ? `rgb(${highlightColor})` : "#e2e8f0",
        }}
      >
        {time}
      </div>
      <div
        className="text-sm"
        style={{
          fontFamily: "'DM Sans', sans-serif",
          color: "#94a3b8",
        }}
      >
        {cost}
      </div>
    </div>
  );
}

function CodeBlock({ code }) {
  return (
    <div
      className="rounded-xl overflow-hidden"
      style={{
        background: "rgba(5, 8, 15, 0.9)",
        border: "1px solid rgba(255,255,255,0.08)",
      }}
    >
      <div
        className="flex items-center gap-2 px-4 py-3"
        style={{
          background: "rgba(255,255,255,0.03)",
          borderBottom: "1px solid rgba(255,255,255,0.06)",
        }}
      >
        <div className="w-3 h-3 rounded-full" style={{ background: "#ff5f57" }} />
        <div className="w-3 h-3 rounded-full" style={{ background: "#ffbd2e" }} />
        <div className="w-3 h-3 rounded-full" style={{ background: "#28c840" }} />
        <span
          className="ml-2 text-xs"
          style={{ fontFamily: "'IBM Plex Mono', monospace", color: "#94a3b8" }}
        >
          agentworkload.yaml
        </span>
      </div>
      <pre
        className="p-5 text-xs leading-relaxed overflow-x-auto"
        style={{
          fontFamily: "'IBM Plex Mono', monospace",
          color: "#e2e8f0",
          margin: 0,
        }}
      >
        <code
          dangerouslySetInnerHTML={{
            __html: (() => {
              // Escape HTML special chars first (but NOT quotes — we need them for step 2)
              const escaped = code
                .replace(/&/g, '&amp;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;');
              // Apply highlighting with quoted strings FIRST, before any <span> tags
              // (with style="color:#...") are injected — prevents the quote regex from
              // accidentally matching hex values inside injected span attributes.
              return escaped
                .replace(/"([^"]*)"/g, '<span style="color:#f59e0b">"$1"</span>')
                .replace(/\b(true|false)\b/g, '<span style="color:#f59e0b">$1</span>')
                .replace(/(apiVersion:|kind:|metadata:|spec:|name:|namespace:|model:|mcpServers:|task:|schedule:|outputBucket:|triggers:|remediation:|alertChannel:|orchestration:|agents:|outputFormat:|timeout:|type:|url:|maxRetries:|dryRun:|parallel:|role:|checkpointer:|format:)/g,
                  '<span style="color:#6366f1">$1</span>')
                .replace(/(claude-3-5-sonnet|AgentWorkload|agentic\.clawdlinux\.org\/v1alpha1)/g,
                  '<span style="color:#00d4aa">$1</span>');
            })(),
          }}
        />
      </pre>
    </div>
  );
}

export default function UseCases() {
  const [activeTab, setActiveTab] = useState(0);
  const ref = useRef(null);

  const tab = tabs[activeTab];

  return (
    <section
      id="use-cases"
      className="py-24 px-4"
      style={{
        background: "linear-gradient(180deg, #05080f 0%, #070b14 100%)",
      }}
    >
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 24 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-60px" }}
          transition={{ duration: 0.6, ease: "easeOut" }}
          className="text-center mb-14"
        >
          <div
            className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-widest mb-6"
            style={{
              background: "rgba(99, 102, 241, 0.08)",
              border: "1px solid rgba(99, 102, 241, 0.2)",
              color: "#6366f1",
              fontFamily: "'DM Sans', sans-serif",
            }}
          >
            Use Cases
          </div>
          <h2
            className="text-3xl sm:text-4xl lg:text-5xl font-bold"
            style={{
              fontFamily: "'Syne', sans-serif",
              color: "#e2e8f0",
            }}
          >
            Three Use Cases,{" "}
            <span
              style={{
                background: "linear-gradient(135deg, #00d4aa, #6366f1)",
                WebkitBackgroundClip: "text",
                WebkitTextFillColor: "transparent",
                backgroundClip: "text",
              }}
            >
              Proven in Production
            </span>
          </h2>
        </motion.div>

        {/* Tab Navigation */}
        <div
          className="flex flex-col sm:flex-row gap-1 p-1 rounded-xl mb-10"
          style={{
            background: "rgba(13, 21, 37, 0.7)",
            border: "1px solid rgba(255,255,255,0.06)",
          }}
        >
          {tabs.map((t, i) => (
            <button
              key={t.id}
              onClick={() => setActiveTab(i)}
              className="flex-1 px-4 py-3 rounded-lg text-sm font-medium transition-all duration-300 text-left sm:text-center relative"
              style={{
                fontFamily: "'DM Sans', sans-serif",
                color: activeTab === i ? "#e2e8f0" : "#94a3b8",
                background: activeTab === i ? "rgba(0, 212, 170, 0.1)" : "transparent",
                border: activeTab === i ? "1px solid rgba(0, 212, 170, 0.25)" : "1px solid transparent",
              }}
            >
              {activeTab === i && (
                <motion.span
                  layoutId="tab-indicator"
                  className="absolute inset-0 rounded-lg"
                  style={{
                    background: "rgba(0, 212, 170, 0.08)",
                  }}
                  transition={{ type: "spring", bounce: 0.2, duration: 0.5 }}
                />
              )}
              <span className="relative z-10">{t.label}</span>
            </button>
          ))}
        </div>

        {/* Tab Content */}
        <AnimatePresence mode="wait">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -16 }}
            transition={{ duration: 0.35, ease: "easeOut" }}
          >
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Left: Problem / Solution / Metrics */}
              <div className="flex flex-col gap-6">
                <div
                  className="rounded-xl p-6"
                  style={{
                    background: "rgba(13, 21, 37, 0.7)",
                    border: "1px solid rgba(255,255,255,0.06)",
                  }}
                >
                  <div
                    className="text-xs font-semibold uppercase tracking-widest mb-3"
                    style={{
                      fontFamily: "'DM Sans', sans-serif",
                      color: "#ef4444",
                    }}
                  >
                    The Problem
                  </div>
                  <p
                    className="text-sm leading-relaxed"
                    style={{
                      fontFamily: "'DM Sans', sans-serif",
                      color: "#94a3b8",
                    }}
                  >
                    {tab.problem}
                  </p>
                </div>

                <div
                  className="rounded-xl p-6"
                  style={{
                    background: "rgba(13, 21, 37, 0.7)",
                    border: "1px solid rgba(0, 212, 170, 0.15)",
                  }}
                >
                  <div
                    className="text-xs font-semibold uppercase tracking-widest mb-3"
                    style={{
                      fontFamily: "'DM Sans', sans-serif",
                      color: "#00d4aa",
                    }}
                  >
                    The Solution
                  </div>
                  <p
                    className="text-sm leading-relaxed"
                    style={{
                      fontFamily: "'DM Sans', sans-serif",
                      color: "#e2e8f0",
                    }}
                  >
                    {tab.solution}
                  </p>
                </div>

                {/* Metrics */}
                <div className="flex gap-4">
                  <MetricCard
                    label="Before"
                    time={tab.metrics.before.time}
                    cost={tab.metrics.before.cost}
                    highlight={false}
                  />
                  <div className="flex items-center">
                    <ArrowRight size={20} color="#94a3b8" />
                  </div>
                  <MetricCard
                    label="After"
                    time={tab.metrics.after.time}
                    cost={tab.metrics.after.cost}
                    highlight={true}
                    highlightColor="0, 212, 170"
                  />
                </div>

                <div
                  className="rounded-xl px-5 py-4 flex items-center gap-3"
                  style={{
                    background: `rgba(${tab.metrics.savingsColor === "#f59e0b" ? "245, 158, 11" : tab.metrics.savingsColor === "#00d4aa" ? "0, 212, 170" : "99, 102, 241"}, 0.08)`,
                    border: `1px solid rgba(${tab.metrics.savingsColor === "#f59e0b" ? "245, 158, 11" : tab.metrics.savingsColor === "#00d4aa" ? "0, 212, 170" : "99, 102, 241"}, 0.2)`,
                  }}
                >
                  <TrendingDown size={18} color={tab.metrics.savingsColor} />
                  <span
                    className="text-sm font-semibold"
                    style={{
                      fontFamily: "'DM Sans', sans-serif",
                      color: tab.metrics.savingsColor,
                    }}
                  >
                    {tab.metrics.savings}
                  </span>
                </div>
              </div>

              {/* Right: Code Snippet */}
              <div>
                <CodeBlock code={tab.yaml} />
              </div>
            </div>
          </motion.div>
        </AnimatePresence>
      </div>
    </section>
  );
}
