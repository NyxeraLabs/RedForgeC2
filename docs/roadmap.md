# RedForgeC2 – Adversary Simulation Platform
Professional Development Roadmap

Stack

Agent: Rust  
Teamserver: Go  
UI: React + Typescript  
Database: PostgreSQL  
Infrastructure: Docker  
Realtime: WebSockets  

---

# Milestone 1 — Project Foundation

## Issue 1 — Repository Initialization

### Dev
- [x] Create GitHub repository
- [x] Configure monorepo structure
- [x] Setup Go modules
- [x] Setup Rust workspace
- [x] Initialize React UI

### QA
- [x] Configure CI pipeline
- [x] Setup linting tools
- [x] Configure automated testing

### Docs
- [x] Create README
- [x] Create CONTRIBUTING
- [x] Create LICENSE

### Commits
- [x] init repository
- [x] add monorepo structure
- [x] configure go module
- [x] configure rust workspace
- [x] initialize react app
- [x] add docker dev environment
- [x] add CI pipeline
- [x] add linting configuration

---

## Issue 2 — Architecture Definition

### Dev
- [x] Define system architecture
- [x] Define agent protocol
- [x] Define API design
- [x] Define transport abstraction

### QA
- [x] Define testing strategy
- [x] Define integration test plan

### Docs
- [x] architecture.md
- [x] threat-model.md
- [x] protocol-spec.md

### Commits
- [x] add architecture documentation
- [x] add protocol specification
- [x] add threat model

---

# Milestone 2 — Teamserver Core

## Issue 3 — Teamserver Skeleton

### Dev
- [x] Implement Go HTTP server
- [x] Implement configuration loader
- [x] Implement logging system
- [x] Implement API routing

### QA
- [x] API unit tests
- [x] Config parsing tests

### Commits
- [x] implement teamserver bootstrap
- [x] add configuration loader
- [x] add structured logging
- [x] add http router
- [x] add health endpoint
- [x] add graceful shutdown

---

## Issue 4 — Authentication System

### Dev
- [x] Implement operator login
- [x] Implement JWT authentication
- [x] Implement RBAC

### QA
- [x] Auth bypass tests
- [x] Token validation tests

### Docs
- [x] auth.md

### Commits
- [x] implement login endpoint
- [x] add jwt authentication
- [x] implement role system
- [x] add middleware auth validation
- [x] add auth tests

---

# Milestone 3 — Agent Development

## Issue 5 — Agent Bootstrap

### Dev
- [x] Implement Rust agent startup
- [x] Implement environment discovery
- [x] Implement metadata collection

### QA
- [x] OS compatibility tests
- [x] metadata validation tests

### Docs
- [x] agent-architecture.md

### Commits
- [x] create rust agent project
- [x] implement agent bootstrap
- [x] implement metadata collection
- [x] add environment discovery
- [x] add serialization library
- [x] add config loader

---

## Issue 6 — Agent Registration

### Dev
- [x] Implement registration protocol
- [x] Implement agent ID generation
- [x] Implement heartbeat mechanism

### QA
- [ ] registration integration test
- [ ] heartbeat tests

### Commits
- [x] implement registration request
- [x] add agent id generation
- [x] add heartbeat protocol
- [x] add reconnect logic
- [x] add agent registry in server

---

# Milestone 4 — Command Execution

## Issue 7 — Tasking System

### Dev
- [x] Implement task queue
- [x] Implement task dispatcher
- [x] Implement result handler

### QA
- [ ] queue reliability tests
- [ ] concurrency tests

### Commits
- [x] implement task queue
- [x] implement dispatcher
- [x] implement result processor
- [x] add task status tracking
- [x] add task history storage

---

## Issue 8 — Core Agent Commands

### Dev
- [x] shell execution
- [x] file upload
- [x] file download
- [x] directory listing
- [x] process enumeration
- [x] working directory command

### QA
- [ ] command reliability tests
- [ ] file transfer tests

### Commits
- [x] implement shell command
- [x] implement file upload
- [x] implement file download
- [x] implement ls command
- [x] implement process list
- [x] add command validation

---

# Milestone 5 — Communication Layer

## Issue 9 — HTTPS Transport

### Dev
- [ ] implement https transport
- [x] implement jitter timing
- [ ] implement retry logic

### QA
- [ ] latency tests
- [ ] network failure tests

### Commits
- [ ] implement https client
- [ ] add transport encryption
- [x] implement jitter algorithm
- [ ] add retry backoff
- [ ] implement listener creation

