package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/config"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/registry"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/users"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// wsHub is a lightweight websocket broadcaster used for real-time UI updates.
// It keeps a set of connected websocket clients and broadcasts JSON payloads to all of them.
type wsHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
}

func newWsHub() *wsHub {
	return &wsHub{conns: make(map[*websocket.Conn]struct{})}
}

func (h *wsHub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[conn] = struct{}{}
}

func (h *wsHub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, conn)
}

func (h *wsHub) broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.conns {
		_ = conn.WriteMessage(websocket.TextMessage, message)
	}
}

// Server represents the teamserver HTTP API.
type Server struct {
	config    *config.Config
	mux       *http.ServeMux
	logger    *log.Logger
	registry  *registry.Registry
	users     *users.Store
	wsHub     *wsHub
	fileStore *FileStore

	hardening    hardeningConfig
	loginLimiter *ipRateLimiter
}

// New creates a new teamserver HTTP server.
func New(cfg *config.Config, logger *log.Logger, pool *pgxpool.Pool) *Server {
	mux := http.NewServeMux()
	userStore := users.New(pool)
	resetAdmin := os.Getenv("REDFORGE_ADMIN_RESET") == "1"
	if err := userStore.EnsureAdmin(context.Background(), cfg.AdminUsername, cfg.AdminPassword, resetAdmin); err != nil {
		logger.Printf("warning: failed to ensure admin user: %v", err)
	}

	// Initialize file store
	fileStore, err := NewFileStore("./data/files")
	if err != nil {
		logger.Printf("warning: failed to initialize file store: %v", err)
		// Create a dummy file store to prevent nil pointer dereference
		fileStore, _ = NewFileStore("/tmp/redforge_files")
	}

	hard := hardeningFromEnv()
	server := &Server{
		config:       cfg,
		mux:          mux,
		logger:       logger,
		registry:     registry.New(pool),
		users:        userStore,
		wsHub:        newWsHub(),
		fileStore:    fileStore,
		hardening:    hard,
		loginLimiter: newIPRateLimiter(hard.loginRPM, time.Minute),
	}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/api/register", server.handleRegister)
	mux.HandleFunc("/api/heartbeat", server.handleHeartbeat)
	mux.HandleFunc("/api/task/result", server.handleTaskResult)
	mux.HandleFunc("/api/files/upload", server.handleFileUpload)
	mux.HandleFunc("/api/files/download", server.handleFileDownload)
	mux.Handle("/api/login", withLoginRateLimit(server.loginLimiter, http.HandlerFunc(server.handleLogin)))
	mux.Handle("/api/me", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, http.HandlerFunc(server.handleMe)))
	mux.Handle("/api/me/profile", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, http.HandlerFunc(server.handleMeProfile)))
	mux.Handle("/api/me/password", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, http.HandlerFunc(server.handleMePassword)))

	operatorReadRoles := []string{string(users.RoleAdmin), string(users.RoleOperator), string(users.RoleObserver)}
	operatorWriteRoles := []string{string(users.RoleAdmin), string(users.RoleOperator)}
	mux.Handle("/api/operator/agents", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleAgentList))))
	mux.Handle("/api/operator/agents/", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorWriteRoles, http.HandlerFunc(server.handleAgentDelete))))
	mux.Handle("/api/operator/task", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorWriteRoles, http.HandlerFunc(server.handleTaskCreate))))
	mux.Handle("/api/operator/results", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleTaskResults))))
	mux.Handle("/api/operator/tasks", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleTasksList))))
	mux.Handle("/api/operator/files", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleFilesList))))
	mux.Handle("/api/operator/files/upload", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorWriteRoles, http.HandlerFunc(server.handleOperatorFileUpload))))
	mux.Handle("/api/operator/files/", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleFileInfo))))

	mux.Handle("/api/admin/users", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireRole(string(users.RoleAdmin), http.HandlerFunc(server.handleAdminUsers))))
	mux.Handle("/api/admin/api-tokens", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, RequireRole(string(users.RoleAdmin), http.HandlerFunc(server.handleAdminAPITokens))))
	mux.Handle("/api/ws", AuthMiddlewareWithAPIToken(cfg.JWTSecret, server.users, http.HandlerFunc(server.handleWebsocket)))

	return server
}

// Listen starts the HTTP server and blocks until stopped.
func (s *Server) Listen(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	handler := http.Handler(s.mux)
	// Audit all requests for operational visibility and compliance.
	handler = withAuditLog(s.logger, handler)
	handler = withBodyLimit(s.hardening.maxBodyBytes, handler)
	handler = withCORS(s.hardening.allowedOrigins, handler)
	handler = withSecurityHeaders(handler)
	handler = withRequestID(handler)

	if s.config.TLSCertFile == "" || s.config.TLSKeyFile == "" {
		return fmt.Errorf("TLS certificates missing: set REDFORGE_TLS_CERT and REDFORGE_TLS_KEY")
	}

	if _, err := os.Stat(s.config.TLSCertFile); err != nil {
		return fmt.Errorf("TLS cert file missing: %s", s.config.TLSCertFile)
	}
	if _, err := os.Stat(s.config.TLSKeyFile); err != nil {
		return fmt.Errorf("TLS key file missing: %s", s.config.TLSKeyFile)
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	s.logger.Printf("teamserver listening on https://0.0.0.0%s", addr)

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctxShutdown)
	}()

	return httpServer.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var reg api.AgentRegistration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.registry.Register(reg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "token": token})

	// Notify connected UIs that agent state has changed.
	s.broadcastState()
}

