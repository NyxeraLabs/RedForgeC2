package server

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
)

// FileStore manages encrypted file storage on the teamserver.
// Files are stored in chunks as received from agents.
type FileStore struct {
	mu         sync.RWMutex
	baseDir    string
	sessions   map[string]*FileSession   // session_id -> FileSession
	filesByID  map[string]*StoredFile    // file_id -> StoredFile
	chunksByID map[string]map[int][]byte // file_id -> chunk_number -> data
}

// FileSession tracks metadata about an ongoing file upload.
type FileSession struct {
	SessionID   string
	FileID      string
	Filename    string
	TotalChunks int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AgentID     string
}

// StoredFile represents a complete uploaded file.
type StoredFile struct {
	FileID         string    `json:"file_id"`
	SessionID      string    `json:"session_id"`
	Filename       string    `json:"filename"`
	TotalChunks    int       `json:"total_chunks"`
	ChunksReceived int       `json:"chunks_received"`
	UploadedAt     time.Time `json:"uploaded_at"`
	AgentID        string    `json:"agent_id"`
	FilePath       string    `json:"-"` // path to the file on disk
}

// NewFileStore creates a new file store with the given base directory.
func NewFileStore(baseDir string) (*FileStore, error) {
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create file store directory: %w", err)
	}

	return &FileStore{
		baseDir:    baseDir,
		sessions:   make(map[string]*FileSession),
		filesByID:  make(map[string]*StoredFile),
		chunksByID: make(map[string]map[int][]byte),
	}, nil
}

// StoreChunk stores an encrypted chunk from an agent.
// Returns the file ID for download requests.
func (fs *FileStore) StoreChunk(chunk *api.FileChunk, agentID string) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Verify this is a known session or start a new one
	session, exists := fs.sessions[chunk.SessionID]
	if !exists {
		// Create new session
		fileID := fmt.Sprintf("%s_%s_%d", agentID, chunk.SessionID[:8], time.Now().Unix())
		session = &FileSession{
			SessionID:   chunk.SessionID,
			FileID:      fileID,
			Filename:    chunk.Filename,
			TotalChunks: chunk.TotalChunks,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AgentID:     agentID,
		}
		fs.sessions[chunk.SessionID] = session
		fs.filesByID[fileID] = &StoredFile{
			FileID:      fileID,
			SessionID:   chunk.SessionID,
			Filename:    chunk.Filename,
			TotalChunks: chunk.TotalChunks,
			UploadedAt:  time.Now(),
			AgentID:     agentID,
		}
		fs.chunksByID[fileID] = make(map[int][]byte)
	}

	// Decode and store the chunk
	chunkData, err := base64.StdEncoding.DecodeString(chunk.Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode chunk data: %w", err)
	}

	fileID := session.FileID
	if _, exists := fs.chunksByID[fileID]; !exists {
		fs.chunksByID[fileID] = make(map[int][]byte)
	}

	fs.chunksByID[fileID][chunk.ChunkNumber] = chunkData
	fs.filesByID[fileID].ChunksReceived = len(fs.chunksByID[fileID])
	session.UpdatedAt = time.Now()

	// If all chunks received, write to disk
	if fs.filesByID[fileID].ChunksReceived == chunk.TotalChunks {
		if err := fs.writeFileToDisk(fileID); err != nil {
			return fileID, fmt.Errorf("failed to write file to disk: %w", err)
		}
	}

	return fileID, nil
}

// GetChunk retrieves an encrypted chunk for download to an agent.
func (fs *FileStore) GetChunk(fileID string, chunkNumber int) (*api.FileChunk, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	storedFile, exists := fs.filesByID[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}

	chunks, exists := fs.chunksByID[fileID]
	if !exists {
		// Try to load from disk
		return nil, fmt.Errorf("chunks not in memory for file: %s", fileID)
	}

	chunkData, exists := chunks[chunkNumber]
	if !exists {
		return nil, fmt.Errorf("chunk %d not found in file %s", chunkNumber, fileID)
	}

	// Reconstruct the FileChunk (without encryption details since we don't store those)
	// In a real implementation, we'd store all chunk metadata
	chunk := &api.FileChunk{
		SessionID:    storedFile.SessionID,
		ChunkNumber:  chunkNumber,
		TotalChunks:  storedFile.TotalChunks,
		Filename:     storedFile.Filename,
		OriginalSize: len(chunkData),
		Data:         base64.StdEncoding.EncodeToString(chunkData),
	}

	return chunk, nil
}

