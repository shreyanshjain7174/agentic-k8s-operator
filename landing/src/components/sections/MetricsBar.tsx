"use client";

import { motion, useInView } from "framer-motion";
import { useRef, useState, useEffect } from "react";
import { Zap, Clock, TrendingUp } from "lucide-react";

interface MetricProps {
  value: string;
  before: string;
  label: string;
  icon: React.ReactNode;
  delay: number;
  color: string;
}

function AnimatedMetric({ value, before, label, icon, delay, color }: MetricProps) {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <motion.div
      ref={ref}
      initial={{ opacity: 0, y: 30 }}
      animate={isInView ? { opacity: 1, y: 0 } : {}}
      transition={{ duration: 0.5, delay }}
      className="glass-card p-6 text-center"
    >
      <div
        className="w-12 h-12 rounded-xl mx-auto mb-4 flex items-center justify-center"
        style={{ backgroundColor: `${color}15` }}
      >
        <div style={{ color }}>{icon}</div>
      </div>
      <div className="text-3xl font-bold text-white mb-1">{value}</div>
      <div className="text-sm text-[#64748B] mb-2">vs {before}</div>
      <div className="text-sm text-[#94A3B8]">{label}</div>
    </motion.div>
  );
}

export function MetricsBar() {
  return (
    <section className="py-16 relative">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <AnimatedMetric
            value="4 min"
            before="8 hours"
            label="Competitive Intelligence"
            icon={<Zap size={24} />}
            delay={0}
            color="#6366F1"
          />
          <AnimatedMetric
            value="47 sec"
            before="22 min"
            label="Autonomous Remediation"
            icon={<Clock size={24} />}
            delay={0.15}
            color="#06B6D4"
          />
          <AnimatedMetric
            value="90 sec"
            before="6 hours"
            label="Research Swarm"
            icon={<TrendingUp size={24} />}
            delay={0.3}
            color="#10B981"
          />
        </div>
      </div>
    </section>
  );
}
