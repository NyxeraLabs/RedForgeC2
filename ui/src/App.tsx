import { useEffect, useMemo, useState } from "react";

const API_URL = import.meta.env.VITE_TEAMSERVER_URL || "http://localhost:9080";

type Agent = {
  agent_id: string;
  hostname: string;
  os: string;
  arch: string;
  version: string;
  last_seen: string;
};

export default function App() {
  const [token, setToken] = useState<string>(() => localStorage.getItem("rf_token") || "");
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("redforge");
  const [agents, setAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<string>("");
  const [taskOutput, setTaskOutput] = useState<string>("");
  const [taskCommand, setTaskCommand] = useState("ls");
  const [taskArgs, setTaskArgs] = useState("/tmp");

  const isAuthenticated = useMemo(() => token.length > 0, [token]);

  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    fetch(`${API_URL}/api/operator/agents`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((r) => r.json())
      .then((data) => setAgents(data))
      .catch(console.error);
  }, [isAuthenticated, token]);

  const login = async () => {
    const res = await fetch(`${API_URL}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (!res.ok) {
      alert("login failed");
      return;
    }

    const body = await res.json();
    setToken(body.token);
    localStorage.setItem("rf_token", body.token);
  };

  const sendTask = async () => {
    if (!selectedAgent) {
      alert("Please select an agent");
      return;
    }

    const res = await fetch(`${API_URL}/api/operator/task`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        agent_id: selectedAgent,
        command: taskCommand,
        args: taskArgs.split(" ").filter(Boolean),
        timeout_seconds: 60,
      }),
    });

    if (!res.ok) {
      alert("failed to send task");
      return;
    }

    const body = await res.json();
    setTaskOutput(`task created: ${body.task_id}`);
  };

  const refreshResults = async () => {
    if (!selectedAgent) return;

    const res = await fetch(
      `${API_URL}/api/operator/results?agent_id=${encodeURIComponent(selectedAgent)}`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
    if (!res.ok) {
      setTaskOutput("failed to fetch results");
      return;
    }
    const list = await res.json();
    setTaskOutput(JSON.stringify(list, null, 2));
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="p-4 border-b border-slate-800">
        <h1 className="text-2xl font-semibold">RedForgeC2 Operator Console</h1>
        <p className="text-sm text-slate-300">Basic operator UI (alpha)</p>
      </header>

      <main className="p-6 grid gap-6 lg:grid-cols-2">
        <section className="rounded-lg border border-slate-800 bg-slate-900 p-6">
          <h2 className="text-xl font-semibold mb-2">Authentication</h2>
          {!isAuthenticated ? (
            <div className="space-y-2">
              <div>
                <label className="block text-sm">Username</label>
                <input
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full rounded border border-slate-700 bg-slate-950 p-2 text-white"
                />
              </div>
              <div>
                <label className="block text-sm">Password</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full rounded border border-slate-700 bg-slate-950 p-2 text-white"
                />
              </div>
              <button
                onClick={login}
                className="mt-3 rounded bg-indigo-600 px-4 py-2 text-sm font-medium hover:bg-indigo-500"
              >
                Login
              </button>
            </div>
          ) : (
            <div className="space-y-2">
              <p className="text-sm">Authenticated. Token stored locally.</p>
              <button
                onClick={() => {
                  setToken("");
                  localStorage.removeItem("rf_token");
                }}
                className="mt-2 rounded bg-red-600 px-4 py-2 text-sm font-medium hover:bg-red-500"
              >
                Logout
              </button>
            </div>
          )}
        </section>

        <section className="rounded-lg border border-slate-800 bg-slate-900 p-6">
          <h2 className="text-xl font-semibold mb-2">Agents</h2>
          {isAuthenticated ? (
            <div className="space-y-3">
              <select
                value={selectedAgent}
                onChange={(e) => setSelectedAgent(e.target.value)}
                className="w-full rounded border border-slate-700 bg-slate-950 p-2"
              >
                <option value="">Select an agent</option>
                {agents.map((a) => (
                  <option key={a.agent_id} value={a.agent_id}>
                    {a.hostname} ({a.agent_id.slice(0, 8)})
                  </option>
                ))}
              </select>
              <button
                onClick={refreshResults}
                className="rounded bg-emerald-600 px-4 py-2 text-sm font-medium hover:bg-emerald-500"
              >
                Refresh Results
              </button>
            </div>
          ) : (
            <p className="text-sm text-slate-400">Login to see agents.</p>
          )}
        </section>

        <section className="rounded-lg border border-slate-800 bg-slate-900 p-6 lg:col-span-2">
          <h2 className="text-xl font-semibold mb-2">Tasking</h2>
          <div className="grid gap-2 md:grid-cols-2">
            <div>
              <label className="block text-sm">Command</label>
              <input
                value={taskCommand}
                onChange={(e) => setTaskCommand(e.target.value)}
                className="w-full rounded border border-slate-700 bg-slate-950 p-2 text-white"
              />
            </div>
            <div>
              <label className="block text-sm">Args (space-separated)</label>
              <input
                value={taskArgs}
                onChange={(e) => setTaskArgs(e.target.value)}
                className="w-full rounded border border-slate-700 bg-slate-950 p-2 text-white"
              />
            </div>
          </div>
          <button
            onClick={sendTask}
            className="mt-4 rounded bg-indigo-600 px-4 py-2 text-sm font-medium hover:bg-indigo-500"
          >
            Send Task
          </button>

          <pre className="mt-4 max-h-56 overflow-y-auto rounded bg-slate-900 p-4 text-xs">
            {taskOutput}
          </pre>
        </section>
      </main>
    </div>
  );
}
