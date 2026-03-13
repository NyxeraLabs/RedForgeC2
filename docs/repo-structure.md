### 1️⃣ `docs/repo-structure.md`

```markdown
# RedForgeC2 — GitHub Repository Structure

This document outlines the recommended folder and file structure for a professional C2 framework.

## Root Layout

```

RedForgeC2/
├─ agent/               # Rust implant code
│  ├─ src/
│  ├─ Cargo.toml
│  ├─ build.rs
├─ teamserver/          # Go backend code
│  ├─ cmd/
│  ├─ internal/
│  ├─ go.mod
├─ ui/                  # React/TypeScript operator console
│  ├─ src/
│  ├─ public/
│  ├─ package.json
├─ docs/                # Documentation
├─ scripts/             # Build, deploy, and test scripts
├─ tests/               # E2E tests & QA scripts
├─ docker/              # Container definitions
├─ .github/             # CI/CD and workflows
│  ├─ workflows/
├─ .gitignore
├─ README.md
├─ LICENSE
└─ Makefile

```

### Recommended Subfolders

- **agent/** → Rust source code, agent bootstrap, command execution, telemetry, transport  
- **teamserver/** → Go backend: task queue, API, database, configs  
- **ui/** → React/TS frontend: dashboard, agent list, tasking, telemetry  
- **scripts/** → Build, deploy, lab simulation scripts  
- **tests/** → Unit, integration, and E2E tests  
- **docker/** → Multi-node lab Dockerfiles  

✅ Best Practices

- Keep `agent` separate to prevent accidental source leaks  
- CI/CD pipelines must include build, test, lint  
- Docs must include architecture, protocol, lab setup, UI guide