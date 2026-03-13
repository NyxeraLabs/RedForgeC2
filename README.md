# RedForgeC2 — Modern Adversary Emulation & C2 Framework

![RedForgeC2 Logo](docs/assets/redforgec2_logo.png)

**RedForgeC2** is a high-fidelity Command & Control (C2) and adversary emulation framework designed for **Red Team, Purple Team, and security R&D**. It combines stealthy Rust agents, a Go teamserver backend, and a modern React/TypeScript operator console. Inspired by **Cobalt Strike, Sliver, Havoc, and Mythic**, RedForgeC2 allows simulation of advanced threat behaviors, secure telemetry collection, lateral movement, and defensive gap testing.

⚠️ **DISCLAIMER**  
For **educational and authorized security testing purposes only**. Unauthorized use against systems you do not own or have explicit permission to test is illegal. Users are responsible for complying with all applicable laws. The authors assume no liability.

---

## 🚀 Features

- Multi-platform Rust agents (Windows, Linux, Android)  
- Secure HTTPS / WebSocket / DNS / ICMP / NTP transport with fallback  
- Modular transport abstraction for stealth & evasion  
- Asynchronous task queue and telemetry collection  
- Artifact harvesting: bash_history, known_hosts, registry, etc.  
- Persistence and survivability mechanisms (Registry, Systemd, mutex)  
- Kill-switch & self-cleanup routines  
- Modern web-based UI for tasking, telemetry, and reporting  
- Real-time alerts and Canary tripwires  
- MITRE ATT&CK mapped simulation activities  
- Advanced EDR/AV evasion techniques (like Cobalt Strike, Havoc, Sliver)  
- Multi-stage build & deploy pipeline  
- Containerized lab environment support  

---

## 🛠 Technology Stack

| Component        | Technology / Language             |
|-----------------|----------------------------------|
| Teamserver       | Go (API, tasking, telemetry DB) |
| Agent            | Rust (async, multi-threaded)     |
| Operator Console | React + TypeScript + TailwindCSS |
| Database         | SQLite / Optional PostgreSQL    |
| Transport        | HTTPS / WSS / DNS / ICMP / NTP  |
| Build & Scripts  | Make, Shell, Docker             |

---

## 📁 Repository Structure

```

RedForgeC2/
├─ agent/          # Rust implant source
├─ teamserver/     # Go backend code
├─ ui/             # React operator console
├─ docs/           # Documentation & UI mockups
├─ scripts/        # Build, deploy, lab scripts
├─ tests/          # Unit & E2E tests
├─ docker/         # Lab containers
├─ .github/        # CI/CD workflows
├─ README.md
├─ LICENSE
└─ Makefile

````

---

## ⚙️ Quick Start

### 1. Install Dependencies

```bash
# Backend
cd teamserver
go mod tidy

# UI
cd ui
npm install
````

### 2. Build Rust Agent

```bash
cd agent
cargo build --release
```

### 3. Launch Teamserver

```bash
cd teamserver
go run main.go
```

### 4. Launch Operator Console

```bash
cd ui
npm start
```

### 5. Deploy Agent in Lab

* Deploy the compiled Rust agent on an isolated VM or container.
* Agent registers automatically with the teamserver.
* Start tasking, telemetry collection, and forensic simulation.

---

## 🧩 Operational Workflow

1. **Agent Registration:** Unique agent UUID + system metadata
2. **Heartbeat Loop:** NHPP-jittered intervals for stealth
3. **Task Assignment:** Pull-based execution from teamserver
4. **Telemetry:** CPU, memory, network, artifact harvesting
5. **Execution:** Remote shell, file transfer, persistence simulation
6. **Reporting:** Live telemetry, forensic PDF output, MITRE ATT&CK mapping
7. **Kill-switch:** Global self-cleanup and trace removal

---

## 📊 UI Screenshots / Mockups

![UI Mockups](docs/ui-mockups.png)
*Mockups for the initial operator experience.*

---

## 📝 Documentation

* [Repo Structure](docs/repo-structure.md)
* [Agent ↔ Teamserver Protocol](docs/protocol.md)
* [UI Mockups](docs/ui-mockups.md)
* [Rust Agent Architecture](docs/agent-architecture.md)
* [Roadmap & Dev Plan](docs/roadmap.md)

---

## 🧪 Testing

* Unit tests: `cargo test` (agent) / `go test ./...` (teamserver)
* Integration tests: `tests/integration/`
* E2E tests: Simulated lab environment via Docker

---

## 🔐 Security & Safety Guidelines

* Run only in isolated lab environments
* Keep sensitive keys out of source control (`.env`, `.pem`, `.db`)
* Use kill-switch to avoid unintended propagation
* Document all simulations for authorized IR / Red Team exercises

---

## 📄 License

Licensed under the **Business Source License 1.1 (BSL 1.1)**.

- Change Date: **2033-03-13**
- Change License: **Apache-2.0**

See `LICENSE` for the full text.

---

**RedForgeC2** — Modern C2, Professional Agent Architecture, Advanced Evasion, Tactical UI.

---
