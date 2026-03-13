### 1️⃣ `docs/repo-structure.md`

```markdown
# RedForgeC2 — Repository Structure (Simulation-Only)

This document outlines the recommended folder and file structure for a professional C2 framework.

## Root Layout

```

RedForgeC2/
├─ agents/              # Simulation agents (Rust + mocks)
│  ├─ rust_agent/
│  └─ mock_agents/
├─ server/              # Teamserver backends (Python/Go)
│  └─ go_backend/
├─ ui/                  # Operator consoles
│  ├─ tui/
│  └─ web/
├─ docs/                # Documentation & UI mockups
├─ scripts/             # Build, deploy, and test scripts
├─ tests/               # Unit, integration, and E2E tests
├─ ci_cd/               # CI/CD notes (workflows in .github/)
├─ docker/              # Container definitions
├─ .github/             # CI/CD workflows
│  ├─ workflows/
├─ .gitignore
├─ README.md
├─ SECURITY.md
├─ LICENSE
└─ Makefile

```

### Recommended Subfolders

- **agents/** → simulation-only agents; mock telemetry + mock tasking  
- **server/** → local-only backend that models task queue + telemetry ingest  
- **ui/** → operator consoles for training (TUI/Web)  
- **scripts/** → Build, deploy, lab simulation scripts  
- **tests/** → Unit, integration, and E2E tests (simulation)  
- **docker/** → Multi-node lab Dockerfiles  

✅ Best Practices

- Bind services to `127.0.0.1` by default  
- CI/CD pipelines must include build, test, lint  
- Docs must include architecture, protocol, lab setup, UI guide
