# Go Backend (Skeleton)

Simulation-only Go teamserver implementation (local lab use only).

Guidelines:
- Bind to `127.0.0.1` by default
- Serve mock agent registration, telemetry ingest, and task polling
- No dangerous “command execution” endpoints; tasks are simulated/logged only

## Run

```bash
cd server/go_backend
go test ./...
go run ./cmd/teamserver -addr 127.0.0.1:8080
```

## Config

Config sources (merged in this order): defaults ← optional JSON file ← env ← flags.

Env:
- `REDFORGE_ADDR`
- `REDFORGE_BASE_URL`
- `REDFORGE_LOG_FILE`
- `REDFORGE_EVENT_LOG`
- `REDFORGE_LOG_JSON`
- `REDFORGE_TLS_CERT`
- `REDFORGE_TLS_KEY`

Flags:
- `-config <path>`
- `-addr <ip:port>`
- `-base-url <url>`
- `-log-file <path>`
- `-event-log <path>`
- `-log-json true|false`
- `-tls-cert <path>` + `-tls-key <path>` (enables HTTPS; TLS 1.2+)

## API (v1)

- `GET /healthz`
- `POST /api/v1/agents/register`
- `POST /api/v1/agents/{agent_id}/telemetry`
- `GET /api/v1/agents/{agent_id}/tasks`
- `POST /api/v1/agents/{agent_id}/tasks/{task_id}/result`
- `POST /api/v1/tasks/enqueue`

## Notes (Safety)

This is a simulation-only backend. It logs and tracks tasks/results but does not provide remote command execution or file transfer.
