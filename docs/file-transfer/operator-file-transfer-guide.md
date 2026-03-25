# File Transfer - Operator Guide

## Accessing File Transfer

1. **Open the Web UI** - Navigate to your RedForgeC2 operator dashboard
2. **Click "File Transfer"** - Located in the main navigation menu
3. **Select an Agent** - Use the dropdown to choose which agent to work with

## Uploading Files from Operator to Agent

1. **Click File Upload Button** - "Choose File" input
2. **Select File** - Pick any file from your computer
3. **Upload Progress** - Watch the progress bar fill as the file uploads in 256 KiB chunks
4. **Confirmation** - Toast notification confirms successful upload

**What's Happening:**
- File is automatically split into 256 KiB chunks
- Each chunk is encrypted with AES-256-GCM (military-grade encryption)
- HMAC-SHA256 authentication ensures chunks aren't tampered with
- Sent over secure HTTPS connection
- Agent automatically decrypts and reassembles the file

## Downloading Files from Agent to Operator

1. **View Files Table** - Shows all files uploaded by selected agent
2. **Locate File** - Find the file you want to download
3. **Click Download Button** - Initiates download
4. **Monitor Progress** - Browser handles the download
5. **File Complete** - Check your downloads folder

**File Details Shown:**
- **Filename** - Name of the file (red highlighted)
- **Chunks** - Shows received/total (e.g., "45/45" = complete)
- **Status Bar** - Visual progress indicator (green = complete, yellow = in-progress)
- **Uploaded** - When the file was transferred
- **Action** - Download button

## Understanding the Progress Bar

```
Status Display: "45 / 45 chunks"
Progress Bar:   [████████████████████████] 100%

Status Display: "32 / 45 chunks"
Progress Bar:   [████████████░░░░░░░░░░░] 71%
```

- **Green Bar**: File transfer is complete and verified
- **Yellow Bar**: File transfer in progress
- **Complete**: All chunks received and HMAC verified

## Security Features (What's Protecting Your Files)

✅ **Encryption**: AES-256-GCM (same as government data encryption)
  - Uses 96-bit random nonce per chunk
  - 256-bit encryption key (derived from agent ID)

✅ **Authentication**: HMAC-SHA256
  - Verifies chunks aren't modified during transfer
  - Computed over both metadata and encrypted data

✅ **Transport Security**: HTTPS TLS 1.2+
  - Additional layer of encryption for network transport
  - Encrypts data in-flight between operator and teamserver

✅ **Agent Authentication**: Token-based
  - Bearer token verified for all file operations
  - Prevents unauthorized access

## Chunk Architecture

Files are automatically split into **256 KiB chunks**:
- **Why?**: Balances memory usage with efficiency
- **How?**: Transparent - you never see the chunks
- **Example**: 10 MB file = 40 chunks

```
File (10 MB) → Split into 40 × 256 KiB chunks
                ↓
           Each chunk encrypted
                ↓
           Each chunk authenticated (HMAC)
                ↓
           All chunks sequentially ordered
                ↓
           Sent over HTTPS to teamserver
```

## Troubleshooting

### "No Files Shown"
- Select an agent from the dropdown
- Agent may not have uploaded any files yet
- Click "Refresh Files" button

### "Upload Not Starting"
- Verify agent is connected to teamserver
- Check agent logs for errors
- Ensure adequate disk space on agent

### "Download Taking Long Time"
- Large files take time to decrypt locally
- Check network connection
- Verify teamserver is responsive

### "Agent Dropdown Empty"
- No agents connected to teamserver
- Check teamserver service is running
- Verify agent authentication credentials

## Best Practices

1. **Large Files**: Split 1GB+ files before uploading
2. **Sensitive Data**: Verify HTTPS connection (lock icon in browser)
3. **Offline Access**: Download important files while agent is connected
4. **Bandwidth**: Large uploads may take time - check network speed
5. **Disk Space**: Ensure sufficient space on agent before uploading

## Technical Details

**Supported**: Any file type and size
**Chunk Size**: 256 KiB (automatically handled)
**Encryption**: AES-256-GCM per chunk
**Authentication**: HMAC-SHA256 per chunk
**Transport**: HTTPS (TLS 1.2+)
**Storage**: In-memory with disk persistence on teamserver

## Support

For issues or questions:
1. Check agent logs: `./redforge_agent --log-level debug`
2. Check teamserver logs: `./teamserver --log-level debug`
3. Review security documentation: [SECURITY.md](../SECURITY.md)
4. File transfer architecture: [implementation.md](implementation.md)
