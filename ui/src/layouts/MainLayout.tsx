import React from "react";
import { Link, NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { getMe } from "../lib/api";
import { setToken } from "../lib/storage";
import logoUrl from "../assets/RedForgeC2-Logo-Transp.png";
import styles from "../demo/DemoLayout.module.css";

function navClass({ isActive }: { isActive: boolean }) {
  return `${styles.navItem} ${isActive ? styles.navItemActive : ""}`;
}

export function MainLayout() {
  const nav = useNavigate();
  const loc = useLocation();
  const [role, setRole] = React.useState<"admin" | "operator" | "observer" | null>(null);
  const [username, setUsername] = React.useState<string | null>(null);

  React.useEffect(() => {
    getMe()
      .then((m) => {
        setRole(m.role);
        setUsername(m.username);
      })
      .catch(() => {
        setRole(null);
        setUsername(null);
      });
  }, []);

  const title = loc.pathname.startsWith("/profile")
    ? "Profile"
    : loc.pathname.startsWith("/admin")
    ? "Administration"
    : "Dashboard";

  function logout() {
    // clear token and navigate to login
    setToken(null);
    nav("/login", { replace: true });
  }

  return (
    <div className={styles.scope}>
      <div className={styles.chrome}>
        <aside className={styles.sidebar}>
          <div className={styles.brand}>
            <img className={styles.logo} src={logoUrl} alt="RedForgeC2" />
            <div className={styles.brandText}>
              <div className={styles.brandTitle}>RedForgeC2</div>
              <div className={styles.brandSub}>Operator Console</div>
            </div>
          </div>

          <nav className={styles.nav} aria-label="Primary navigation">
            <NavLink to="/" end className={navClass}>
              Dashboard
            </NavLink>
            <NavLink to="/telemetry" className={navClass}>
              Telemetry
            </NavLink>
            <NavLink to="/map" className={navClass}>
              Map
            </NavLink>
            <NavLink to="/reports" className={navClass}>
              Reports
            </NavLink>
            <NavLink to="/settings" className={navClass}>
              Settings
            </NavLink>
            <NavLink to="/profile" className={navClass}>
              Profile
            </NavLink>
            {role === "admin" ? (
              <NavLink to="/admin/users" className={navClass}>
                Users
              </NavLink>
            ) : null}
          </nav>

          <div className={styles.hint}>
            <div style={{ marginBottom: 8 }}>
              {username ? `${username} (${role ?? "?"})` : "Loading profile..."}
            </div>
            <button className={styles.navItem} onClick={logout} style={{ marginTop: 8 }}>
              Logout
            </button>
          </div>
        </aside>

        <main className={styles.main}>
          <header className={styles.head}>
            <div>
              <h1 className={styles.headTitle}>{title}</h1>
              <div className={styles.headMeta}>Real-time data from your teamserver</div>
            </div>
            <div className={styles.chips}>
              <span className={styles.chip}>MODE: LIVE</span>
              <span className={styles.chip}>DATA: API</span>
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
