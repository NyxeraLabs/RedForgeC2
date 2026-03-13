# RedForgeC2 UI & UX Mockups

This document provides a comprehensive blueprint for the **User Interfaces** of RedForgeC2, covering both **Terminal UI (TUI)** and **Web UI**, including mockups, layouts, interactions, and operator workflows.  

⚠️ Simulation-only: all “commands”, “tasks”, “alerts”, and “telemetry” are **mocked** for training and UI testing. No host command execution, file transfer, persistence, lateral movement, or exfiltration.

The goal is to provide **high-fidelity training interfaces** inspired by modern operator consoles, while keeping all behaviors **lab-safe**.

---

## 1️⃣ Terminal User Interface (TUI)

### 1.1 Overview

The TUI is designed for **local lab training**, enabling operators to:

- Interact with multiple agents and sessions
- Monitor telemetry in real time
- Trigger **mock tasks** and simulated scenarios
- Perform automated reporting and logging

**Tech Stack:** Python (`rich`, `textual`) or Rust (`tui-rs`, `crossterm`).

---

### 1.2 Layout

```

+----------------------------------------------------------+
| REDFORGE-C2 TUI v1.0                                    |
| [Connected: 3 Agents | Alerts: 1 RED | Tasks: 5]       |
+----------------------+-------------------+---------------+

| Sessions Panel                                               | Telemetry Panel | Alerts Panel   |               |
| ------------------------------------------------------------ | --------------- | -------------- | ------------- |
| ID                                                           | Hostname        | CPU: 32%       | [RED] Exfil   |
| 001                                                          | ALPHA-LAB       | RAM: 4.2GB     | [YEL] Shell   |
| 002                                                          | BRAVO-LAB       | NET: 12 Mbps   | [GRN] Persist |
| 003                                                          | GAMMA-LAB       | STATUS: ACTIVE |               |
| +----------------------+------------------+---------------+  |                 |                |               |
| Command Input: >                                             |                 |                |               |
| +----------------------------------------------------------+ |                 |                |               |

```

**Panels:**

1. **Sessions Panel**  
   - List all connected agents, platform, status, last check-in  
   - Supports `select <ID>` or `use <ID>` to switch session focus  

2. **Telemetry Panel**  
   - Real-time metrics: CPU, RAM, network, battery (mobile), OS version  
   - Heartbeat timestamps and last command executed  

3. **Alerts Panel**  
   - Color-coded alerts: RED (critical), YELLOW (warning), GREEN (info)  
   - Event logs: simulated alerts (e.g., “policy violation”, “anomaly detected”)  

4. **Command Input**  
   - Auto-complete commands (tab support)  
   - Context-sensitive help: `help <command>`  
   - Multi-line input for scripts  

---

### 1.3 Key Commands (TUI)

| Command | Description |
|---------|------------|
| `list` | List all connected agents with metadata |
| `use <ID>` | Focus session on specific agent |
| `task <type> [json]` | Enqueue a mock task for the focused agent |
| `telemetry` | Show recent telemetry for focused agent |
| `alerts` | Show alert log history |
| `tasks` | Show pending C2 tasks |
| `kill <ID>` | Terminate agent safely |

---

### 1.4 Mockup Features

- **Tabbed Session View:** Multiple agents visible, switchable with `Ctrl+Tab`  
- **Dynamic Color Alerts:** Critical events flash in RED with bold text  
- **Mini Graphs:** CPU/Memory/Network as sparkline ASCII charts  
- **Auto-Scrolling Logs:** Tail-like behavior for telemetry updates  
- **Command History:** Scroll with `Up/Down` keys  

---

## 2️⃣ Web User Interface (Web UI)

### 2.1 Overview

The Web UI provides a **rich visual dashboard** for operators, analysts, and purple team members:

- Multi-agent map & telemetry
- Network topology and scenario visualization (simulated)
- Real-time alerts and forensic data
- Configurable dashboards and MITRE ATT&CK integration

