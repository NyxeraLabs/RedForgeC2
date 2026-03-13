import { useMemo } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import { Select } from "../components/Select";
import type { Agent } from "../lib/api";

function agentStatus(lastSeen: string): { label: string; cls: string } {
  const ts = Date.parse(lastSeen);
  if (Number.isNaN(ts)) return { label: "unknown", cls: "text-slate-300" };

  const ageMs = Date.now() - ts;
  if (ageMs < 60_000) return { label: "active", cls: "text-emerald-300" };
  if (ageMs < 5 * 60_000) return { label: "stale", cls: "text-yellow-300" };
  return { label: "offline", cls: "text-red-300" };
}

export function AgentsCard({
  isAuthenticated,
  agents,
  selectedAgentId,
  onSelectedAgentId,
  onRefreshAgents,
}: {
  isAuthenticated: boolean;
  agents: Agent[];
  selectedAgentId: string;
  onSelectedAgentId: (next: string) => void;
  onRefreshAgents: () => void;
}) {
  const selected = useMemo(
    () => agents.find((a) => a.agent_id === selectedAgentId),
    [agents, selectedAgentId]
  );

  return (
    <Card title="Agents">
      {!isAuthenticated ? (
        <p className="text-sm text-slate-400">Login to see agents.</p>
      ) : (
        <div className="space-y-3">
          <Select value={selectedAgentId} onChange={(e) => onSelectedAgentId(e.target.value)}>
            <option value="">Select an agent</option>
            {agents.map((a) => {
              const s = agentStatus(a.last_seen);
              return (
                <option key={a.agent_id} value={a.agent_id}>
                  {a.hostname} ({a.agent_id.slice(0, 8)}) — {s.label}
                </option>
              );
            })}
          </Select>
          <div className="flex gap-2">
            <Button onClick={onRefreshAgents}>Refresh Agents</Button>
          </div>

          {selected ? (
            <div className="rounded border border-slate-800 bg-slate-950 p-3 text-sm">
              <div className="flex items-center justify-between">
                <div className="font-medium">{selected.hostname}</div>
                <div className={agentStatus(selected.last_seen).cls}>
                  {agentStatus(selected.last_seen).label}
                </div>
              </div>
              <div className="mt-2 grid grid-cols-2 gap-2 text-slate-200">
                <div>
                  <div className="text-slate-400">Agent ID</div>
                  <div className="font-mono text-xs">{selected.agent_id}</div>
                </div>
                <div>
                  <div className="text-slate-400">Last Seen</div>
                  <div className="font-mono text-xs">{selected.last_seen}</div>
                </div>
                <div>
                  <div className="text-slate-400">OS</div>
                  <div className="font-mono text-xs">
                    {selected.os} / {selected.arch}
                  </div>
                </div>
                <div>
                  <div className="text-slate-400">Version</div>
                  <div className="font-mono text-xs">{selected.version}</div>
                </div>
              </div>
            </div>
          ) : null}
        </div>
      )}
    </Card>
  );
}

