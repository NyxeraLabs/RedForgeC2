# RedForgeC2 – Adversary Simulation Platform
Professional Development Roadmap

Stack

Agent: Rust  
Teamserver: Go  
UI: React + Typescript  
Database: PostgreSQL  
Infrastructure: Docker  
Realtime: WebSockets  

---

# Milestone 1 — Project Foundation

## Issue 1 — Repository Initialization

### Dev
- [x] Create GitHub repository
- [x] Configure monorepo structure
- [x] Setup Go modules
- [x] Setup Rust workspace
- [x] Initialize React UI

### QA
- [x] Configure CI pipeline
- [x] Setup linting tools
- [x] Configure automated testing

### Docs
- [x] Create README
- [x] Create CONTRIBUTING
- [x] Create LICENSE

### Commits
- [x] init repository
- [x] add monorepo structure
- [x] configure go module
- [x] configure rust workspace
- [x] initialize react app
- [x] add docker dev environment
- [x] add CI pipeline
- [x] add linting configuration

---

## Issue 2 — Architecture Definition

### Dev
- [x] Define system architecture
- [x] Define agent protocol
- [x] Define API design
- [x] Define transport abstraction

### QA
- [x] Define testing strategy
- [x] Define integration test plan

### Docs
- [x] architecture.md
- [x] threat-model.md
- [x] protocol-spec.md

### Commits
- [x] add architecture documentation
- [x] add protocol specification
- [x] add threat model

---

# Milestone 2 — Teamserver Core

## Issue 3 — Teamserver Skeleton

### Dev
- [x] Implement Go HTTP server
- [x] Implement configuration loader
- [x] Implement logging system
- [x] Implement API routing

### QA
- [ ] API unit tests
- [ ] Config parsing tests

### Commits
- [x] implement teamserver bootstrap
- [x] add configuration loader
- [x] add structured logging
- [x] add http router
- [x] add health endpoint
- [x] add graceful shutdown

---

## Issue 4 — Authentication System

### Dev
- [x] Implement operator login
- [x] Implement JWT authentication
- [x] Implement RBAC

### QA
- [ ] Auth bypass tests
- [ ] Token validation tests

### Docs
- [x] auth.md

### Commits
- [x] implement login endpoint
- [x] add jwt authentication
- [x] implement role system
- [x] add middleware auth validation
- [ ] add auth tests

---

# Milestone 3 — Agent Development

## Issue 5 — Agent Bootstrap

### Dev
- [x] Implement Rust agent startup
- [x] Implement environment discovery
- [x] Implement metadata collection

### QA
- [ ] OS compatibility tests
- [ ] metadata validation tests

### Docs
- [x] agent-architecture.md

### Commits
- [x] create rust agent project
- [x] implement agent bootstrap
- [x] implement metadata collection
- [x] add environment discovery
- [x] add serialization library
- [x] add config loader

---

## Issue 6 — Agent Registration

### Dev
- [x] Implement registration protocol
- [x] Implement agent ID generation
- [x] Implement heartbeat mechanism

### QA
- [ ] registration integration test
- [ ] heartbeat tests

### Commits
- [x] implement registration request
- [x] add agent id generation
- [x] add heartbeat protocol
- [x] add reconnect logic
- [x] add agent registry in server

---

# Milestone 4 — Command Execution

## Issue 7 — Tasking System

### Dev
- [x] Implement task queue
- [x] Implement task dispatcher
- [x] Implement result handler

### QA
- [ ] queue reliability tests
- [ ] concurrency tests

### Commits
- [x] implement task queue
- [x] implement dispatcher
- [x] implement result processor
- [x] add task status tracking
- [x] add task history storage

---

## Issue 8 — Core Agent Commands

### Dev
- [x] shell execution
- [x] file upload
- [x] file download
- [x] directory listing
- [x] process enumeration
- [x] working directory command

### QA
- [ ] command reliability tests
- [ ] file transfer tests

### Commits
- [x] implement shell command
- [x] implement file upload
- [x] implement file download
- [x] implement ls command
- [x] implement process list
- [x] add command validation

---

# Milestone 5 — Communication Layer

## Issue 9 — HTTPS Transport

