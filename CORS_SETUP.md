# CORS Configuration for RedForgeC2

## Issue: NetworkError when attempting to fetch resource during login

If you're getting a "NetworkError" or "TypeError: NetworkError when attempting to fetch resource" during login, it's likely a CORS (Cross-Origin Resource Sharing) issue.

### Root Cause

By default, the teamserver only allows HTTPS requests from:
- `https://localhost:5174`
- `https://localhost:5173`

If your development UI is running on:
- `http://localhost:3000` (different port or HTTP instead of HTTPS)
- `http://127.0.0.1:5173`
- Any other origin

The browser will block the request due to CORS policy.

### Solution

Set the `REDFORGE_CORS_ORIGINS` environment variable before starting the teamserver:

#### For Local Development (HTTP)

```bash
export REDFORGE_CORS_ORIGINS="http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173"
# Then start your teamserver
cargo run -p redforge-teamserver
```

#### For Production (HTTPS)

```bash
export REDFORGE_CORS_ORIGINS="https://yourui.domain.com,https://another.domain.com"
# Then start your teamserver
```

#### Allow Any Origin (Development Only - NOT SECURE)

```bash
export REDFORGE_CORS_ORIGINS="*"
# This allows requests from any origin
```

### Docker Compose Setup

If using Docker, set the environment variable in your `docker-compose.yml`:

```yaml
teamserver:
  environment:
    - REDFORGE_CORS_ORIGINS=http://localhost:5173,http://localhost:3000
```

### Vite Dev Server

When running the UI with Vite:

```bash
cd ui
npm run dev
# Usually starts on http://localhost:5173
```

Add that origin to your `REDFORGE_CORS_ORIGINS` before starting the teamserver.

### API Base URL Configuration

You can also specify a custom teamserver URL in the login page:

1. If you've set `REDFORGE_CORS_ORIGINS` correctly, the login page should auto-detect the right URL
2. You can manually enter the teamserver URL in the "Teamserver URL" field on the login page
3. Or set `VITE_TEAMSERVER_URL` environment variable when building the UI:

```bash
export VITE_TEAMSERVER_URL="https://192.168.1.100:9080"
npm run build
```

### Debugging

If you still get a NetworkError:

1. Open browser DevTools (F12)
2. Go to Network tab
3. Try to login
4. Look for the `/api/login` request
5. Check the response headers for `Access-Control-Allow-Origin`
6. If it's missing, the origin is not in the allowed list

### Related Environment Variables

- `REDFORGE_CORS_ORIGINS` - Comma-separated list of allowed origins (default: `https://localhost:5173,https://localhost:5174`)
- `REDFORGE_MAX_BODY_BYTES` - Maximum request body size in bytes (default: 25 MiB)
- `REDFORGE_LOGIN_RPM` - Login rate limit in requests per minute (default: 20)
