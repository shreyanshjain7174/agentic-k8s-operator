"use client";

import { useCurrentFrame, useVideoConfig, interpolate, spring } from "remotion";

const NODE_COLORS = ["#6366F1", "#8B5CF6", "#06B6D4", "#10B981", "#3B82F6", "#D946EF"];

function KubeNode({
  cx,
  cy,
  radius,
  color,
  delay,
  label,
}: {
  cx: number;
  cy: number;
  radius: number;
  color: string;
  delay: number;
  label: string;
}) {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const scale = spring({ frame: frame - delay, fps, config: { damping: 12, mass: 0.5 } });
  const opacity = interpolate(frame - delay, [0, 15], [0, 1], { extrapolateLeft: "clamp", extrapolateRight: "clamp" });
  const floatY = Math.sin((frame + delay * 5) / 30) * 3;

  return (
    <g transform={`translate(${cx}, ${cy + floatY}) scale(${scale})`} opacity={opacity}>
      <circle r={radius} fill={color} opacity={0.15} />
      <circle r={radius * 0.7} fill={color} opacity={0.3} />
      <circle r={radius * 0.4} fill={color} opacity={0.8} />
      <text
        y={radius + 16}
        textAnchor="middle"
        fill="#94A3B8"
        fontSize={11}
        fontFamily="Inter, sans-serif"
      >
        {label}
      </text>
    </g>
  );
}

function Connection({
  x1,
  y1,
  x2,
  y2,
  delay,
  color,
}: {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  delay: number;
  color: string;
}) {
  const frame = useCurrentFrame();
  const progress = interpolate(frame - delay, [0, 40], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const dx = x2 - x1;
  const dy = y2 - y1;
  const endX = x1 + dx * progress;
  const endY = y1 + dy * progress;

  const dotProgress = ((frame - delay) % 60) / 60;
  const dotX = x1 + dx * dotProgress;
  const dotY = y1 + dy * dotProgress;

  return (
    <g>
      <line
        x1={x1}
        y1={y1}
        x2={endX}
        y2={endY}
        stroke={color}
        strokeWidth={1.5}
        strokeOpacity={0.3}
        strokeDasharray="4 4"
      />
      {progress > 0.5 && (
        <circle cx={dotX} cy={dotY} r={3} fill={color} opacity={0.8}>
          <animate attributeName="r" values="2;4;2" dur="1s" repeatCount="indefinite" />
        </circle>
      )}
    </g>
  );
}

export function HeroAnimation() {
  const frame = useCurrentFrame();

  const centerX = 400;
  const centerY = 300;
  const orbitRadius = 160;

  const agents = [
    { angle: 0, label: "Agent 1", color: NODE_COLORS[0] },
    { angle: 60, label: "Agent 2", color: NODE_COLORS[1] },
    { angle: 120, label: "MCP", color: NODE_COLORS[2] },
    { angle: 180, label: "Agent 3", color: NODE_COLORS[3] },
    { angle: 240, label: "Argo", color: NODE_COLORS[4] },
    { angle: 300, label: "Agent 4", color: NODE_COLORS[5] },
  ];

  const rotationOffset = frame * 0.3;

  return (
    <svg viewBox="0 0 800 600" style={{ width: "100%", height: "100%" }}>
      <defs>
        <radialGradient id="bgGlow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#6366F1" stopOpacity={0.08} />
          <stop offset="100%" stopColor="#0A0B14" stopOpacity={0} />
        </radialGradient>
      </defs>

      <rect width="800" height="600" fill="transparent" />
      <circle cx={centerX} cy={centerY} r={250} fill="url(#bgGlow)" />

      {/* Orbit ring */}
      <circle
        cx={centerX}
        cy={centerY}
        r={orbitRadius}
        fill="none"
        stroke="#6366F1"
        strokeWidth={1}
        strokeOpacity={0.15}
        strokeDasharray="8 8"
      />

      {/* Connections from agents to center */}
      {agents.map((agent, i) => {
        const rad = ((agent.angle + rotationOffset) * Math.PI) / 180;
        const ax = centerX + Math.cos(rad) * orbitRadius;
        const ay = centerY + Math.sin(rad) * orbitRadius;
        return (
          <Connection
            key={`conn-${i}`}
            x1={centerX}
            y1={centerY}
            x2={ax}
            y2={ay}
            delay={i * 8}
            color={agent.color}
          />
        );
      })}

      {/* Center operator node */}
      <KubeNode
        cx={centerX}
        cy={centerY}
        radius={35}
        color="#6366F1"
        delay={0}
        label="Operator"
      />

      {/* Orbiting agent nodes */}
      {agents.map((agent, i) => {
        const rad = ((agent.angle + rotationOffset) * Math.PI) / 180;
        const ax = centerX + Math.cos(rad) * orbitRadius;
        const ay = centerY + Math.sin(rad) * orbitRadius;
        return (
          <KubeNode
            key={i}
            cx={ax}
            cy={ay}
            radius={22}
            color={agent.color}
            delay={i * 10 + 5}
            label={agent.label}
          />
        );
      })}

      {/* Status text */}
      <g opacity={interpolate(frame, [60, 80], [0, 1], { extrapolateLeft: "clamp", extrapolateRight: "clamp" })}>
        <text x={centerX} y={520} textAnchor="middle" fill="#10B981" fontSize={13} fontFamily="Inter, sans-serif">
          47/47 pods healthy
        </text>
        <circle cx={centerX - 80} cy={516} r={4} fill="#10B981">
          <animate attributeName="opacity" values="0.4;1;0.4" dur="2s" repeatCount="indefinite" />
        </circle>
      </g>
    </svg>
  );
}
