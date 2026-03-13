import { useState } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import type { TelemetryPayload } from "../lib/api";
import { getLatestTelemetry } from "../lib/api";

export function TelemetryCard({
  isAuthenticated,
  token,
  selectedAgentId,
}: {
  isAuthenticated: boolean;
  token: string;
  selectedAgentId: string;
}) {
  const [telemetry, setTelemetry] = useState<TelemetryPayload | null>(null);
  const [error, setError] = useState<string>("");

  return (
    <Card title="Telemetry (latest)">
      {!isAuthenticated ? (
        <p className="text-sm text-slate-400">Login to view telemetry.</p>
      ) : !selectedAgentId ? (
        <p className="text-sm text-slate-400">Select an agent to view telemetry.</p>
      ) : (
        <div className="space-y-3">
          {error ? <p className="text-sm text-red-300">{error}</p> : null}
          <div className="flex gap-2">
            <Button
              onClick={async () => {
                setError("");
                try {
                  setTelemetry(await getLatestTelemetry(token, selectedAgentId));
                } catch (e) {
                  setTelemetry(null);
                  setError(e instanceof Error ? e.message : "failed to load telemetry");
                }
              }}
            >
              Refresh
            </Button>
          </div>
          {telemetry ? (
            <div className="rounded border border-slate-800 bg-slate-950 p-3 text-sm">
              <div className="grid grid-cols-2 gap-2 text-slate-200">
                <div>
                  <div className="text-slate-400">CPU</div>
                  <div className="font-mono text-xs">{telemetry.cpu.toFixed(2)}%</div>
                </div>
                <div>
                  <div className="text-slate-400">Memory</div>
                  <div className="font-mono text-xs">{telemetry.memory} bytes</div>
                </div>
                <div>
                  <div className="text-slate-400">Uptime</div>
                  <div className="font-mono text-xs">{telemetry.uptime} s</div>
                </div>
                <div>
                  <div className="text-slate-400">Timestamp</div>
                  <div className="font-mono text-xs">{telemetry.timestamp}</div>
                </div>
              </div>
            </div>
          ) : (
            <p className="text-sm text-slate-400">No telemetry yet (agent must heartbeat).</p>
          )}
        </div>
      )}
    </Card>
  );
}

