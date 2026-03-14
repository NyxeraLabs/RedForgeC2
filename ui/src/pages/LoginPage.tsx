import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import logoUrl from "../assets/redforgec2_logo.png";
import { login } from "../lib/api";
import { setApiBase, setToken } from "../lib/storage";
import { useToasts } from "../components/ToastProvider";

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
    <div className="auth">
      <div className="auth-card">
        <div className="auth-brand">
          <img className="brand-logo" src={logoUrl} alt="RedForgeC2" />
          <div>
            <div className="auth-title">RedForgeC2</div>
            <div className="auth-sub">Operator Console</div>
          </div>
        </div>

        <div className="auth-form">
          <label>
            Teamserver URL
            <input value={apiBase} onChange={(e) => setApiBaseState(e.target.value)} placeholder="http://127.0.0.1:9080" />
          </label>
          <label>
            Username
            <input value={username} onChange={(e) => setUsername(e.target.value)} />
          </label>
          <label>
            Password
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          </label>
          <button className="btn btn-primary" disabled={busy || !username || !password || !apiBase} onClick={onLogin}>
            {busy ? "Authenticating..." : "Login"}
          </button>
        </div>

        <div className="auth-hint">Default admin credentials are controlled by server env vars.</div>
      </div>
    </div>
  );
}
