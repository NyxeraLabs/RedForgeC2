package server

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/registry"
)

func TestHandleTaskResult_AllowsNullError(t *testing.T) {
	reg := registry.New(nil)
	token, err := reg.Register(api.AgentRegistration{
		AgentID:  "agent-1",
		OS:       "linux",
		Arch:     "amd64",
		Hostname: "host",
		Version:  "v0",
	})
	if err != nil {
		t.Fatalf("register agent: %v", err)
	}

	s := &Server{registry: reg, wsHub: newWsHub()}

	ts := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)
	body := map[string]any{
		"agent_id":  "agent-1",
		"token":     token,
		"task_id":   "task-1",
		"status":    "success",
		"output":    "hello\nworld\n",
		"error":     nil,
		"timestamp": ts,
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/task/result", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	s.handleTaskResult(rr, req)

	if rr.Code != 200 {
		t.Fatalf("expected 200, got %d (%s)", rr.Code, rr.Body.String())
	}

	results := reg.GetResults("agent-1")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Output != "hello\nworld\n" {
		t.Fatalf("unexpected output: %q", results[0].Output)
	}
	if results[0].Error != "" {
		t.Fatalf("expected empty error, got %q", results[0].Error)
	}
	if !results[0].Timestamp.Equal(ts) {
		t.Fatalf("expected timestamp %s, got %s", ts.Format(time.RFC3339), results[0].Timestamp.Format(time.RFC3339))
	}
}
