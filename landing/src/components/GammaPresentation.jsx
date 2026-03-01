import { motion } from 'framer-motion';

const containerVariants = {
  hidden: { opacity: 0, y: 40 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.7, ease: 'easeOut', staggerChildren: 0.2 },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 24 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6, ease: 'easeOut' } },
};

const stats = [
  { label: '4 min to first report' },
  { label: 'Zero infra changes required' },
];

export default function GammaPresentation() {
  return (
    <section
      id="demo"
      style={{ background: '#05080f' }}
      className="py-24 px-6 overflow-hidden"
    >
      <motion.div
        className="max-w-5xl mx-auto"
        variants={containerVariants}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, amount: 0.2 }}
      >
        {/* Heading */}
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

        {/* Subtitle */}
        <motion.p
          variants={itemVariants}
          className="text-center text-lg mb-12"
          style={{ color: 'rgba(255,255,255,0.55)', fontFamily: "'DM Sans', sans-serif" }}
        >
          Watch the agentic operator handle real-world Kubernetes workloads
        </motion.p>

        {/* iframe container */}
        <motion.div
          variants={itemVariants}
          style={{ position: 'relative', display: 'flex', justifyContent: 'center' }}
        >
          {/* Ambient glow */}
          <div
            className="absolute inset-0 rounded-full blur-3xl pointer-events-none"
            style={{ background: 'rgba(0,212,170,0.08)' }}
          />
          {/* Outer ring glow */}
          <div
            className="absolute rounded-2xl pointer-events-none"
            style={{
              inset: '-2px',
              background: 'linear-gradient(135deg, rgba(0,212,170,0.3) 0%, rgba(99,102,241,0.15) 50%, rgba(0,212,170,0.1) 100%)',
              borderRadius: '18px',
              filter: 'blur(1px)',
            }}
          />
          <iframe
            src="https://gamma.app/embed/g53pjztg8z13h71"
            style={{
              width: '100%',
              maxWidth: 800,
              height: 500,
              borderRadius: 16,
              border: '1px solid rgba(0,212,170,0.2)',
              position: 'relative',
              zIndex: 10,
              display: 'block',
            }}
            allow="fullscreen"
            title="Agentic Operator"
          />
        </motion.div>

        {/* Quick stats */}
        <motion.div
          variants={itemVariants}
          className="flex flex-col sm:flex-row items-center justify-center gap-4 mt-10"
        >
          {stats.map((stat, i) => (
            <div key={i} className="flex items-center gap-2">
              {i > 0 && (
                <span
                  className="hidden sm:block w-1 h-1 rounded-full"
                  style={{ background: 'rgba(0,212,170,0.4)' }}
                />
              )}
              <span
                className="text-sm font-medium"
                style={{
                  color: '#00d4aa',
                  fontFamily: "'IBM Plex Mono', monospace",
                }}
              >
                {stat.label}
              </span>
            </div>
          ))}
        </motion.div>
      </motion.div>
    </section>
  );
}
