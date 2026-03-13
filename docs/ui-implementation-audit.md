# UI/TUI Implementation Audit (vs Mockups)

This repo contains a **UI/TUI mockup blueprint** in `docs/ui-mockups.md` (sometimes referenced as “UI-MOCKUPS.md”), plus PNG assets like `docs/ui-mockups.png`.

This document answers: **what is actually implemented today** in the repo, and **what’s missing** to reach the mockups.

## Source of truth (current code)

**Implemented Web UI (current)**
- Single-page React UI: `ui/src/App.tsx`
- UI README explicitly says placeholder: `ui/README.md`

**Implemented Teamserver operator APIs (current)**
- HTTP endpoints are registered in: `teamserver/internal/server/server.go`
  - `POST /api/login`
  - `GET  /api/operator/agents`
  - `POST /api/operator/task`
  - `GET  /api/operator/results?agent_id=...`

**Agent capabilities relevant to “File Manager”**
- Agent command handlers are in: `agent/src/bootstrap.rs` (e.g., `upload`, `download`, `ls`, `ps`, `pwd/cwd`, and arbitrary command execution)

**Notably absent**
- No WebSocket/SSE endpoints in `teamserver/` (only HTTP polling-style APIs exist).
- No TUI application directory/entrypoint in the repo.
- Telemetry is sent in the heartbeat payload (`teamserver/internal/api/models.go`) but is not persisted/returned to operators.

---

## Checklist: Web UI (mockups → implementation)

Mockup reference: `docs/ui-mockups.md` section “2️⃣ Web User Interface (Web UI)”.

### Authentication
- Mockup: “Logged in as Operator1”, JWT auth.
- Status: **Implemented (minimal)**.
- Evidence:
  - API: `POST /api/login` in `teamserver/internal/server/server.go`
  - UI: login form + token in localStorage in `ui/src/App.tsx`

### Navigation + multi-tab layout (Sessions / Telemetry / File Manager / Reports)
- Mockup: explicit tabs + distinct panels.
- Status: **Missing** (single page only, no routing, no tabs).
- Evidence: only `ui/src/App.tsx` exists; no router/components under `ui/src/`.

### Sessions view (agent list + metadata + status indicators)
- Mockup: sessions metadata and status.
- Status: **Partial**.
- Evidence:
  - API: `GET /api/operator/agents`
  - UI: dropdown list of agents (no table/details view) in `ui/src/App.tsx`
  - Data model in UI includes `last_seen`, but server registry returns `last_seen` only if populated; there is no “status” computation in UI.

### Interactive console (command console + history + streaming results)
- Mockup: “Command Console”, autocomplete, multiline scripts, result streaming.
- Status: **Partial**.
- Evidence:
  - API supports enqueueing a task and later fetching results (`/api/operator/task`, `/api/operator/results`) in `teamserver/internal/server/server.go`
  - UI supports sending a command + args and manually “Refresh Results” in `ui/src/App.tsx`
- Missing pieces:
  - Real-time result streaming (no WebSocket/SSE)
  - Command history UI, multiline editor, autocomplete, per-agent console transcript

### Telemetry dashboard (charts/sparklines, uptime, last check-in)
- Mockup: dedicated Telemetry panel and dashboards.
- Status: **Missing (data path incomplete)**.
- Evidence:
  - Heartbeat payload contains telemetry (`teamserver/internal/api/models.go`)
  - Teamserver ignores telemetry content in `teamserver/internal/server/server.go` (only updates last-seen and returns tasks)
- Missing pieces:
  - Persist telemetry in DB
  - Operator API to retrieve latest/series telemetry
  - UI components/charts

### Alerts panel (severity feed, filters, notifications)
- Mockup: alert feed, filters, toasts.
- Status: **Missing**.
- Evidence:
  - No alert models/endpoints/tables in `teamserver/`
  - UI has no alerts panel in `ui/src/App.tsx`

### File Manager (upload/download, browsing)
- Mockup: File Manager tab with upload/download.
- Status: **Missing in UI; partially supported by agent commands**.
- Evidence:
  - Agent supports `upload` and `download` commands in `agent/src/bootstrap.rs`
  - Operator can enqueue tasks; results are stored in DB (`teamserver/internal/registry/registry.go`)
- Missing pieces:
  - UI workflows to:
    - Read local file → base64 → send `upload` task
    - Decode base64 output from `download` task into a browser download
  - Remote browsing UX (can be approximated via `ls` initially)

### Network topology / map visualization
- Mockup: map + D3-style node graph of lateral movement.
- Status: **Missing (no data + no UI)**.
- Evidence:
  - No topology/lateral movement data model in `teamserver/`
  - No D3/Recharts components in `ui/src/`

