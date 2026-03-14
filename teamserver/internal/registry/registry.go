package registry

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Agent represents a registered agent.
type Agent struct {
	AgentID              string            `json:"agent_id"`
	Token                string            `json:"token"`
	OS                   string            `json:"os"`
	Arch                 string            `json:"arch"`
	Hostname             string            `json:"hostname"`
	Version              string            `json:"version"`
	Metadata             map[string]string `json:"metadata,omitempty"`
	LastSeen             time.Time         `json:"last_seen"`
	Registered           time.Time         `json:"registered"`
	HeartbeatFailures    int               `json:"heartbeat_failures,omitempty"`
	HeartbeatLastError   string            `json:"heartbeat_last_error,omitempty"`
	HeartbeatLastAttempt time.Time         `json:"heartbeat_last_attempt,omitempty"`
	HeartbeatLastBackoff int64             `json:"heartbeat_last_backoff_ms,omitempty"`
}

// Registry manages registered agents, tasks, and results.
// If a database pool is provided, it persists data.
type Registry struct {
	mu      sync.RWMutex
	agents  map[string]*Agent
	tasks   map[string][]*api.TaskMessage
	results map[string][]*api.TaskResult
	db      *pgxpool.Pool
}

// New creates an agent registry.
// If db is nil, a memory-backed registry is used.
func New(db *pgxpool.Pool) *Registry {
	return &Registry{
		agents:  make(map[string]*Agent),
		tasks:   make(map[string][]*api.TaskMessage),
		results: make(map[string][]*api.TaskResult),
		db:      db,
	}
}

