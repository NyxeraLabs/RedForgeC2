# RedForgeC2 Threat Model

This threat model describes the assumed adversary capabilities, system trust boundaries, and mitigations for RedForgeC2.

## Assumptions

- **Environment:** RedForgeC2 is executed in controlled lab environments (virtual machines, containers, isolated networks).
- **Adversary:** The operator is trusted; threat modeling focuses on unintentionally enabling escalation or persistence beyond the lab.
- **Scope:** The system is not intended for use against production systems without explicit authorization.

## Assets

- **Agent Identity:** UUID and secret token used to authenticate to the teamserver.
- **Telemetry & Task Data:** Commands, output, and harvested artifacts.
- **Operator Credentials:** JWT or session tokens used in the UI.
- **Transport Keys:** Encryption keys and TLS certificates.

## Threats

### T1 — Unauthorized Access
- **Description:** An attacker gains access to the teamserver API or UI.
- **Mitigation:** Enforce authentication, MFA (future), and role-based access controls.

### T2 — Agent Impersonation
- **Description:** A malicious actor emulates an agent and sends fabricated telemetry or task results.
- **Mitigation:** Agents authenticate with a unique token and verify message integrity (HMAC + sequence tracking).

### T3 — Data Exfiltration Outside Lab
- **Description:** Telemetry or artifacts get transmitted to external networks.
- **Mitigation:** Use network isolation, allow-listing, and enforce lab-only configuration.

### T4 — Persistence Abuse
- **Description:** Agent persistence mechanisms (registry, systemd, cron) persist beyond intended scope.
- **Mitigation:** Provide strong kill-switch + cleanup routines; document safe lab practices.

## Trust Boundaries

1. **Agent ↔ Teamserver**
   - Secured by transport encryption (TLS + AES/HMAC).
   - Authentication at connection setup.

2. **Operator UI ↔ Teamserver**
   - Secured by HTTPS and JWT authentication.
   - UI never stores secrets in cleartext.

3. **Storage Layer (DB/files)**
   - Sensitive artifacts are stored encrypted (future phase).

## Security Controls (Current & Planned)

- **Secure Defaults:** TLS enforced, strict CORS, minimal exposure.
- **Logging & Auditing:** Record actions for traceability.
- **Secrets Management:** No hardcoded secrets; use environment variables or vaults.
- **Hardening:** Minimal container images, regular dependency updates.

---

> Note: This threat model is continuously updated as the system evolves. Always validate against the latest implementation.
