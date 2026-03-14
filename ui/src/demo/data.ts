export type DemoSession = {
  id: string;
  hostname: string;
  user: string;
  osarch: string;
  lastBeacon: string;
  sev: "ok" | "warn" | "bad";
  cpu: number;
  ramGb: number;
  netMbps: number;
};

export const demoSessions: DemoSession[] = [
  { id: "001", hostname: "ALPHA-LAB", user: "SYSTEM", osarch: "WIN64", lastBeacon: "3 sec ago", sev: "ok", cpu: 32, ramGb: 4.2, netMbps: 23 },
  { id: "002", hostname: "BRAVO-LAB", user: "root", osarch: "LINUX_AMD64", lastBeacon: "12 sec ago", sev: "warn", cpu: 48, ramGb: 7.9, netMbps: 8 },
  { id: "003", hostname: "CHARLIE-LAB", user: "user", osarch: "WIN64", lastBeacon: "41 sec ago", sev: "ok", cpu: 19, ramGb: 3.6, netMbps: 11 },
  { id: "004", hostname: "DELTA-LAB", user: "svc", osarch: "WIN64", lastBeacon: "5 min ago", sev: "bad", cpu: 88, ramGb: 12.6, netMbps: 0.3 },
];

export type DemoAlert = { name: string; sev: "ok" | "warn" | "bad"; detail: string };
export const demoAlerts: DemoAlert[] = [
  { name: "SEID", sev: "bad", detail: "High-confidence signal" },
  { name: "METTY", sev: "warn", detail: "Anomalous pattern" },
  { name: "BEEB", sev: "ok", detail: "Benign indicator" },
  { name: "HOLLOW", sev: "warn", detail: "Heuristic hit" },
];

export function demoSeries(n = 36, seed = 0.35) {
  const out: number[] = [];
  let v = seed;
  for (let i = 0; i < n; i++) {
    v += (Math.sin(i / 2.7) * 0.03 + Math.cos(i / 6.3) * 0.02);
    v += (i % 9 === 0 ? 0.06 : 0) - (i % 13 === 0 ? 0.05 : 0);
    v = Math.max(0.05, Math.min(0.95, v));
    out.push(v);
  }
  return out;
}

