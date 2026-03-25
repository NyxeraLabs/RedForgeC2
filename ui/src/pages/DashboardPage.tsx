import React, { useEffect, useMemo, useRef, useState } from "react";
import { Agent, apiBase, createTask, deleteAgent, getMe, listAgents, listResults, listTasks, TaskResult, TaskSummary, wsBase } from "../lib/api";
import { useToasts } from "../components/ToastProvider";
import { getSelectedAgent, getToken, setSelectedAgent as persistSelectedAgent } from "../lib/storage";

function ageSeconds(iso: string): number {
  const t = new Date(iso).getTime();
  return Math.max(0, Math.floor((Date.now() - t) / 1000));
}

function heartbeatClass(lastSeenIso: string) {
  const age = ageSeconds(lastSeenIso);
  if (age <= 90) return "hb good";
  if (age <= 240) return "hb warn";
  return "hb bad";
}

export function DashboardPage() {
  const toasts = useToasts();
  const [role, setRole] = useState<"admin" | "operator" | "observer" | null>(null);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<string>(() => getSelectedAgent() || "");
  const [command, setCommand] = useState("ls");
  const [args, setArgs] = useState("");
  const [timeoutSec, setTimeoutSec] = useState(30);
  const [transport, setTransport] = useState("https");
  const [results, setResults] = useState<TaskResult[]>([]);
  const [tasks, setTasks] = useState<TaskSummary[]>([]);
  const [status, setStatus] = useState<string>("idle");

  const knownAgents = useRef<Set<string>>(new Set());
  const selectedAgentRef = useRef<string>(selectedAgent);
  const api = useMemo(() => apiBase(), []);

  useEffect(() => {
    selectedAgentRef.current = selectedAgent;
    persistSelectedAgent(selectedAgent || null);
  }, [selectedAgent]);

  useEffect(() => {
    const current = selectedAgentRef.current;
    if (agents.length === 0) {
      if (current) setSelectedAgent("");
      return;
    }
    if (!current || !agents.some((a) => a.agent_id === current)) {
      setSelectedAgent(agents[0].agent_id);
    }
  }, [agents]);

  async function removeAgent(agentId: string) {
    if (role === "observer") return;
    const ok = window.confirm(`Delete agent ${agentId.slice(0, 8)}? This will also delete tasks/results for the agent.`);
    if (!ok) return;
    setStatus("deleting agent...");
    try {
      await deleteAgent(agentId);
      toasts.push({ kind: "good", title: "Agent deleted", message: agentId.slice(0, 8) });
      await refreshAgents(false);
      await refreshTasks();
      await refreshResults();
      setStatus("agent deleted");
    } catch (e) {
      setStatus(`delete failed: ${e}`);
      toasts.push({ kind: "bad", title: "Delete failed", message: String(e) });
    }
  }

  async function refreshAgents(showToast = false) {
    setStatus("loading agents...");
    try {
      const list = await listAgents();
      setAgents(list);
      const current = selectedAgentRef.current;
      if (list.length === 0) {
        if (current) setSelectedAgent("");
      } else if (!current || !list.some((a) => a.agent_id === current)) {
        setSelectedAgent(list[0].agent_id);
      }
      setStatus(`loaded ${list.length} agents`);

      const newOnes = list.filter((a) => !knownAgents.current.has(a.agent_id));
      newOnes.forEach((a) => knownAgents.current.add(a.agent_id));
      if (showToast) {
        newOnes.forEach((a) =>
          toasts.push({ kind: "info", title: "New agent", message: `${a.hostname} (${a.agent_id.slice(0, 8)})` })
        );
      }
    } catch (e) {
      setStatus(`agents failed: ${e}`);
      toasts.push({ kind: "bad", title: "Agents refresh failed", message: String(e) });
    }
  }

  async function refreshTasks() {
    try {
      const list = await listTasks(selectedAgentRef.current || undefined);
      setTasks(list);
    } catch (e) {
      toasts.push({ kind: "warn", title: "Tasks refresh failed", message: String(e) });
    }
  }

  async function refreshResults() {
    const agentId = selectedAgentRef.current;
    if (!agentId) return;
    try {
      const list = await listResults(agentId);
      setResults(list);
    } catch (e) {
      toasts.push({ kind: "warn", title: "Results refresh failed", message: String(e) });
    }
  }

  useEffect(() => {
    getMe()
      .then((m) => setRole(m.role))
      .catch(() => setRole(null));
    refreshAgents(true);

    const id = window.setInterval(() => {
      refreshAgents(true);
      refreshTasks();
      refreshResults();
    }, 5000);

    // Real-time websocket updates (fallbacks to polling)
    const token = getToken();
    const wsUrl = token ? `${wsBase()}/api/ws?token=${encodeURIComponent(token)}` : `${wsBase()}/api/ws`;
    const ws = new WebSocket(wsUrl);
    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as { type: string; agents?: Agent[]; tasks?: TaskSummary[] };
        if (msg.type === "state") {
          if (msg.agents) setAgents(msg.agents);
          if (msg.tasks) {
            const current = selectedAgentRef.current;
            setTasks(current ? msg.tasks.filter((t) => t.agent_id === current) : msg.tasks);
          }
        }
      } catch {
        // ignore bad payloads
      }
    };

    return () => {
      window.clearInterval(id);
      ws.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    refreshTasks();
    refreshResults();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAgent]);

  async function sendTask() {
    if (!selectedAgent) {
      toasts.push({ kind: "warn", title: "Select an agent first" });
      return;
    }
    setStatus("submitting task...");
    try {
      const rawCommand = command.trim();
      if (!rawCommand) {
        toasts.push({ kind: "warn", title: "Command is required" });
        setStatus("idle");
        return;
      }

      // Accept "ls -la" typed into the Command box by splitting on whitespace and
      // combining with the Args field. This prevents trying to exec "ls -la" as
      // a single binary (which fails with ENOENT).
      const cmdParts = rawCommand.split(/\s+/).filter(Boolean);
      const cmd = cmdParts[0];
      const inlineArgs = cmdParts.slice(1);
      const extraArgs = args.trim() === "" ? [] : args.trim().split(/\s+/).filter(Boolean);
      const argv = [...inlineArgs, ...extraArgs];

      await createTask(selectedAgent, cmd, argv, timeoutSec, transport);
      setStatus("task submitted");
      toasts.push({ kind: "good", title: "Task queued", message: `${cmd} ${argv.join(" ")}`.trim() });
      await refreshTasks();
      await refreshResults();
    } catch (e) {
      setStatus(`task failed: ${e}`);
      toasts.push({ kind: "bad", title: "Task failed", message: String(e) });
    }
  }

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Dashboard</div>
          <div className="page-sub">Teamserver: {api}</div>
        </div>
        <div className="chips">
          <span className="chip">AGENTS: {agents.length}</span>
          <span className="chip">TASKS: {tasks.length}</span>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>Agents</h2>
            <button onClick={() => refreshAgents(true)}>Refresh</button>
          </div>
          <div className="panel-body">
            <div className="table-wrap">
              <table className="table">
                <thead>
                  <tr>
                    <th>Agent</th>
                    <th>OS</th>
                    <th>Arch</th>
                    <th>HB</th>
                    <th>Fails</th>
                    <th>Last Err</th>
                    <th>Last Seen</th>
                  </tr>
                </thead>
                <tbody>
                  {agents.map((a) => (
                    <tr
                      key={a.agent_id}
                      className={a.agent_id === selectedAgent ? "selected" : ""}
                      onClick={() => setSelectedAgent(a.agent_id)}
                    >
                      <td>
                        <div className="agent-line">
                          <span className="agent-host" title={a.hostname}>
                            {a.hostname}
                          </span>
                          <span className="agent-id">{a.agent_id.slice(0, 8)}</span>
                          {role === "observer" ? null : (
                            <span className="agent-actions">
                              <button
                                className="btn-mini btn-danger"
                                title="Delete agent"
                                onClick={(e) => {
                                  e.preventDefault();
                                  e.stopPropagation();
                                  void removeAgent(a.agent_id);
                                }}
                              >
                                Delete
                              </button>
                            </span>
                          )}
                        </div>
                      </td>
                      <td>{a.os}</td>
                      <td>{a.arch}</td>
                      <td>
                        <span className={heartbeatClass(a.last_seen)}>{ageSeconds(a.last_seen)}s</span>
                      </td>
                      <td>{a.heartbeat_failures ?? 0}</td>
                      <td title={a.heartbeat_last_error || ""}>
                        {a.heartbeat_last_error ? a.heartbeat_last_error.slice(0, 24) : "-"}
                      </td>
                      <td>{new Date(a.last_seen).toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Tasking</h2>
            <button onClick={sendTask} disabled={!selectedAgent || role === "observer"}>
              Send
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
              <label>Timeout (sec)</label>
              <input type="number" value={timeoutSec} onChange={(e) => setTimeoutSec(Number(e.target.value))} />
            </div>

            <div className="form-row">
              <label>Transport</label>
              <select value={transport} onChange={(e) => setTransport(e.target.value)}>
                <option value="https">HTTPS</option>
                <option value="dns">DNS (Fallback)</option>
                <option value="icmp">ICMP (Signaling)</option>
              </select>
            </div>

            <div className="split">
              <div className="split-head">
                <div className="split-title">Queue</div>
                <button onClick={refreshTasks}>Refresh</button>
              </div>
              <div className="split-body">
                {tasks.length === 0 ? (
                  <div className="empty">No tasks</div>
                ) : (
                  tasks.slice(0, 25).map((t) => (
                    <div key={t.task_id} className={`task ${t.status}`}>
                      <div className="task-top">
                        <span className="tag">{t.status}</span>
                        <span className="task-id">{t.task_id.slice(0, 8)}</span>
                        <span className="task-when">{new Date(t.updated_at).toLocaleString()}</span>
                      </div>
                      <div className="task-cmd">
                        <span className="mono">{t.command}</span>{" "}
                        <span className="mono muted">{(t.args || []).join(" ")}</span>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        </section>

        <section className="panel panel-span">
          <div className="panel-header">
            <h2>Results</h2>
            <button onClick={refreshResults} disabled={!selectedAgent}>
              Refresh
            </button>
          </div>
          <div className="panel-body results">
            {results.length === 0 ? (
              <div className="empty">No results yet</div>
            ) : (
              results.slice(0, 30).map((r) => (
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

      <footer className="footer">
        <div className="status">Status: {status}</div>
        <div className="hint">Active agent: {selectedAgent ? selectedAgent.slice(0, 8) : "--"}</div>
      </footer>
    </div>
  );
}
