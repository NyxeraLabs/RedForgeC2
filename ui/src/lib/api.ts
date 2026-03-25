import { getApiBase, getToken, setToken } from "./storage";

function defaultApiBaseFromLocation(): string {
  const hostname = window.location.hostname;
  return `https://${hostname}:9080`;
}

export function apiBase(): string {
  // Allow runtime overrides (login page Teamserver URL) to win over build-time defaults.
  return getApiBase() || import.meta.env.VITE_TEAMSERVER_URL || defaultApiBaseFromLocation();
}

export class AuthError extends Error {
  constructor(message = "unauthorized") {
    super(message);
    this.name = "AuthError";
  }
}

async function request(path: string, init: RequestInit = {}) {
  const base = apiBase();
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type") && init.body) headers.set("Content-Type", "application/json");
  const token = getToken();
  if (token && !headers.has("Authorization")) headers.set("Authorization", `Bearer ${token}`);

  try {
    const res = await fetch(`${base}${path}`, { ...init, headers });
    if (res.status === 401 || res.status === 403) {
      setToken(null);
      throw new AuthError();
    }
    return res;
  } catch (err) {
    if (err instanceof TypeError) {
      throw new Error(
        `Network error: check that the teamserver is reachable at ${base} and CORS is configured for your UI origin. If ${base} uses a self-signed TLS cert, open ${base}/healthz in your browser and accept/trust the certificate.`,
      );
    }
    throw err;
  }
}

export async function login(username: string, password: string): Promise<string> {
  const res = await request("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  if (!res.ok) {
    if (res.status === 401) throw new Error("invalid credentials");
    throw new Error(`${res.status} ${res.statusText}`);
  }
  const data = (await res.json()) as { token: string };
  return data.token;
}

export type Me = {
  username: string;
  role: "admin" | "operator" | "observer";
  display_name?: string;
};

export async function getMe(): Promise<Me> {
  const res = await request("/api/me");
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  return (await res.json()) as Me;
}

export async function updateProfile(displayName: string) {
  const res = await request("/api/me/profile", { method: "POST", body: JSON.stringify({ display_name: displayName }) });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
}

export async function changePassword(oldPassword: string, newPassword: string) {
  const res = await request("/api/me/password", {
    method: "POST",
    body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
  });
  if (!res.ok) {
    const data = (await res.json().catch(() => null)) as any;
    throw new Error(data?.error || `${res.status} ${res.statusText}`);
  }
}

export type Agent = {
  agent_id: string;
  os: string;
  arch: string;
  hostname: string;
  version: string;
  last_seen: string;
  registered: string;
  heartbeat_failures?: number;
  heartbeat_last_error?: string;
  heartbeat_last_attempt?: string;
  heartbeat_last_backoff_ms?: number;
};

export async function listAgents(): Promise<Agent[]> {
  const res = await request("/api/operator/agents");
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  const raw = (await res.json()) as unknown;
  return Array.isArray(raw) ? (raw as Agent[]) : [];
}

export async function deleteAgent(agentId: string) {
  const res = await request(`/api/operator/agents/${encodeURIComponent(agentId)}`, { method: "DELETE" });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
}

export type TaskSummary = {
  task_id: string;
  agent_id: string;
  command: string;
  args: string[];
  timeout_seconds: number;
  status: string;
  created_at: string;
  updated_at: string;
};

export async function listTasks(agentId?: string): Promise<TaskSummary[]> {
  const qs = agentId ? `?agent_id=${encodeURIComponent(agentId)}` : "";
  const res = await request(`/api/operator/tasks${qs}`);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  const raw = (await res.json()) as unknown;
  return Array.isArray(raw) ? (raw as TaskSummary[]) : [];
}

export async function createTask(agentId: string, command: string, args: string[], timeoutSeconds: number, transport?: string) {
  const res = await request("/api/operator/task", {
    method: "POST",
    body: JSON.stringify({ agent_id: agentId, command, args, timeout_seconds: timeoutSeconds, transport }),
  });
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
}

export type TaskResult = {
  agent_id: string;
  task_id: string;
  status: string;
  output: string;
  error?: string;
  timestamp: string;
};

export function wsBase(): string {
  const base = apiBase();
  return base.replace(/^https:\/\//, "wss://").replace(/^http:\/\//, "ws://");
}

export async function listResults(agentId: string): Promise<TaskResult[]> {
  const res = await request(`/api/operator/results?agent_id=${encodeURIComponent(agentId)}`);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  const raw = (await res.json()) as unknown;
  return Array.isArray(raw) ? (raw as TaskResult[]) : [];
}

export type AdminUser = {
  username: string;
  role: "admin" | "operator" | "observer";
  display_name?: string;
  created_at: string;
  updated_at: string;
};

export async function adminListUsers(): Promise<AdminUser[]> {
  const res = await request("/api/admin/users");
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  const raw = (await res.json()) as unknown;
  return Array.isArray(raw) ? (raw as AdminUser[]) : [];
}

export async function adminCreateUser(body: { username: string; password: string; role: string; display_name?: string }) {
  const res = await request("/api/admin/users", { method: "POST", body: JSON.stringify(body) });
  if (!res.ok) {
    const data = (await res.json().catch(() => null)) as any;
    throw new Error(data?.error || `${res.status} ${res.statusText}`);
  }
}
// File Transfer API types and functions

export type StoredFile = {
  file_id: string;
  session_id: string;
  filename: string;
  total_chunks: number;
  chunks_received: number;
  uploaded_at: string;
  agent_id: string;
};

export async function listFiles(agentId?: string): Promise<StoredFile[]> {
  const qs = agentId ? `?agent_id=${encodeURIComponent(agentId)}` : "";
  const res = await request(`/api/operator/files${qs}`);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  const data = (await res.json()) as { files: StoredFile[] } | null;
  return data?.files || [];
}

export async function getFileInfo(fileId: string): Promise<StoredFile> {
  const res = await request(`/api/operator/files/${encodeURIComponent(fileId)}`);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  return (await res.json()) as StoredFile;
}

export async function downloadFile(fileId: string, fileName: string): Promise<void> {
  const token = getToken();
  const base = apiBase();
  
  // Create a temporary element to trigger download
  const link = document.createElement("a");
  link.href = `${base}/api/operator/files/${encodeURIComponent(fileId)}/download?token=${encodeURIComponent(token || "")}`;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}