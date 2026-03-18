import React from 'react';
import { motion } from 'framer-motion';
import { CheckCircle2, XCircle, Zap } from 'lucide-react';

const Comparison = () => {
  const containerVariants = {
    hidden: {},
    visible: {
      transition: {
        staggerChildren: 0.12,
        delayChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 28 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.55, ease: 'easeOut' },
    },
  };

  const headerVariants = {
    hidden: { opacity: 0, y: -20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.6, ease: 'easeOut' },
    },
  };

  const features = [
    {
      name: 'Activation Speed',
      agentic: '5-min quickstart',
      orkes: '15+ min setup',
      agenticWins: true,
      description: 'Time to first agent',
    },
    {
      name: 'Cost Transparency',
      agentic: 'Cost-aware routing + quotas',
      orkes: 'Opaque usage meters',
      agenticWins: true,
      description: 'Budget control',
    },
    {
      name: 'Deployment Model',
      agentic: 'Pure open-source self-managed',
      orkes: 'SaaS-first proprietary',
      agenticWins: true,
      description: 'Full control & ownership',
    },
    {
      name: 'Agent Control',
      agentic: 'Full policy isolation + DAG design',
      orkes: 'Limited workflow controls',
      agenticWins: true,
      description: 'Advanced orchestration',
    },
    {
      name: 'Enterprise Support',
      agentic: 'Open community + enterprise package',
      orkes: 'SaaS lock-in model',
      agenticWins: true,
      description: 'Flexibility & opt-in support',
    },
  ];

  return (
    <section id="comparison" className="py-24 px-6 bg-gradient-to-b from-[#05080f] to-[#0a0e1a] overflow-hidden">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <motion.div
          className="mb-16 text-center"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.3 }}
          variants={headerVariants}
        >
          <div className="inline-flex items-center gap-2 mb-4 px-3 py-1.5 rounded-full bg-gradient-to-r from-[#00d4aa]/20 to-[#6366f1]/20 border border-[#00d4aa]/40">
            <Zap size={16} className="text-[#00d4aa]" />
            <span className="text-sm font-semibold text-[#00d4aa]">Competitive Advantage</span>
          </div>
          <h2 className="text-4xl md:text-5xl font-bold text-[#e2e8f0] mb-4 leading-tight">
            Why enterprises choose{' '}
            <span className="bg-gradient-to-r from-[#00d4aa] to-[#6366f1] bg-clip-text text-transparent">
              Agentic Operator
            </span>
          </h2>
          <p className="text-[#94a3b8] text-lg max-w-2xl mx-auto">
            Open-source flexibility with enterprise-grade control. No lock-in, no surprises.
          </p>
        </motion.div>

        {/* Comparison Table */}
        <motion.div
          className="space-y-4 mb-12"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.2 }}
          variants={containerVariants}
        >
          {/* Header Row */}
          <motion.div 
            className="grid grid-cols-[1fr_1fr_1fr] gap-6 mb-6 px-6 py-4 rounded-lg bg-[#0f1419]/50 border border-[#1e293b]/50"
            variants={itemVariants}
          >
            <div className="text-[#94a3b8] font-semibold text-sm">Feature</div>
            <div className="text-center">
              <div className="text-[#00d4aa] font-bold text-base">Agentic Operator</div>
              <div className="text-[#6a7d93] text-xs mt-1">Open Source</div>
            </div>
            <div className="text-center">
              <div className="text-[#94a3b8] font-bold text-base">Orkes</div>
              <div className="text-[#6a7d93] text-xs mt-1">Proprietary</div>
            </div>
          </motion.div>

          {/* Feature Rows */}
          {features.map((feature, idx) => (
            <motion.div
              key={idx}
              className="grid grid-cols-[1fr_1fr_1fr] gap-6 px-6 py-5 rounded-lg bg-gradient-to-r from-[#0f1419]/40 to-[#0f1419]/20 border border-[#1e293b]/30 hover:border-[#1e293b]/60 transition-colors"
              variants={itemVariants}
            >
              {/* Feature Name */}
              <div>
                <div className="text-[#e2e8f0] font-semibold text-sm">{feature.name}</div>
                <div className="text-[#6a7d93] text-xs mt-1.5">{feature.description}</div>
              </div>

              {/* Agentic Column */}
              <div className="flex items-center gap-3">
                <CheckCircle2 size={20} className="text-[#00d4aa] flex-shrink-0" />
                <span className="text-[#e2e8f0] text-sm">{feature.agentic}</span>
              </div>

              {/* Orkes Column */}
              <div className="flex items-center gap-3">
                <XCircle size={20} className="text-[#94a3b8] flex-shrink-0" />
                <span className="text-[#94a3b8] text-sm">{feature.orkes}</span>
              </div>
            </motion.div>
          ))}
        </motion.div>

        {/* CTA Section */}
        <motion.div
          className="mt-16 flex flex-col sm:flex-row items-center justify-center gap-6"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.3 }}
          variants={{
            hidden: { opacity: 0, y: 20 },
            visible: {
              opacity: 1,
              y: 0,
              transition: { duration: 0.6, ease: 'easeOut' },
            },
          }}
        >
          <a
            href="#quickstart"
            className="px-8 py-3.5 rounded-lg bg-gradient-to-r from-[#00d4aa] to-[#00b88d] text-[#05080f] font-semibold text-base hover:shadow-lg hover:shadow-[#00d4aa]/20 transition-all duration-300 hover:scale-105"
          >
            Start in 5 Minutes
          </a>
          <a
            href="#"
            className="px-8 py-3.5 rounded-lg border-2 border-[#1e293b] text-[#e2e8f0] font-semibold text-base hover:border-[#00d4aa]/60 hover:text-[#00d4aa] transition-all duration-300"
          >
            View Full Comparison
          </a>
        </motion.div>

        {/* Trust Statement */}
        <motion.div
          className="mt-12 text-center"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.3 }}
          variants={{
            hidden: { opacity: 0 },
            visible: {
              opacity: 1,
              transition: { duration: 0.8, delay: 0.3 },
            },
          }}
        >
          <p className="text-[#6a7d93] text-sm">
            Join teams running <span className="text-[#00d4aa] font-semibold">1000+</span> agentic workflows in production
          </p>
        </motion.div>
      </div>
    </section>
  );
};

export default Comparison;
