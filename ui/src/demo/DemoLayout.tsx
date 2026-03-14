import React from "react";
import { Link, NavLink, Outlet, useLocation } from "react-router-dom";
import logoUrl from "../assets/RedForgeC2-Logo-Transp.png";
import styles from "./DemoLayout.module.css";

function navClass({ isActive }: { isActive: boolean }) {
  return `${styles.navItem} ${isActive ? styles.navItemActive : ""}`;
}

export function DemoLayout() {
  const loc = useLocation();
  const title = loc.pathname.startsWith("/demo/map")
    ? "Map"
    : loc.pathname.startsWith("/demo/telemetry")
      ? "Telemetry"
      : loc.pathname.startsWith("/demo/reports")
        ? "Reports"
        : loc.pathname.startsWith("/demo/settings")
          ? "Settings"
          : "Dashboard";

  return (
    <div className={styles.scope}>
      <div className={styles.chrome}>
        <aside className={styles.sidebar}>
          <div className={styles.brand}>
            <img className={styles.logo} src={logoUrl} alt="RedForgeC2" />
            <div className={styles.brandText}>
              <div className={styles.brandTitle}>RedForgeC2</div>
              <div className={styles.brandSub}>Demo mockup (static)</div>
            </div>
          </div>

          <nav className={styles.nav} aria-label="Demo navigation">
            <NavLink to="/demo" end className={navClass}>
              Dashboard
            </NavLink>
            <NavLink to="/demo/telemetry" className={navClass}>
              Telemetry
            </NavLink>
            <NavLink to="/demo/map" className={navClass}>
              Map
            </NavLink>
            <NavLink to="/demo/reports" className={navClass}>
              Reports
            </NavLink>
            <NavLink to="/demo/settings" className={navClass}>
              Settings
            </NavLink>
          </nav>

          <div className={styles.hint}>
            Demo-only: mock data, no API calls, no operational controls.
            <div style={{ marginTop: 10 }}>
              <Link to="/login" className={styles.navItem}>
                Back to Login
              </Link>
            </div>
          </div>
        </aside>

        <main className={styles.main}>
          <header className={styles.head}>
            <div>
              <h1 className={styles.headTitle}>{title}</h1>
              <div className={styles.headMeta}>Mockup skin inspired by `docs/ui-mockups.png`</div>
            </div>
            <div className={styles.chips}>
              <span className={styles.chip}>MODE: DEMO</span>
              <span className={styles.chip}>DATA: STATIC</span>
              <span className={styles.chip}>BUILD: UI</span>
            </div>
          </header>

          <div style={{ marginTop: 16 }}>
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}

