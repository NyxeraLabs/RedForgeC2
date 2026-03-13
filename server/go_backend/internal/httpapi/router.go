package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"redforgec2/server/go_backend/internal/protocol"
	"redforgec2/server/go_backend/internal/store"
)

type Router struct {
	store *store.Store
	mux   *http.ServeMux
}

func NewRouter(st *store.Store) http.Handler {
	r := &Router{
		store: st,
		mux:   http.NewServeMux(),
	}

	r.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	r.mux.HandleFunc("/api/v1/", r.handleV1)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) handleV1(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}

	switch parts[0] {
	case "agents":
		r.handleAgents(w, req, parts[1:])
		return
	case "tasks":
		r.handleTasks(w, req, parts[1:])
		return
	default:
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
}

func (r *Router) handleAgents(w http.ResponseWriter, req *http.Request, parts []string) {
	// POST /api/v1/agents/register
	if len(parts) == 1 && parts[0] == "register" && req.Method == http.MethodPost {
		var rr protocol.RegisterRequest
		if err := decodeJSON(req, &rr); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		agentID := r.store.Register(rr)
		writeJSON(w, http.StatusOK, protocol.RegisterResponse{
			AgentID:    agentID,
			ServerTime: protocol.NowRFC3339Nano(),
		})
		return
	}

	// /api/v1/agents/{agent_id}/...
	if len(parts) < 2 {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
	agentID := parts[0]

	// POST /api/v1/agents/{agent_id}/telemetry
	if len(parts) == 2 && parts[1] == "telemetry" && req.Method == http.MethodPost {
		var tr protocol.TelemetryRequest
		if err := decodeJSON(req, &tr); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := r.store.AddTelemetry(agentID, tr.Events); err != nil {
			if errors.Is(err, store.ErrAgentNotFound) {
				writeErr(w, http.StatusNotFound, "agent not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]any{"ok": true})
		return
	}

	// GET /api/v1/agents/{agent_id}/tasks
	if len(parts) == 2 && parts[1] == "tasks" && req.Method == http.MethodGet {
		tasks, err := r.store.PollTasks(agentID)
		if err != nil {
			if errors.Is(err, store.ErrAgentNotFound) {
				writeErr(w, http.StatusNotFound, "agent not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, protocol.PollTasksResponse{Tasks: tasks})
		return
	}

	// POST /api/v1/agents/{agent_id}/tasks/{task_id}/result
	if len(parts) == 4 && parts[1] == "tasks" && parts[3] == "result" && req.Method == http.MethodPost {
		taskID := parts[2]
		var sr protocol.SubmitResultRequest
		if err := decodeJSON(req, &sr); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := r.store.SubmitResult(agentID, taskID, sr); err != nil {
			if errors.Is(err, store.ErrAgentNotFound) {
				writeErr(w, http.StatusNotFound, "agent not found")
				return
			}
			if errors.Is(err, store.ErrTaskNotFound) {
				writeErr(w, http.StatusNotFound, "task not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, protocol.SubmitResultResponse{Ok: true, ServerTime: protocol.NowRFC3339Nano()})
		return
	}

	writeErr(w, http.StatusNotFound, "not found")
}

func (r *Router) handleTasks(w http.ResponseWriter, req *http.Request, parts []string) {
	// POST /api/v1/tasks/enqueue
	if len(parts) == 1 && parts[0] == "enqueue" && req.Method == http.MethodPost {
		var er protocol.EnqueueTaskRequest
		if err := decodeJSON(req, &er); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if er.AgentID == "" || er.TaskType == "" {
			writeErr(w, http.StatusBadRequest, "agent_id and task_type are required")
			return
		}
		task, err := r.store.EnqueueTask(er.AgentID, er.TaskType, er.Params)
		if err != nil {
			if errors.Is(err, store.ErrAgentNotFound) {
				writeErr(w, http.StatusNotFound, "agent not found")
				return
			}
			writeErr(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, protocol.EnqueueTaskResponse{Task: task})
		return
	}

	writeErr(w, http.StatusNotFound, "not found")
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

