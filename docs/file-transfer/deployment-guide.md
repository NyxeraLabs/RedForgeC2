# File Transfer - Deployment & Testing Guide

## Pre-Deployment Checklist

### Code Compilation
- [x] Agent (Rust): `cargo check` - No errors, 6 warnings (unused code for future integration)
- [x] Teamserver (Go): `go build ./cmd/teamserver` - No errors
- [x] Web UI: `npm run build` - 339 modules transformed, production ready

### Dependency Verification
**Agent (Cargo.toml):**
```toml
aes-gcm = "0.10"    # AES-256-GCM encryption
hmac = "0.12"       # HMAC authentication
sha2 = "0.10"       # SHA-256 hashing
tempfile = "3.8"    # Temporary file handling
```

**Teamserver (go.mod):**
- Standard library only (net, crypto, encoding, etc.)
- No external dependencies required

**Web UI (package.json):**
- React 18.3
- TypeScript 5.2
- React Router 6
- Vite 5.4

## Deployment Steps

### 1. Build Agent Binary
```bash
cd /home/xoce/Workspace/RedForgeC2/agent
cargo build --release
# Output: target/release/redforge_agent
```

### 2. Build Teamserver Binary
```bash
cd /home/xoce/Workspace/RedForgeC2/teamserver
go build -o teamserver ./cmd/teamserver
# Output: ./teamserver
```

### 3. Build Web UI
```bash
cd /home/xoce/Workspace/RedForgeC2/ui
npm install  # If not already done
npm run build
# Output: dist/ directory
```

### 4. Set Up File Storage Directory
```bash
mkdir -p ./data/files
chmod 755 ./data/files
```

### 5. Deploy Teamserver
```bash
./teamserver \
  --listen 0.0.0.0:9080 \
  --cert certs/server.crt \
  --key certs/server.key
```

### 6. Deploy Web UI
```bash
# Option A: Serve with built-in server
cd dist
python3 -m http.server 8080

# Option B: Use production web server (nginx, Apache)
# Copy dist/* to web root
```

### 7. Configure Agent
```bash
./redforge_agent \
  --teamserver https://your-teamserver:9080 \
  --id my-agent-01 \
  --token <auth-token>
```

## Testing Scenarios

### Test 1: Basic File Upload (Small File)
**Setup:**
1. Start teamserver
2. Connect agent
3. Open Web UI
4. Select agent from dropdown

**Test Steps:**
1. Click "Choose File"
2. Select a 1 MB test file
3. Observe progress bar
4. Wait for completion

**Expected Results:**
- ✅ Progress bar reaches 100%
- ✅ Status shows "X / X chunks"
- ✅ File appears in the table
- ✅ "Complete" status shown in green
- ✅ Download button is available

**Verification:**
```bash
# Check teamserver logs
./teamserver --log-level debug 2>&1 | grep -i "upload"

# Check teamserver storage
ls -la ./data/files/
```

### Test 2: Large File Upload (>1 GB)
**Setup:** Same as Test 1

**Test Steps:**
1. Upload 1 GB test file
2. Observe chunk progress
3. Monitor system resources

**Expected Results:**
- ✅ Chunks upload sequentially (progress 0% → 100%)
- ✅ File appears in UI during upload
- ✅ UI remains responsive
- ✅ File storage directory grows
- ✅ No data corruption on disk

**Verification:**
```bash
# Verify file size
du -h ./data/files/*/
```

### Test 3: File Download
**Setup:** File already uploaded and visible in UI

**Test Steps:**
1. Click download button
2. Browser initiates download
3. Verify file in downloads folder

**Expected Results:**
- ✅ Download starts immediately
- ✅ File saves to downloads folder
- ✅ Correct filename preserved
- ✅ File is accessible (not encrypted)

### Test 4: Encryption Verification
**Objective:** Verify chunks are encrypted on disk

**Setup:** File uploaded, storage location known

**Test Steps:**
```bash
# Check raw file data (should be binary, not readable)
hexdump -C ./data/files/agent-session-id/chunk_0000 | head -20

# Expected: Random binary data, not recognizable text
00000000  d4 73 b2 1a 9c f4 8e 2f  42 19 e3 c2 7a 8b 5f d1
00000010  a9 12 4c 68 3f 19 c7 e2  51 a3 92 d4 19 f2 8b e0
```

**Expected Results:**
- ✅ Chunk files are binary
- ✅ No readable text from original file
- ✅ Random-looking data confirms encryption
- ✅ File can be reconstructed from encrypted chunks

### Test 5: Multiple Agents
**Setup:** Multiple agents connected to teamserver

**Test Steps:**
1. Agent A uploads file "doc1.pdf"
2. Switch to Agent B in dropdown
3. Verify doc1.pdf is not shown
4. Agent B uploads file "doc2.pdf"
5. Switch to Agent A
6. Verify doc1.pdf is shown, not doc2.pdf

**Expected Results:**
- ✅ Each agent has isolated file list
- ✅ No cross-agent file visibility
- ✅ Agent filter works correctly

### Test 6: Concurrent Operations
**Setup:** Multiple browser tabs with file transfer UI

**Test Steps:**
1. Tab 1: Start uploading large file
2. Tab 2: Download different file
3. Tab 3: Refresh file list
4. Tab 4: Select different agent

