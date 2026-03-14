package api

import "time"

// AgentRegistration is the payload sent by an agent to register with the teamserver.
type AgentRegistration struct {
	AgentID  string            `json:"agent_id"`
	OS       string            `json:"os"`
	Arch     string            `json:"arch"`
	Hostname string            `json:"hostname"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Token    string            `json:"token,omitempty"`
}

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
	HeartbeatLastAttempt string            `json:"heartbeat_last_attempt,omitempty"`
	HeartbeatLastBackoff int64             `json:"heartbeat_last_backoff_ms,omitempty"`
}

type TaskMessage struct {
	TaskID  string   `json:"task_id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Timeout int      `json:"timeout_seconds"`
}

type TaskResult struct {
	AgentID   string    `json:"agent_id"`
	Token     string    `json:"token"`
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// TelemetryPayload contains periodic metrics from agents.
type TransportStatus struct {
	ConsecutiveFailures int    `json:"consecutive_failures"`
	LastError           string `json:"last_error,omitempty"`
	LastBackoffMs       int64  `json:"last_backoff_ms,omitempty"`
	LastAttemptAt       string `json:"last_attempt_at,omitempty"`
}

type TelemetryPayload struct {
	AgentID   string    `json:"agent_id"`
	CPU       float64   `json:"cpu"`
	Memory    uint64    `json:"memory"`
	Uptime    uint64    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

// HeartbeatRequest is sent by an agent to report telemetry and request tasks.
type HeartbeatRequest struct {
	AgentID   string           `json:"agent_id"`
	Token     string           `json:"token"`
	Telemetry TelemetryPayload `json:"telemetry"`
	Transport *TransportStatus `json:"transport,omitempty"`
}

// HeartbeatResponse is returned by teamserver to an agent.
type HeartbeatResponse struct {
	Status string        `json:"status"`
	Tasks  []TaskMessage `json:"tasks"`
}
