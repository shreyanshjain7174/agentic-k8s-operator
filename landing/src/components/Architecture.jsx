import { useRef } from "react";
import { motion, useInView } from "framer-motion";
import { FileCode, RefreshCw, Zap, CheckCircle } from "lucide-react";

const steps = [
  {
    number: "01",
    icon: FileCode,
    title: "Apply YAML",
    description:
      "kubectl apply your AgentWorkload manifest. The Kubernetes API server stores it in etcd.",
    command: "kubectl apply -f agent.yaml",
    color: "#00d4aa",
    colorAlpha: "rgba(0, 212, 170, 0.12)",
    colorBorder: "rgba(0, 212, 170, 0.25)",
  },
  {
    number: "02",
    icon: RefreshCw,
    title: "Operator Reconciles",
    description:
      "The Go controller watches for AgentWorkload events, validates policy with OPA, and creates an Argo Workflow.",
    command: "// controller-runtime reconcile loop",
    color: "#6366f1",
    colorAlpha: "rgba(99, 102, 241, 0.12)",
    colorBorder: "rgba(99, 102, 241, 0.25)",
  },
  {
    number: "03",
    icon: Zap,
    title: "Agent Launches",
    description:
      "Argo Workflow spins up a pod running your LangGraph agent with all MCP tools mounted and ready.",
    command: "argo submit --from=wftmpl/agent",
    color: "#f59e0b",
    colorAlpha: "rgba(245, 158, 11, 0.12)",
    colorBorder: "rgba(245, 158, 11, 0.25)",
  },
  {
    number: "04",
    icon: CheckCircle,
    title: "Report Delivered",
    description:
      "Agent completes its task, uploads artifacts to MinIO, and updates the AgentWorkload status with results.",
    command: "status: phase: Succeeded",
    color: "#22c55e",
    colorAlpha: "rgba(34, 197, 94, 0.12)",
    colorBorder: "rgba(34, 197, 94, 0.25)",
  },
];

function GridPattern() {
  return (
    <svg
      className="absolute inset-0 w-full h-full"
      style={{ opacity: 0.025 }}
      xmlns="http://www.w3.org/2000/svg"
    >
      <defs>
        <pattern
          id="arch-grid"
          width="40"
          height="40"
          patternUnits="userSpaceOnUse"
        >
          <path
            d="M 40 0 L 0 0 0 40"
            fill="none"
            stroke="white"
            strokeWidth="1"
          />
        </pattern>
      </defs>
      <rect width="100%" height="100%" fill="url(#arch-grid)" />
    </svg>
  );
}

function ConnectorArrow({ color }) {
  return (
    <div className="hidden lg:flex items-center justify-center w-16 flex-shrink-0 relative">
      <motion.div
        initial={{ scaleX: 0, opacity: 0 }}
        whileInView={{ scaleX: 1, opacity: 1 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6, ease: "easeOut", delay: 0.4 }}
        className="w-full h-0.5 relative"
        style={{
          background: `linear-gradient(90deg, ${color}40, ${color}90)`,
          transformOrigin: "left center",
        }}
      >
        <div
          className="absolute right-0 top-1/2 -translate-y-1/2"
          style={{
            width: 0,
            height: 0,
            borderTop: "5px solid transparent",
            borderBottom: "5px solid transparent",
            borderLeft: `7px solid ${color}90`,
          }}
        />
      </motion.div>
    </div>
  );
}

