package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
)

// handleFileUpload handles file chunk uploads from agents.
// POST /api/files/upload
func (s *Server) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req api.FileUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate agent and token
	if !s.registry.ValidateToken(req.AgentID, req.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Store the chunk
	_, err := s.fileStore.StoreChunk(&req.Chunk, req.AgentID)
	if err != nil {
		s.logger.Printf("failed to store chunk: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "failed to store chunk",
		})
		return
	}

	// Update agent state for UI notification
	s.broadcastState()

	response := api.FileUploadResponse{
		Status:      "ok",
		ChunkNumber: req.Chunk.ChunkNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// handleFileDownload handles file chunk downloads from the operator to agents.
// POST /api/files/download
func (s *Server) handleFileDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req api.FileDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate agent and token
	if !s.registry.ValidateToken(req.AgentID, req.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the chunk from file store
	chunk, err := s.fileStore.GetChunk(req.FileID, req.ChunkNumber)
	if err != nil {
		s.logger.Printf("failed to retrieve chunk: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := api.FileDownloadResponse{
		Chunk: *chunk,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// handleFilesList returns a list of uploaded files (operator only).
// GET /api/files/list
func (s *Server) handleFilesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Query parameter for agent_id (optional)
	agentID := r.URL.Query().Get("agent_id")

	files := s.fileStore.ListFiles(agentID)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"files": files,
	})
}

// handleOperatorFileUpload allows operators to upload files to the teamserver.
// POST /api/operator/files/upload
func (s *Server) handleOperatorFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form data with 100MB limit
	err := r.ParseMultipartForm(100 * 1024 * 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to parse form data",
		})
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "no file provided",
		})
		return
	}
	defer file.Close()

	// Get agent_id from form
	agentID := r.FormValue("agent_id")
	if agentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "agent_id is required",
		})
		return
	}

	// Read file content
	fileBytes := make([]byte, handler.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to read file",
		})
		return
	}

	// Store the file (without encryption for operator uploads)
	fileID, err := s.fileStore.StoreOperatorFile(handler.Filename, agentID, fileBytes)
	if err != nil {
		s.logger.Printf("failed to store operator file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to store file",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"file_id":  fileID,
		"filename": handler.Filename,
	})
}

// handleFileInfo returns info about a specific file (operator only).
// GET /api/operator/files/{file_id}
func (s *Server) handleFileInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract file_id from path: /api/operator/files/{file_id}
	// Since the route is "/api/operator/files/", the remaining path is after that prefix
	path := r.URL.Path[len("/api/operator/files/"):]

	// Check if this is a download request
	if strings.HasSuffix(path, "/download") {
		// Remove /download suffix to get the file_id
		fileID := path[:len(path)-len("/download")]
		s.handleFileDownloadForOperator(w, r, fileID)
		return
	}

	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := s.fileStore.GetFile(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(file)
}

// handleFileDownloadForOperator returns a file for download by the operator.
// GET /api/operator/files/{file_id}/download?token={token}
func (s *Server) handleFileDownloadForOperator(w http.ResponseWriter, r *http.Request, fileID string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if fileID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the file metadata
	file, err := s.fileStore.GetFile(fileID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if file is complete
	if file.ChunksReceived < file.TotalChunks {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "file upload incomplete",
		})
		return
	}

	// Check if file exists on disk
	if file.FilePath == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "file not found on disk",
		})
		return
	}

	// Send file as binary download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+file.Filename+`"`)
	http.ServeFile(w, r, file.FilePath)
}