**Expected Results:**
- ✅ All operations complete successfully
- ✅ No race conditions
- ✅ File integrity maintained
- ✅ UI remains responsive

### Test 7: Network Interruption
**Setup:** File transfer in progress

**Test Steps:**
1. Start uploading large file
2. When upload ~50% complete:
   - Disconnect network
   - Wait 30 seconds
   - Reconnect network
3. Observe retry mechanism

**Expected Results:**
- ✅ Upload resumes after reconnection
- ✅ No partial/corrupted files
- ✅ Chunks already uploaded preserved
- ✅ Final file integrity maintained

### Test 8: HMAC Authentication
**Objective:** Verify HMAC prevents tampering

**Setup:** File uploaded, stored on teamserver

**Test Steps:**
```bash
# Simulate tampering with encrypted chunk
CHUNK=$(ls ./data/files/*/chunk_0000 | head -1)

# Backup original
cp "$CHUNK" "${CHUNK}.bak"

# Modify chunk (flip a bit)
dd if="${CHUNK}.bak" of="$CHUNK" bs=1 count=50 2>/dev/null
echo -n "TAMPERED" | dd of="$CHUNK" seek=100 bs=1 conv=notrunc 2>/dev/null

# Try to download
# Expected: Agent detects HMAC mismatch
```

**Expected Results:**
- ✅ Agent detects HMAC mismatch
- ✅ Agent rejects corrupted chunk
- ✅ Error logged in agent
- ✅ File marked as corrupted in UI

## Performance Benchmarks

### Target Metrics
| Metric | Target | Notes |
|--------|--------|-------|
| Upload Speed | 10-50 MBps | Network dependent |
| Download Speed | 10-50 MBps | Network dependent |
| Decryption Overhead | <5% | CPU limited |
| UI Responsiveness | <100ms | React updates |
| Memory per 256KiB chunk | ~5 MB | During processing |

### Load Testing
```bash
# Generate test file
dd if=/dev/zero of=test-1gb.bin bs=1M count=1024

# Upload via agent
time ./redforge_agent upload test-1gb.bin

# Monitor system resources
watch -n 1 'ps aux | grep redforge'
```

## Monitoring Checklist

### Teamserver Monitoring
- [x] File upload request rate
- [x] Chunk storage rate
- [x] Disk space usage
- [x] Memory usage (in-memory chunks)
- [x] API response times
- [x] Authentication failures

### Agent Monitoring
- [x] Encryption/Decryption time
- [x] Network errors/retries
- [x] Memory usage (chunk buffer)
- [x] CPU usage (AES-256-GCM)
- [x] File I/O performance

### Web UI Monitoring
- [x] Page load time
- [x] API request latency
- [x] File list refresh rate
- [x] Download initiation time
- [x] Browser memory usage

## Logging Configuration

### Agent Debug Logging
```bash
./redforge_agent --log-level debug 2>&1 | tee agent.log
```

Expected log patterns:
```
[DEBUG] Splitting file into chunks: 40 chunks
[DEBUG] Chunk 0: Size=262144, Nonce=0x[random], HMAC=0x[hash]
[DEBUG] Uploading chunk 0/40
[INFO] File transfer complete: chunks=40, total_size=10485760
```

### Teamserver Debug Logging
```bash
./teamserver --log-level debug 2>&1 | tee server.log
```

Expected log patterns:
```
[DEBUG] FileUpload: Agent=agent-01, Chunks=40, Size=10485760
[DEBUG] StoreChunk: session=xyz, chunk=0, encrypted_size=262175
[INFO] File complete: session=xyz, file_id=abc123, chunks=40
```

## Cleanup and Reset

### Clear All File Storage
```bash
rm -rf ./data/files/*
# Recreate empty directory
mkdir -p ./data/files
```

### Clear Agent Cache
```bash
# On agent system
rm -rf ~/.redforge/cache/
```

### Browser Cache
```
# In browser dev tools
Application → Storage → Clear All
```

## Success Criteria

✅ **All Tests Pass**: All 8 test scenarios complete successfully
✅ **No Errors**: No errors in agent, teamserver, or UI logs
✅ **File Integrity**: Downloaded files match originals (checksums)
✅ **Encryption**: Encrypted chunks confirmed on disk
✅ **Performance**: Within target benchmarks
✅ **Security**: HMAC verification prevents tampering
✅ **Scalability**: Multiple concurrent operations work
✅ **Recovery**: Network interruptions handled gracefully

## Rollback Plan

If issues occur post-deployment:

1. **Stop new uploads**: Disable file transfer in UI (remove nav link)
2. **Investigate logs**: Check agent and server logs
3. **Backup data**: `cp -r ./data/files ./data/files.backup`
4. **Rollback binaries**: Use previous working versions
5. **Clear cache**: Reset in-memory storage
6. **Resume**: Re-enable after fixes

## Post-Deployment Checklist

- [ ] All three binaries deployed (agent, teamserver, ui)
- [ ] File storage directory created and writable
- [ ] HTTPS certificates configured
- [ ] Agent authentication tokens configured
- [ ] All 8 test scenarios passed
- [ ] Monitoring/logging configured
- [ ] Backup strategy in place
- [ ] Documentation reviewed by operations team
- [ ] Incident response procedures documented
- [ ] Performance baselines recorded
