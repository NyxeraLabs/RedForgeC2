#!/usr/bin/env bash
set -euo pipefail

addr="127.0.0.1:18080"
base="http://${addr}"
agent_id="agent-smoke-1"

cleanup() {
  if [[ -n "${server_pid:-}" ]] && kill -0 "${server_pid}" 2>/dev/null; then
    kill "${server_pid}" 2>/dev/null || true
    wait "${server_pid}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

( cd server/go_backend && go run ./cmd/teamserver -addr "${addr}" >/dev/null 2>&1 ) &
server_pid="$!"

for _ in {1..50}; do
  if curl -sf "${base}/healthz" >/dev/null; then
    break
  fi
  sleep 0.1
done
curl -sf "${base}/healthz" >/dev/null

# Pre-register the agent and enqueue a mock task.
curl -sf -X POST "${base}/api/v1/agents/register" \
  -H 'Content-Type: application/json' \
  -d "{\"agent_id\":\"${agent_id}\",\"os\":\"linux\",\"arch\":\"amd64\",\"hostname\":\"lab\"}" >/dev/null

curl -sf -X POST "${base}/api/v1/tasks/enqueue" \
  -H 'Content-Type: application/json' \
  -d "{\"agent_id\":\"${agent_id}\",\"task_type\":\"echo\",\"params\":{\"msg\":\"hello\"}}" >/dev/null

( cd agents/rust_agent && \
  REDFORGE_SERVER_URL="${base}" \
  REDFORGE_AGENT_ID="${agent_id}" \
  REDFORGE_ONCE=1 \
  cargo run --quiet >/dev/null 2>&1 )

echo "OK: agent↔server simulation smoke test passed"