**Tech Stack:**  
- Backend: Python (FastAPI/Quart) or Go (Gin/Echo)  
- Frontend: React + TailwindCSS + D3.js/Recharts  
- Real-time: Websockets (async updates)  
- Auth: JWT, optional 2FA  

---

### 2.2 Layout

```

+------------------------------------------------------------+
| RedForgeC2 Web UI - Dashboard                               |
| Logged in as: Operator1                                    |
+-------------------+----------------------+----------------+

| Agents Map Panel                                               | Alerts Panel       | Scenario Panel|         |
| -------------------------------------------------------------- | ------------------ | ------------- | ------- |
| [Map of lab hosts]                                             | [RED] Sim Alert    | Node Graph    |         |
| Node 001: ALPHA                                                | [YEL] Task Alert   | Gamma->Bravo  |         |
| Node 002: BRAVO                                                | [GRN] Policy Alert | Bravo->Alpha  |         |
| Node 003: GAMMA                                                |                    |               |         |
| +-------------------+---------------------+----------------+   |                    |               |         |
| Command Console: >                                             |                    |               |         |
| +------------------------------------------------------------+ |                    |               |         |
| Tabs: Sessions                                                 | Telemetry          | Tasks         | Reports |
| +------------------------------------------------------------+ |                    |               |         |

```

---

### 2.3 Panels

1. **Agents Map Panel**  
   - Shows all agents on a network map (2D or 3D)  
   - Color-coded by status: RED (alert), GREEN (idle), BLUE (controlled)  

2. **Network Topology Panel**  
   - Graph view of **simulated** relationships and lab scenarios  
   - Live edge updates when scenarios change (simulated)  

3. **Alerts Panel**  
   - Live feed of critical events  
   - Filters for severity, agent, or MITRE ATT&CK tactic  

4. **Command Console**  
   - Web-based tasking (simulation)  
   - Autocomplete task types  
   - Optional JSON param editor for training scenarios  

5. **Tabs**  
   - **Sessions:** All agent sessions and metadata  
   - **Telemetry:** CPU, memory, OS, network, battery  
   - **Tasks:** Create/review simulated tasks and outcomes  
   - **Reports:** Generate forensic/IR reports mapped to MITRE ATT&CK  

---

### 2.4 Key Web UI Features

- **Real-Time Dashboards:** Dynamic charts and sparkline graphs for telemetry  
- **Interactive Network Graph:** D3.js visualization for simulated scenarios  
- **Alerts & Notifications:** Toast pop-ups and persistent feed  
- **Agent Health Overview:** Uptime, last check-in, payload version  
- **Reporting Module:** PDF or JSON export of sessions, telemetry, alerts, and simulated events  
- **Role-Based Access Control:** Operator, Analyst, QA, Viewer  
- **Dark/Light Mode:** For usability during long engagements  

---

### 2.5 UX Considerations

- **Consistency:** TUI and Web UI use same command naming and alert colors  
- **Immediate Feedback:** Every operator action triggers clear logs or alerts  
- **Customization:** Adjustable layouts, color themes, and refresh intervals  
- **Auditability:** All commands and agent events logged for QA, testing, and training  

---

### 3️⃣ UI Mockup Assets

- **Logo:** `docs/assets/redforgec2_logo.png`  
- **Banner:** `docs/assets/redforgec2_banner.png`  
- **TUI Wireframe:** ASCII mockups included above  
- **Web UI Wireframe:** See layout section for panel mockups  

---

### 4️⃣ Next Steps for Dev

- [ ] Implement TUI with `rich` or `tui-rs`  
- [ ] Implement Web UI with React + WebSocket telemetry  
- [ ] Connect TUI/Web UI to unified Agent ↔ Teamserver protocol  
- [ ] Add real-time alert rendering and MITRE ATT&CK tagging  
- [ ] Conduct UX testing with internal lab simulations  

---

**RedForgeC2 UI Mockups:** Designed for high-fidelity C2 operations with dual TUI & Web UI, modern visualization, and professional usability for Red Team exercises.
```

---
