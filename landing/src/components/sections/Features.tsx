"use client";

import { AnimatedSection } from "../ui/AnimatedSection";
import {
  Boxes,
  Plug,
  Network,
  Database,
  Shield,
  Terminal,
} from "lucide-react";

const features = [
  {
    icon: <Boxes size={24} />,
    title: "AgentWorkload CRD",
    description:
      "Declare AI agent workflows as native Kubernetes resources. One YAML, full orchestration with status tracking.",
    color: "#6366F1",
  },
  {
    icon: <Plug size={24} />,
    title: "MCP Integration",
    description:
      "Tool-agnostic Model Context Protocol support. Connect Browserless, LLMs, databases, and custom tools.",
    color: "#8B5CF6",
  },
  {
    icon: <Network size={24} />,
    title: "Multi-Agent Orchestration",
    description:
      "LangGraph-powered agent coordination with ReAct patterns. Agents collaborate on complex tasks autonomously.",
    color: "#06B6D4",
  },
  {
    icon: <Database size={24} />,
    title: "Durable State",
    description:
      "PostgreSQL checkpointing survives pod preemption. Your agents resume exactly where they left off.",
    color: "#10B981",
  },
  {
    icon: <Shield size={24} />,
    title: "Safety-First Architecture",
    description:
      "SSRF protection, OPA policies, validating webhooks, credential sanitization. 8 critical security layers.",
    color: "#F59E0B",
  },
  {
    icon: <Terminal size={24} />,
    title: "One-Line Deploy",
    description:
      "helm install and you're running. 47 pods, full observability stack with Prometheus and Grafana, production-ready.",
    color: "#D946EF",
  },
];

export function Features() {
  return (
    <section id="features" className="py-24 relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <AnimatedSection>
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              Everything You Need to Run{" "}
              <span className="gradient-text">AI Agents on K8s</span>
            </h2>
            <p className="text-[#94A3B8] text-lg max-w-2xl mx-auto">
              A complete platform for deploying, orchestrating, and monitoring
              autonomous AI agents inside your Kubernetes cluster.
            </p>
          </div>
        </AnimatedSection>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, i) => (
            <AnimatedSection key={i} delay={i * 0.1}>
              <div className="glass-card gradient-border p-6 h-full">
                <div
                  className="w-12 h-12 rounded-xl mb-4 flex items-center justify-center"
                  style={{ backgroundColor: `${feature.color}15` }}
                >
                  <div style={{ color: feature.color }}>{feature.icon}</div>
                </div>
                <h3 className="text-lg font-semibold text-white mb-2">
                  {feature.title}
                </h3>
                <p className="text-sm text-[#94A3B8] leading-relaxed">
                  {feature.description}
                </p>
              </div>
            </AnimatedSection>
          ))}
        </div>
      </div>
    </section>
  );
}
