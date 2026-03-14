import React from "react";

type Node = { x: number; y: number; label: string; sev: "ok" | "warn" | "bad" };

function color(sev: Node["sev"]) {
  if (sev === "ok") return "rgba(34,197,94,0.95)";
  if (sev === "warn") return "rgba(245,158,11,0.95)";
  return "rgba(255,0,51,0.95)";
}

export function TacticalMap({
  width = 740,
  height = 420,
  nodes,
}: {
  width?: number;
  height?: number;
  nodes: Node[];
}) {
  const w = Math.max(200, width);
  const h = Math.max(180, height);

  const links = nodes.slice(1).map((n, i) => ({ a: nodes[i], b: n }));

  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} role="img" aria-label="Tactical map (demo)">
      <defs>
        <radialGradient id="glow" cx="50%" cy="0%" r="90%">
          <stop offset="0" stopColor="rgba(255,0,51,0.18)" />
          <stop offset="1" stopColor="rgba(0,0,0,0)" />
        </radialGradient>
        <pattern id="grid" width="28" height="28" patternUnits="userSpaceOnUse">
          <path d="M 28 0 L 0 0 0 28" fill="none" stroke="rgba(255,255,255,0.06)" strokeWidth="1" />
        </pattern>
      </defs>

      <rect x="0" y="0" width={w} height={h} rx="16" fill="rgba(0,0,0,0.18)" stroke="rgba(255,0,51,0.18)" />
      <rect x="0" y="0" width={w} height={h} rx="16" fill="url(#glow)" />
      <rect x="12" y="12" width={w - 24} height={h - 24} rx="12" fill="url(#grid)" />

      {links.map((l, idx) => (
        <line
          key={idx}
          x1={l.a.x * w}
          y1={l.a.y * h}
          x2={l.b.x * w}
          y2={l.b.y * h}
          stroke="rgba(255,0,51,0.25)"
          strokeWidth="1.2"
        />
      ))}

      {nodes.map((n) => (
        <g key={n.label}>
          <circle cx={n.x * w} cy={n.y * h} r="8" fill={color(n.sev)} opacity="0.92" />
          <circle cx={n.x * w} cy={n.y * h} r="18" fill={color(n.sev)} opacity="0.10" />
          <text
            x={n.x * w + 12}
            y={n.y * h - 10}
            fill="rgba(255,255,255,0.82)"
            fontFamily="ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace"
            fontSize="12"
          >
            {n.label}
          </text>
        </g>
      ))}
    </svg>
  );
}

