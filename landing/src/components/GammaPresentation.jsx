import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Play, Presentation } from 'lucide-react';

const VIDEOS = [
  {
    id: 'operator-demo',
    title: 'Operator Demo',
    label: '01',
    description: 'Watch the AgentWorkload CRD spin up a live hedge-fund intelligence pipeline',
    src: '/videos/agentic-operator-demo.mp4',
    badge: 'Live Run',
  },
  {
    id: 'full-walkthrough',
    title: 'Full Walkthrough',
    label: '02',
    description: 'End-to-end tour — CRD authoring, MCP integration, Argo scheduling, OPA policies',
    src: '/videos/demo-video.mp4',
    badge: 'Architecture',
  },
];

const containerVariants = {
  hidden: { opacity: 0, y: 40 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.7, ease: 'easeOut', staggerChildren: 0.18 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 24 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: 'easeOut' } },
};

function VideoPlayer({ video, isActive }) {
  return (
    <AnimatePresence mode="wait">
      {isActive && (
        <motion.div
          key={video.id}
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -12 }}
          transition={{ duration: 0.35, ease: 'easeOut' }}
          className="rounded-xl overflow-hidden"
          style={{
            border: '1px solid rgba(0,212,170,0.22)',
            boxShadow: '0 0 60px rgba(0,212,170,0.06), 0 8px 40px rgba(0,0,0,0.5)',
          }}
        >
          {/* Terminal-style title bar */}
          <div
            className="flex items-center gap-2 px-4 py-2.5 border-b"
            style={{
              background: '#0d1117',
              borderColor: 'rgba(255,255,255,0.05)',
            }}
          >
            <div className="w-3 h-3 rounded-full" style={{ background: '#ff5f57' }} />
            <div className="w-3 h-3 rounded-full" style={{ background: '#febc2e' }} />
            <div className="w-3 h-3 rounded-full" style={{ background: '#28c840' }} />
            <span
              className="ml-2 text-xs text-slate-500 tracking-wider"
              style={{ fontFamily: "'IBM Plex Mono', monospace" }}
            >
              agentic-operator — {video.title.toLowerCase().replace(/ /g, '-')}.mp4
            </span>
            <span
              className="ml-auto text-xs px-2 py-0.5 rounded-full"
              style={{
                color: '#00d4aa',
                background: 'rgba(0,212,170,0.1)',
                border: '1px solid rgba(0,212,170,0.2)',
                fontFamily: "'IBM Plex Mono', monospace",
              }}
            >
              {video.badge}
            </span>
          </div>

          {/* Video */}
          <div style={{ background: '#000' }}>
            <video
              key={video.src}
              src={video.src}
              controls
              preload="metadata"
              playsInline
              className="w-full block"
              style={{ maxHeight: 420, display: 'block' }}
            />
          </div>

          {/* Description bar */}
          <div
            className="px-4 py-3"
            style={{
              background: 'rgba(10,14,26,0.95)',
              borderTop: '1px solid rgba(255,255,255,0.04)',
            }}
          >
            <p
              className="text-sm"
              style={{
                color: 'rgba(255,255,255,0.5)',
                fontFamily: "'DM Sans', sans-serif",
              }}
            >
              {video.description}
            </p>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}

export default function GammaPresentation() {
  const [activeVideo, setActiveVideo] = useState(0);
  const [showDeck, setShowDeck] = useState(false);

  return (
    <section
      id="demo"
      style={{ background: '#05080f' }}
      className="py-24 px-6 overflow-hidden"
    >
      <motion.div
        className="max-w-4xl mx-auto"
        variants={containerVariants}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, amount: 0.15 }}
      >
        {/* ── Header ── */}
        <motion.div variants={itemVariants} className="text-center mb-4">
          <span
            className="inline-block text-xs font-semibold tracking-widest uppercase mb-4 px-3 py-1 rounded-full"
            style={{
              color: '#00d4aa',
              background: 'rgba(0,212,170,0.08)',
              border: '1px solid rgba(0,212,170,0.2)',
              fontFamily: "'IBM Plex Mono', monospace",
            }}
          >
            Live Demo
          </span>
          <h2
            className="text-4xl md:text-5xl font-bold text-white"
            style={{ fontFamily: "'Syne', sans-serif" }}
          >
            See It In Action
          </h2>
        </motion.div>

        <motion.p
          variants={itemVariants}
          className="text-center text-lg mb-10"
          style={{ color: 'rgba(255,255,255,0.5)', fontFamily: "'DM Sans', sans-serif" }}
        >
          Real kubectl output. Real cluster. Real agents.
        </motion.p>

        {/* ── Video tabs ── */}
        <motion.div variants={itemVariants} className="mb-4">
          <div
            className="inline-flex rounded-xl p-1 mb-6 w-full"
            style={{
              background: 'rgba(13,21,37,0.7)',
              border: '1px solid rgba(255,255,255,0.07)',
            }}
          >
            {VIDEOS.map((video, i) => (
              <button
                key={video.id}
                onClick={() => setActiveVideo(i)}
                className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition-all duration-250"
                style={{
                  background: activeVideo === i ? 'rgba(0,212,170,0.12)' : 'transparent',
                  border: activeVideo === i ? '1px solid rgba(0,212,170,0.25)' : '1px solid transparent',
                  color: activeVideo === i ? '#00d4aa' : 'rgba(255,255,255,0.4)',
                  fontFamily: "'DM Sans', sans-serif",
                  cursor: 'pointer',
                }}
              >
                <Play
                  size={13}
                  style={{ opacity: activeVideo === i ? 1 : 0.4 }}
                />
                <span
                  className="mr-1 text-xs opacity-50"
                  style={{ fontFamily: "'IBM Plex Mono', monospace" }}
                >
                  {video.label}
                </span>
                {video.title}
              </button>
            ))}
          </div>

          {/* Video players */}
          <div className="relative min-h-[300px]">
            {VIDEOS.map((video, i) => (
              <VideoPlayer key={video.id} video={video} isActive={activeVideo === i} />
            ))}
          </div>
        </motion.div>

        {/* ── Quick stats ── */}
        <motion.div
          variants={itemVariants}
          className="flex flex-wrap items-center justify-center gap-6 mt-8 mb-16"
        >
          {[
            { value: '4 min', label: 'first report' },
            { value: 'Zero', label: 'infra changes' },
            { value: 'Native', label: 'kubectl UX' },
          ].map((s) => (
            <div key={s.label} className="flex items-center gap-2">
              <span
                className="text-sm font-semibold"
                style={{ color: '#00d4aa', fontFamily: "'IBM Plex Mono', monospace" }}
              >
                {s.value}
              </span>
              <span
                className="text-sm"
                style={{ color: 'rgba(255,255,255,0.35)', fontFamily: "'DM Sans', sans-serif" }}
              >
                {s.label}
              </span>
            </div>
          ))}
        </motion.div>

        {/* ── Pitch deck (collapsible) ── */}
        <motion.div variants={itemVariants}>
          <button
            onClick={() => setShowDeck((v) => !v)}
            className="w-full flex items-center justify-between px-5 py-4 rounded-xl transition-all duration-200 group"
            style={{
              background: 'rgba(13,21,37,0.5)',
              border: '1px solid rgba(0,212,170,0.12)',
              cursor: 'pointer',
            }}
          >
            <div className="flex items-center gap-3">
              <div
                className="w-8 h-8 rounded-lg flex items-center justify-center"
                style={{ background: 'rgba(0,212,170,0.1)', border: '1px solid rgba(0,212,170,0.2)' }}
              >
                <Presentation size={15} style={{ color: '#00d4aa' }} />
              </div>
              <div className="text-left">
                <div
                  className="text-sm font-semibold text-white"
                  style={{ fontFamily: "'DM Sans', sans-serif" }}
                >
                  View Full Pitch Deck
                </div>
                <div
                  className="text-xs mt-0.5"
                  style={{ color: 'rgba(255,255,255,0.4)', fontFamily: "'DM Sans', sans-serif" }}
                >
                  Architecture · Use Cases · Roadmap
                </div>
              </div>
            </div>
            <motion.div
              animate={{ rotate: showDeck ? 180 : 0 }}
              transition={{ duration: 0.25 }}
              style={{ color: 'rgba(0,212,170,0.6)' }}
            >
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M4 6l4 4 4-4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </motion.div>
          </button>

          <AnimatePresence initial={false}>
            {showDeck && (
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: 'auto', opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.35, ease: 'easeInOut' }}
                style={{ overflow: 'hidden' }}
              >
                <div className="pt-4 relative">
                  {/* Ambient glow behind iframe */}
                  <div
                    className="absolute inset-0 rounded-full blur-3xl pointer-events-none"
                    style={{ background: 'rgba(0,212,170,0.05)' }}
                  />
                  <div
                    className="rounded-xl overflow-hidden"
                    style={{
                      position: 'relative',
                      border: '1px solid rgba(0,212,170,0.18)',
                      boxShadow: '0 8px 40px rgba(0,0,0,0.4)',
                    }}
                  >
                    <iframe
                      src="https://gamma.app/embed/g53pjztg8z13h71"
                      style={{
                        width: '100%',
                        height: 500,
                        display: 'block',
                        border: 'none',
                      }}
                      allow="fullscreen"
                      title="Agentic Operator — Pitch Deck"
                      loading="lazy"
                    />
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
      </motion.div>
    </section>
  );
}
