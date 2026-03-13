package registry

import (
	"sync"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/google/uuid"
)

// Agent represents a registered agent.
type Agent struct {
	AgentID    string
	Token      string
	OS         string
	Arch       string
	Hostname   string
	Version    string
	Metadata   map[string]string
	LastSeen   time.Time
	Registered time.Time
}

// Registry manages registered agents in memory.
type Registry struct {
	mu      sync.RWMutex
	agents  map[string]*Agent
	tasks   map[string][]*api.TaskMessage
	results map[string][]*api.TaskResult
}

// New creates an empty agent registry.
func New() *Registry {
	return &Registry{
		agents:  make(map[string]*Agent),
		tasks:   make(map[string][]*api.TaskMessage),
		results: make(map[string][]*api.TaskResult),
	}
}

// Register registers or updates an agent and returns a token.
func (r *Registry) Register(reg api.AgentRegistration) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, found := r.agents[reg.AgentID]
	if !found {
		agent = &Agent{AgentID: reg.AgentID, Registered: time.Now()}
		r.agents[reg.AgentID] = agent
	}

	if agent.Token == "" {
		agent.Token = uuid.NewString()
	}

	agent.OS = reg.OS
	agent.Arch = reg.Arch
	agent.Hostname = reg.Hostname
	agent.Version = reg.Version
	agent.Metadata = reg.Metadata
	agent.LastSeen = time.Now()

	return agent.Token, nil
}

// AddTask adds a task to an agent's queue.
func (r *Registry) AddTask(agentID string, task api.TaskMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[agentID] = append(r.tasks[agentID], &task)
}

// GetTasks returns and clears queued tasks for an agent.
func (r *Registry) GetTasks(agentID string) []api.TaskMessage {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := r.tasks[agentID]
	out := make([]api.TaskMessage, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, *t)
	}
	delete(r.tasks, agentID)
	return out
}

// AddResult stores a task result for an agent.
func (r *Registry) AddResult(agentID string, result api.TaskResult) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results[agentID] = append(r.results[agentID], &result)
}

// GetResults returns stored task results for an agent.
func (r *Registry) GetResults(agentID string) []api.TaskResult {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := r.results[agentID]
	out := make([]api.TaskResult, 0, len(results))
	for _, r := range results {
		out = append(out, *r)
	}
	return out
}

// ValidateToken checks if the given token matches the agent's token.
func (r *Registry) ValidateToken(agentID, token string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, ok := r.agents[agentID]
	return ok && agent.Token == token
}

// UpdateHeartbeat updates the last-seen timestamp for an agent.
func (r *Registry) UpdateHeartbeat(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if agent, ok := r.agents[agentID]; ok {
		agent.LastSeen = time.Now()
	}
}

// ListAgents returns a snapshot of registered agents.
func (r *Registry) ListAgents() []*Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Agent, 0, len(r.agents))
	for _, a := range r.agents {
		out = append(out, a)
	}
	return out
}