### Dev
- [x] implement https transport
- [x] implement jitter timing
- [x] implement retry logic

### QA
- [x] latency tests
- [x] network failure tests

### Commits
- [x] implement https client
- [x] add transport encryption
- [x] implement jitter algorithm
- [x] add retry backoff
- [x] implement listener creation

---

## Issue 10 — Transport Abstraction

### Dev
- [x] create transport interface
- [x] implement tcp transport
- [x] implement dns placeholder

### QA
- [x] failover tests
- [x] transport switching tests

### Commits
- [x] add transport interface
- [x] implement tcp transport
- [x] add dns transport skeleton
- [x] implement transport registry

---

# Milestone 6 — Operator UI

## Issue 11 — UI Foundation

### Dev
- [x] dashboard layout
- [x] authentication page
- [x] navigation system

### QA
- [x] UI navigation tests
- [x] authentication flow tests

### Commits
- [x] initialize react ui
- [x] implement login page
- [x] implement dashboard
- [x] add sidebar navigation
- [x] add api client

---

## Issue 12 — Agent Dashboard

### Dev
- [x] agent list view
- [x] agent metadata display
- [x] agent status indicators

### QA
- [x] real-time update tests

### Commits
- [x] implement agent table
- [x] add agent details panel
- [x] add real-time websocket updates

---

## Issue 13 — Interactive Console

### Dev
- [x] agent console terminal
- [x] command submission
- [x] result streaming

### QA
- [x] console reliability tests

### Commits
- [x] implement console UI
- [x] add websocket command streaming
- [x] add command history

---

# Milestone 7 — Advanced Operator Features

## Issue 14 — Listener Manager

### Dev
- [x] create listener UI
- [x] delete listener
- [x] configure listeners

### QA
- [x] listener creation tests

### Commits
- [x] implement listener API
- [x] add listener UI
- [x] add listener status monitoring

---

## Issue 15 — Loot Manager

### Dev
- [x] implement file storage
- [x] implement screenshot storage
- [x] implement credential storage

### QA
- [x] large file tests
- [x] storage reliability tests

### Commits
- [x] implement loot storage
- [x] add loot search
- [x] add loot tagging
- [x] implement loot download

---

# Milestone 8 — Adversary Simulation

## Issue 16 — Playbook Engine

### Dev
- [x] yaml playbook format
- [x] playbook execution engine
- [x] task sequencing

### QA
- [x] playbook execution tests

### Commits
- [x] implement playbook parser
- [x] implement playbook executor
- [x] add playbook scheduling

---

## Issue 17 — Detection Telemetry

### Dev
- [x] technique tracking
- [x] detection metadata
- [x] telemetry dashboard

### QA
- [x] telemetry accuracy tests

### Commits
- [x] implement telemetry collector
- [x] add telemetry database
- [x] implement detection dashboard

---

# Milestone 9 — Security Hardening

### Dev
- [ ] TLS enforcement
- [ ] audit logging
- [ ] API token system

### QA
- [ ] auth bypass tests
- [ ] session security tests

### Commits
- [ ] enforce TLS
- [ ] add audit logging
- [ ] implement api tokens
- [ ] add security tests

---

# Milestone 10 — End-to-End Testing

### Dev
- [x] build docker lab
- [ ] deploy test infrastructure
- [ ] automate agent deployment

### QA
- [ ] multi-agent simulation
- [ ] network failure testing
- [ ] server restart testing

### Commits
- [x] add docker lab
- [ ] implement integration tests
- [ ] add multi-agent simulation

---

# Milestone 11 — Documentation & Release

### Docs
- [ ] getting-started.md
- [ ] operator-guide.md
- [ ] developer-guide.md
- [ ] architecture diagrams
- [ ] lab setup guide

### Release
- [ ] record demo video
- [ ] publish GitHub release
- [ ] write technical blog

### Commits
- [x] add documentation
- [x] add diagrams
- [x] prepare v1 release

---

# Milestone 12 — UI/TUI Mockup Parity (Actual)

Bring the implemented UI/TUI up to the feature set described in `docs/ui-mockups.md` and validated in `docs/ui-implementation-audit.md`.

## Issue 18 — Web UI Refactor & Navigation