### Reports (export PDF/JSON; MITRE mapping)
- Mockup: reporting module with exports and ATT&CK mapping.
- Status: **Missing**.
- Evidence:
  - No report endpoints or MITRE tagging model in `teamserver/`
  - No Reports UI in `ui/src/`

### RBAC (Operator/Analyst/QA/Viewer) + 2FA
- Mockup: role-based access control and optional 2FA.
- Status: **Missing / simplified**.
- Evidence:
  - Teamserver issues tokens with a single hard-coded role “admin” in `teamserver/internal/server/server.go`
  - Middleware requires role “admin” for operator endpoints (see `teamserver/internal/server/server.go`)

### Dark/Light mode
- Mockup: dark/light theme.
- Status: **Partial** (dark-only styling).
- Evidence: UI uses dark palette classes in `ui/src/App.tsx` with no toggle.

---

## Checklist: TUI (mockups → implementation)

Mockup reference: `docs/ui-mockups.md` section “1️⃣ Terminal User Interface (TUI)”.

- Status: **Not implemented** (no TUI app present).
- Evidence:
  - No `tui/` directory, no Python `textual` app, no Rust `ratatui/tui-rs` crate in workspace.
  - `docs/ui-mockups.md` explicitly lists TUI as a “Next Step” (`docs/ui-mockups.md` section “4️⃣ Next Steps for Dev”).

What exists that a future TUI can use:
- Operator APIs: `docs/operator-guide.md` + `teamserver/internal/server/server.go`

---

## Notes on `docs/roadmap.md`

`docs/roadmap.md` contains milestones that mark many UI items as “done” (dashboard, websocket updates, listener UI, loot manager, etc.). Those features **do not appear in the current code** under `ui/src/` or `teamserver/` (no websocket code, no listener/loot/report endpoints).

Treat the roadmap as **aspirational / stale** unless updated to match the codebase.

---

## Smallest plan to reach “mockups parity” (MVP-first)

This is the smallest sequence that gets you close to the mockup experience while reusing what already exists.

### Web UI MVP (build on current APIs)
1. **Split UI into pages/components** (Login, Agents, Console, Results) and add basic navigation (React Router).
2. **Agents page:** table view (hostname/os/arch/version/last_seen) + status pill computed from `last_seen`.
3. **Console page:** per-agent console transcript + command history (client-side) + results view.
4. **File Manager (v0):**
   - “Browse” via `ls <path>`
   - Upload: local file → base64 → task `upload <remote_path> <base64>`
   - Download: task `download <remote_path>` then decode base64 output to a browser download
5. **Polling-based “real-time” (v0):** poll `/api/operator/agents` and `/api/operator/results` on intervals; defer WebSockets.

### Teamserver API additions needed for telemetry/alerts
1. **Persist telemetry** on `/api/heartbeat` (store latest + optional time-series).
2. Add `GET /api/operator/telemetry/latest?agent_id=...` and optionally `/api/operator/telemetry/series?...`.
3. Define an **alerts/events** table + model (severity, type, message, agent_id, timestamp).
4. Add `GET /api/operator/alerts` (+ filters) and optionally `POST /api/operator/alerts/ack`.

### Web UI “mockup features” (after MVP)
1. **Telemetry charts:** render latest + time-series.
2. **Alerts panel + notifications:** toast + persistent feed.
3. **Reports export:** `GET /api/operator/reports/session.json` (start with JSON), then add PDF later.
4. **WebSockets/SSE:** replace polling for agent list, results, alerts, telemetry.
5. **Topology/Map:** only after you have a real data source for movement/edges.

### TUI MVP (terminal operator experience)
1. Choose stack (Python `textual` is fastest to iterate; Rust `ratatui` for a single-binary feel).
2. Implement:
   - Sessions panel (agents list)
   - Tasking prompt (enqueue `/api/operator/task`)
   - Results view (poll `/api/operator/results`)
3. Add telemetry + alerts panels once the server exposes them (API additions above).

---

## GitFlow: “after each task complete”

If you want a consistent GitFlow-style delivery loop for each task (UI page, API endpoint, etc.), use:

1. Start work:
   - `git checkout develop`
   - `git pull`
   - `git checkout -b feature/<short-task-name>`
2. Commit as you go:
   - `git add -A`
   - `git commit -m "feat(ui): <what changed>"`
3. Finish task:
   - `git push -u origin feature/<short-task-name>`
   - Open PR → merge into `develop` (squash or merge per your preference)
4. Release (when a set of tasks is ready):
   - `git checkout -b release/<version>`
   - bump versions/changelog
   - merge `release/<version>` → `main` and tag
   - merge `main` → `develop`

Hotfix (production only):
- `git checkout main && git checkout -b hotfix/<issue>`
- merge back into `main` and `develop`

