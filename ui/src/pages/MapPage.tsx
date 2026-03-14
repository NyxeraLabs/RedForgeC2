import React from "react";

type AgentNode = {
  id: string;
  label: string;
  lat: number;
  lon: number;
  status: "ok" | "warn" | "bad";
  connectedTo: string[];
};

const agents: AgentNode[] = [
  { id: "A1", label: "NYC-01", lat: 40.7128, lon: -74.006, status: "ok", connectedTo: ["B2", "C3"] },
  { id: "B2", label: "LON-01", lat: 51.5074, lon: -0.1278, status: "warn", connectedTo: ["A1"] },
  { id: "C3", label: "SYD-01", lat: -33.8688, lon: 151.2093, status: "ok", connectedTo: ["A1", "D4"] },
  { id: "D4", label: "SFO-01", lat: 37.7749, lon: -122.4194, status: "bad", connectedTo: ["C3"] },
];

function statusColor(status: "ok" | "warn" | "bad") {
  if (status === "ok") return "#22c55e";
  if (status === "warn") return "#f59e0b";
  return "#ef4444";
}

function project(lat: number, lon: number, width: number, height: number) {
  // Simple equirectangular projection for mock UI.
  const x = ((lon + 180) / 360) * width;
  const y = ((90 - lat) / 180) * height;
  return { x, y };
}

export function MapPage() {
  const width = 940;
  const height = 520;

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Map</div>
          <div className="page-sub">Global agent topology (mocked). Location + connectivity shown.</div>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>Global View</h2>
          </div>
          <div className="panel-body" style={{ padding: 18, display: "grid", gap: 14 }}>
            <div
              style={{
                borderRadius: 14,
                background: "rgba(5, 5, 10, 0.55)",
                border: "1px solid rgba(255, 255, 255, 0.12)",
                boxShadow: "0 18px 60px rgba(0,0,0,0.55)",
                overflow: "hidden",
              }}
            >
              <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`}>
                <defs>
                  <linearGradient id="ocean" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="#05101f" />
                    <stop offset="100%" stopColor="#000" />
                  </linearGradient>
                  <linearGradient id="land" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="#1b1f28" />
                    <stop offset="100%" stopColor="#0e0f14" />
                  </linearGradient>
                </defs>

                {/* background */}
                <rect width={width} height={height} fill="url(#ocean)" />
                {/* simplified grid */}
                {[...Array(9)].map((_, i) => (
                  <line
                    key={`lon-${i}`}
                    x1={(i * width) / 8}
                    y1={0}
                    x2={(i * width) / 8}
                    y2={height}
                    stroke="rgba(255,255,255,0.08)"
                    strokeWidth={1}
                  />
                ))}
                {[...Array(5)].map((_, i) => (
                  <line
                    key={`lat-${i}`}
                    x1={0}
                    y1={(i * height) / 4}
                    x2={width}
                    y2={(i * height) / 4}
                    stroke="rgba(255,255,255,0.08)"
                    strokeWidth={1}
                  />
                ))}

                {/* agent edges */}
                {agents.map((agent) => {
                  const start = project(agent.lat, agent.lon, width, height);
                  return agent.connectedTo.map((to) => {
                    const target = agents.find((a) => a.id === to);
                    if (!target) return null;
                    const end = project(target.lat, target.lon, width, height);
                    return (
                      <line
                        key={`${agent.id}-${to}`}
                        x1={start.x}
                        y1={start.y}
                        x2={end.x}
                        y2={end.y}
                        stroke="rgba(255,255,255,0.25)"
                        strokeWidth={2}
                        strokeDasharray="6 4"
                      />
                    );
                  });
                })}

                {/* agent nodes */}
                {agents.map((agent) => {
                  const { x, y } = project(agent.lat, agent.lon, width, height);
                  return (
                    <g key={agent.id}>
                      <circle cx={x} cy={y} r={10} fill={statusColor(agent.status)} stroke="rgba(0,0,0,0.5)" strokeWidth={2} />
                      <text
                        x={x + 14}
                        y={y + 4}
                        fontSize={12}
                        fill="rgba(255,255,255,0.9)"
                        fontFamily="ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif"
                      >
                        {agent.label}
                      </text>
                    </g>
                  );
                })}
              </svg>

              <div style={{ padding: 14, color: "rgba(235, 235, 245, 0.75)" }}>
                Mocked global topology: each node represents an agent location (geo-coded) and dashed lines show known connections.
                In production, this will be populated by real beaconing telemetry and topology correlation.
              </div>
            </div>

            <section className="panel">
              <div className="panel-header">
                <h3>Agents (mock data)</h3>
              </div>
              <div className="panel-body" style={{ padding: 12 }}>
                <table className="table">
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>Location</th>
                      <th>Status</th>
                      <th>Peers</th>
                    </tr>
                  </thead>
                  <tbody>
                    {agents.map((agent) => (
                      <tr key={agent.id}>
                        <td>{agent.label}</td>
                        <td>
                          {agent.lat.toFixed(2)}, {agent.lon.toFixed(2)}
                        </td>
                        <td>
                          <span style={{ color: statusColor(agent.status), fontWeight: 700 }}>{agent.status}</span>
                        </td>
                        <td>{agent.connectedTo.join(", ")}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                <div style={{ marginTop: 10, color: "rgba(235, 235, 245, 0.7)" }}>
                  Note: Coordinates are used for visualization only and do not imply real geolocation data.
                </div>
              </div>
            </section>
          </div>
        </section>
      </main>
    </div>
  );
}
