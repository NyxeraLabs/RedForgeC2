# RedForgeC2 — Rust Agent Architecture

## Overview

Professional-grade implant design following industry standards (Sliver/Havoc/Cobalt Strike style)

- Self-contained Rust binary  
- Multi-threaded: transport, task execution, telemetry  
- Memory-safe with `zeroize` for sensitive data  
- Modular transport abstraction  

## Core Modules

### 1. Bootstrap

- Detect OS/Arch  
- Load configuration  
- Initialize logging  

### 2. Transport Layer

- Trait: `Transport { send(); receive(); }`  
- Implementations: HTTPS/WSS, DNS tunneling, ICMP signaling  
- Dynamic transport switching  

### 3. Task Executor

- Queue-based, async  
- Shell execution, file upload/download  
- Optional PowerShell in-memory execution (Windows)  

### 4. Telemetry

- System metrics: CPU, memory, uptime  
- Artifact harvest: bash_history, known_hosts, registry  
- Zeroize sensitive data after sending  

### 5. Persistence

- Registry run keys (Windows)  
- Systemd service (Linux)  
- Single instance mutex  

### 6. Safety & Kill-Switch

- Global kill-switch: terminate agent + clean traces  
- Rate-limiting propagation  
- Optional sandbox detection  

### 7. Encryption & Security

- AES-256-CBC per message  
- HMAC verification  
- Replay protection & sequence tracking  
- Multi-layer transport fallback  

## Diagram

```

+----------------------+
| Bootstrap            |
+----------------------+
|
+----------------------+
| Transport Layer      |
|  - HTTPS/WSS         |
|  - DNS fallback      |
|  - ICMP fallback     |
+----------------------+
|
+----------------------+
| Task Executor        |
+----------------------+
|
+----------------------+
| Telemetry Collector  |
+----------------------+
|
+----------------------+
| Persistence / Safety |
+----------------------+

```

✅ Notes

- Highly modular (hot-swappable transports)  
- Fully async loop  
- Memory-safe with zeroize  
- Designed for containerized lab & isolated environments