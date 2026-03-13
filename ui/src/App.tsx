import { useEffect, useMemo, useState } from "react";
import logoUrl from "./assets/redforgec2_logo.png";

function defaultApiBaseFromLocation(): string {
  if (typeof window === "undefined") return "http://localhost:9080";
  const protocol = window.location.protocol;
  const hostname = window.location.hostname;
  return `${protocol}//${hostname}:9080`;
}

const savedApiBase = typeof window !== "undefined" ? window.localStorage.getItem("redforge_api_base") : null;

const API_BASE = import.meta.env.VITE_TEAMSERVER_URL || savedApiBase || defaultApiBaseFromLocation();

type Agent = {
  agent_id: string;
  os: string;
  arch: string;
  hostname: string;
  version: string;
  last_seen: string;
  registered: string;
};

type TaskResult = {
  agent_id: string;
  task_id: string;
  status: string;
  output: string;
  error?: string;
  timestamp: string;
};

export default function App() {
  const [token, setToken] = useState<string | null>(null);
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("redforge");
  const [agents, setAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<string>("");
  const [command, setCommand] = useState("ls");
  const [args, setArgs] = useState("/tmp");
  const [timeoutSec, setTimeoutSec] = useState(30);
  const [results, setResults] = useState<TaskResult[]>([]);
  const [status, setStatus] = useState<string>("idle");

  const authHeader = useMemo(() => (token ? { Authorization: `Bearer ${token}` } : {}), [token]);

  function clearSession(reason: string) {
    window.localStorage.removeItem("redforge_token");
    setToken(null);
    setAgents([]);
    setSelectedAgent("");
    setResults([]);
    setStatus(reason);
  }

  useEffect(() => {
    const saved = window.localStorage.getItem("redforge_token");
    if (saved) {
      setStatus("restoring session...");
      setToken(saved);
    }
  }, []);

  useEffect(() => {
    if (!token) return;
    fetchAgents();
  }, [token]);

  async function login() {
    setStatus("logging in...");
    try {
      const r = await fetch(`${API_BASE}/api/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });
      if (r.status === 401) throw new Error("invalid credentials");
      if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
      const data = await r.json();
      window.localStorage.setItem("redforge_token", data.token);
      setToken(data.token);
      setStatus("logged in");
    } catch (err) {
      setStatus(`login failed: ${err}`);
    }
  }

  async function fetchAgents() {
    if (!token) return;
    setStatus("loading agents...");
    try {
      const r = await fetch(`${API_BASE}/api/operator/agents`, {
        headers: { ...authHeader, "Content-Type": "application/json" },
      });
      if (r.status === 401 || r.status === 403) {
        clearSession("session expired (agents): please log in again");
        return;
      }
      if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
      const raw = (await r.json()) as unknown;
      const data = Array.isArray(raw) ? (raw as Agent[]) : [];
      setAgents(data);
      if (data.length > 0) setSelectedAgent(data[0].agent_id);
      setStatus(`loaded ${data.length} agents`);
    } catch (err) {
      setStatus(`failed to load agents: ${err}`);
    }
  }

  async function sendTask() {
    if (!token) return;
    if (!selectedAgent) {
      setStatus("select an agent first");
      return;
    }

    setStatus("submitting task...");

    const body = {
      agent_id: selectedAgent,
      command,
      args: args.trim() === "" ? [] : args.split(" "),
      timeout_seconds: timeoutSec,
    };

    try {
      const r = await fetch(`${API_BASE}/api/operator/task`, {
        method: "POST",
        headers: { ...authHeader, "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (r.status === 401 || r.status === 403) {
        clearSession("session expired (task): please log in again");
        return;
      }
      if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
      setStatus("task submitted");
      await fetchResults();
    } catch (err) {
      setStatus(`task failed: ${err}`);
    }
  }

  async function fetchResults() {
    if (!token || !selectedAgent) return;
    setStatus("loading results...");
    try {
      const r = await fetch(
        `${API_BASE}/api/operator/results?agent_id=${encodeURIComponent(selectedAgent)}`,
        {
          headers: { ...authHeader, "Content-Type": "application/json" },
        }
      );
      if (r.status === 401 || r.status === 403) {
        clearSession("session expired (results): please log in again");
        return;
      }
      if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
      const raw = (await r.json()) as unknown;
      const data = Array.isArray(raw) ? (raw as TaskResult[]) : [];
      setResults(data);
      setStatus(`loaded ${data.length} results`);
    } catch (err) {
      setStatus(`failed to load results: ${err}`);
    }
  }

  return (
    <div className="app">
      <header className="header">
        <div className="title">
          <div className="brand">
            <img className="brand-logo" src={logoUrl} alt="RedForgeC2" />
            <div className="brand-text">
              <h1>RedForgeC2</h1>
              <p className="subtitle">Operator Console</p>
            </div>
          </div>
          <div className="chips">
            <span className={`chip ${token ? "ok" : "warn"}`}>{token ? "AUTH: OK" : "AUTH: NONE"}</span>
            <span className="chip">AGENTS: {agents.length}</span>
            <span className="chip">TARGET: {selectedAgent ? selectedAgent.slice(0, 8) : "--"}</span>
          </div>
        </div>
        <div className="login">
          <div className="login-row">
            <label>Username</label>
            <input value={username} onChange={(e) => setUsername(e.target.value)} />
          </div>
          <div className="login-row">
            <label>Password</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          </div>
          <div className="login-actions">
            <button className="btn btn-primary" onClick={login} disabled={!username || !password}>
            Login
            </button>
            <button className="btn btn-ghost" onClick={() => clearSession("logged out")} disabled={!token}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {!token ? (
        <main className="main main-single">
          <section className="panel">
            <div className="panel-header">
              <h2>Authentication Required</h2>
            </div>
            <div className="panel-body">
              <div className="empty">Log in to view agents, send tasking, and view results.</div>
            </div>
          </section>
        </main>
      ) : (
      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>Agents</h2>
            <button onClick={fetchAgents} disabled={!token}>
              Refresh
            </button>
          </div>

          <div className="panel-body">
            <table className="table">
              <thead>
                <tr>
                  <th>Agent ID</th>
                  <th>Host</th>
                  <th>OS</th>
                  <th>Arch</th>
                  <th>Last Seen</th>
                </tr>
              </thead>
              <tbody>
                {agents.map((agent) => (
                  <tr
                    key={agent.agent_id}
                    className={agent.agent_id === selectedAgent ? "selected" : ""}
                    onClick={() => setSelectedAgent(agent.agent_id)}
                  >
                    <td>{agent.agent_id}</td>
                    <td>{agent.hostname}</td>
                    <td>{agent.os}</td>
                    <td>{agent.arch}</td>
                    <td>{new Date(agent.last_seen).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Tasking</h2>
            <button onClick={sendTask} disabled={!token || !selectedAgent}>
              Send Task
            </button>
          </div>

          <div className="panel-body">
            <div className="form-row">
              <label>Agent</label>
              <select value={selectedAgent} onChange={(e) => setSelectedAgent(e.target.value)}>
                <option value="">-- select agent --</option>
                {agents.map((a) => (
                  <option key={a.agent_id} value={a.agent_id}>
                    {a.hostname} ({a.agent_id.slice(0, 8)})
                  </option>
                ))}
              </select>
            </div>
            <div className="form-row">
              <label>Command</label>
              <input value={command} onChange={(e) => setCommand(e.target.value)} />
            </div>
            <div className="form-row">
              <label>Args</label>
              <input value={args} onChange={(e) => setArgs(e.target.value)} />
            </div>
            <div className="form-row">
              <label>Timeout</label>
              <input
                type="number"
                value={timeoutSec}
                onChange={(e) => setTimeoutSec(Number(e.target.value))}
              />
            </div>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Results</h2>
            <button onClick={fetchResults} disabled={!token || !selectedAgent}>
              Refresh
            </button>
          </div>

          <div className="panel-body results">
            {results.length === 0 ? (
              <div className="empty">No results yet</div>
            ) : (
              results.map((r) => (
                <div key={r.task_id} className="result">
                  <div className="result-meta">
                    <span className="tag">{r.status}</span>
                    <span>{new Date(r.timestamp).toLocaleString()}</span>
                    <span>task: {r.task_id.slice(0, 8)}</span>
                  </div>
                  {r.error ? <pre className="error">{r.error}</pre> : null}
                  <pre className="output">{r.output}</pre>
                </div>
              ))
            )}
          </div>
        </section>
      </main>
      )}

      <footer className="footer">
        <div className="status">Status: {status}</div>
        <div className="hint">Teamserver: {API_BASE}</div>
      </footer>
    </div>
  );
}
