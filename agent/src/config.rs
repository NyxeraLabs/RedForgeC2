use std::env;

#[derive(Debug)]
pub struct AgentConfig {
    /// Teamserver base URL (e.g. https://localhost:8080)
    pub server_url: String,
    /// Agent heartbeat interval in seconds.
    pub heartbeat_interval: u64,
}

impl AgentConfig {
    pub fn load() -> Self {
        let server_url =
            env::var("REDFORGE_SERVER_URL").unwrap_or_else(|_| "http://localhost:8080".to_string());
        let heartbeat_interval = env::var("REDFORGE_HEARTBEAT_INTERVAL")
            .ok()
            .and_then(|v| v.parse().ok())
            .unwrap_or(60);

        AgentConfig {
            server_url,
            heartbeat_interval,
        }
    }
}
