# RedForgeC2 Architecture (Simulation-Only)

This document defines the high-level architecture for **RedForgeC2 as a lab-safe simulation platform**.

Non-goals (explicitly out of scope):
- Remote command execution (RCE) against hosts
- Arbitrary file upload/download
- Persistence, evasion, exploitation, lateral movement

## Components

### 1) Teamserver (Go)

Responsibilities:
- Accept agent registration and mock telemetry
- Maintain a per-agent task queue (simulation tasks only)
- Persist an append-only event log for auditability (`jsonl`)
- Provide a read-only API for operator UIs (agent list, task history, telemetry)

Default safety:
- Binds to `127.0.0.1` by default
- Can enable TLS (minimum TLS 1.2) when cert/key are provided

### 2) Agent Simulator (Rust)

Responsibilities:
- Register with the teamserver and maintain a heartbeat loop
- Emit synthetic telemetry events (CPU/mem-like samples, labels)
- Poll tasks and return simulated outcomes

Default safety:
- Refuses non-local server URLs unless explicitly overridden for isolated Docker/VM labs
- Executes no host commands and does not mutate the filesystem

### 3) Operator UIs

#### Web UI (React)

Responsibilities:
- Show agents list + status (derived from last-seen timestamps)
- Show task queue/history and outcomes
- Show telemetry charts/logs (simulated)

#### TUI (planned)

Responsibilities:
- Fast interactive listing and selection of agents
- Issue simulation tasks
- Tail logs and telemetry streams

## Data Flow (Simulation)

1. Agent → Teamserver: `register`
2. Agent → Teamserver: `telemetry` events
3. Operator UI → Teamserver: enqueue simulation tasks
4. Agent → Teamserver: poll tasks → submit results
5. Teamserver → Event log: append events for audit and replay

## Storage Model

- **Event log**: append-only `jsonl` (primary audit trail)
- **In-memory state**: current agents + queued/dispatched tasks (can be replaced later by a DB)

## Security Model (Lab-Only)

- TLS optional; when enabled, minimum TLS 1.2
- Strict input validation (`DisallowUnknownFields` for JSON where practical)
- No sensitive secrets or keys in repo; use env/config files outside source control

