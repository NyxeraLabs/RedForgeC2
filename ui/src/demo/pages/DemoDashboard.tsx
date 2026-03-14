import React from "react";
import { demoAlerts, demoSeries, demoSessions } from "../data";
import { Sparkline } from "../widgets/Sparkline";
import styles from "./DemoDashboard.module.css";

function sevClass(sev: "ok" | "warn" | "bad") {
  if (sev === "ok") return styles.ok;
  if (sev === "warn") return styles.warn;
  return styles.bad;
}

export function DemoDashboard() {
  const cpu = demoSeries(42, 0.42);
  const net = demoSeries(42, 0.28);

  return (
    <div>
      <div className={styles.grid}>
        <section className={styles.panel}>
          <div className={styles.panelHead}>
            <div className={styles.panelTitle}>Sessions</div>
            <div className={styles.small}>static</div>
          </div>
          <div className={styles.panelBody}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Hostname</th>
                  <th>User</th>
                  <th>OS/Arch</th>
                  <th>Last Beacon</th>
                </tr>
              </thead>
              <tbody>
                {demoSessions.map((s) => (
                  <tr key={s.id}>
                    <td className={styles.id}>{s.id}</td>
                    <td>{s.hostname}</td>
                    <td>{s.user}</td>
                    <td>{s.osarch}</td>
                    <td className={sevClass(s.sev)}>{s.lastBeacon}</td>
                  </tr>
                ))}
              </tbody>
            </table>

            <div style={{ marginTop: 12 }} className={styles.panelTitle}>
              Telemetry
            </div>
            <div style={{ marginTop: 10 }}>
              <Sparkline values={cpu} />
            </div>
          </div>
        </section>

        <section className={styles.panel}>
          <div className={styles.panelHead}>
            <div className={styles.panelTitle}>Alerts</div>
            <div className={styles.small}>static</div>
          </div>
          <div className={styles.panelBody}>
            <div className={styles.alerts}>
              {demoAlerts.map((a) => (
                <div key={a.name} className={styles.alert}>
                  <div className={`${styles.alertName} ${sevClass(a.sev)}`}>{a.name}</div>
                  <div className={styles.small}>{a.detail}</div>
                </div>
              ))}
            </div>

            <div style={{ marginTop: 14 }} className={styles.panelTitle}>
              Net
            </div>
            <div style={{ marginTop: 10 }}>
              <Sparkline values={net} stroke="rgba(255,255,255,0.65)" />
            </div>
          </div>
        </section>
      </div>

      <div className={styles.row2}>
        <section className={styles.panel}>
          <div className={styles.panelHead}>
            <div className={styles.panelTitle}>Notes</div>
            <div className={styles.small}>demo</div>
          </div>
          <div className={styles.panelBody}>
            <div className={styles.small}>
              This demo is a mockup skin. It intentionally does not connect to the teamserver and does not provide
              operational controls.
            </div>
          </div>
        </section>
        <section className={styles.panel}>
          <div className={styles.panelHead}>
            <div className={styles.panelTitle}>Graphs</div>
            <div className={styles.small}>demo</div>
          </div>
          <div className={styles.panelBody}>
            <Sparkline values={demoSeries(48, 0.55)} stroke="rgba(255,0,51,0.82)" />
          </div>
        </section>
      </div>
    </div>
  );
}

