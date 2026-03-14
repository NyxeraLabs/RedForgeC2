import React from "react";

const nodes = [
  { id: "A1", label: "Gateway", status: "ok" },
  { id: "B2", label: "Win10-01", status: "warn" },
  { id: "C3", label: "Linux-02", status: "ok" },
  { id: "D4", label: "DC01", status: "bad" },
];

function statusColor(status: string) {
  if (status === "ok") return "#22c55e";
  if (status === "warn") return "#f59e0b";
  return "#ef4444";
}

export function MapPage() {
  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Map</div>
          <div className="page-sub">Network topology (mocked). Real map data in future releases.</div>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>Topology</h2>
          </div>
          <div className="panel-body" style={{ display: "grid", gap: 14, padding: 18 }}>
            <div
              style={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fit, minmax(160px, 1fr))",
                gap: 12,
              }}
            >
              {nodes.map((n) => (
                <div
                  key={n.id}
                  style={{
                    padding: 14,
                    borderRadius: 12,
                    border: `1px solid rgba(255,255,255,0.12)`,
                    background: "rgba(15, 15, 20, 0.55)",
                    display: "flex",
                    flexDirection: "column",
                    gap: 10,
                  }}
                >
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div style={{ fontWeight: 700 }}>{n.label}</div>
                    <span
                      style={{
                        borderRadius: 999,
                        padding: "2px 10px",
                        fontSize: "0.75rem",
                        background: "rgba(0,0,0,0.4)",
                        color: statusColor(n.status),
                        border: `1px solid ${statusColor(n.status)}33`,
                      }}
                    >
                      {n.status.toUpperCase()}
                    </span>
                  </div>
                  <div style={{ color: "rgba(235, 235, 245, 0.7)" }}>
                    This is a placeholder node card representing host state. In a full build, this will be replaced
                    with a real-time topology graph with drilldowns.
                  </div>
                </div>
              ))}
            </div>

            <div style={{ color: "rgba(235, 235, 245, 0.7)" }}>
              <strong>Note:</strong> This map is currently mocked. Integration with the Teamserver topology/graph engine is planned.
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
