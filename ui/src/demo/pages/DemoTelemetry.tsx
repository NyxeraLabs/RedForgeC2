import React from "react";
import { demoSeries } from "../data";
import { Sparkline } from "../widgets/Sparkline";

export function DemoTelemetry() {
  return (
    <div style={{ display: "grid", gap: 16 }}>
      <Sparkline values={demoSeries(80, 0.35)} width={860} height={130} stroke="rgba(255,0,51,0.82)" />
      <Sparkline values={demoSeries(80, 0.55)} width={860} height={130} stroke="rgba(245,158,11,0.82)" />
      <Sparkline values={demoSeries(80, 0.25)} width={860} height={130} stroke="rgba(34,197,94,0.82)" />
    </div>
  );
}

