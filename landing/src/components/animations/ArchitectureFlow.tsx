"use client";

import { useCurrentFrame, interpolate, spring, useVideoConfig } from "remotion";

const STEPS = [
  { label: "AgentWorkload\nYAML", color: "#6366F1", x: 80 },
  { label: "Operator", color: "#8B5CF6", x: 260 },
  { label: "Argo\nWorkflows", color: "#06B6D4", x: 440 },
  { label: "Agent\nPods", color: "#10B981", x: 620 },
  { label: "Results", color: "#D946EF", x: 800 },
];

function FlowArrow({ x1, x2, delay }: { x1: number; x2: number; delay: number }) {
  const frame = useCurrentFrame();
  const progress = interpolate(frame - delay, [0, 30], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const y = 150;
  const endX = x1 + (x2 - x1) * progress;

  const dotProgress = ((frame - delay) % 40) / 40;
  const dotX = x1 + (x2 - x1) * dotProgress;

  return (
    <g>
      <line x1={x1 + 50} y1={y} x2={endX - 10} y2={y} stroke="#334155" strokeWidth={2} strokeOpacity={progress} />
      {progress > 0.8 && (
        <>
          <polygon
            points={`${x2 - 50},${y - 6} ${x2 - 40},${y} ${x2 - 50},${y + 6}`}
            fill="#334155"
            opacity={progress}
          />
          <circle cx={dotX} cy={y} r={3} fill="#6366F1" opacity={0.8}>
            <animate attributeName="opacity" values="0.4;1;0.4" dur="0.8s" repeatCount="indefinite" />
          </circle>
        </>
      )}
    </g>
  );
}

export function ArchitectureFlow() {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  return (
    <svg viewBox="0 0 900 300" style={{ width: "100%", height: "100%" }}>
      <defs>
        <linearGradient id="archGlow" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stopColor="#6366F1" stopOpacity={0.05} />
          <stop offset="50%" stopColor="#8B5CF6" stopOpacity={0.08} />
          <stop offset="100%" stopColor="#D946EF" stopOpacity={0.05} />
        </linearGradient>
      </defs>

      <rect width="900" height="300" fill="url(#archGlow)" rx={16} />

      {/* Arrows between steps */}
      {STEPS.slice(0, -1).map((step, i) => (
        <FlowArrow key={i} x1={step.x} x2={STEPS[i + 1].x} delay={i * 20 + 10} />
      ))}

      {/* Step nodes */}
      {STEPS.map((step, i) => {
        const scale = spring({ frame: frame - i * 15, fps, config: { damping: 12 } });
        const opacity = interpolate(frame - i * 15, [0, 10], [0, 1], {
          extrapolateLeft: "clamp",
          extrapolateRight: "clamp",
        });

        const lines = step.label.split("\n");

        return (
          <g key={i} transform={`translate(${step.x}, 150) scale(${scale})`} opacity={opacity}>
            <rect
              x={-45}
              y={-35}
              width={90}
              height={70}
              rx={12}
              fill="rgba(20, 22, 40, 0.8)"
              stroke={step.color}
              strokeWidth={1.5}
              strokeOpacity={0.5}
            />
            {lines.map((line, li) => (
              <text
                key={li}
                y={li * 16 - (lines.length - 1) * 8}
                textAnchor="middle"
                fill="#F8FAFC"
                fontSize={12}
                fontFamily="Inter, sans-serif"
                fontWeight={500}
              >
                {line}
              </text>
            ))}
            <circle cx={0} cy={-35} r={4} fill={step.color} opacity={0.8}>
              <animate attributeName="r" values="3;5;3" dur="2s" repeatCount="indefinite" />
            </circle>
          </g>
        );
      })}
    </svg>
  );
}
