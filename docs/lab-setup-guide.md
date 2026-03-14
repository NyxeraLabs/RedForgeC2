# Lab Setup Guide

This guide describes setting up a local lab environment for RedForgeC2.

## Prerequisites

- Docker & Docker Compose
- A machine capable of running multiple containers

## Running the Lab

From the repository root:

```sh
cp .env.example .env
# edit .env and set: REDFORGE_DB_PASS, REDFORGE_JWT_SECRET, REDFORGE_ADMIN_PASS
make up
```

This starts:

- `postgres`: database
- `teamserver`: C2 server on port 9080
- `ui`: operator UI on port 5174

## Custom Configuration

Environment variables are loaded from a repo-root `.env` file.

Hardening-related environment variables (optional):
- `REDFORGE_CORS_ORIGINS` (comma-separated allowlist)
- `REDFORGE_LOGIN_RPM` (per-IP login requests/minute)
- `REDFORGE_MAX_BODY_BYTES` (request body limit for non-GET endpoints)

## Resetting the Lab

To rebuild and reset state:

```sh
make db-reset
make up
```
