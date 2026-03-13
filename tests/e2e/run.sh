#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT_DIR"

compose_file="docker/docker-compose.yml"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

require_cmd docker
require_cmd curl
require_cmd jq
require_cmd go
require_cmd cargo
require_cmd python3

rand_hex() {
  local n="$1"
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex "$n" 2>/dev/null
    return
  fi
  python3 - <<PY
import secrets
print(secrets.token_hex(${n}))
PY
}

JWT_SECRET="${REDFORGE_JWT_SECRET:-$(rand_hex 32)}"
if [[ -z "${JWT_SECRET:-}" ]]; then
  echo "failed to generate JWT secret; set REDFORGE_JWT_SECRET" >&2
  exit 1
fi

ADMIN_USER="${REDFORGE_ADMIN_USER:-admin}"
ADMIN_PASS="${REDFORGE_ADMIN_PASS:-$(rand_hex 12)}"
if [[ -z "${ADMIN_PASS:-}" ]]; then
  echo "failed to generate admin password; set REDFORGE_ADMIN_PASS" >&2
  exit 1
fi

PORT="${REDFORGE_PORT:-9080}"

DB_USER="${REDFORGE_DB_USER:-redforge}"
DB_PASS="${REDFORGE_DB_PASS:-$(rand_hex 12)}"
DB_NAME="${REDFORGE_DB_NAME:-redforge}"

export REDFORGE_JWT_SECRET="$JWT_SECRET"
export REDFORGE_ADMIN_USER="$ADMIN_USER"
export REDFORGE_ADMIN_PASS="$ADMIN_PASS"
export REDFORGE_PORT="$PORT"
export REDFORGE_DB_USER="$DB_USER"
export REDFORGE_DB_PASS="$DB_PASS"
export REDFORGE_DB_NAME="$DB_NAME"

TEAMSERVER_URL="http://localhost:${PORT}"
DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@localhost:5433/${DB_NAME}?sslmode=disable"
export DATABASE_URL

cleanup() {
  set +e
  if [[ -n "${AGENT_PID:-}" ]]; then kill "${AGENT_PID}" >/dev/null 2>&1 || true; fi
  if [[ -n "${TEAMSERVER_PID:-}" ]]; then kill "${TEAMSERVER_PID}" >/dev/null 2>&1 || true; fi
  docker compose -f "$compose_file" down >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "[e2e] starting postgres via docker compose"
docker compose -f "$compose_file" up -d postgres >/dev/null

echo "[e2e] waiting for postgres to accept connections"
for _ in $(seq 1 60); do
  if docker compose -f "$compose_file" exec -T postgres pg_isready -U "$DB_USER" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

echo "[e2e] starting teamserver locally on :${PORT}"
(cd teamserver && go run ./cmd/teamserver) >/tmp/redforge.teamserver.log 2>&1 &
TEAMSERVER_PID=$!

echo "[e2e] waiting for /healthz"
for _ in $(seq 1 60); do
  if curl -fsS "${TEAMSERVER_URL}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

echo "[e2e] login as operator"
TOKEN="$(
  curl -fsS "${TEAMSERVER_URL}/api/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}" \
    | jq -r .token
)"
if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "[e2e] failed to obtain operator token" >&2
  exit 1
fi

echo "[e2e] starting agent locally"
export REDFORGE_SERVER_URL="${TEAMSERVER_URL}"
(cd agent && cargo run) >/tmp/redforge.agent.log 2>&1 &
AGENT_PID=$!

echo "[e2e] waiting for agent registration (agents list non-empty)"
AGENT_ID=""
for _ in $(seq 1 90); do
  agents_json="$(curl -fsS "${TEAMSERVER_URL}/api/operator/agents" -H "Authorization: Bearer ${TOKEN}")"
  count="$(echo "$agents_json" | jq 'length')"
  if [[ "$count" -ge 1 ]]; then
    AGENT_ID="$(echo "$agents_json" | jq -r '.[0].agent_id')"
    break
  fi
  sleep 1
done
if [[ -z "$AGENT_ID" || "$AGENT_ID" == "null" ]]; then
  echo "[e2e] agent did not register in time" >&2
  echo "[e2e] teamserver log: /tmp/redforge.teamserver.log" >&2
  echo "[e2e] agent log: /tmp/redforge.agent.log" >&2
  exit 1
fi

echo "[e2e] enqueue a safe task (pwd)"
TASK_ID="$(
  curl -fsS "${TEAMSERVER_URL}/api/operator/task" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"agent_id\":\"${AGENT_ID}\",\"command\":\"pwd\",\"args\":[],\"timeout_seconds\":20}" \
    | jq -r .task_id
)"
if [[ -z "$TASK_ID" || "$TASK_ID" == "null" ]]; then
  echo "[e2e] failed to create task" >&2
  exit 1
fi

echo "[e2e] waiting for task result"
RESULT_STATUS=""
for _ in $(seq 1 90); do
  results_json="$(curl -fsS "${TEAMSERVER_URL}/api/operator/results?agent_id=${AGENT_ID}" -H "Authorization: Bearer ${TOKEN}")"
  RESULT_STATUS="$(echo "$results_json" | jq -r --arg id "$TASK_ID" '.[] | select(.task_id==$id) | .status' | head -n1 || true)"
  if [[ -n "$RESULT_STATUS" && "$RESULT_STATUS" != "null" ]]; then
    break
  fi
  sleep 1
done
if [[ "$RESULT_STATUS" != "success" ]]; then
  echo "[e2e] task did not succeed; status=${RESULT_STATUS:-missing}" >&2
  echo "[e2e] teamserver log: /tmp/redforge.teamserver.log" >&2
  echo "[e2e] agent log: /tmp/redforge.agent.log" >&2
  exit 1
fi

echo "[e2e] verify telemetry endpoint returns data"
curl -fsS "${TEAMSERVER_URL}/api/operator/telemetry/latest?agent_id=${AGENT_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  | jq -e '.agent_id and .timestamp' >/dev/null

echo "[e2e] PASS"
echo "[e2e] logs: /tmp/redforge.teamserver.log /tmp/redforge.agent.log"
