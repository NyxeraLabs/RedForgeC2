# RedForgeC2 — Agent ↔ Teamserver Protocol

## Overview

Defines secure communication, tasking, and telemetry between Rust Agents and the Go Teamserver.

- **Transport:** HTTPS / WebSocket primary, fallback to DNS/ICMP  
- **Serialization:** JSON for commands, optional CBOR for telemetry  
- **Encryption:** AES-256-CBC per message  
- **Authentication:** Agent registration token (UUID-based)  

## Message Types

| Direction | Type   | Payload                          | Description |
|-----------|--------|---------------------------------|-------------|
| Agent → Server | REG    | {agent_id, os, arch, meta}       | Initial registration |
| Server → Agent | TASK   | {task_id, command, args}          | Commands to execute |
| Agent → Server | RESULT | {task_id, output, status}         | Execution results |
| Agent → Server | TELEMETRY | {metrics, heartbeat}           | Status update |
| Server → Agent | CONFIG | {settings}                        | Update agent configuration |

## Heartbeat & Tasking Flow

```

Agent boot → REGISTER → heartbeat loop (NHPP-jittered interval)
↓
Server receives heartbeat
↓
Server responds with TASKS (if any)
↓
Agent executes task → sends RESULT

```

### Security Measures

- AES-256-CBC encrypted messages  
- HMAC verification  
- Replay protection (sequence numbers)  
- Optional multi-layer transport fallback (HTTP → ICMP → DNS)  

### Future Enhancements

- P2P agent relay mesh  
- Multi-stage command pipelines  
- Transport polymorphism & dynamic port hopping