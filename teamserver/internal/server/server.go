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
)

// Server represents the teamserver HTTP API.
type Server struct {
	config *config.Config
	mux    *http.ServeMux
	logger *log.Logger
}

// New creates a new teamserver HTTP server.
func New(cfg *config.Config, logger *log.Logger) *Server {
	mux := http.NewServeMux()
	server := &Server{config: cfg, mux: mux, logger: logger}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/api/register", server.handleRegister)
	mux.HandleFunc("/api/heartbeat", server.handleHeartbeat)
	mux.HandleFunc("/api/login", server.handleLogin)
	mux.Handle("/api/operator/agents", AuthMiddleware(cfg.JWTSecret, RequireRole("admin", http.HandlerFunc(server.handleAgentList))))

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

	// TODO: Validate and store agent registration.

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// TODO: Parse telemetry payload, assign tasks.

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
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
	// TODO: Return list of registered agents.
	_ = json.NewEncoder(w).Encode([]interface{}{})
}