### Dev
- [ ] split `App.tsx` into pages/components
- [ ] add navigation (tabs/sidebar) and basic routing
- [ ] improve agents view (table + details panel)
- [ ] compute agent status from `last_seen`

### QA
- [ ] navigation smoke tests
- [ ] auth persistence tests (localStorage token)

### Commits
- [ ] refactor ui structure
- [ ] add routing/navigation
- [ ] implement agent table + details

---

## Issue 19 — Interactive Console UX

### Dev
- [ ] per-agent console transcript view
- [ ] command history (client-side)
- [ ] results view with refresh + basic filtering
- [ ] nicer task submission UX (timeouts, validation)

### QA
- [ ] command submission reliability tests
- [ ] results rendering tests (large output)

### Commits
- [ ] add console page
- [ ] add command history
- [ ] improve results view

---

## Issue 20 — File Manager (v0 via Agent Commands)

### Dev
- [ ] browse remote filesystem via `ls <path>`
- [ ] upload: local file → base64 → `upload <remote_path> <base64>`
- [ ] download: `download <remote_path>` → base64 decode → browser download

### QA
- [ ] upload/download e2e tests (small + large files)
- [ ] path traversal and invalid path tests

### Commits
- [ ] implement file manager page (v0)
- [ ] add upload/download helpers

---

## Issue 21 — Telemetry Persistence + Operator APIs

### Dev
- [ ] persist telemetry from `/api/heartbeat` to database
- [ ] add operator endpoints for telemetry (latest + optional series)
- [ ] add Web UI telemetry panel (latest metrics first)

### QA
- [ ] telemetry ingest tests
- [ ] operator telemetry API tests

### Commits
- [ ] add telemetry tables/migrations
- [ ] store telemetry on heartbeat
- [ ] add telemetry operator endpoints
- [ ] implement telemetry UI panel

---

## Issue 22 — Alerts / Event Feed

### Dev
- [ ] define alert/event model (severity, type, agent, timestamp, message)
- [ ] persist alerts/events to database
- [ ] operator endpoints for alerts (+ basic filtering)
- [ ] Web UI alerts panel + notifications

### QA
- [ ] alert ingestion tests
- [ ] alert filtering tests

### Commits
- [ ] add alerts tables/migrations
- [ ] add alerts operator endpoints
- [ ] implement alerts UI panel

---

## Issue 23 — Real-Time Updates (SSE/WebSockets)

### Dev
- [ ] add real-time stream for agent updates, results, telemetry, and alerts (SSE or WebSockets)
- [ ] update UI to subscribe instead of polling (fallback to polling)

### QA
- [ ] disconnect/reconnect tests
- [ ] multi-agent update ordering tests

### Commits
- [ ] add realtime transport endpoints
- [ ] wire UI realtime subscriptions

---

## Issue 24 — Reports (JSON Export → PDF Later)

### Dev
- [ ] JSON export of sessions/tasks/results/telemetry/alerts
- [ ] Web UI reports page to download exports
- [ ] (later) PDF export

### QA
- [ ] export correctness tests
- [ ] large dataset export tests

### Commits
- [ ] add reports endpoints (json)
- [ ] implement reports UI

---

## Issue 25 — Network Topology / Map Visualization

### Dev
- [ ] define topology graph model (nodes/edges)
- [ ] persist topology events/edges
- [ ] Web UI topology panel (D3/Recharts)

### QA
- [ ] topology model validation tests
- [ ] UI render performance smoke tests

### Commits
- [ ] add topology schema
- [ ] add topology endpoints
- [ ] implement topology UI

---

## Issue 26 — TUI MVP

### Dev
- [ ] choose TUI stack (Python `textual`/`rich` or Rust `ratatui`/`crossterm`)
- [ ] implement sessions panel + agent selection
- [ ] implement command input + task enqueue
- [ ] implement results view (polling first)
- [ ] add telemetry + alerts panels after Issue 21/22

### QA
- [ ] TUI smoke tests (login, list agents, run command, view result)
- [ ] cross-platform checks (linux/mac/windows where applicable)

### Commits
- [ ] add tui project skeleton
- [ ] implement tui sessions + tasking
- [ ] implement tui results view
