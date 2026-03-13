package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"redforgec2/server/go_backend/internal/protocol"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
	ErrTaskNotFound  = errors.New("task not found")
)

type Store struct {
	mu     sync.RWMutex
	agents map[string]*agentState
}

type agentState struct {
	id        string
	os        string
	arch      string
	hostname  string
	createdAt time.Time
	lastSeen  time.Time
	meta      map[string]string
	telemetry []protocol.TelemetryEvent
	tasks     []protocol.Task
}

func New() *Store {
	return &Store{agents: map[string]*agentState{}}
}

func (s *Store) Register(req protocol.RegisterRequest) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.AgentID
	if id == "" {
		id = newID()
	}
	now := time.Now().UTC()
	if a, ok := s.agents[id]; ok {
		a.lastSeen = now
		return id
	}
	s.agents[id] = &agentState{
		id:        id,
		os:        req.OS,
		arch:      req.Arch,
		hostname:  req.Hostname,
		createdAt: now,
		lastSeen:  now,
		meta:      req.Meta,
	}
	return id
}

func (s *Store) AddTelemetry(agentID string, events []protocol.TelemetryEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.agents[agentID]
	if !ok {
		return ErrAgentNotFound
	}
	a.lastSeen = time.Now().UTC()
	a.telemetry = append(a.telemetry, events...)
	return nil
}

func (s *Store) EnqueueTask(agentID, taskType string, params map[string]any) (protocol.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.agents[agentID]
	if !ok {
		return protocol.Task{}, ErrAgentNotFound
	}
	a.lastSeen = time.Now().UTC()

	task := protocol.Task{
		TaskID:    newID(),
		TaskType:  taskType,
		Params:    params,
		CreatedAt: protocol.NowRFC3339Nano(),
		Status:    protocol.TaskQueued,
	}
	a.tasks = append(a.tasks, task)
	return task, nil
}

func (s *Store) PollTasks(agentID string) ([]protocol.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.agents[agentID]
	if !ok {
		return nil, ErrAgentNotFound
	}
	a.lastSeen = time.Now().UTC()

	var out []protocol.Task
	for i := range a.tasks {
		if a.tasks[i].Status == protocol.TaskQueued {
			a.tasks[i].Status = protocol.TaskDispatched
			out = append(out, a.tasks[i])
		}
	}
	return out, nil
}

func (s *Store) SubmitResult(agentID, taskID string, result protocol.SubmitResultRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.agents[agentID]
	if !ok {
		return ErrAgentNotFound
	}
	a.lastSeen = time.Now().UTC()
	for i := range a.tasks {
		if a.tasks[i].TaskID == taskID {
			a.tasks[i].Status = protocol.TaskCompleted
			return nil
		}
	}
	return ErrTaskNotFound
}

func (s *Store) ListAgents() []protocol.AgentSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]protocol.AgentSummary, 0, len(s.agents))
	for _, a := range s.agents {
		out = append(out, protocol.AgentSummary{
			AgentID:   a.id,
			OS:        a.os,
			Arch:      a.arch,
			Hostname:  a.hostname,
			Meta:      a.meta,
			CreatedAt: a.createdAt.Format(time.RFC3339Nano),
			LastSeen:  a.lastSeen.Format(time.RFC3339Nano),
		})
	}
	return out
}

func newID() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
