# Authentication (Operator) — Simulation-Only

RedForgeC2 implements **operator authentication** for the teamserver API used by the UIs.

Scope:
- Protects operator endpoints like agent listing and task enqueueing.
- Agent endpoints (register/telemetry/task poll/result) remain lab-scoped and should run on `127.0.0.1` or isolated Docker/VM networks.

## Enable auth (dev/admin)

Set:
- `REDFORGE_AUTH_ENABLED=true`
- `REDFORGE_JWT_SECRET` (>= 32 bytes)
- `REDFORGE_ADMIN_USER`
- `REDFORGE_ADMIN_PASS`

Example:

```bash
export REDFORGE_AUTH_ENABLED=true
export REDFORGE_JWT_SECRET="01234567890123456789012345678901"
export REDFORGE_ADMIN_USER=admin
export REDFORGE_ADMIN_PASS='change-me'

cd server/go_backend
go run ./cmd/teamserver -addr 127.0.0.1:8080
```

## Login

`POST /api/v1/auth/login`

Body:
```json
{"username":"admin","password":"change-me"}
```

Response:
```json
{"token":"<jwt>","expires_at":"...","role":"admin"}
```

Use the token:
`Authorization: Bearer <jwt>`

## RBAC (current)

Roles:
- `admin` (highest)
- `operator`
- `analyst`
- `viewer` (lowest)

Protected endpoints (current):
- `GET /api/v1/agents`: requires `viewer+`
- `POST /api/v1/tasks/enqueue`: requires `operator+`