func (s *Server) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var hb api.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&hb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.registry.ValidateToken(hb.AgentID, hb.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.registry.UpdateHeartbeat(hb.AgentID, hb.Transport)

	tasks := s.registry.GetTasks(hb.AgentID)
	if tasks == nil {
		tasks = []api.TaskMessage{}
	}
	resp := api.HeartbeatResponse{
		Status: "ok",
		Tasks:  tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)

	// Notify connected UIs that agent/task state may have changed.
	s.broadcastState()
}

func (s *Server) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AgentID        string   `json:"agent_id"`
		Command        string   `json:"command"`
		Args           []string `json:"args"`
		TimeoutSeconds int      `json:"timeout_seconds"`
		Transport      string   `json:"transport"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := api.TaskMessage{
		TaskID:  uuid.NewString(),
		Command: req.Command,
		Args:    req.Args,
		Timeout: req.TimeoutSeconds,
	}
	s.registry.AddTask(req.AgentID, task)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"task_id": task.TaskID})
}

func (s *Server) handleTaskResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// The agent serializes `error` as JSON null on success. Decode it as a pointer
	// so we accept both `"error": "..."` and `"error": null`.
	var in struct {
		AgentID   string     `json:"agent_id"`
		Token     string     `json:"token"`
		TaskID    string     `json:"task_id"`
		Status    string     `json:"status"`
		Output    string     `json:"output"`
		Error     *string    `json:"error"`
		Timestamp *time.Time `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.registry.ValidateToken(in.AgentID, in.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ts := time.Now().UTC()
	if in.Timestamp != nil {
		ts = in.Timestamp.UTC()
	}
	out := api.TaskResult{
		AgentID:   in.AgentID,
		Token:     in.Token,
		TaskID:    in.TaskID,
		Status:    in.Status,
		Output:    in.Output,
		Error:     "",
		Timestamp: ts,
	}
	if in.Error != nil {
		out.Error = *in.Error
	}

	s.registry.AddResult(in.AgentID, out)
	// Notify UIs that task status and results changed.
	s.broadcastState()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleTaskResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	agentID := r.URL.Query().Get("agent_id")
	if agentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	results := s.registry.GetResults(agentID)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := s.users.Authenticate(r.Context(), creds.Username, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(s.config.JWTSecret, u.Username, string(u.Role), s.config.TokenExpiryMins)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) handleAgentList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.registry.ListAgents())
}

func (s *Server) handleAgentDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	agentID := strings.TrimPrefix(r.URL.Path, "/api/operator/agents/")
	if agentID == "" || strings.Contains(agentID, "/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := s.registry.DeleteAgent(agentID); err != nil {
		if errors.Is(err, registry.ErrAgentNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Notify connected UIs that agent state has changed.
	s.broadcastState()

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u, err := s.users.GetByUsername(r.Context(), claims.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"username":      u.Username,
		"role":          u.Role,
		"display_name":  u.DisplayName,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
		"token_expires": claims.ExpiresAt,
	})
}

func (s *Server) handleMeProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	claims, ok := ClaimsFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var body struct {
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.users.UpdateProfile(r.Context(), claims.Username, body.DisplayName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleMePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	claims, ok := ClaimsFromRequest(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var body struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := s.users.ChangePassword(r.Context(), claims.Username, body.OldPassword, body.NewPassword); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usersList, err := s.users.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(usersList)
	case http.MethodPost:
		var body struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			Role        string `json:"role"`
			DisplayName string `json:"display_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.users.Create(r.Context(), body.Username, body.Password, users.Role(body.Role), body.DisplayName); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAdminAPITokens(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tokens, err := s.users.ListAPITokens(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tokens)
	case http.MethodPost:
		var body struct {
			Username       string `json:"username"`
			Description    string `json:"description"`
			ExpiresMinutes int    `json:"expires_minutes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var expiresAt *time.Time
		if body.ExpiresMinutes > 0 {
			t := time.Now().Add(time.Duration(body.ExpiresMinutes) * time.Minute)
			expiresAt = &t
		}

		token, err := s.users.CreateAPIToken(r.Context(), body.Username, body.Description, expiresAt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTasksList(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.registry.ListTasks(agentID))
}

// wsStatePayload is used to send full state snapshots over websocket.
type wsStatePayload struct {
	Type   string                 `json:"type"`
	Agents []*registry.Agent      `json:"agents,omitempty"`
	Tasks  []registry.TaskSummary `json:"tasks,omitempty"`
}

func (s *Server) broadcastState() {
	payload := wsStatePayload{
		Type:   "state",
		Agents: s.registry.ListAgents(),
		Tasks:  s.registry.ListTasks(""),
	}
	b, _ := json.Marshal(payload)
	s.wsHub.broadcast(b)
}

func (s *Server) sendState(conn *websocket.Conn) error {
	payload := wsStatePayload{
		Type:   "state",
		Agents: s.registry.ListAgents(),
		Tasks:  s.registry.ListTasks(""),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, b)
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	s.wsHub.add(conn)
	defer s.wsHub.remove(conn)

	_ = s.sendState(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