// Register registers or updates an agent and returns a token.
func (r *Registry) Register(reg api.AgentRegistration) (string, error) {
	if r.db != nil {
		ctx := context.Background()
		var token string
		if err := r.db.QueryRow(ctx, "SELECT token FROM agents WHERE agent_id=$1", reg.AgentID).Scan(&token); err != nil {
			// ignore not found
			if !errors.Is(err, pgx.ErrNoRows) {
				return "", err
			}
		}

		if token == "" {
			token = uuid.NewString()
		}

		meta, _ := json.Marshal(reg.Metadata)
		now := time.Now().UTC()
		_, err := r.db.Exec(ctx, `
			INSERT INTO agents (agent_id, token, os, arch, hostname, version, metadata, last_seen, registered)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			ON CONFLICT (agent_id) DO UPDATE SET
				token = EXCLUDED.token,
				os = EXCLUDED.os,
				arch = EXCLUDED.arch,
				hostname = EXCLUDED.hostname,
				version = EXCLUDED.version,
				metadata = EXCLUDED.metadata,
				last_seen = EXCLUDED.last_seen
		`, reg.AgentID, token, reg.OS, reg.Arch, reg.Hostname, reg.Version, meta, now, now)
		if err != nil {
			return "", err
		}
		return token, nil
	}

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
	if r.db != nil {
		ctx := context.Background()
		args, _ := json.Marshal(task.Args)
		now := time.Now().UTC()
		_, _ = r.db.Exec(ctx, `
			INSERT INTO tasks (task_id, agent_id, command, args, timeout_seconds, status, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,'pending',$6,$6)
		`, task.TaskID, agentID, task.Command, args, task.Timeout, now)
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[agentID] = append(r.tasks[agentID], &task)
}

type TaskSummary struct {
	TaskID         string    `json:"task_id"`
	AgentID        string    `json:"agent_id"`
	Command        string    `json:"command"`
	Args           []string  `json:"args"`
	TimeoutSeconds int       `json:"timeout_seconds"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListTasks returns tasks (optionally for a specific agent).
func (r *Registry) ListTasks(agentID string) []TaskSummary {
	if r.db != nil {
		ctx := context.Background()
		var rows pgx.Rows
		var err error
		if agentID != "" {
			rows, err = r.db.Query(ctx, "SELECT task_id, agent_id, command, args, timeout_seconds, status, created_at, updated_at FROM tasks WHERE agent_id=$1 ORDER BY updated_at DESC", agentID)
		} else {
			rows, err = r.db.Query(ctx, "SELECT task_id, agent_id, command, args, timeout_seconds, status, created_at, updated_at FROM tasks ORDER BY updated_at DESC LIMIT 200")
		}
		if err != nil {
			return []TaskSummary{}
		}
		defer rows.Close()

		out := make([]TaskSummary, 0)
		for rows.Next() {
			var t TaskSummary
			var argsBytes []byte
			if err := rows.Scan(&t.TaskID, &t.AgentID, &t.Command, &argsBytes, &t.TimeoutSeconds, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
				continue
			}
			_ = json.Unmarshal(argsBytes, &t.Args)
			out = append(out, t)
		}
		return out
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]TaskSummary, 0)
	if agentID != "" {
		for _, t := range r.tasks[agentID] {
			out = append(out, TaskSummary{
				TaskID:         t.TaskID,
				AgentID:        agentID,
				Command:        t.Command,
				Args:           t.Args,
				TimeoutSeconds: t.Timeout,
				Status:         "pending",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
		}
		return out
	}
	for aid, tasks := range r.tasks {
		for _, t := range tasks {
			out = append(out, TaskSummary{
				TaskID:         t.TaskID,
				AgentID:        aid,
				Command:        t.Command,
				Args:           t.Args,
				TimeoutSeconds: t.Timeout,
				Status:         "pending",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			})
		}
	}
	return out
}

// GetTasks returns and clears queued tasks for an agent.
func (r *Registry) GetTasks(agentID string) []api.TaskMessage {
	if r.db != nil {
		ctx := context.Background()
		rows, err := r.db.Query(ctx, "SELECT task_id, command, args, timeout_seconds FROM tasks WHERE agent_id=$1 AND status='pending'", agentID)
		if err != nil {
			return nil
		}
		defer rows.Close()

		var out []api.TaskMessage
		for rows.Next() {
			var t api.TaskMessage
			var argsBytes []byte
			if err := rows.Scan(&t.TaskID, &t.Command, &argsBytes, &t.Timeout); err != nil {
				continue
			}
			_ = json.Unmarshal(argsBytes, &t.Args)
			out = append(out, t)
		}
		_, _ = r.db.Exec(ctx, "UPDATE tasks SET status='in-progress', updated_at=$2 WHERE agent_id=$1 AND status='pending'", agentID, time.Now().UTC())
		return out
	}

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
	if r.db != nil {
		ctx := context.Background()
		_, _ = r.db.Exec(ctx, `
			INSERT INTO task_results (task_id, status, output, error, timestamp)
			VALUES ($1,$2,$3,$4,$5)
		`, result.TaskID, result.Status, result.Output, result.Error, result.Timestamp)
		_, _ = r.db.Exec(ctx, "UPDATE tasks SET status=$2, updated_at=$3 WHERE task_id=$1", result.TaskID, result.Status, time.Now().UTC())
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.results[agentID] = append(r.results[agentID], &result)
}

// GetResults returns stored task results for an agent.
func (r *Registry) GetResults(agentID string) []api.TaskResult {
	if r.db != nil {
		ctx := context.Background()
		rows, err := r.db.Query(ctx, "SELECT task_id, status, output, error, timestamp FROM task_results WHERE task_id IN (SELECT task_id FROM tasks WHERE agent_id=$1) ORDER BY timestamp DESC", agentID)
		if err != nil {
			return []api.TaskResult{}
		}
		defer rows.Close()

		out := make([]api.TaskResult, 0)
		for rows.Next() {
			var r api.TaskResult
			if err := rows.Scan(&r.TaskID, &r.Status, &r.Output, &r.Error, &r.Timestamp); err != nil {
				continue
			}
			out = append(out, r)
		}
		return out
	}

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
	if r.db != nil {
		ctx := context.Background()
		var stored string
		if err := r.db.QueryRow(ctx, "SELECT token FROM agents WHERE agent_id=$1", agentID).Scan(&stored); err != nil {
			return false
		}
		return stored == token
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, ok := r.agents[agentID]
	return ok && agent.Token == token
}

// UpdateHeartbeat updates the last-seen timestamp for an agent.
func (r *Registry) UpdateHeartbeat(agentID string, transport *api.TransportStatus) {
	if r.db != nil {
		ctx := context.Background()
		_, _ = r.db.Exec(ctx, "UPDATE agents SET last_seen=$2 WHERE agent_id=$1", agentID, time.Now().UTC())
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if agent, ok := r.agents[agentID]; ok {
		agent.LastSeen = time.Now()
		if transport != nil {
			agent.HeartbeatFailures = transport.ConsecutiveFailures
			agent.HeartbeatLastError = transport.LastError
			// Parse attempt timestamp if present, otherwise leave as zero.
			if transport.LastAttemptAt != "" {
				if t, err := time.Parse(time.RFC3339, transport.LastAttemptAt); err == nil {
					agent.HeartbeatLastAttempt = t
				}
			}
			agent.HeartbeatLastBackoff = transport.LastBackoffMs
		}
	}
}

// ListAgents returns a snapshot of registered agents.
func (r *Registry) ListAgents() []*Agent {
	if r.db != nil {
		ctx := context.Background()
		rows, err := r.db.Query(ctx, "SELECT agent_id, token, os, arch, hostname, version, metadata, last_seen, registered FROM agents")
		if err != nil {
			return []*Agent{}
		}
		defer rows.Close()

		out := make([]*Agent, 0)
		for rows.Next() {
			var a Agent
			var metaBytes []byte
			if err := rows.Scan(&a.AgentID, &a.Token, &a.OS, &a.Arch, &a.Hostname, &a.Version, &metaBytes, &a.LastSeen, &a.Registered); err != nil {
				continue
			}
			_ = json.Unmarshal(metaBytes, &a.Metadata)
			out = append(out, &a)
		}
		return out
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Agent, 0, len(r.agents))
	for _, a := range r.agents {
		out = append(out, a)
	}
	return out
}
