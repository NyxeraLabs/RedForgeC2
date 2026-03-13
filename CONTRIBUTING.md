# Contributing to RedForgeC2 (Simulation-Only)

RedForgeC2 is a **lab-safe simulation project**. Contributions must not add real malware/C2 capabilities.

## Rules

- No remote command execution or arbitrary file transfer features.
- No persistence, evasion, lateral movement, or exploitation modules.
- Default networking must bind to `127.0.0.1` / `localhost` unless explicitly documented for isolated Docker/VM labs.
- Keep changes small and well-tested; update docs and `docs/roadmap.md` checkboxes.

## Development workflow (manual Gitflow)

```bash
git checkout -b feat/<feature>
git add .
git commit -S -m "feat: <description>"
git push origin feat/<feature>
git checkout dev
git merge feat/<feature>
git push origin dev
```

## Testing

- Go teamserver: `cd server/go_backend && go test ./...`
- Rust agent: `cd agents/rust_agent && cargo test`
- E2E smoke: `tests/e2e/agent_server_smoke.sh`

