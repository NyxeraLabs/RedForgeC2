# RedForgeC2 Protocol Specification

This document defines the message formats and communication patterns used between Agents and the Teamserver.

## 1. Transport

### Primary Transport
- **HTTPS** (REST / WebSocket)
- TLS 1.2+ enforced
- Payloads are encrypted at the application layer (AES-256-CBC + HMAC) in addition to TLS.

### Fallback Transports
- **DNS tunneling** (future)
- **ICMP signaling** (future)

## 2. Message Envelope

All messages share a common envelope:

```json
{
  "agent_id": "<uuid>",
  "sequence": 123,
  "type": "REG|TASK|RESULT|TELEMETRY|CONFIG",
  "payload": { /* type-specific payload */ },
  "hmac": "<hex>"
}
```

- `agent_id`: UUID identifying the agent.
- `sequence`: Monotonic message counter for replay protection.
- `type`: Message type.
- `payload`: Message-specific data structure.
- `hmac`: HMAC-SHA256 over the serialized message fields.

## 3. Message Types

### 3.1 REG (Registration)
**Direction:** Agent → Server

```json
{
  "agent_id": "<uuid>",
  "os": "linux",
  "arch": "x86_64",
  "hostname": "lab-host",
  "version": "0.1.0",
  "metadata": {
    "user": "redteam",
    "kernel": "5.15.0"
  }
}
```

### 3.2 TASK (Task Delivery)
**Direction:** Server → Agent

```json
{
  "task_id": "<uuid>",
  "command": "shell",
  "args": ["whoami"],
  "timeout": 120
}
```

### 3.3 RESULT (Task Result)
**Direction:** Agent → Server

```json
{
  "task_id": "<uuid>",
  "status": "success",
  "output": "redteam\n",
  "error": "",
  "timestamp": "2026-03-13T12:00:00Z"
}
```

### 3.4 TELEMETRY (Heartbeat)
**Direction:** Agent → Server

```json
{
  "cpu": 12.4,
  "memory": 4253120,
  "uptime": 12345,
  "timestamp": "2026-03-13T12:00:00Z"
}
```

### 3.5 CONFIG (Agent Configuration)
**Direction:** Server → Agent

```json
{
  "heartbeat_interval": 60,
  "max_jitter": 20,
  "allowed_transports": ["https", "dns"]
}
```

## 4. Security

- **Encryption:** AES-256-CBC per message with random IV.
- **Integrity:** HMAC-SHA256 over the serialized envelope.
- **Replay Protection:** Sequence numbers and timestamp windows.

## 5. Versioning

Every message includes an optional `protocol_version` field. When the teamserver or agent detects a version mismatch, it returns a `CONFIG` message with the required protocol version.

---

> This document is the authoritative protocol spec used for implementation and testing.
