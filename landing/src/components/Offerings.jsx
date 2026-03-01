import { motion } from "framer-motion";
import {
  Layers,
  Plug,
  GitBranch,
  Shield,
  Package,
  BarChart2,
} from "lucide-react";

const features = [
  {
    icon: Layers,
    title: "AgentWorkload CRD",
    description:
      "Define AI agents as native Kubernetes resources. Apply YAML, get autonomous agents.",
  },
  {
    icon: Plug,
    title: "MCP Tool Integration",
    description:
      "Connect any MCP server. Browserless, LiteLLM, PostgreSQL — zero lock-in.",
  },
  {
    icon: GitBranch,
    title: "LangGraph Orchestration",
    description:
      "Multi-agent state machines with PostgreSQL checkpointing. Survives pod evictions.",
  },
  {
    icon: Shield,
    title: "OPA Policy Governance",
    description:
      "Enforce strict agent action boundaries. SSRF protection, allowlist validation.",
  },
  {
    icon: Package,
    title: "Helm Umbrella Chart",
    description:
      "One command deploys the full stack: operator, Argo, MinIO, PostgreSQL, monitoring.",
  },
  {
    icon: BarChart2,
    title: "Full Observability",
    description:
      "Prometheus metrics, Grafana dashboards, Loki logs, AlertManager — batteries included.",
  },
];

const containerVariants = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const cardVariants = {
  hidden: { opacity: 0, y: 32 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, ease: "easeOut" },
  },
};

function FeatureCard({ feature }) {
  const Icon = feature.icon;

  return (
    <motion.div
      variants={cardVariants}
      whileHover={{
        borderColor: "rgba(0, 212, 170, 0.45)",
        boxShadow: "0 0 28px rgba(0, 212, 170, 0.12), 0 4px 24px rgba(0,0,0,0.4)",
        y: -4,
      }}
      className="rounded-xl p-6 flex flex-col gap-4 cursor-default transition-colors duration-300"
      style={{
        background: "rgba(13, 21, 37, 0.7)",
        border: "1px solid rgba(255,255,255,0.06)",
        backdropFilter: "blur(12px)",
      }}
    >
      <div
        className="w-11 h-11 rounded-xl flex items-center justify-center flex-shrink-0"
        style={{
          background: "rgba(0, 212, 170, 0.12)",
          border: "1px solid rgba(0, 212, 170, 0.2)",
        }}
      >
        <Icon size={22} color="#00d4aa" strokeWidth={1.75} />
      </div>

      <div>
        <h3
          className="text-base font-semibold mb-2"
          style={{
            fontFamily: "'Syne', sans-serif",
            color: "#e2e8f0",
          }}
        >
          {feature.title}
        </h3>
        <p
          className="text-sm leading-relaxed"
          style={{
            fontFamily: "'DM Sans', sans-serif",
            color: "#94a3b8",
          }}
        >
          {feature.description}
        </p>
      </div>
    </motion.div>
  );
}

export default function Offerings() {
  return (
    <section
      id="features"
      className="py-24 px-4"
      style={{ background: "#05080f" }}
    >
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 24 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-60px" }}
          transition={{ duration: 0.6, ease: "easeOut" }}
          className="text-center mb-16"
        >
          <div
            className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-widest mb-6"
            style={{
              background: "rgba(0, 212, 170, 0.08)",
              border: "1px solid rgba(0, 212, 170, 0.2)",
              color: "#00d4aa",
              fontFamily: "'DM Sans', sans-serif",
            }}
          >
            Platform Features
          </div>
          <h2
            className="text-3xl sm:text-4xl lg:text-5xl font-bold leading-tight"
            style={{
              fontFamily: "'Syne', sans-serif",
              color: "#e2e8f0",
            }}
          >
            Everything You Need to Run{" "}
            <span
              style={{
                background: "linear-gradient(135deg, #00d4aa, #6366f1)",
                WebkitBackgroundClip: "text",
                WebkitTextFillColor: "transparent",
                backgroundClip: "text",
              }}
            >
              AI Agents
            </span>{" "}
            in Production
          </h2>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-60px" }}
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5"
        >
          {features.map((feature) => (
            <FeatureCard key={feature.title} feature={feature} />
          ))}
        </motion.div>
      </div>
    </section>
  );
}
