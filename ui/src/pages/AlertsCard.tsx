import { useState } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import type { Alert } from "../lib/api";
import { listAlerts } from "../lib/api";

export function AlertsCard({ isAuthenticated, token }: { isAuthenticated: boolean; token: string }) {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [error, setError] = useState<string>("");

  return (
    <Card title="Alerts (derived)">
      {!isAuthenticated ? (
        <p className="text-sm text-slate-400">Login to view alerts.</p>
      ) : (
        <div className="space-y-3">
          {error ? <p className="text-sm text-red-300">{error}</p> : null}
          <div className="flex gap-2">
            <Button
              onClick={async () => {
                setError("");
                try {
                  setAlerts(await listAlerts(token, 50));
                } catch (e) {
                  setAlerts([]);
                  setError(e instanceof Error ? e.message : "failed to load alerts");
                }
              }}
            >
              Refresh
            </Button>
          </div>

          <div className="overflow-x-auto rounded border border-slate-800">
            <table className="min-w-full text-sm">
              <thead className="bg-slate-900 text-slate-200">
                <tr>
                  <th className="p-2 text-left">Severity</th>
                  <th className="p-2 text-left">Agent</th>
                  <th className="p-2 text-left">Command</th>
                  <th className="p-2 text-left">Status</th>
                  <th className="p-2 text-left">Time</th>
                </tr>
              </thead>
              <tbody className="bg-slate-950 text-slate-200">
                {alerts.length === 0 ? (
                  <tr>
                    <td className="p-2 text-slate-400" colSpan={5}>
                      No alerts (alerts are currently derived from task results with status error/timeout).
                    </td>
                  </tr>
                ) : (
                  alerts.map((a) => (
                    <tr key={`${a.agent_id}:${a.task_id}`} className="border-t border-slate-800">
                      <td className="p-2 font-mono">{a.severity}</td>
                      <td className="p-2 font-mono">{a.agent_id.slice(0, 8)}</td>
                      <td className="p-2 font-mono">{a.command}</td>
                      <td className="p-2 font-mono">{a.status}</td>
                      <td className="p-2 font-mono text-xs">{a.timestamp}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </Card>
  );
}

