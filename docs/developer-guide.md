# Developer Guide

This guide covers the basic architecture and how to build components.

## Repository Layout

- `agent/`: Rust-based agent that runs on target hosts.
- `teamserver/`: Go-based C2 server and API.
- `ui/`: React-based operator UI.

## Building

### Teamserver

```sh
cd teamserver
go build ./cmd/teamserver
```

### Agent

```sh
cd agent
cargo build --release
```

### UI

```sh
cd ui
npm install
npm run dev
```

## Adding Features

- For teamserver changes, work in `teamserver/internal/*` and add routes in `teamserver/internal/server`.
- For agent changes, update `agent/src/bootstrap.rs` for task execution or `agent/src/protocol.rs` for message formats.

## Configuration

Environment variables are used for configuration. See `teamserver/internal/config/config.go` for the full list.
