package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/auth"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/config"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/registry"
	"github.com/google/uuid"
)

// Server represents the teamserver HTTP API.
type Server struct {
	config   *config.Config
	mux      *http.ServeMux
	logger   *log.Logger
	registry *registry.Registry
}

// New creates a new teamserver HTTP server.
func New(cfg *config.Config, logger *log.Logger) *Server {
	mux := http.NewServeMux()
	server := &Server{config: cfg, mux: mux, logger: logger, registry: registry.New()}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/api/register", server.handleRegister)
	mux.HandleFunc("/api/heartbeat", server.handleHeartbeat)
	mux.HandleFunc("/api/task/result", server.handleTaskResult)
	mux.HandleFunc("/api/login", server.handleLogin)
	mux.Handle("/api/operator/agents", AuthMiddleware(cfg.JWTSecret, RequireRole("admin", http.HandlerFunc(server.handleAgentList))))
	mux.Handle("/api/operator/task", AuthMiddleware(cfg.JWTSecret, RequireRole("admin", http.HandlerFunc(server.handleTaskCreate))))
	mux.Handle("/api/operator/results", AuthMiddleware(cfg.JWTSecret, RequireRole("admin", http.HandlerFunc(server.handleTaskResults))))

	return server
}

// Listen starts the HTTP server and blocks until stopped.
func (s *Server) Listen(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	s.logger.Printf("teamserver listening on %s", addr)

	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctxShutdown)
	}()

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
		TaskID:        uuid.NewString(),
		Command:       req.Command,
		Args:          req.Args,
		Timeout:       req.TimeoutSecond,
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

	if creds.Username != s.config.AdminUsername || creds.Password != s.config.AdminPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(s.config.JWTSecret, creds.Username, "admin", s.config.TokenExpiryMins)
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
