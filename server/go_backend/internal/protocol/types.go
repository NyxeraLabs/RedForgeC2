package protocol

import "time"

type RegisterRequest struct {
	AgentID  string            `json:"agent_id,omitempty"`
	OS       string            `json:"os,omitempty"`
	Arch     string            `json:"arch,omitempty"`
	Hostname string            `json:"hostname,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
}

type RegisterResponse struct {
	AgentID    string `json:"agent_id"`
	ServerTime string `json:"server_time"`
}

type TelemetryEvent struct {
	Timestamp string             `json:"timestamp"`
	Metrics   map[string]float64 `json:"metrics,omitempty"`
	Labels    map[string]string  `json:"labels,omitempty"`
}

type TelemetryRequest struct {
	Events []TelemetryEvent `json:"events"`
}

type TaskStatus string

const (
	TaskQueued     TaskStatus = "queued"
	TaskDispatched TaskStatus = "dispatched"
	TaskCompleted  TaskStatus = "completed"
)

type Task struct {
	TaskID    string         `json:"task_id"`
	TaskType  string         `json:"task_type"`
	Params    map[string]any `json:"params,omitempty"`
	CreatedAt string         `json:"created_at"`
	Status    TaskStatus     `json:"status"`
}

type PollTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type EnqueueTaskRequest struct {
	AgentID  string         `json:"agent_id"`
	TaskType string         `json:"task_type"`
	Params   map[string]any `json:"params,omitempty"`
}

type EnqueueTaskResponse struct {
	Task Task `json:"task"`
}

type SubmitResultRequest struct {
	Outcome map[string]any `json:"outcome,omitempty"`
	Status  string        `json:"status,omitempty"`
}

type SubmitResultResponse struct {
	Ok         bool   `json:"ok"`
	ServerTime string `json:"server_time"`
}

func NowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

