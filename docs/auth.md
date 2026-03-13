# Authentication

RedForgeC2 uses a simple JWT-based authentication scheme for operator access and a shared token approach for agent authentication.

## Operator Authentication (JWT)

### Login Flow
1. Operator POSTs credentials to `/api/login`:
   - `username`
   - `password`

2. Server validates credentials against environment variables:
   - `REDFORGE_ADMIN_USER`
   - `REDFORGE_ADMIN_PASS`

3. On success, the server returns a signed JWT:
   - `alg: HS256`
   - `exp` based on the `REDFORGE_TOKEN_EXPIRY_MIN` configuration
   - includes `username` and `role` claims

### Usage
- All operator API endpoints under `/api/operator/*` require an `Authorization: Bearer <token>` header.
- The server validates the token signature and checks the required role (`admin`).

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
