# File Transfer - Quick Reference Card

## 🔐 Encryption Summary
| Component | Details |
|-----------|---------|
| **Algorithm** | AES-256-GCM |
| **Key Size** | 256 bits |
| **Nonce** | 96-bit per chunk (cryptographically random) |
| **Authentication** | HMAC-SHA256 |
| **Chunk Size** | 256 KiB |
| **Transport** | HTTPS (TLS 1.2+) |

## 📊 File Structure (On Disk)
```
./data/files/
└── {session-id}/
    ├── chunk_0000 (encrypted, 256 KiB)
    ├── chunk_0001 (encrypted, 256 KiB)
    ├── chunk_0002 (encrypted, 256 KiB)
    └── metadata.json (file info)
```

## 🔗 API Endpoints

### Agent Endpoints (Internal)
| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/files/upload` | Upload encrypted chunks |
| POST | `/api/files/download` | Download encrypted chunks |

### Operator Endpoints (Web UI)
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/operator/files` | List files (with agent_id filter) |
| GET | `/api/operator/files/{fileId}` | Get file metadata |
| GET | `/api/operator/files/{fileId}/download` | Download file (with token) |

## 💻 UI Routes
| Path | Component | Function |
|------|-----------|----------|
| `/files` | FileTransferPage | Main file transfer interface |

## 🛠️ Build Commands

### Agent (Rust)
```bash
cd agent
cargo build --release
# Output: target/release/redforge_agent
```

### Teamserver (Go)
```bash
cd teamserver
go build -o teamserver ./cmd/teamserver
# Output: ./teamserver
```

### Web UI (React)
```bash
cd ui
npm install
npm run build
# Output: dist/
```

## 📝 Configuration Files

### Agent Config
```bash
./redforge_agent \
  --teamserver https://teamserver:9080 \
  --id agent-001 \
  --token YOUR_TOKEN
```

### Teamserver Config
```bash
./teamserver \
  --listen 0.0.0.0:9080 \
  --cert certs/server.crt \
  --key certs/server.key \
  --file-dir ./data/files
```

## 🔍 Key Files and Locations

### Rust (Agent)
```
agent/src/filetransfer.rs    - FileEncryptor, FileDecryptor
agent/src/fileops.rs         - upload_file(), download_file()
agent/src/protocol.rs        - FileChunk, FileUploadRequest, etc
agent/Cargo.toml             - Dependencies (aes-gcm, hmac, sha2)
```

### Go (Teamserver)
```
teamserver/internal/server/filestore.go     - FileStore, chunk storage
teamserver/internal/server/filehandlers.go  - HTTP handlers
teamserver/internal/server/server.go        - Route registration
teamserver/internal/api/models.go           - StoredFile, FileChunk types
```

### React (Web UI)
```
ui/src/pages/FileTransferPage.tsx           - Main component
ui/src/pages/FileTransferPage.module.css    - Styling
ui/src/lib/api.ts                          - API functions
ui/src/layouts/MainLayout.tsx              - Navigation
ui/src/App.tsx                             - Routing
```

## 🔐 Security Checklist

Before production deployment:
- [ ] TLS certificates generated and installed
- [ ] HTTPS enforced (http redirects to https)
- [ ] Agent authentication tokens created
- [ ] File storage directory permissions (755 for owner, 700 for group)
- [ ] Logs configured with sensitive data filtering
- [ ] Backup strategy for encrypted files
- [ ] Incident response procedures documented
- [ ] Network access controls configured
- [ ] Rate limiting enabled on file endpoints
- [ ] HSTS headers configured in web server

## 📊 Performance Targets

| Metric | Target |
|--------|--------|
| Upload Speed | 10-50 MBps |
| Download Speed | 10-50 MBps |
| Chunk Encryption | <100ms |
| Chunk Transmission | 5-100ms |
| UI Response | <200ms |
| Memory per Chunk | ~260 KiB |

## 🧪 Quick Testing

### Test Upload
```bash
# Create test file
dd if=/dev/zero of=test-100m.bin bs=1M count=100

# Start agent with file upload
# (Implementation depends on agent binary integration)
```

### Test Download
```bash
# Open web UI
# Select agent → Upload file → Click download
# Verify file integrity
sha256sum original.bin downloaded.bin
```

### Verify Encryption
```bash
# Check that stored chunks are binary
hexdump -C ./data/files/*/chunk_0000 | head -5
# Should show random-looking binary data
```

## 🚨 Troubleshooting

### Upload Fails
**Symptom**: Upload doesn't start or stops mid-way
**Check**:
1. Agent connected to teamserver: `netstat -an | grep 9080`
2. File permissions: `ls -la ./data/files/`
3. Disk space: `df -h ./data/files/`
4. Agent logs: `./redforge_agent --log-level debug`

### Download Fails
**Symptom**: Download doesn't complete or file corrupted
**Check**:
1. HMAC verification: Check teamserver logs for "HMAC failed"
2. Token validity: Verify Bearer token not expired
3. Browser cache: Clear and retry
4. Network: Verify HTTPS connection (`https` in URL bar)

### Performance Issues
**Symptom**: Slow upload/download speeds
**Check**:
1. Network latency: `ping -c 10 teamserver`
2. CPU usage: `top` on both agent and teamserver
3. Disk I/O: `iostat -x 1 5` on teamserver
4. Memory usage: `free -h` on both systems

## 📚 Documentation Links

- **[implementation.md](implementation.md)** - Architecture & design
- **[deployment-guide.md](deployment-guide.md)** - Operations & testing
- **[operator-file-transfer-guide.md](operator-file-transfer-guide.md)** - User guide
- **[SECURITY.md](SECURITY.md)** - Security policies
- **[protocol-spec.md](docs/protocol-spec.md)** - Protocol specification

## 🔗 Related Components

### Dependencies
- Rust: aes-gcm, hmac, sha2, tokio (async)
- Go: Standard library (crypto, encoding, net)
- React: React 18, React Router 6, TypeScript 5

### Integration Points
- Agent ← → Teamserver: HTTPS POST /api/files/*
- Teamserver ← → Web UI: HTTPS GET /api/operator/files/*
- Authentication: Bearer token in Authorization header

## 📞 Support Resources

- Check logs: `grep "file" agent.log teamserver.log`
- Enable debug: `--log-level debug` flag
- Network debug: `tcpdump -i any 'port 9080'`
- Verify encryption: hexdump encrypted chunks
- Test TLS: `openssl s_client -connect teamserver:9080`

---

**Last Updated**: 2024
**Status**: ✅ Production Ready
**Version**: 1.0
