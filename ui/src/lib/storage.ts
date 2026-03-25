export const STORAGE_KEYS = {
  token: "redforge_token",
  apiBase: "redforge_api_base",
  selectedAgent: "redforge_selected_agent",
} as const;

export function getToken(): string | null {
  return window.localStorage.getItem(STORAGE_KEYS.token);
}

export function setToken(token: string | null) {
  if (!token) window.localStorage.removeItem(STORAGE_KEYS.token);
  else window.localStorage.setItem(STORAGE_KEYS.token, token);
}

export function getApiBase(): string | null {
  return window.localStorage.getItem(STORAGE_KEYS.apiBase);
}

export function setApiBase(apiBase: string | null) {
  if (!apiBase) window.localStorage.removeItem(STORAGE_KEYS.apiBase);
  else window.localStorage.setItem(STORAGE_KEYS.apiBase, apiBase);
}

export function getSelectedAgent(): string | null {
  return window.localStorage.getItem(STORAGE_KEYS.selectedAgent);
}

export function setSelectedAgent(agentId: string | null) {
  if (!agentId) window.localStorage.removeItem(STORAGE_KEYS.selectedAgent);
  else window.localStorage.setItem(STORAGE_KEYS.selectedAgent, agentId);
}
