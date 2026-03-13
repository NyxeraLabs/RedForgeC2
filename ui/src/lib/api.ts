const API_URL = import.meta.env.VITE_TEAMSERVER_URL || "http://localhost:9080";

export type Agent = {
  agent_id: string;
  hostname: string;
  os: string;
  arch: string;
  version: string;
  last_seen: string;
};

export type TaskResult = {
  task_id: string;
  status: string;
  output: string;
  error?: string;
  timestamp: string;
};

export type TelemetryPayload = {
  agent_id: string;
  cpu: number;
  memory: number;
  uptime: number;
  timestamp: string;
};

export type Alert = {
  agent_id: string;
  task_id: string;
  severity: "RED" | "YEL" | "GRN";
  status: string;
  command: string;
  error?: string;
  timestamp: string;
};

export async function login(username: string, password: string): Promise<string> {
  const res = await fetch(`${API_URL}/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) throw new Error(`login failed: ${res.status}`);
  const body = (await res.json()) as { token: string };
  return body.token;
}

function authHeaders(token: string): HeadersInit {
  return { Authorization: `Bearer ${token}` };
}

export async function listAgents(token: string): Promise<Agent[]> {
  const res = await fetch(`${API_URL}/api/operator/agents`, { headers: authHeaders(token) });
  if (!res.ok) throw new Error(`list agents failed: ${res.status}`);
  return (await res.json()) as Agent[];
}

export async function createTask(
  token: string,
  agentId: string,
  command: string,
  args: string[],
  timeoutSeconds = 60
): Promise<{ task_id: string }> {
  const res = await fetch(`${API_URL}/api/operator/task`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders(token) },
    body: JSON.stringify({
      agent_id: agentId,
      command,
      args,
      timeout_seconds: timeoutSeconds,
    }),
  });
  if (!res.ok) throw new Error(`create task failed: ${res.status}`);
  return (await res.json()) as { task_id: string };
}

export async function listResults(token: string, agentId: string): Promise<TaskResult[]> {
  const res = await fetch(
    `${API_URL}/api/operator/results?agent_id=${encodeURIComponent(agentId)}`,
    { headers: authHeaders(token) }
  );
  if (!res.ok) throw new Error(`list results failed: ${res.status}`);
  return (await res.json()) as TaskResult[];
}

export async function getLatestTelemetry(token: string, agentId: string): Promise<TelemetryPayload> {
  const res = await fetch(
    `${API_URL}/api/operator/telemetry/latest?agent_id=${encodeURIComponent(agentId)}`,
    { headers: authHeaders(token) }
  );
  if (!res.ok) throw new Error(`telemetry failed: ${res.status}`);
  return (await res.json()) as TelemetryPayload;
}

export async function listAlerts(token: string, limit = 50): Promise<Alert[]> {
  const res = await fetch(`${API_URL}/api/operator/alerts?limit=${limit}`, {
    headers: authHeaders(token),
  });
  if (!res.ok) throw new Error(`alerts failed: ${res.status}`);
  return (await res.json()) as Alert[];
}

