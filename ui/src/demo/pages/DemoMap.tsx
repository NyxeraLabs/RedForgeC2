import React from "react";
import { TacticalMap } from "../widgets/TacticalMap";

export function DemoMap() {
  const nodes = [
    { x: 0.22, y: 0.28, label: "ALPHA", sev: "ok" as const },
    { x: 0.46, y: 0.42, label: "BRAVO", sev: "warn" as const },
    { x: 0.70, y: 0.36, label: "CHARLIE", sev: "ok" as const },
    { x: 0.62, y: 0.70, label: "DELTA", sev: "bad" as const },
  ];
  return <TacticalMap nodes={nodes} width={940} height={520} />;
}

