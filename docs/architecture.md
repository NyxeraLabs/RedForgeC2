# RedForgeC2 — System Architecture

This document captures the high-level architecture of RedForgeC2 and describes the major subsystems, data flows, and design principles.

## High-Level Components

### 1. Teamserver (Go)
- **Role:** Central command-and-control server.
- **Responsibilities:**
  - Agent registration & authentication
  - Task queue management
  - Telemetry ingest & storage
  - API endpoints for operator UI
  - Transport abstraction for multiple C2 channels

### 2. Agent (Rust)
- **Role:** Implant running on target systems.
- **Responsibilities:**
  - Bootstrap and environment discovery
  - Secure registration & heartbeat
  - Task execution (shell, file operations, etc.)
  - Telemetry collection & reporting
  - Transport abstraction (HTTPS/WSS, DNS, ICMP, etc.)
  - Persistence and kill-switch

### 3. Operator UI (React)
- **Role:** Operator console for managing agents and tasks.
- **Responsibilities:**
  - Authentication / session management
  - Agent inventory and status
  - Tasking console (shell, file transfer, etc.)
  - Telemetry visualization and alerts
  - Reports and threat tracking

## Data Flow Overview

1. Agent starts and performs **registration** over the primary transport.
2. Teamserver stores agent metadata and returns an authentication token.
3. Agent enters a **heartbeat loop**, polling for tasks.
4. Teamserver assigns tasks via the task queue, and agent reports results.
5. Telemetry is ingested and recorded for operator analysis.

## Modular Architecture

### Transport Abstraction
- Both agent and teamserver expose a **transport interface**.
- The agent can switch between transports (HTTPS → DNS → ICMP) based on reachability.
- The teamserver can accept incoming messages over different transports but normalizes them into a common internal pipeline.

### Tasking Engine
- The teamserver maintains per-agent task queues (in-memory / persistent).
- Tasks are represented as structured objects (command, args, ID, and metadata).
- Results are validated and stored for operator retrieval.

### Security Boundary
- All C2 messages are encrypted with **AES-256-CBC** and authenticated via **HMAC**.
- Agents authenticate via a unique agent token and ID.
- Operators authenticate via JWT / session tokens (future phases).

## Deployment Patterns

- **Lab Deployment:** Docker Compose lab with PostgreSQL, teamserver, UI, and agent runners.
- **Production-Style Deployment:** Separate services, TLS termination, and hardened network segmentation.

---

## References
- Protocol: `docs/protocol-spec.md`
- Transport design: `docs/protocol-spec.md#transport-abstraction`
- Agent architecture: `docs/agent-architecture.md`
