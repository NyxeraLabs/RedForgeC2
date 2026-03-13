# RedForgeC2 — Rust Agent Architecture (Simulation-Only)

## Overview

This document describes a **lab-safe agent simulator** design for training and UI development.

- Self-contained Rust binary (simulator)  
- Multi-threaded/async: transport, task simulation, telemetry generation  
- Modular transport (local HTTP/WebSocket only)  

## Core Modules

### 1. Bootstrap

- Detect OS/Arch  
- Load configuration  
- Initialize logging  

### 2. Transport Layer

- Trait: `Transport { send(); receive(); }`  
- Implementations: HTTP / WebSocket (local-only)  

### 3. Task Executor

- Queue-based, async  
- **Simulation-only task runner**: maps `task_type` → synthetic outcomes  
- Hard denylist: no host command execution, no filesystem mutation  

### 4. Telemetry

- Synthetic metrics: CPU, memory, uptime, “network-like” samples  
- Deterministic/fuzzable generators for tests  

### 5. Safety

- Default-bind to `127.0.0.1` targets only  
- No propagation logic  
- Strict config validation (reject non-local endpoints unless explicitly allowed for Docker lab)  

## Diagram

```

+----------------------+
| Bootstrap            |
+----------------------+
|
+----------------------+
| Transport Layer      |
|  - HTTP/WSS (local)  |
+----------------------+
|
+----------------------+
| Task Simulator       |
+----------------------+
|
+----------------------+
| Telemetry Collector  |
+----------------------+
|
+----------------------+
| Safety               |
+----------------------+

```

✅ Notes

- Designed for containerized lab & isolated environments  
- Tasking and telemetry are **mocked** for education, QA, and UI demos
