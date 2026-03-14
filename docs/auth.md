# Authentication

RedForgeC2 uses a simple JWT-based authentication scheme for operator access and a shared token approach for agent authentication.

## Operator Authentication (JWT)

### Login Flow
1. Operator POSTs credentials to `/api/login`:
   - `username`
   - `password`

2. Server validates credentials against the `users` table (Postgres).
   - On startup, the teamserver ensures an admin user exists from:
     - `REDFORGE_ADMIN_USER` (default: `admin`)
    - `REDFORGE_ADMIN_PASS` (no default; required)
   - Passwords are stored using bcrypt hashes (never plaintext).

3. On success, the server returns a signed JWT:
   - `alg: HS256`
   - `exp` based on the `REDFORGE_TOKEN_EXPIRY_MIN` configuration
   - includes `username` and `role` claims

### Usage
- All operator API endpoints under `/api/operator/*` require an `Authorization: Bearer <token>` header.
- The server validates the token signature and checks role permissions (Admin/Operator/Observer).

## User Management

- Admin-only endpoint: `GET/POST /api/admin/users`
  - Create users with roles: `admin`, `operator`, `observer`
- Self-service endpoints:
  - `GET /api/me`
  - `POST /api/me/profile`
  - `POST /api/me/password`

## Agent Authentication (Token)

- Agents register themselves at `/api/register`.
- The server issues a per-agent token that is stored locally on the agent.
- Agents include this token in every heartbeat and task result payload.

### Validation
- The teamserver verifies the token against the in-memory agent registry.
- Invalid tokens result in HTTP 401 responses.

## Configuration

- `REDFORGE_JWT_SECRET`: secret used for signing and validating tokens.
- `REDFORGE_TOKEN_EXPIRY_MIN`: JWT expiration in minutes (default: 60).
- `REDFORGE_ADMIN_RESET`: if set to `1`, updates the admin password on startup (dev only).

## Hardening (defensive)

- `REDFORGE_CORS_ORIGINS`: comma-separated allowlist for browser Origins (default allows local UI dev ports).
- `REDFORGE_LOGIN_RPM`: per-IP rate limit for `/api/login` (default: 20/min).
- `REDFORGE_MAX_BODY_BYTES`: max request body size enforced for non-GET endpoints (default: 25 MiB).
