# Rust Agent (Simulation Skeleton)

This is a **safe educational skeleton**.

Planned behavior (lab-safe):
- Registers to a local teamserver (`localhost` / Docker network)
- Sends **mock telemetry** (JSON)
- Receives **mock tasks** (JSON) and records them as simulated actions

Non-goals:
- No real remote shell, file operations, persistence, or evasion.

## Run

```bash
cd agents/rust_agent
cargo test

# Terminal 1: run the Go teamserver
# (from repo root)
#   cd server/go_backend && go run ./cmd/teamserver -addr 127.0.0.1:8080

# Terminal 2: run the agent (one loop)
REDFORGE_SERVER_URL=http://127.0.0.1:8080 REDFORGE_ONCE=1 cargo run
```
