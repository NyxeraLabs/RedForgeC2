package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

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
	if _, ok := s.agents[id]; ok {
		return id
	}
	s.agents[id] = &agentState{
		id:   id,
		meta: req.Meta,
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
	for i := range a.tasks {
		if a.tasks[i].TaskID == taskID {
			a.tasks[i].Status = protocol.TaskCompleted
			return nil
		}
	}
	return ErrTaskNotFound
}

func newID() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}

