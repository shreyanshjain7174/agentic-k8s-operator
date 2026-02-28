"use client";

import dynamic from "next/dynamic";
import { motion } from "framer-motion";
import { ArrowRight, Github } from "lucide-react";

const RemotionPlayer = dynamic(
  () => import("./HeroPlayer"),
  { ssr: false, loading: () => <div className="w-full h-full bg-[#0F1120] rounded-2xl animate-pulse" /> }
);

export function Hero() {
  return (
    <section className="relative min-h-screen flex items-center pt-16 overflow-hidden">
      {/* Background gradient mesh */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-[#6366F1]/10 rounded-full blur-[120px]" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-[#8B5CF6]/10 rounded-full blur-[120px]" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-[#06B6D4]/5 rounded-full blur-[150px]" />
      </div>

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left column - Text */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.7, delay: 0.2 }}
          >
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full border border-[#10B981]/30 bg-[#10B981]/10 mb-6">
              <span className="w-2 h-2 rounded-full bg-[#10B981] animate-pulse" />
              <span className="text-xs text-[#10B981] font-medium">
                Apache 2.0 Open Source
              </span>
            </div>

            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold leading-tight mb-6">
              Deploy{" "}
              <span className="gradient-text">AI Agents</span>{" "}
              Inside Your Kubernetes Cluster
            </h1>

            <p className="text-lg text-[#94A3B8] mb-8 max-w-lg">
              Production-grade operator for orchestrating autonomous AI agent
              workloads. One Helm command. Zero data leaves your infrastructure.
            </p>

            <div className="flex flex-wrap gap-4 mb-10">
              <a
                href="#waitlist"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-full bg-gradient-to-r from-[#6366F1] to-[#8B5CF6] text-white font-medium hover:shadow-lg hover:shadow-[#6366F1]/25 transition-all"
              >
                Join the Waitlist
                <ArrowRight size={18} />
              </a>
              <a
                href="https://github.com/shreyanshjain7174/agentic-k8s-operator"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 rounded-full border border-[#334155] text-[#94A3B8] hover:text-white hover:border-[#6366F1]/50 transition-all"
              >
                <Github size={18} />
                View on GitHub
              </a>
            </div>

            {/* Micro stats */}
            <div className="flex flex-wrap gap-6 text-sm">
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-[#10B981]" />
                <span className="text-[#94A3B8]">47/47 pods healthy</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-[#3B82F6]" />
                <span className="text-[#94A3B8]">46 tests passing</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-[#8B5CF6]" />
                <span className="text-[#94A3B8]">Production ready</span>
              </div>
            </div>
          </motion.div>

          {/* Right column - Animation */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.7, delay: 0.4 }}
            className="hidden lg:block"
          >
            <div className="relative w-full aspect-[4/3] rounded-2xl overflow-hidden glow-effect">
              <RemotionPlayer />
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
