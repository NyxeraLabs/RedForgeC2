use std::env;

#[derive(Debug)]
pub struct AgentConfig {
    /// Teamserver base URL (e.g. https://localhost:9080)
    pub server_url: String,
    /// Optional path to a CA certificate to trust for mTLS/HTTPS.
    pub ca_cert_path: Option<String>,
    /// Agent heartbeat interval in seconds.
    pub heartbeat_interval: u64,
    /// Maximum number of retries for transport operations.
    pub transport_max_retries: usize,
    /// Initial backoff duration in milliseconds.
    pub transport_backoff_ms: u64,
}

impl AgentConfig {
    pub fn load() -> Self {
        let server_url = env::var("REDFORGE_SERVER_URL")
            .unwrap_or_else(|_| "https://localhost:9080".to_string());
        let ca_cert_path = env::var("REDFORGE_CA_CERT").ok();
        let heartbeat_interval = env::var("REDFORGE_HEARTBEAT_INTERVAL")
            .ok()
            .and_then(|v| v.parse().ok())
            .unwrap_or(60);
        let transport_max_retries = env::var("REDFORGE_TRANSPORT_RETRIES")
            .ok()
            .and_then(|v| v.parse().ok())
            .unwrap_or(3);
        let transport_backoff_ms = env::var("REDFORGE_TRANSPORT_BACKOFF_MS")
            .ok()
            .and_then(|v| v.parse().ok())
            .unwrap_or(250);

        AgentConfig {
            server_url,
            ca_cert_path,
            heartbeat_interval,
            transport_max_retries,
            transport_backoff_ms,
        }
    }
}
