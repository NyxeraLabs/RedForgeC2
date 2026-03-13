# RedForgeC2 Protocol Specification (Simulation-Only)

This file is the **canonical spec** for the agent ↔ teamserver protocol in this repository.

See also:
- `protocol.md` for an overview and examples

## Transport

- HTTP (JSON) on `127.0.0.1` by default
- Optional HTTPS (TLS 1.2+) when cert/key are provided (lab-only)

## Endpoints (v1)

- `GET /healthz`
- `POST /api/v1/agents/register`
- `GET /api/v1/agents`
- `POST /api/v1/agents/{agent_id}/telemetry`
- `GET /api/v1/agents/{agent_id}/tasks`
- `POST /api/v1/agents/{agent_id}/tasks/{task_id}/result`
- `POST /api/v1/tasks/enqueue`

## Message Schemas (JSON)

### Register

Request: `RegisterRequest`
- `agent_id` (string, optional)
- `os` (string)
- `arch` (string)
- `hostname` (string)
- `meta` (object<string,string>, optional)

Response: `RegisterResponse`
- `agent_id` (string)
- `server_time` (RFC3339Nano string)

### Telemetry

Request: `TelemetryRequest`
- `events` (array of `TelemetryEvent`)

`TelemetryEvent`
- `timestamp` (string; agent may send RFC3339; server stores as-is)
- `metrics` (object<string,number>, optional)
- `labels` (object<string,string>, optional)

### Tasking (Simulation)

`Task`
- `task_id` (string)
- `task_type` (string; simulation-only types)
- `params` (object, optional)
- `created_at` (RFC3339Nano string)
- `status` (`queued` | `dispatched` | `completed`)

Poll response: `PollTasksResponse`
- `tasks` (array of `Task`)

Enqueue request: `EnqueueTaskRequest`
- `agent_id` (string)
- `task_type` (string)
- `params` (object, optional)

Result submission: `SubmitResultRequest`
- `outcome` (object, optional)
- `status` (string, optional; e.g., `"simulated"`)

## Allowed Task Types

Simulation-only (non-exhaustive):
- `echo`
- `sleep` (clamped to safe max in agent)
- `collect_telemetry`

Any task types that would execute host commands or move files are out of scope.

