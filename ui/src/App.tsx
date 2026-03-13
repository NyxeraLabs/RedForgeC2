import { useEffect, useMemo, useState } from "react";
import { Button } from "./components/Button";
import type { Agent, TaskResult } from "./lib/api";
import { listAgents, listResults } from "./lib/api";
import { AgentsCard } from "./pages/AgentsCard";
import { AlertsCard } from "./pages/AlertsCard";
import { AuthCard } from "./pages/AuthCard";
import { ConsoleCard } from "./pages/ConsoleCard";
import { FilesCard } from "./pages/FilesCard";
import { TelemetryCard } from "./pages/TelemetryCard";

export default function App() {
  const [token, setToken] = useState<string>(() => localStorage.getItem("rf_token") || "");
  const [agents, setAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<string>("");
  const [results, setResults] = useState<TaskResult[]>([]);
  const [active, setActive] = useState<"console" | "files" | "telemetry" | "alerts">("console");
  const [error, setError] = useState<string>("");

  const isAuthenticated = useMemo(() => token.length > 0, [token]);

  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    listAgents(token)
      .then(setAgents)
      .catch((e) => setError(e instanceof Error ? e.message : "failed to load agents"));
  }, [isAuthenticated, token]);

  const setTokenAndPersist = (next: string) => {
    setToken(next);
    if (next) localStorage.setItem("rf_token", next);
    else localStorage.removeItem("rf_token");
  };

  const refreshAgents = async () => {
    setError("");
    if (!isAuthenticated) return;
    try {
      setAgents(await listAgents(token));
    } catch (e) {
      setError(e instanceof Error ? e.message : "failed to load agents");
    }
  };

  const refreshResults = async () => {
    setError("");
    if (!isAuthenticated || !selectedAgent) return;
    try {
      setResults(await listResults(token, selectedAgent));
    } catch (e) {
      setError(e instanceof Error ? e.message : "failed to load results");
    }
  };

  useEffect(() => {
    setResults([]);
    if (selectedAgent && isAuthenticated) {
      refreshResults().catch(() => {});
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAgent]);

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="p-4 border-b border-slate-800">
        <h1 className="text-2xl font-semibold">RedForgeC2 Operator Console</h1>
        <p className="text-sm text-slate-300">MVP operator UI (agents, console, files)</p>
      </header>

      <main className="p-6 grid gap-6 lg:grid-cols-2">
        <AuthCard token={token} onToken={setTokenAndPersist} />

        <AgentsCard
          isAuthenticated={isAuthenticated}
          agents={agents}
          selectedAgentId={selectedAgent}
          onSelectedAgentId={setSelectedAgent}
          onRefreshAgents={refreshAgents}
        />

        <section className="rounded-lg border border-slate-800 bg-slate-900 p-4 lg:col-span-2">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex gap-2">
              <Button
                className={active === "console" ? "bg-indigo-600 hover:bg-indigo-500" : ""}
                onClick={() => setActive("console")}
              >
                Console
              </Button>
              <Button
                className={active === "files" ? "bg-indigo-600 hover:bg-indigo-500" : ""}
                onClick={() => setActive("files")}
              >
                Files
              </Button>
              <Button
                className={active === "telemetry" ? "bg-indigo-600 hover:bg-indigo-500" : ""}
                onClick={() => setActive("telemetry")}
              >
                Telemetry
              </Button>
              <Button
                className={active === "alerts" ? "bg-indigo-600 hover:bg-indigo-500" : ""}
                onClick={() => setActive("alerts")}
              >
                Alerts
              </Button>
            </div>

            <div className="flex gap-2">
              <Button onClick={refreshAgents}>Refresh Agents</Button>
              <Button onClick={refreshResults}>Refresh Results</Button>
            </div>
          </div>
          {error ? <p className="mt-3 text-sm text-red-300">{error}</p> : null}
        </section>

        <div className="lg:col-span-2">
          {active === "console" ? (
            <ConsoleCard
              isAuthenticated={isAuthenticated}
              token={token}
              selectedAgentId={selectedAgent}
              results={results}
              onRefreshResults={refreshResults}
            />
          ) : active === "files" ? (
            <FilesCard
              isAuthenticated={isAuthenticated}
              token={token}
              selectedAgentId={selectedAgent}
              results={results}
              onRefreshResults={refreshResults}
            />
          ) : active === "telemetry" ? (
            <TelemetryCard
              isAuthenticated={isAuthenticated}
              token={token}
              selectedAgentId={selectedAgent}
            />
          ) : (
            <AlertsCard isAuthenticated={isAuthenticated} token={token} />
          )}
        </div>
      </main>
    </div>
  );
}
