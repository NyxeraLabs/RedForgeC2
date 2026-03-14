package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Init creates a Postgres connection pool and ensures the schema exists.
func Init(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err := ensureSchema(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			display_name TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS api_tokens (
			token_hash TEXT PRIMARY KEY,
			username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
			description TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMPTZ,
			last_used TIMESTAMPTZ
		);`,
		`CREATE TABLE IF NOT EXISTS agents (
			agent_id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			os TEXT,
			arch TEXT,
			hostname TEXT,
			version TEXT,
			metadata JSONB,
			last_seen TIMESTAMPTZ NOT NULL,
			registered TIMESTAMPTZ NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_agents_last_seen ON agents(last_seen);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			task_id TEXT PRIMARY KEY,
			agent_id TEXT NOT NULL REFERENCES agents(agent_id) ON DELETE CASCADE,
			command TEXT NOT NULL,
			args JSONB NOT NULL,
			timeout_seconds INT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_agent_status ON tasks(agent_id, status);`,
		`CREATE TABLE IF NOT EXISTS task_results (
			result_id SERIAL PRIMARY KEY,
			task_id TEXT NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
			status TEXT NOT NULL,
			output TEXT NOT NULL,
			error TEXT,
			timestamp TIMESTAMPTZ NOT NULL
		);`,
	}

	for _, q := range schema {
		if _, err := pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("migrate schema: %w", err)
		}
	}

	return nil
}

func encodeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func now() time.Time {
	return time.Now().UTC()
}
