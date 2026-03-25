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

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct TransportStatus {
    pub consecutive_failures: u64,
    pub last_error: Option<String>,
    pub last_backoff_ms: Option<u64>,
    pub last_attempt_at: Option<String>,
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
    #[serde(skip_serializing_if = "Option::is_none")]
    pub transport: Option<TransportStatus>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HeartbeatResponse {
    pub status: String,
    pub tasks: Vec<TaskMessage>,
}

/// Encrypted file chunk for transfer between agent and operator.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FileChunk {
    /// Unique identifier for this file transfer session
    pub session_id: String,
    /// Sequential chunk number (0-indexed)
    pub chunk_number: usize,
    /// Total number of chunks in this file
    pub total_chunks: usize,
    /// Original filename (for reference)
    pub filename: String,
    /// Size of original (unencrypted) data in this chunk
    pub original_size: usize,
    /// Nonce used for this chunk's encryption (base64-encoded)
    pub nonce: String,
    /// Encrypted chunk data (base64-encoded)
    pub data: String,
    /// HMAC-SHA256 for authentication (base64-encoded)
    pub hmac: String,
}

/// Upload request from agent to operator (sends encrypted file chunks).
#[derive(Debug, Serialize, Deserialize)]
pub struct FileUploadRequest {
    pub agent_id: String,
    pub token: String,
    pub chunk: FileChunk,
}

/// Upload response from operator confirming chunk receipt.
#[derive(Debug, Serialize, Deserialize)]
pub struct FileUploadResponse {
    pub status: String,
    pub message: Option<String>,
    pub chunk_number: usize,
}

/// Download request from agent to operator (requests encrypted file chunks).
#[derive(Debug, Serialize, Deserialize)]
pub struct FileDownloadRequest {
    pub agent_id: String,
    pub token: String,
    pub file_id: String,
    pub chunk_number: usize,
}

/// Download response from operator (sends encrypted file chunks).
#[derive(Debug, Serialize, Deserialize)]
pub struct FileDownloadResponse {
    pub chunk: FileChunk,
}

