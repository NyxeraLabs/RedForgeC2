package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/config"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/registry"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/users"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server represents the teamserver HTTP API.
type Server struct {
	config   *config.Config
	mux      *http.ServeMux
	logger   *log.Logger
	registry *registry.Registry
	users    *users.Store

	hardening   hardeningConfig
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

	hard := hardeningFromEnv()
	server := &Server{
		config:        cfg,
		mux:           mux,
		logger:        logger,
		registry:      registry.New(pool),
		users:         userStore,
		hardening:     hard,
		loginLimiter:  newIPRateLimiter(hard.loginRPM, time.Minute),
	}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/api/register", server.handleRegister)
	mux.HandleFunc("/api/heartbeat", server.handleHeartbeat)
	mux.HandleFunc("/api/task/result", server.handleTaskResult)
	mux.Handle("/api/login", withLoginRateLimit(server.loginLimiter, http.HandlerFunc(server.handleLogin)))
	mux.Handle("/api/me", AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(server.handleMe)))
	mux.Handle("/api/me/profile", AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(server.handleMeProfile)))
	mux.Handle("/api/me/password", AuthMiddleware(cfg.JWTSecret, http.HandlerFunc(server.handleMePassword)))

	operatorReadRoles := []string{string(users.RoleAdmin), string(users.RoleOperator), string(users.RoleObserver)}
	operatorWriteRoles := []string{string(users.RoleAdmin), string(users.RoleOperator)}
	mux.Handle("/api/operator/agents", AuthMiddleware(cfg.JWTSecret, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleAgentList))))
	mux.Handle("/api/operator/task", AuthMiddleware(cfg.JWTSecret, RequireAnyRole(operatorWriteRoles, http.HandlerFunc(server.handleTaskCreate))))
	mux.Handle("/api/operator/results", AuthMiddleware(cfg.JWTSecret, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleTaskResults))))
	mux.Handle("/api/operator/tasks", AuthMiddleware(cfg.JWTSecret, RequireAnyRole(operatorReadRoles, http.HandlerFunc(server.handleTasksList))))

	mux.Handle("/api/admin/users", AuthMiddleware(cfg.JWTSecret, RequireRole(string(users.RoleAdmin), http.HandlerFunc(server.handleAdminUsers))))

	return server
}

// Listen starts the HTTP server and blocks until stopped.
func (s *Server) Listen(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	handler := http.Handler(s.mux)
	handler = withBodyLimit(s.hardening.maxBodyBytes, handler)
	handler = withCORS(s.hardening.allowedOrigins, handler)
	handler = withSecurityHeaders(handler)
	handler = withRequestID(handler)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	if s.config.TLSCertFile != "" && s.config.TLSKeyFile != "" {
		s.logger.Printf("teamserver listening on https://0.0.0.0%s", addr)
	} else {
		s.logger.Printf("teamserver listening on http://0.0.0.0%s", addr)
	}

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctxShutdown)
	}()

	if s.config.TLSCertFile != "" && s.config.TLSKeyFile != "" {
		return httpServer.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
	}
	return httpServer.ListenAndServe()
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

	s.registry.UpdateHeartbeat(hb.AgentID)

	tasks := s.registry.GetTasks(hb.AgentID)
	resp := api.HeartbeatResponse{
		Status: "ok",
		Tasks:  tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AgentID       string   `json:"agent_id"`
		Command       string   `json:"command"`
		Args          []string `json:"args"`
		TimeoutSecond int      `json:"timeout_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task := api.TaskMessage{
		TaskID:  uuid.NewString(),
		Command: req.Command,
		Args:    req.Args,
		Timeout: req.TimeoutSecond,
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

	var result api.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.registry.ValidateToken(result.AgentID, result.Token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.registry.AddResult(result.AgentID, result)

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

func (s *Server) handleTasksList(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.registry.ListTasks(agentID))
}