// ListFiles returns all uploaded files.
func (fs *FileStore) ListFiles(agentID string) []*StoredFile {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var files []*StoredFile
	for _, f := range fs.filesByID {
		if agentID == "" || f.AgentID == agentID {
			files = append(files, f)
		}
	}
	return files
}

// GetFile returns metadata about a stored file.
func (fs *FileStore) GetFile(fileID string) (*StoredFile, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	file, exists := fs.filesByID[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}
	return file, nil
}

// writeFileToDisk writes all chunks of a file to disk once complete.
func (fs *FileStore) writeFileToDisk(fileID string) error {
	storedFile := fs.filesByID[fileID]
	chunks := fs.chunksByID[fileID]

	// Create directory structure: baseDir/agent_id/uploads/
	uploadDir := filepath.Join(fs.baseDir, storedFile.AgentID, "uploads")
	if err := os.MkdirAll(uploadDir, 0700); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Write chunks in order to file
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%s_%s", fileID, storedFile.Filename))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for i := 0; i < storedFile.TotalChunks; i++ {
		chunkData, exists := chunks[i]
		if !exists {
			return fmt.Errorf("chunk %d missing for file %s", i, fileID)
		}
		if _, err := file.Write(chunkData); err != nil {
			return fmt.Errorf("failed to write chunk %d: %w", i, err)
		}
	}

	storedFile.FilePath = filePath
	return nil
}

// DeleteFile removes a file from storage (both metadata and chunks).
func (fs *FileStore) DeleteFile(fileID string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	storedFile, exists := fs.filesByID[fileID]
	if !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	// Remove from disk if it exists
	if storedFile.FilePath != "" {
		_ = os.Remove(storedFile.FilePath)
	}

	// Remove from memory
	delete(fs.filesByID, fileID)
	delete(fs.chunksByID, fileID)
	delete(fs.sessions, storedFile.SessionID)

	return nil
}

// CleanupExpiredChunks removes incomplete uploads older than maxAge.
func (fs *FileStore) CleanupExpiredChunks(maxAge time.Duration) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	now := time.Now()
	for fileID, file := range fs.filesByID {
		if file.ChunksReceived < file.TotalChunks && now.Sub(file.UploadedAt) > maxAge {
			delete(fs.filesByID, fileID)
			delete(fs.chunksByID, fileID)
			delete(fs.sessions, file.SessionID)
		}
	}
}

// StoreOperatorFile stores a file uploaded by an operator for agents to download.
// These files are NOT encrypted (encryption happens on the agent side).
func (fs *FileStore) StoreOperatorFile(filename string, agentID string, fileBytes []byte) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Generate file ID
	fileID := fmt.Sprintf("op_%s_%s_%d", agentID, filename, time.Now().Unix())

	// Create directory structure: baseDir/agent_id/downloads/
	downloadDir := filepath.Join(fs.baseDir, agentID, "downloads")
	if err := os.MkdirAll(downloadDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	// Write file to disk
	filePath := filepath.Join(downloadDir, filename)
	if err := os.WriteFile(filePath, fileBytes, 0600); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Store metadata in memory
	storedFile := &StoredFile{
		FileID:         fileID,
		SessionID:      fileID, // Use same as file ID for operator uploads
		Filename:       filename,
		TotalChunks:    1, // Not chunked for operator uploads
		ChunksReceived: 1,
		UploadedAt:     time.Now(),
		AgentID:        agentID,
		FilePath:       filePath,
	}

	fs.filesByID[fileID] = storedFile
	fs.chunksByID[fileID] = make(map[int][]byte)
	fs.chunksByID[fileID][0] = fileBytes

	return fileID, nil
}
