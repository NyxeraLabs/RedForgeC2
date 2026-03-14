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

## TLS Enforcement (Optional but recommended)

RedForgeC2 supports mTLS and can be configured to refuse non-TLS connections.

- Set `REDFORGE_TLS_CERT_FILE` and `REDFORGE_TLS_KEY_FILE` to point to your certificate and key.
- Enable strict TLS enforcement by setting `REDFORGE_REQUIRE_TLS=1`.

If `REDFORGE_REQUIRE_TLS` is enabled, the teamserver will fail to start unless both the cert and key are configured.

The teamserver emits audit-style logs for every API request (user, IP, path, status, duration) via standard output for easy log collection.

## API Tokens

The teamserver supports long-lived API tokens for automation and non-interactive access. Tokens can be created by an administrator via the `/api/admin/api-tokens` endpoint and are only shown once when created. Store them securely.

Example request (as an admin):

```sh
curl -X POST \
  -H "Authorization: Bearer <admin-jwt>" \
  -H "Content-Type: application/json" \
  -d '{"username":"operator","description":"CI runner","expires_minutes":1440}' \
  http://localhost:9080/api/admin/api-tokens
```

## Running the agent locally

```sh
cd agent
cargo run
```

## Next Steps

- Log in via the UI (default credentials are set via environment variables)
- Register an agent via the agent binary
- Send tasking via the operator API