---

## Issue 10 — Transport Abstraction

### Dev
- [x] create transport interface
- [ ] implement tcp transport
- [ ] implement dns placeholder

### QA
- [ ] failover tests
- [ ] transport switching tests

### Commits
- [x] add transport interface
- [ ] implement tcp transport
- [ ] add dns transport skeleton
- [ ] implement transport registry

---

# Milestone 6 — Operator UI

## Issue 11 — UI Foundation

### Dev
- [x] dashboard layout
- [x] authentication page
- [x] navigation system

### QA
- [x] UI navigation tests
- [x] authentication flow tests

### Commits
- [x] initialize react ui
- [x] implement login page
- [x] implement dashboard
- [x] add sidebar navigation
- [x] add api client

---

## Issue 12 — Agent Dashboard

### Dev
- [x] agent list view
- [x] agent metadata display
- [x] agent status indicators

### QA
- [ ] real-time update tests

### Commits
- [x] implement agent table
- [x] add agent details panel
- [ ] add real-time websocket updates

---

## Issue 13 — Interactive Console

### Dev
- [ ] agent console terminal
- [ ] command submission
- [ ] result streaming

### QA
- [ ] console reliability tests

### Commits
- [ ] implement console UI
- [ ] add websocket command streaming
- [ ] add command history

---

# Milestone 7 — Advanced Operator Features

## Issue 14 — Listener Manager

### Dev
- [ ] create listener UI
- [ ] delete listener
- [ ] configure listeners

### QA
- [ ] listener creation tests

### Commits
- [ ] implement listener API
- [ ] add listener UI
- [ ] add listener status monitoring

---

## Issue 15 — Loot Manager

### Dev
- [ ] implement file storage
- [ ] implement screenshot storage
- [ ] implement credential storage

### QA
- [ ] large file tests
- [ ] storage reliability tests

### Commits
- [ ] implement loot storage
- [ ] add loot search
- [ ] add loot tagging
- [ ] implement loot download

---

# Milestone 8 — Adversary Simulation

## Issue 16 — Playbook Engine

### Dev
- [ ] yaml playbook format
- [ ] playbook execution engine
- [ ] task sequencing

### QA
- [ ] playbook execution tests

### Commits
- [ ] implement playbook parser
- [ ] implement playbook executor
- [ ] add playbook scheduling

---

## Issue 17 — Detection Telemetry

### Dev
- [ ] technique tracking
- [ ] detection metadata
- [ ] telemetry dashboard

### QA
- [ ] telemetry accuracy tests

### Commits
- [ ] implement telemetry collector
- [ ] add telemetry database
- [ ] implement detection dashboard

---

# Milestone 9 — Security Hardening

### Dev
- [ ] TLS enforcement
- [ ] audit logging
- [ ] API token system

### QA
- [ ] auth bypass tests
- [ ] session security tests

### Commits
- [ ] enforce TLS
- [ ] add audit logging
- [ ] implement api tokens
- [ ] add security tests

---

# Milestone 10 — End-to-End Testing

### Dev
- [x] build docker lab
- [ ] deploy test infrastructure
- [ ] automate agent deployment

### QA
- [ ] multi-agent simulation
- [ ] network failure testing
- [ ] server restart testing

### Commits
- [x] add docker lab
- [ ] implement integration tests
- [ ] add multi-agent simulation

---

# Milestone 11 — Documentation & Release

### Docs
- [ ] getting-started.md
- [ ] operator-guide.md
- [ ] developer-guide.md
- [ ] architecture diagrams
- [ ] lab setup guide

### Release
- [ ] record demo video
- [ ] publish GitHub release
- [ ] write technical blog

### Commits
- [x] add documentation
- [x] add diagrams
- [x] prepare v1 release

---

# Milestone 12 — Agent Capabilities (Priority)

## Issue 18 — File Operations

### Dev
- [ ] implement file upload (agent receives file bytes and writes to disk)
- [ ] implement file download (agent reads file and returns base64 payload)
- [ ] handle large file streaming / chunking

### QA
- [ ] file integrity tests
- [ ] partial transfer recovery tests

### Commits
- [ ] add file operation task handlers
- [ ] add UI support for upload/download tasks

---

## Issue 19 — Keylogging & Persistence

### Dev
- [ ] implement keylogger task (capture keystrokes)
- [ ] persist keystrokes securely
- [ ] query keystrokes from server

