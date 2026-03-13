# Agents

This folder contains **simulation-only** agent implementations used for local lab exercises.

- `rust_agent/`: Rust agent skeleton that **only** exchanges mock telemetry and mock tasks with a local teamserver.
- `mock_agents/`: Non-native/mock agents used for UI and test fixtures.

Safety:
- No real exploitation, persistence, lateral movement, or host command execution.
- Networking must be limited to `localhost` or an isolated lab network (Docker/VM).

