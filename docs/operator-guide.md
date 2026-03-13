# Operator Guide

This guide describes how operators interact with the RedForgeC2 teamserver.

## Authentication

Operators authenticate using username/password and receive a JWT token.

- Default username/password are controlled via environment variables in the teamserver.
- Use the `/api/login` endpoint to obtain a token.

Example request:

```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"redforge"}'
```

Use the returned token for subsequent requests:

```sh
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/operator/agents
```

## Tasking

Operators can enqueue tasks for agents using the `/api/operator/task` endpoint.

Example:

```sh
curl -X POST http://localhost:8080/api/operator/task \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"agent_id":"<agent-id>","command":"ls","args":["/tmp"],"timeout_seconds":60}'
```

## Viewing Task Results

Use `/api/operator/results` to retrieve stored task results for an agent:

```sh
curl -H "Authorization: Bearer <token>" "http://localhost:8080/api/operator/results?agent_id=<agent-id>"
```
