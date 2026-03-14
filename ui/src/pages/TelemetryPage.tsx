import React from "react";
import { demoSeries } from "../demo/data";
import { Sparkline } from "../demo/widgets/Sparkline";

export function TelemetryPage() {
  const cpu = React.useMemo(() => demoSeries(48, 0.42), []);
  const net = React.useMemo(() => demoSeries(48, 0.28), []);

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Telemetry</div>
          <div className="page-sub">Mock data (real telemetry coming soon)</div>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>CPU</h2>
          </div>
          <div className="panel-body">
            <Sparkline values={cpu} />
            <div style={{ marginTop: 14, color: "rgba(235, 235, 245, 0.7)" }}>
              This view is mocked for the current prototype. Live agent telemetry will be streamed here when the
              telemetry pipeline is implemented.
            </div>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Network</h2>
          </div>
          <div className="panel-body">
            <Sparkline values={net} stroke="rgba(255,255,255,0.65)" />
            <div style={{ marginTop: 14, color: "rgba(235, 235, 245, 0.7)" }}>
              Mock network usage based on demo data. When implemented, this will show per-agent traffic and transfer rates.
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