### QA
- [ ] keylogging accuracy tests
- [ ] privacy / access control tests

### Commits
- [ ] add keylogger task handler
- [ ] add keylogging viewer in UI

---

## Issue 20 — Pivoting & Lateral Movement

### Dev
- [ ] implement pivot proxy (agent acts as SOCKS/HTTP proxy)
- [ ] implement port forwarding task
- [ ] implement remote command execution via pivot

### QA
- [ ] pivot stability tests
- [ ] port forwarding stress tests

### Commits
- [ ] add pivoting task handlers
- [ ] add UI pivot configuration controls

# Appendix A: Military-Grade Security Hardening & Infrastructure Resilience
### Phase 4: NyxeraLabs Sovereign Infrastructure Protection

## Milestone 19 — Cryptographic Rigor & Protocol Supremacy
*Eradicating legacy cryptography and transitioning to Quantum-Resistant, Authenticated Encryption with Associated Data (AEAD) to ensure absolute forward secrecy and payload integrity.*

- [ ] **19.1 Post-Quantum & AEAD Cryptography**
  - [ ] Deprecate AES-256-CBC in favor of `XChaCha20-Poly1305` or `AES-256-GCM` for all Agent ↔ Teamserver and Client ↔ Teamserver data layers.
  - [ ] Implement Elliptic-Curve Diffie-Hellman (`X25519`) for ephemeral session key exchange to guarantee Perfect Forward Secrecy (PFS).
  - [ ] Integrate optional Kyber (NIST PQC) key encapsulation mechanisms for post-quantum forward secrecy on Tier-1 (Rust) agents targeting highly classified networks.
- [ ] **19.2 Hardware Security Module (HSM) Integration**
  - [ ] Implement PKCS#11 support within the Go Teamserver to offload Root CA and private key storage to an HSM (e.g., YubiHSM 2 or AWS CloudHSM).
  - [ ] Ensure the Teamserver *never* holds the private signing keys in memory; all JWT signing and mTLS certificate issuance must occur within the HSM boundary.
- [ ] **19.3 Cryptographic Task Watermarking & Non-Repudiation**
  - [ ] Require the Operator Client to cryptographically sign every task request (e.g., via Ed25519 operator-specific keys).
  -[ ] Implement signature verification at both the Teamserver (for RBAC authorization) and the Agent (for execution authorization), ensuring a compromised Teamserver cannot issue arbitrary commands to agents.

## Milestone 20 — Zero-Trust Operator Infrastructure
*Enforcing biometric, phishing-resistant authentication and strict network isolation for NyxeraLabs operators accessing the C2.*

- [ ] **20.1 Phishing-Resistant MFA (FIDO2/WebAuthn)**
  -[ ] Eradicate password-only logins. Enforce hardware-backed FIDO2 security keys (YubiKey) for all Operator UI (Web/TUI) access.
  - [ ] Implement short-lived, biometrically gated session tokens using Time-to-Live (TTL) micro-sessions (max 15 minutes of idle time before re-authentication).
- [ ] **20.2 Bastion Hosts & Network Segmentation**
  - [ ] Restrict Operator Client ↔ Teamserver gRPC communication exclusively through a WireGuard or Nebula encrypted overlay network.
  - [ ] Drop all external ingress traffic to the Teamserver management ports. Only authenticated IPsec/WireGuard tunnels from designated jump boxes are permitted.
- [ ] **20.3 Strict Role-Based Access Control (RBAC) & Four-Eyes Principle**
  - [ ] Implement "Four-Eyes" execution authorization for critical tasks (e.g., `Domain Admin` escalation, widespread ransomware simulation). Operator A tasks; Admin B must approve via cryptographic signature.

## Milestone 21 — Ephemeral & Untraceable Teamserver Footprint
*Hardening the Teamserver host OS to resist physical and forensic analysis. The C2 must operate as a "Ghost Server" that self-destructs upon tampering.*

- [ ] **21.1 RAM-Only Execution & Anti-Swapping**
  - [ ] Deploy the Teamserver operating system and binaries strictly via stateless Live CD/PXE boot mechanisms into a `tmpfs` (RAM disk).
  - [ ] Enforce `mlockall()` / `mlock()` in the Go Teamserver to pin all cryptographic keys and sensitive variables in memory, preventing them from being paged to disk/swap.
