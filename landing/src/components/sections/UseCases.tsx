"use client";

import { AnimatedSection } from "../ui/AnimatedSection";
import { BarChart3, Wrench, Users } from "lucide-react";

const useCases = [
  {
    icon: <BarChart3 size={28} />,
    title: "Competitive Intelligence",
    category: "Financial Intelligence",
    description:
      "Deploy autonomous web scraping agents that gather competitive data from dozens of sources, synthesize findings through Claude, and deliver executive PDF reports -- all inside your cluster.",
    before: { value: "8 hours", label: "Manual research" },
    after: { value: "4 min", label: "Automated pipeline" },
    savings: "$320 saved per run",
    color: "#6366F1",
  },
  {
    icon: <Wrench size={28} />,
    title: "Autonomous K8s Remediation",
    category: "Platform Engineering",
    description:
      "AI agents detect cluster issues, diagnose root causes, and apply fixes automatically. CrashLoopBackOff, resource exhaustion, and config drift resolved without human intervention.",
    before: { value: "22 min", label: "MTTR" },
    after: { value: "47 sec", label: "Automated fix" },
    savings: "~$28K downtime protected",
    color: "#06B6D4",
  },
  {
    icon: <Users size={28} />,
    title: "Multi-Agent Research Swarm",
    category: "Research & Analysis",
    description:
      "Coordinate a swarm of specialized agents -- data collectors, analyzers, and synthesizers -- that work in parallel to produce comprehensive research reports.",
    before: { value: "6 hours", label: "Analyst time" },
    after: { value: "90 sec", label: "Agent swarm" },
    savings: "$400 analyst cost saved",
    color: "#10B981",
  },
];

function ProgressBar({ percentage, color }: { percentage: number; color: string }) {
  return (
    <div className="w-full h-2 bg-[#1E293B] rounded-full overflow-hidden">
      <div
        className="h-full rounded-full transition-all duration-1000"
        style={{ width: `${percentage}%`, backgroundColor: color }}
      />
    </div>
  );
}

export function UseCases() {
  return (
    <section id="use-cases" className="py-24 bg-[#0F1120]/50 relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <AnimatedSection>
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              Real-World{" "}
              <span className="gradient-text-blue">Use Cases</span>
            </h2>
            <p className="text-[#94A3B8] text-lg max-w-2xl mx-auto">
              See how teams are using the Agentic K8s Operator to automate
              complex workflows and save thousands of hours.
            </p>
          </div>
        </AnimatedSection>

        <div className="space-y-8">
          {useCases.map((uc, i) => (
            <AnimatedSection key={i} delay={i * 0.15}>
              <div className="glass-card p-8">
                <div className="grid lg:grid-cols-2 gap-8 items-center">
                  {/* Content */}
                  <div>
                    <div
                      className="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium mb-4"
                      style={{
                        backgroundColor: `${uc.color}15`,
                        color: uc.color,
                      }}
                    >
                      {uc.category}
                    </div>
                    <h3 className="text-2xl font-bold text-white mb-3 flex items-center gap-3">
                      <span style={{ color: uc.color }}>{uc.icon}</span>
                      {uc.title}
                    </h3>
                    <p className="text-[#94A3B8] leading-relaxed mb-6">
                      {uc.description}
                    </p>
                  </div>

                  {/* Metrics */}
                  <div className="space-y-6">
                    <div className="grid grid-cols-2 gap-4">
                      <div className="bg-[#0A0B14] rounded-xl p-4 border border-[#1E293B]">
                        <div className="text-xs text-[#64748B] mb-1">Before</div>
                        <div className="text-2xl font-bold text-[#EF4444]">
                          {uc.before.value}
                        </div>
                        <div className="text-xs text-[#64748B]">{uc.before.label}</div>
                      </div>
                      <div className="bg-[#0A0B14] rounded-xl p-4 border border-[#1E293B]">
                        <div className="text-xs text-[#64748B] mb-1">After</div>
                        <div className="text-2xl font-bold" style={{ color: uc.color }}>
                          {uc.after.value}
                        </div>
                        <div className="text-xs text-[#64748B]">{uc.after.label}</div>
                      </div>
                    </div>

                    <div>
                      <div className="flex justify-between text-sm mb-2">
                        <span className="text-[#64748B]">Efficiency Gain</span>
                        <span style={{ color: uc.color }}>99%+</span>
                      </div>
                      <ProgressBar percentage={99} color={uc.color} />
                    </div>

                    <div
                      className="flex items-center justify-center gap-2 py-2 rounded-lg text-sm font-medium"
                      style={{
                        backgroundColor: `${uc.color}10`,
                        color: uc.color,
                      }}
                    >
                      {uc.savings}
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          ))}
        </div>
      </div>
    </section>
  );
}
