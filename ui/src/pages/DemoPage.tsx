import React from "react";
import logoUrl from "../assets/RedForgeC2-Logo-Transp.png";
import styles from "./DemoPage.module.css";

type Session = {
  id: string;
  hostname: string;
  user: string;
  osarch: string;
  last: string;
  sev: "ok" | "warn" | "bad";
};

const sessions: Session[] = [
  { id: "001", hostname: "ALPHA-LAB", user: "SYSTEM", osarch: "WIN64", last: "3 sec ago", sev: "ok" },
  { id: "002", hostname: "BRAVO-LAB", user: "root", osarch: "LINUX_AMD64", last: "12 sec ago", sev: "warn" },
  { id: "003", hostname: "CHARLIE-LAB", user: "user", osarch: "WIN64", last: "41 sec ago", sev: "ok" },
];

const alerts = [
  { name: "SEID", sev: "bad" as const },
  { name: "METTY", sev: "warn" as const },
  { name: "BEEB", sev: "ok" as const },
];

function sevClass(sev: "ok" | "warn" | "bad") {
  if (sev === "ok") return styles.ok;
  if (sev === "warn") return styles.warn;
  return styles.bad;
}

export function DemoPage() {
  return (
    <div className={styles.scope}>
      <div className={styles.chrome}>
        <aside className={styles.sidebar}>
          <div className={styles.brand}>
            <img className={styles.logo} src={logoUrl} alt="RedForgeC2" />
            <div className={styles.brandText}>
              <div className={styles.brandTitle}>RedForgeC2</div>
              <div className={styles.brandSub}>Visual demo (static)</div>
            </div>
          </div>

          <nav className={styles.nav} aria-label="Demo navigation">
            <div className={`${styles.navItem} ${styles.navItemActive}`}>Dashboard</div>
            <div className={styles.navItem}>Sessions</div>
            <div className={styles.navItem}>Tasks</div>
            <div className={styles.navItem}>Reports</div>
            <div className={styles.navItem}>Settings</div>
          </nav>

          <div className={styles.small}>
            This page is intentionally non-operational: no API calls, no command submission, mock data only.
          </div>
        </aside>

        <main className={styles.main}>
          <header className={styles.head}>
            <div>
              <h1 className={styles.headTitle}>Operator Dashboard</h1>
              <div className={styles.headMeta}>Obsidian / Crimson high-density mockup skin</div>
            </div>
            <div className={styles.chips}>
              <span className={styles.chip}>SESSIONS: {sessions.length}</span>
              <span className={styles.chip}>ALERTS: {alerts.length}</span>
              <span className={styles.chip}>MODE: DEMO</span>
            </div>
          </header>

          <div className={styles.grid}>
            <section className={styles.panel}>
              <div className={styles.panelHead}>
                <div className={styles.panelTitle}>Sessions</div>
                <div className={styles.kbd}>R refresh</div>
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
                    {sessions.map((s) => (
                      <tr key={s.id}>
                        <td className={styles.id}>{s.id}</td>
                        <td>{s.hostname}</td>
                        <td>{s.user}</td>
                        <td>{s.osarch}</td>
                        <td className={sevClass(s.sev)}>{s.last}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>

                <div style={{ marginTop: 12 }} className={styles.panelTitle}>
                  Telemetry
                </div>
                <div style={{ marginTop: 10 }} className={styles.spark} />
              </div>
            </section>

            <section className={styles.panel}>
              <div className={styles.panelHead}>
                <div className={styles.panelTitle}>Alerts</div>
                <div className={styles.kbd}>Acknowledge</div>
              </div>
              <div className={styles.panelBody}>
                <div className={styles.alerts}>
                  {alerts.map((a) => (
                    <div key={a.name} className={styles.alert}>
                      <div className={`${styles.alertName} ${sevClass(a.sev)}`}>{a.name}</div>
                      <div className={styles.small}>signal</div>
                    </div>
                  ))}
                </div>
              </div>
            </section>
          </div>

          <section className={styles.cmd} style={{ marginTop: 16 }}>
            <div className={styles.panelHead}>
              <div className={styles.panelTitle}>Command Input (visual only)</div>
              <div className={styles.kbd}>Enter</div>
            </div>
            <div className={styles.cmdBody}>
              <div className={styles.cmdLine}>
                <span className={styles.prompt}>&gt;</span>
                <input className={styles.cmdInput} value={'help'} readOnly aria-label="Demo command input" />
              </div>
              <div className={styles.small}>list · use &lt;id&gt; · shell &lt;cmd&gt; · upload · download · persist</div>
            </div>
          </section>
        </main>
      </div>
    </div>
  );
}