- [ ] **21.2 Database Field-Level Encryption & LUKS**
  - [ ] Implement AES-GCM Field-Level Encryption (FLE) within PostgreSQL for all captured loot, keystrokes, and task outputs before they are written to the database.
  - [ ] Mandate LUKS2 Full Disk Encryption for any persistent storage volumes, utilizing Tang network-bound disk encryption (NBDE) to prevent decryption if the server is physically seized and isolated.
- [ ] **21.3 Dead Man's Switch & Self-Destruct Sequence**
  - [ ] Implement an eBPF-based tamper detection module tracking unauthorized SSH attempts, physical chassis intrusion, or unexpected network isolation.
  - [ ] Develop a "Scorched Earth" subroutine: upon tamper threshold violation, instantly execute `WipeFile()` on database volumes, overwrite memory with `/dev/urandom`, and trigger an immediate kernel panic (`sysrq-trigger`).

## Milestone 22 — Advanced Redirector & Traffic Laundering
*Isolating the Teamserver from the target network via a distributed, highly expendable, and intelligent redirection tier.*

- [ ] **22.1 JA3/JA4 Fingerprint Filtering & Active Defense**
  - [ ] Deploy Nginx/HAProxy edge redirectors equipped with eBPF/Lua modules to inspect incoming TLS ClientHello packets (JA3/JA4 signatures).
  - [ ] Automatically drop or tarpit connections matching known Blue Team scanners (e.g., Shodan, Censys, Palo Alto Cortex, CrowdStrike) or Python/Go default HTTP libraries.
- [ ] **22.2 Dynamic Traffic Routing & Payload Proxying**
  - [ ] Implement Domain Fronting and CDN pivoting (via Cloudflare, Fastly, CloudFront) to mask the true IP of the edge redirectors.
  - [ ] Configure edge redirectors to perform SNI routing: valid agent traffic (verified via custom headers or TLS SNI) is proxied via reverse SSH tunnels back to the Teamserver. Invalid traffic is seamlessly redirected to a benign corporate webpage (e.g., `https://www.microsoft.com`).
- [ ] **22.3 Fast-Flux Infrastructure Rotation**
  - [ ] Develop a Terraform/Pulumi automation pipeline integrated into the Teamserver to spin up, rotate, and burn redirector VPS instances (DigitalOcean, Linode, AWS) every 4–8 hours automatically to burn threat intel IoCs (Indicators of Compromise).

# Appendix B: Direct Competitor Supremacy & Extensibility
### Phase 5: Advanced Extensibility & Modern Perimeter Domination

## Milestone 23 — The "Anvil" Engine (In-Memory Object Execution)
*To compete with Cobalt Strike's BOF and Havoc's module loading, the Rust Tier-1 Agent must be able to execute unlinked C/C++ object files and .NET assemblies entirely in memory, without spawning child processes.*

- [ ] **23.1 COFF / ELF Object Loader (BOF Compatibility)**
  - [ ] Implement a custom Common Object File Format (COFF) and Executable and Linkable Format (ELF) loader within the Rust agent to map, relocate, and execute C/C++ object files directly in the agent's memory space.
  - [ ] Develop a compatibility layer wrapper to natively support existing Cobalt Strike BOFs and TrustedSec's SA (Situation Awareness) tools without requiring recompilation.
  - [ ] Ensure all object file memory allocations are routed through the Milestone 15 Indirect Syscall engine to bypass user-land API hooks.
- [ ] **23.2 Inline .NET Assembly Execution (Windows)**
  - [ ] Implement a CLR (Common Language Runtime) hosting interface via COM (Component Object Model) to load and execute C# binaries (`.exe`/`.dll`) entirely in memory (e.g., BloodHound, Seatbelt, Rubeus).
  - [ ] Develop an AMSI/ETW-TI bypass that dynamically patches the local CLR instance immediately before assembly invocation, restoring original bytes post-execution to avoid memory scanning alerts.
- [ ] **23.3 Reflective DLL Injection & Shellcode Orchestration**
  - [ ] Implement an sRDI (Shellcode Reflective DLL Injection) module to seamlessly convert arbitrary native DLLs into position-independent shellcode for thread-hijack execution.

## Milestone 24 — Cloud Native & Identity Graph Operations
*Modern breaches rarely rely solely on Active Directory; they target Entra ID (Azure AD), AWS IAM, and Okta. RedForgeC2 must treat Cloud C2 as a first-class citizen.*

