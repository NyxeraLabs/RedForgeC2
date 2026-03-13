package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"redforgec2/server/go_backend/internal/protocol"
	"redforgec2/server/go_backend/internal/store"
)

func TestRegisterEnqueuePollSubmitTelemetry(t *testing.T) {
	st := store.New()
	srv := httptest.NewServer(NewRouter(Deps{Store: st}))
	t.Cleanup(srv.Close)

	// Register
	regBody, _ := json.Marshal(protocol.RegisterRequest{OS: "linux", Arch: "amd64", Hostname: "lab"})
	regRes, err := http.Post(srv.URL+"/api/v1/agents/register", "application/json", bytes.NewReader(regBody))
	if err != nil {
		t.Fatalf("register post: %v", err)
	}
	defer regRes.Body.Close()
	if regRes.StatusCode != http.StatusOK {
		t.Fatalf("register status: %s", regRes.Status)
	}
	var rr protocol.RegisterResponse
	if err := json.NewDecoder(regRes.Body).Decode(&rr); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if rr.AgentID == "" {
		t.Fatalf("expected agent_id")
	}

	// Enqueue a task
	enqBody, _ := json.Marshal(protocol.EnqueueTaskRequest{AgentID: rr.AgentID, TaskType: "echo", Params: map[string]any{"msg": "hello"}})
	enqRes, err := http.Post(srv.URL+"/api/v1/tasks/enqueue", "application/json", bytes.NewReader(enqBody))
	if err != nil {
		t.Fatalf("enqueue post: %v", err)
	}
	defer enqRes.Body.Close()
	if enqRes.StatusCode != http.StatusOK {
		t.Fatalf("enqueue status: %s", enqRes.Status)
	}
	var er protocol.EnqueueTaskResponse
	if err := json.NewDecoder(enqRes.Body).Decode(&er); err != nil {
		t.Fatalf("decode enqueue response: %v", err)
	}
	if er.Task.TaskID == "" {
		t.Fatalf("expected task_id")
	}

	// Poll tasks
	pollRes, err := http.Get(srv.URL + "/api/v1/agents/" + rr.AgentID + "/tasks")
	if err != nil {
		t.Fatalf("poll get: %v", err)
	}
	defer pollRes.Body.Close()
	if pollRes.StatusCode != http.StatusOK {
		t.Fatalf("poll status: %s", pollRes.Status)
	}
	var pr protocol.PollTasksResponse
	if err := json.NewDecoder(pollRes.Body).Decode(&pr); err != nil {
		t.Fatalf("decode poll response: %v", err)
	}
	if len(pr.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(pr.Tasks))
	}
	if pr.Tasks[0].Status != protocol.TaskDispatched {
		t.Fatalf("expected dispatched, got %q", pr.Tasks[0].Status)
	}

	// Submit result
	subBody, _ := json.Marshal(protocol.SubmitResultRequest{Outcome: map[string]any{"ok": true}, Status: "simulated"})
	subRes, err := http.Post(srv.URL+"/api/v1/agents/"+rr.AgentID+"/tasks/"+er.Task.TaskID+"/result", "application/json", bytes.NewReader(subBody))
	if err != nil {
		t.Fatalf("submit post: %v", err)
	}
	defer subRes.Body.Close()
	if subRes.StatusCode != http.StatusOK {
		t.Fatalf("submit status: %s", subRes.Status)
	}

	// Telemetry
	telBody, _ := json.Marshal(protocol.TelemetryRequest{
		Events: []protocol.TelemetryEvent{{Timestamp: protocol.NowRFC3339Nano(), Metrics: map[string]float64{"cpu": 0.42}}},
	})
	telRes, err := http.Post(srv.URL+"/api/v1/agents/"+rr.AgentID+"/telemetry", "application/json", bytes.NewReader(telBody))
	if err != nil {
		t.Fatalf("telemetry post: %v", err)
	}
	defer telRes.Body.Close()
	if telRes.StatusCode != http.StatusAccepted {
		t.Fatalf("telemetry status: %s", telRes.Status)
	}
}
