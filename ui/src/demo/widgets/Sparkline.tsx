import React from "react";

export function Sparkline({
  values,
  width = 420,
  height = 92,
  stroke = "rgba(255,0,51,0.85)",
}: {
  values: number[];
  width?: number;
  height?: number;
  stroke?: string;
}) {
  const pad = 6;
  const w = Math.max(10, width);
  const h = Math.max(10, height);
  const max = Math.max(...values, 1);
  const min = Math.min(...values, 0);
  const span = Math.max(1e-6, max - min);

  const pts = values
    .map((v, i) => {
      const x = pad + (i * (w - pad * 2)) / Math.max(1, values.length - 1);
      const y = pad + ((max - v) * (h - pad * 2)) / span;
      return `${x.toFixed(2)},${y.toFixed(2)}`;
    })
    .join(" ");

  return (
    <svg width={w} height={h} viewBox={`0 0 ${w} ${h}`} role="img" aria-label="Telemetry sparkline">
      <defs>
        <linearGradient id="rg" x1="0" x2="0" y1="0" y2="1">
          <stop offset="0" stopColor="rgba(255,0,51,0.25)" />
          <stop offset="1" stopColor="rgba(0,0,0,0)" />
        </linearGradient>
      </defs>
      <rect x="0" y="0" width={w} height={h} fill="rgba(0,0,0,0.12)" rx="12" />
      <path d={`M ${pts}`} fill="none" stroke={stroke} strokeWidth="2.25" />
      <path d={`M ${pts} L ${w - pad} ${h - pad} L ${pad} ${h - pad} Z`} fill="url(#rg)" />
    </svg>
  );
}