- [ ] **24.1 Primary Refresh Token (PRT) Extraction & Forgery**
  - [ ] Develop native Rust modules to interface with the Windows CloudAP plugin and TPM (Trusted Platform Module) to extract or request Entra ID PRTs for seamless Azure lateral movement.
  - [ ] Implement a local proxy within the Agent to tunnel operator browser traffic directly through the victim's authenticated PRT session, bypassing conditional access policies (Device Compliance/IP fencing).
- [ ] **24.2 Cloud Metadata API Pivoting**
  - [ ] Equip Tier-3 (Go) infrastructure agents with automated AWS IMDSv2, Azure IMDS, and GCP metadata extraction modules.
  - [ ] Implement automated STS (Security Token Service) assumption and temporary credential generation for immediate cloud control plane escalation.
- [ ] **24.3 OAuth & Device Code Phishing Workflows**
  - [ ] Integrate a Teamserver module to generate, track, and weaponize Microsoft/Google Device Code authentication flows, pushing the authentication prompts directly to the Operator UI.

## Milestone 25 — macOS / Apple Silicon Supremacy
*Sliver is currently the standard for macOS C2. We must dethrone it by building a Tier-1 implant explicitly designed for Apple Silicon (ARM64) and modern macOS (Sonoma/Sequoia) security frameworks.*

- [ ] **25.1 Mach-O & Swift Native Implant**
  - [ ] Develop a dedicated macOS implant written in Swift/Objective-C to interface natively with Apple's private APIs, bypassing the heavy signature footprint of cross-compiled Go/Rust binaries.
  - [ ] Implement Mach-O memory execution techniques to reflectively load dylibs (Dynamic Libraries) without touching the APFS filesystem.
- [ ] **25.2 Endpoint Security Framework (ESF) Evasion**
  -[ ] Develop memory-safe techniques to unhook or blind ESF sensors (e.g., Jamf Protect, CrowdStrike Falcon on Mac).
  - [ ] Implement TCC (Transparency, Consent, and Control) database manipulation and zero-click bypasses to grant the agent Full Disk Access and Screen Recording permissions without user prompts.
- [ ] **25.3 Keychain & Secure Enclave Subversion**
  - [ ] Implement native API calls (`SecItemCopyMatching`) to dump cleartext passwords, cryptographic keys, and Safari cookies directly from the macOS Keychain.

## Milestone 26 — The "ForgeScript" Ecosystem (API & Armory)
*A C2 lives and dies by its community and extensibility. We must provide a decentralized package manager and a robust scripting API for operators.*

- [ ] **26.1 Operator Scripting Engine (Python/Lua via gRPC)**
  - [ ] Expose a 100% coverage gRPC API for the Teamserver, allowing operators to write headless Python or Lua scripts to automate tasks (e.g., "If new agent checks in as SYSTEM, auto-execute BOF Seatbelt and dump LSASS").
  - [ ] Implement Webhook and Slack/Discord/Mattermost native integrations for automated alerting on high-value check-ins or lateral movement successes.
- [ ] **26.2 Decentralized "Forge Armory" Package Manager**
  - [ ] Build an in-UI module repository (like Sliver's Armory) that securely pulls vetted Red Team tools, BOFs, and custom lateral movement modules from a NyxeraLabs-signed GitHub repository.
  - [ ] Implement automatic server-side compilation of these tools via the Milestone 14 "Forge" engine, ensuring payloads are obfuscated specifically for the current engagement before deployment.

## Milestone 27 — Sub-System & Kernel Domination
*To truly outmaneuver modern XDR platforms, RedForgeC2 must push operations below Ring 3 (Userland).*

- [ ] **27.1 BYOVD (Bring Your Own Vulnerable Driver) Automation**
  - [ ] Integrate a database of signed, vulnerable Windows drivers (e.g., `RTCore64.sys`, `gdrv.sys`).
  - [ ] Develop an automated agent task to drop a vulnerable driver, exploit it to gain Ring 0 execution, and forcibly remove EDR kernel callbacks (`PspCreateProcessNotifyRoutine`, `ObRegisterCallbacks`) to effectively blind the EDR system-wide.
- [ ] **27.2 Linux eBPF Rootkit Capabilities**
  - [ ] For Linux Tier-1 agents, implement an eBPF (Extended Berkeley Packet Filter) module to intercept syscalls at the kernel level.
  - [ ] Use eBPF to hide the agent's PID from `ps`/`top`, hide network sockets from `netstat`/`ss`, and intercept/modify SSH credentials in memory during authentication without patching `sshd`.
