# Security Policy for RedForgeC2

## ⚠️ Important Disclaimer

**RedForgeC2** is a research and development tool intended **strictly for authorized Red Team, Purple Team, and incident response exercises**. Unauthorized use, deployment, or testing on live systems without explicit consent is illegal and strictly prohibited. The authors disclaim all liability for misuse, damages, or legal consequences.

---

## 🛡 Rules of Engagement (RoE)

1. **Authorized Environments Only**
   - Execute RedForgeC2 in isolated lab environments (VMs, containers, or physically segmented networks).
   - Ensure no production systems are connected to test networks.

2. **Explicit Consent Required**
   - Obtain written authorization from system owners for any testing or simulation.
   - Keep a signed engagement agreement on record.

3. **Propagation Limits**
   - Multi-node testing is permitted only in sandboxed lab setups.
   - Lateral movement, persistence, or exfiltration modules must be restricted to lab hosts.
   - Avoid uncontrolled propagation to prevent unintended Denial of Service (DoS).

4. **Kill-Switch & Cleanup**
   - Use the built-in **Kill-Switch** to terminate all agents immediately.
   - Ensure automated cleanup scripts remove binaries, logs, and temporary data after each test.
   - Always verify the lab state post-simulation.

5. **Telemetry and Logging**
   - RedForgeC2 telemetry is for simulation visibility only.
   - Avoid sending logs or artifacts to production networks.

---

## 🕵️‍♂️ Evasion and Sensitive Techniques

RedForgeC2 implements advanced evasion techniques for realism in simulations:

- Anti-forensic memory clearing
- In-memory execution and payload injection
- Dynamic protocol polymorphism for network transport
- Credential simulation for lateral movement exercises

> **Note:** These features are for **training, research, and testing only**. Misuse may violate local and international law.

---

## 🔒 Secure Development Practices

- All source code undergoes **code review** and **static analysis** for unintentional vulnerabilities.
- Secrets (keys, tokens, credentials) are **never hardcoded** in source.
- Agent communication uses **TLS 1.2+** with certificate verification.
- Rust agents follow **memory-safe patterns** and implement zeroize for sensitive buffers.

---

## 📢 Reporting Security Issues

If you discover a vulnerability, misconfiguration, or unsafe behavior:

1. Report via GitHub Issues with a clear subject line:
   - `[SECURITY] <short description>`
2. Include reproduction steps and environment details.
3. Avoid publicly disclosing exploits; maintain responsible disclosure.
4. Maintainers will respond within **72 hours** for acknowledgment and follow-up.

---

## 📝 Acknowledgments

This security policy is inspired by professional offensive security frameworks including **Cobalt Strike**, **Sliver**, **Havoc**, and **Mythic**, adapted for safe research and educational use.

---

**Always operate RedForgeC2 within the legal and ethical boundaries of your jurisdiction.**