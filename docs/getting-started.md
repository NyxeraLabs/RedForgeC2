# Getting Started

This document provides a quick overview for getting RedForgeC2 running locally.

## Prerequisites

- Docker & Docker Compose (for full stack lab)
- Rust toolchain (for building the agent)
- Go toolchain (for building the teamserver)
- Node.js + npm/yarn (for building the UI)

## Running with Docker

1. From the repo root:

```sh
cp .env.example .env
# edit .env and set: REDFORGE_DB_PASS, REDFORGE_JWT_SECRET, REDFORGE_ADMIN_PASS
make up
```

2. Open the UI at `http://localhost:5174`.
   - Demo-only visual mock: `http://localhost:5174/demo`
   - The login page shows a small `UI:<sha>@<epoch>` build stamp so you can confirm rebuilds are live.

3. The teamserver API is available at `http://localhost:9080`.

### CORS (UI access)

By default the teamserver only allows browser requests from common local UI origins:
- `http://localhost:5174`
- `http://localhost:5173`

To override, set `REDFORGE_CORS_ORIGINS` (comma-separated) in `.env`.

## Running the teamserver locally

```sh
cd teamserver
go run ./cmd/teamserver
```

## Running the agent locally
See `agent/README.md`.

## Next Steps

- Log in via the UI (default credentials are set via environment variables)
- Register an agent via the agent binary
- Send tasking via the operator API
