# RedForgeC2 Threat Model (Simulation-Only)

This threat model assumes RedForgeC2 is operated **only** in an isolated lab environment (VMs or Docker networks).

## Assets

- Event log integrity (audit trail)
- Operator credentials/tokens (when auth is implemented)
- Agent identity (agent_id) and registration data
- Telemetry and task history (training data)

## Trust Boundaries

- Operator ↔ Teamserver (UI clients to HTTP API)
- Agent ↔ Teamserver (agent simulator to HTTP API)
- Local filesystem (event logs, config files)

## Threats and Mitigations

### 1) Unauthorized API access

Threat:
- An attacker on the same lab network calls tasking endpoints or scrapes telemetry.

Mitigations:
- Bind to `127.0.0.1` by default
- Add authentication/authorization (Issue 4)
- Consider allowlists for lab subnets when Docker networking is used

### 2) Event log tampering

Threat:
- Event log is modified to hide activity or corrupt training results.

Mitigations:
- Append-only `jsonl` with restrictive permissions (`0600`)
- Future: log signing / hash-chaining for tamper evidence

### 3) Input validation and JSON poisoning

Threat:
- Crafted JSON causes crashes or unexpected behavior.

Mitigations:
- Strict JSON decoding (disallow unknown fields)
- Size limits and timeouts (future)

### 4) Credential leakage

Threat:
- Tokens/certs committed to repo or leaked in logs.

Mitigations:
- `.gitignore` for secrets
- Document env-based configuration
- Keep logs structured; avoid dumping secrets

## Safety Constraints

To avoid dual-use misuse, RedForgeC2 (this version) must not implement:
- Remote command execution
- Arbitrary file transfer
- Persistence/evasion/exploitation capabilities

