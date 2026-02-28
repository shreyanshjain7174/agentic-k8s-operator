"use client";

import dynamic from "next/dynamic";
import { AnimatedSection } from "../ui/AnimatedSection";

const ArchPlayer = dynamic(
  () => import("./ArchitecturePlayer"),
  { ssr: false, loading: () => <div className="w-full h-48 bg-[#141628] rounded-2xl animate-pulse" /> }
);

const techStack = [
  { name: "Kubernetes", color: "#326CE5" },
  { name: "Argo Workflows", color: "#E7553B" },
  { name: "PostgreSQL", color: "#336791" },
  { name: "MinIO", color: "#C72C48" },
  { name: "Prometheus", color: "#E6522C" },
  { name: "Grafana", color: "#F46800" },
  { name: "LangGraph", color: "#10B981" },
  { name: "Cloudflare AI", color: "#F6821F" },
];

export function Architecture() {
  return (
    <section id="architecture" className="py-24 relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <AnimatedSection>
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              How It <span className="gradient-text">Works</span>
            </h2>
            <p className="text-[#94A3B8] text-lg max-w-2xl mx-auto">
              From YAML declaration to intelligent execution -- a fully
              orchestrated pipeline running inside your cluster.
            </p>
          </div>
        </AnimatedSection>

        <AnimatedSection delay={0.2}>
          <div className="glass-card p-8 mb-12 overflow-hidden">
            <ArchPlayer />
          </div>
        </AnimatedSection>

        {/* Tech Stack */}
        <AnimatedSection delay={0.3}>
          <div className="text-center mb-8">
            <h3 className="text-lg font-semibold text-white mb-6">
              Powered by Battle-Tested Infrastructure
            </h3>
            <div className="flex flex-wrap justify-center gap-4">
              {techStack.map((tech) => (
                <div
                  key={tech.name}
                  className="flex items-center gap-2 px-4 py-2 rounded-full bg-[#141628] border border-[#1E293B] hover:border-[#334155] transition-colors"
                >
                  <div
                    className="w-3 h-3 rounded-full"
                    style={{ backgroundColor: tech.color }}
                  />
                  <span className="text-sm text-[#94A3B8]">{tech.name}</span>
                </div>
              ))}
            </div>
          </div>
        </AnimatedSection>
      </div>
    </section>
  );
}
