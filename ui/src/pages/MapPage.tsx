import React, { useMemo, useState } from "react";
import { ComposableMap, Geographies, Geography, Marker } from "react-simple-maps";
import { feature } from "topojson-client";
import countries110m from "world-atlas/countries-110m.json";

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

export function MapPage() {
  const [hovered, setHovered] = useState<AgentNode | null>(null);

  const geoFeatures = useMemo(() => {
    const geo = feature(countries110m as any, (countries110m as any).objects.countries);
    return (geo as any).features || [];
  }, []);

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
                position: "relative",
                borderRadius: 14,
                background: "rgba(5, 5, 10, 0.70)",
                border: "1px solid rgba(255, 255, 255, 0.12)",
                boxShadow: "0 18px 60px rgba(0,0,0,0.55)",
                overflow: "hidden",
              }}
            >
              <ComposableMap
                projection="geoNaturalEarth1"
                width={940}
                height={520}
                projectionConfig={{ scale: 180 }}
                style={{ width: "100%", height: "auto" }}
              >
                <Geographies geography={{ type: "FeatureCollection", features: geoFeatures }}>
                  {({ geographies }) =>
                    geographies.map((geo) => (
                      <Geography
                        key={geo.rsmKey}
                        geography={geo}
                        fill="rgba(255,255,255,0.04)"
                        stroke="rgba(255,255,255,0.15)"
                        strokeWidth={0.5}
                      />
                    ))
                  }
                </Geographies>

                {agents.map((agent) => (
                  <Marker
                    key={agent.id}
                    coordinates={[agent.lon, agent.lat]}
                    onMouseEnter={() => setHovered(agent)}
                    onMouseLeave={() => setHovered(null)}
                  >
                    <circle cx={0} cy={0} r={8} fill={statusColor(agent.status)} stroke="#000" strokeWidth={2} />
                  </Marker>
                ))}
              </ComposableMap>

              {hovered ? (
                <div
                  style={{
                    position: "absolute",
                    top: 18,
                    right: 18,
                    padding: "10px 12px",
                    borderRadius: 12,
                    background: "rgba(0,0,0,0.75)",
                    border: "1px solid rgba(255,255,255,0.18)",
                    color: "rgba(255,255,255,0.9)",
                    fontSize: "0.9rem",
                    maxWidth: 240,
                  }}
                >
                  <div style={{ fontWeight: 700, marginBottom: 4 }}>{hovered.label}</div>
                  <div>Location: {hovered.lat.toFixed(2)}, {hovered.lon.toFixed(2)}</div>
                  <div>Status: {hovered.status.toUpperCase()}</div>
                </div>
              ) : null}

              <div style={{ padding: 14, color: "rgba(235, 235, 245, 0.75)" }}>
                World map is rendered with a dark theme; node markers represent agent locations and hover reveals metadata. Connections are mocked for demo.
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
