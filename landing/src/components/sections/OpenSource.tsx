"use client";

import { AnimatedSection } from "../ui/AnimatedSection";
import { Github, Star, GitFork, CheckCircle, Scale } from "lucide-react";

const yamlSnippet = `apiVersion: agentic.clawdlinux.org/v1alpha1
kind: AgentWorkload
metadata:
  name: competitive-intel
spec:
  agentImage: "ghcr.io/shreyansh/agent:latest"
  mcpServers:
    - name: browserless
      endpoint: "http://browserless:3000"
  tools:
    - web_scrape
    - synthesize_report
  schedule: "0 6 * * 1-5"
  checkpoint:
    enabled: true
    storage: postgresql`;

export function OpenSource() {
  return (
    <section id="open-source" className="py-24 bg-[#0F1120]/50 relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <AnimatedSection>
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              100% <span className="gradient-text-blue">Open Source</span>
            </h2>
            <p className="text-[#94A3B8] text-lg max-w-2xl mx-auto">
              Built in the open with Apache 2.0 -- the same license as
              Kubernetes itself. Inspect, modify, and contribute.
            </p>
          </div>
        </AnimatedSection>

        <div className="grid lg:grid-cols-2 gap-8">
          {/* GitHub Card */}
          <AnimatedSection delay={0.1}>
            <div className="glass-card p-8 h-full">
              <div className="flex items-center gap-3 mb-6">
                <Github size={28} className="text-white" />
                <div>
                  <h3 className="font-semibold text-white">
                    shreyanshjain7174/agentic-k8s-operator
                  </h3>
                  <p className="text-sm text-[#64748B]">
                    Production-grade K8s operator for AI agent workloads
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-6">
                <div className="flex items-center gap-2 text-sm text-[#94A3B8]">
                  <CheckCircle size={16} className="text-[#10B981]" />
                  46/46 tests passing
                </div>
                <div className="flex items-center gap-2 text-sm text-[#94A3B8]">
                  <span className="w-3 h-3 rounded-full bg-[#00ADD8]" />
                  Go + Python
                </div>
                <div className="flex items-center gap-2 text-sm text-[#94A3B8]">
                  <Scale size={16} className="text-[#F59E0B]" />
                  Apache 2.0
                </div>
                <div className="flex items-center gap-2 text-sm text-[#94A3B8]">
                  <span className="w-2 h-2 rounded-full bg-[#10B981]" />
                  Production Ready
                </div>
              </div>

              <a
                href="https://github.com/shreyanshjain7174/agentic-k8s-operator"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-full bg-[#1E293B] text-white hover:bg-[#334155] transition-colors w-full justify-center"
              >
                <Star size={18} />
                Star on GitHub
              </a>
            </div>
          </AnimatedSection>

          {/* Code Snippet */}
          <AnimatedSection delay={0.2}>
            <div className="glass-card overflow-hidden h-full">
              <div className="flex items-center gap-2 px-4 py-3 bg-[#0A0B14] border-b border-[#1E293B]">
                <div className="w-3 h-3 rounded-full bg-[#EF4444]" />
                <div className="w-3 h-3 rounded-full bg-[#F59E0B]" />
                <div className="w-3 h-3 rounded-full bg-[#10B981]" />
                <span className="ml-2 text-xs text-[#64748B]">
                  agentworkload.yaml
                </span>
              </div>
              <pre className="p-6 text-sm font-mono text-[#94A3B8] overflow-x-auto leading-relaxed">
                <code>{yamlSnippet}</code>
              </pre>
            </div>
          </AnimatedSection>
        </div>
      </div>
    </section>
  );
}
