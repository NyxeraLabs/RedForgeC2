import React from "react";

const reports = [
  { id: "R-001", title: "Recon Summary", updated: "2026-03-10T14:19:00Z", status: "complete" },
  { id: "R-002", title: "Initial Access", updated: "2026-03-12T09:33:00Z", status: "pending" },
  { id: "R-003", title: "Lateral Movement", updated: "2026-03-12T17:05:00Z", status: "pending" },
];

function statusLabel(status: string) {
  if (status === "complete") return "Complete";
  if (status === "pending") return "In progress";
  return status;
}

export function ReportsPage() {
  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Reports</div>
          <div className="page-sub">Mocked report catalog. Reports will be generated from live operations.</div>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>Report Catalog</h2>
          </div>
          <div className="panel-body">
            <table className="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Title</th>
                  <th>Updated</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {reports.map((r) => (
                  <tr key={r.id}>
                    <td>{r.id}</td>
                    <td>{r.title}</td>
                    <td>{new Date(r.updated).toLocaleString()}</td>
                    <td>{statusLabel(r.status)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div style={{ marginTop: 12, color: "rgba(235, 235, 245, 0.7)" }}>
              Report generation is currently mocked and will be replaced with data driven by real tasking and telemetry.
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
