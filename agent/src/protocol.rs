use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentRegistration {
    pub agent_id: String,
    pub os: String,
    pub arch: String,
    pub hostname: String,
    pub version: String,
    #[serde(default)]
    pub metadata: HashMap<String, String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TaskMessage {
    pub task_id: String,
    pub command: String,
    pub args: Vec<String>,
    pub timeout_seconds: u64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TaskResult {
    pub agent_id: String,
    pub token: String,
    pub task_id: String,
    pub status: String,
    pub output: String,
    pub error: Option<String>,
    pub timestamp: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TelemetryPayload {
    pub agent_id: String,
    pub cpu: f64,
    pub memory: u64,
    pub uptime: u64,
    pub timestamp: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HeartbeatRequest {
    pub agent_id: String,
    pub token: String,
    pub telemetry: TelemetryPayload,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HeartbeatResponse {
    pub status: String,
    pub tasks: Vec<TaskMessage>,
}