function StepCard({ step, index }) {
  const Icon = step.icon;

  return (
    <motion.div
      initial={{ opacity: 0, y: 32 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: "-40px" }}
      transition={{ duration: 0.55, ease: "easeOut", delay: index * 0.2 }}
      className="flex-1 min-w-0 flex flex-col"
    >
      <div
        className="rounded-2xl p-6 flex flex-col gap-5 h-full relative overflow-hidden group"
        style={{
          background: "rgba(13, 21, 37, 0.8)",
          border: `1px solid ${step.colorBorder}`,
          backdropFilter: "blur(12px)",
        }}
      >
        {/* Glow effect on hover */}
        <div
          className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
          style={{
            background: `radial-gradient(ellipse at top left, ${step.colorAlpha} 0%, transparent 70%)`,
          }}
        />

        {/* Number badge */}
        <div className="flex items-center justify-between">
          <div
            className="text-xs font-bold tracking-widest"
            style={{
              fontFamily: "'IBM Plex Mono', monospace",
              color: step.color,
              opacity: 0.6,
            }}
          >
            {step.number}
          </div>
          <div
            className="w-10 h-10 rounded-xl flex items-center justify-center"
            style={{
              background: step.colorAlpha,
              border: `1px solid ${step.colorBorder}`,
            }}
          >
            <Icon size={20} color={step.color} strokeWidth={1.75} />
          </div>
        </div>

        <div>
          <h3
            className="text-base font-bold mb-2"
            style={{
              fontFamily: "'Syne', sans-serif",
              color: "#e2e8f0",
            }}
          >
            {step.title}
          </h3>
          <p
            className="text-sm leading-relaxed"
            style={{
              fontFamily: "'DM Sans', sans-serif",
              color: "#94a3b8",
            }}
          >
            {step.description}
          </p>
        </div>

        {/* Command snippet */}
        <div
          className="rounded-lg px-4 py-3 mt-auto"
          style={{
            background: "rgba(5, 8, 15, 0.8)",
            border: "1px solid rgba(255,255,255,0.06)",
          }}
        >
          <span
            className="text-xs"
            style={{
              fontFamily: "'IBM Plex Mono', monospace",
              color: step.color,
            }}
          >
            $ {step.command}
          </span>
        </div>
      </div>
    </motion.div>
  );
}

export default function Architecture() {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-80px" });

  return (
    <section
      id="architecture"
      ref={ref}
      className="py-24 px-4 relative overflow-hidden"
      style={{ background: "#05080f" }}
    >
      <GridPattern />

      {/* Ambient glow */}
      <div
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-2/3 h-64 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse, rgba(0, 212, 170, 0.04) 0%, transparent 70%)",
        }}
      />

      <div className="max-w-6xl mx-auto relative z-10">
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
            Architecture
          </div>
          <h2
            className="text-3xl sm:text-4xl lg:text-5xl font-bold"
            style={{
              fontFamily: "'Syne', sans-serif",
              color: "#e2e8f0",
            }}
          >
            How It{" "}
            <span
              style={{
                background: "linear-gradient(135deg, #00d4aa, #6366f1)",
                WebkitBackgroundClip: "text",
                WebkitTextFillColor: "transparent",
                backgroundClip: "text",
              }}
            >
              Works
            </span>
          </h2>
          <p
            className="mt-4 text-base max-w-xl mx-auto"
            style={{
              fontFamily: "'DM Sans', sans-serif",
              color: "#94a3b8",
            }}
          >
            From a single kubectl command to a running AI agent — four steps, fully automated.
          </p>
        </motion.div>

        {/* Steps — desktop: horizontal row with arrows, mobile: vertical stack */}
        <div className="flex flex-col lg:flex-row items-stretch gap-4 lg:gap-0">
          {steps.map((step, index) => (
            <div key={step.number} className="flex flex-col lg:flex-row items-stretch flex-1 min-w-0">
              <StepCard step={step} index={index} />
              {index < steps.length - 1 && (
                <ConnectorArrow color={step.color} />
              )}
            </div>
          ))}
        </div>

        {/* Mobile: vertical connector lines */}
        <div className="flex lg:hidden flex-col items-center mt-0">
          {steps.slice(0, -1).map((step, i) => (
            <motion.div
              key={i}
              initial={{ scaleY: 0, opacity: 0 }}
              whileInView={{ scaleY: 1, opacity: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4, delay: i * 0.15 }}
              className="w-0.5 h-6 my-1"
              style={{
                background: `linear-gradient(to bottom, ${step.color}60, ${steps[i + 1].color}60)`,
                transformOrigin: "top center",
              }}
            />
          ))}
        </div>

        {/* Bottom CTA */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6, delay: 0.8 }}
          className="text-center mt-16"
        >
          <div
            className="inline-flex items-center gap-3 px-6 py-4 rounded-xl"
            style={{
              background: "rgba(13, 21, 37, 0.8)",
              border: "1px solid rgba(255,255,255,0.08)",
            }}
          >
            <span
              className="text-sm"
              style={{
                fontFamily: "'IBM Plex Mono', monospace",
                color: "#94a3b8",
              }}
            >
              Total time from apply to output:
            </span>
            <span
              className="text-sm font-bold"
              style={{
                fontFamily: "'IBM Plex Mono', monospace",
                color: "#00d4aa",
              }}
            >
              under 4 minutes
            </span>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
