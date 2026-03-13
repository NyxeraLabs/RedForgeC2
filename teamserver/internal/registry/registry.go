package registry

import (
	"context"
	"encoding/json"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Agent represents a registered agent.
type Agent struct {
	AgentID    string            `json:"agent_id"`
	Token      string            `json:"token"`
	OS         string            `json:"os"`
	Arch       string            `json:"arch"`
	Hostname   string            `json:"hostname"`
	Version    string            `json:"version"`
	Metadata   map[string]string `json:"metadata"`
	LastSeen   time.Time         `json:"last_seen"`
	Registered time.Time         `json:"registered"`
}

// Registry manages agents/tasks via Postgres.
type Registry struct {
	pool *pgxpool.Pool
}

// New creates a new registry backed by the given DB pool.
func New(pool *pgxpool.Pool) *Registry {
	return &Registry{pool: pool}
}

// Register registers or updates an agent and returns a token.
func (r *Registry) Register(reg api.AgentRegistration) (string, error) {
	ctx := context.Background()
	metadata, _ := json.Marshal(reg.Metadata)

	// Try to reuse existing token if it exists
	var token string
	err := r.pool.QueryRow(ctx, `SELECT token FROM agents WHERE agent_id = $1`, reg.AgentID).Scan(&token)
	if err != nil && err != pgx.ErrNoRows {
		return "", err
	}

	if token == "" {
		token = uuid.NewString()
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO agents (agent_id, token, os, arch, hostname, version, metadata, last_seen, registered)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now(), now())
		ON CONFLICT (agent_id) DO UPDATE
		SET token = $2,
		    os = EXCLUDED.os,
		    arch = EXCLUDED.arch,
		    hostname = EXCLUDED.hostname,
		    version = EXCLUDED.version,
		    metadata = EXCLUDED.metadata,
		    last_seen = now()
	`, reg.AgentID, token, reg.OS, reg.Arch, reg.Hostname, reg.Version, metadata)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken checks if the given token matches the agent's token.
func (r *Registry) ValidateToken(agentID, token string) bool {
	ctx := context.Background()
	var dbToken string
	if err := r.pool.QueryRow(ctx, `SELECT token FROM agents WHERE agent_id=$1`, agentID).Scan(&dbToken); err != nil {
		return false
	}
	return dbToken == token
}

// UpdateHeartbeat updates the last-seen timestamp for an agent.
func (r *Registry) UpdateHeartbeat(agentID string) {
	ctx := context.Background()
	_, _ = r.pool.Exec(ctx, `UPDATE agents SET last_seen = now() WHERE agent_id=$1`, agentID)
}

// ListAgents returns a snapshot of registered agents.
func (r *Registry) ListAgents() ([]*Agent, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `
		SELECT agent_id, token, os, arch, hostname, version, metadata, last_seen, registered
		FROM agents
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*Agent, 0)
	for rows.Next() {
		var a Agent
		var meta []byte
		if err := rows.Scan(&a.AgentID, &a.Token, &a.OS, &a.Arch, &a.Hostname, &a.Version, &meta, &a.LastSeen, &a.Registered); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(meta, &a.Metadata)
		out = append(out, &a)
	}
	return out, nil
}

// AddTask queues a task for an agent.
func (r *Registry) AddTask(agentID string, task api.TaskMessage) error {
	ctx := context.Background()
	args, _ := json.Marshal(task.Args)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tasks (task_id, agent_id, command, args, timeout_seconds, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())
	`, task.TaskID, agentID, task.Command, args, task.Timeout, "queued")
	return err
}

// GetTasks returns queued tasks for an agent (and marks them as sent).
func (r *Registry) GetTasks(agentID string) ([]api.TaskMessage, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `
		SELECT task_id, command, args, timeout_seconds
		FROM tasks
		WHERE agent_id = $1 AND status = 'queued'
		ORDER BY created_at ASC
	`, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]api.TaskMessage, 0)
	for rows.Next() {
		var t api.TaskMessage
		var args []byte
		if err := rows.Scan(&t.TaskID, &t.Command, &args, &t.Timeout); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(args, &t.Args)
		tasks = append(tasks, t)
	}

	if len(tasks) == 0 {
		return tasks, nil
	}

	// mark as sent
	ids := make([]any, 0, len(tasks))
	for _, t := range tasks {
		ids = append(ids, t.TaskID)
	}
	_, err = r.pool.Exec(ctx, `
		UPDATE tasks SET status = 'sent', updated_at = now() WHERE task_id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddResult stores a task result for an agent.
func (r *Registry) AddResult(agentID string, result api.TaskResult) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx, `
		UPDATE tasks
		SET status = $1, output = $2, error = $3, updated_at = now()
		WHERE task_id = $4 AND agent_id = $5
	`, result.Status, result.Output, result.Error, result.TaskID, agentID)
	return err
}

// GetResults returns stored task results for an agent.
func (r *Registry) GetResults(agentID string) ([]api.TaskResult, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, `
		SELECT task_id, status, output, error, updated_at
		FROM tasks
		WHERE agent_id = $1
		ORDER BY updated_at DESC
	`, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]api.TaskResult, 0)
	for rows.Next() {
		var r api.TaskResult
		if err := rows.Scan(&r.TaskID, &r.Status, &r.Output, &r.Error, &r.Timestamp); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
