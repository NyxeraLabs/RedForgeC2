import { useMemo, useState } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import { TextInput } from "../components/TextInput";
import type { TaskResult } from "../lib/api";
import { createTask } from "../lib/api";

export function ConsoleCard({
  isAuthenticated,
  token,
  selectedAgentId,
  results,
  onRefreshResults,
}: {
  isAuthenticated: boolean;
  token: string;
  selectedAgentId: string;
  results: TaskResult[];
  onRefreshResults: () => void;
}) {
  const [command, setCommand] = useState("ls");
  const [args, setArgs] = useState("/tmp");
  const [lastTaskId, setLastTaskId] = useState<string>("");
  const [error, setError] = useState<string>("");

  const lastResult = useMemo(() => {
    if (!lastTaskId) return undefined;
    return results.find((r) => r.task_id === lastTaskId);
  }, [results, lastTaskId]);

  return (
    <Card title="Console">
      {!isAuthenticated ? (
        <p className="text-sm text-slate-400">Login to run commands.</p>
      ) : !selectedAgentId ? (
        <p className="text-sm text-slate-400">Select an agent to run commands.</p>
      ) : (
        <div className="space-y-3">
          <div className="grid gap-2 md:grid-cols-2">
            <div>
              <label className="block text-sm">Command</label>
              <TextInput value={command} onChange={(e) => setCommand(e.target.value)} />
            </div>
            <div>
              <label className="block text-sm">Args (space-separated)</label>
              <TextInput value={args} onChange={(e) => setArgs(e.target.value)} />
            </div>
          </div>

          {error ? <p className="text-sm text-red-300">{error}</p> : null}

          <div className="flex flex-wrap gap-2">
            <Button
              variant="primary"
              onClick={async () => {
                setError("");
                try {
                  const resp = await createTask(
                    token,
                    selectedAgentId,
                    command,
                    args.split(" ").filter(Boolean),
                    60
                  );
                  setLastTaskId(resp.task_id);
                } catch (e) {
                  setError(e instanceof Error ? e.message : "failed to send task");
                }
              }}
            >
              Send
            </Button>
            <Button onClick={onRefreshResults}>Refresh Results</Button>
          </div>

          <div className="rounded border border-slate-800 bg-slate-950 p-3">
            <div className="text-sm text-slate-300">Last task</div>
            <div className="font-mono text-xs text-slate-200">
              {lastTaskId ? lastTaskId : "(none)"}
            </div>
            {lastResult ? (
              <pre className="mt-2 max-h-56 overflow-y-auto rounded bg-slate-900 p-3 text-xs">
                {lastResult.error ? `ERROR: ${lastResult.error}\n\n` : ""}
                {lastResult.output}
              </pre>
            ) : (
              <p className="mt-2 text-sm text-slate-400">
                {lastTaskId ? "Waiting for result (refresh)..." : "Send a task to see output."}
              </p>
            )}
          </div>

          <details className="rounded border border-slate-800 bg-slate-950 p-3">
            <summary className="cursor-pointer text-sm text-slate-200">
              Recent results ({results.length})
            </summary>
            <pre className="mt-2 max-h-56 overflow-y-auto rounded bg-slate-900 p-3 text-xs">
              {JSON.stringify(results, null, 2)}
            </pre>
          </details>
        </div>
      )}
    </Card>
  );
}

