# RedForgeC2 — Agent ↔ Teamserver Protocol (Simulation-Only)

## Overview

Defines **lab-safe**, **local-only** communication for simulated agents and a simulated teamserver.

Non-goals:
- No remote code execution, shells, file transfer, persistence, or lateral movement.
- No “fallback” transports (DNS/ICMP/etc.). Keep traffic on `localhost` or an isolated Docker/VM network.

- **Transport:** HTTP and/or WebSocket on `127.0.0.1` / Docker network  
- **Serialization:** JSON  
- **Authentication (optional):** simple dev token for lab sessions (not for real-world use)

## Message Types

| Direction | Type   | Payload                          | Description |
|-----------|--------|---------------------------------|-------------|
| Agent → Server | REG    | {agent_id, os, arch, meta}       | Initial registration |
| Server → Agent | TASK   | {task_id, task_type, params}      | Mock tasks to simulate |
| Agent → Server | RESULT | {task_id, outcome, status}        | Simulated results |
| Agent → Server | TELEMETRY | {metrics, heartbeat}           | Mock status update |
| Server → Agent | CONFIG | {settings}                        | Update simulation settings |

## Heartbeat & Tasking Flow

```

Agent boot → REGISTER → heartbeat loop (fixed/jittered interval for UI demo)
↓
Server receives heartbeat
↓
Server responds with TASKS (if any)
↓
Agent simulates task → sends RESULT

```

### Example Mock Tasks

- `collect_telemetry`: return synthetic CPU/RAM/NET samples
- `sleep`: simulate a delayed heartbeat
- `set_tag`: update a label for the agent (e.g., `lab-group-a`)
- `emit_alert`: generate a synthetic alert event for UI testing
- `echo`: return a provided string as a simulated outcome

### Future Enhancements

- Web UI live updates via local WebSocket
- More telemetry schemas for training scenarios (still simulated)
