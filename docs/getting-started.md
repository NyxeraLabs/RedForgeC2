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
docker-compose up --build
```

2. Open the UI at `http://localhost:5173`.

## Running the teamserver locally

```sh
cd teamserver
go run ./cmd/teamserver
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
