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

## API (v1)

- `GET /healthz`
- `POST /api/v1/agents/register`
- `POST /api/v1/agents/{agent_id}/telemetry`
- `GET /api/v1/agents/{agent_id}/tasks`
- `POST /api/v1/agents/{agent_id}/tasks/{task_id}/result`
- `POST /api/v1/tasks/enqueue`
