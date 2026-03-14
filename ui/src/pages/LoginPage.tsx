import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import logoUrl from "../assets/RedForgeC2-Logo-Transp.png";
import { login } from "../lib/api";
import { setApiBase, setToken } from "../lib/storage";
import { useToasts } from "../components/ToastProvider";
import { BuildStamp } from "../components/BuildStamp";
import styles from "./LoginPage.module.css";

function defaultApiBaseFromLocation(): string {
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  return `${protocol}//${hostname}:9080`;
}

export function LoginPage() {
  const nav = useNavigate();
  const toasts = useToasts();

  const [apiBase, setApiBaseState] = useState(() => window.localStorage.getItem("redforge_api_base") || defaultApiBaseFromLocation());
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);

  async function onLogin() {
    setBusy(true);
    try {
      setApiBase(apiBase);
      const token = await login(username, password);
      setToken(token);
      toasts.push({ kind: "good", title: "Authenticated", message: `Welcome, ${username}` });
      nav("/", { replace: true });
    } catch (e) {
      toasts.push({ kind: "bad", title: "Login failed", message: String(e) });
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className={styles.scope}>
      <div className={styles.card}>
        <div className={styles.brandRow}>
          <img className={styles.logo} src={logoUrl} alt="RedForgeC2" />
          <div>
            <h1 className={styles.title}>RedForgeC2</h1>
            <div className={styles.sub}>Operator Console</div>
          </div>
        </div>

        <div className={styles.form}>
          <label className={styles.fieldLabel}>
            Teamserver URL
            <input
              className={styles.input}
              value={apiBase}
              onChange={(e) => setApiBaseState(e.target.value)}
              placeholder="http://127.0.0.1:9080"
            />
          </label>
          <label className={styles.fieldLabel}>
            Username
            <input className={styles.input} value={username} onChange={(e) => setUsername(e.target.value)} />
          </label>
          <label className={styles.fieldLabel}>
            Password
            <input className={styles.input} type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          </label>
          <div className={styles.actions}>
            <button className={styles.primary} disabled={busy || !username || !password || !apiBase} onClick={onLogin}>
              {busy ? "Authenticating..." : "Login"}
            </button>
          </div>
        </div>

        <div className={styles.hint}>
          Default admin credentials are controlled by server env vars. <BuildStamp className={styles.buildStamp} />
        </div>
      </div>
    </div>
  );
}
