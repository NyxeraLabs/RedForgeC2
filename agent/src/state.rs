use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

const STATE_FILE: &str = ".redforge_state.json";

#[derive(Debug, Serialize, Deserialize)]
pub struct AgentState {
    pub agent_id: String,
    pub token: Option<String>,
}

impl AgentState {
    pub fn load() -> Option<Self> {
        let path = Self::path();
        let data = fs::read_to_string(path).ok()?;
        serde_json::from_str(&data).ok()
    }

    pub fn save(&self) -> std::io::Result<()> {
        let path = Self::path();
        let data = serde_json::to_string_pretty(self)?;
        fs::write(path, data)
    }

    fn path() -> PathBuf {
        std::env::current_dir().unwrap_or_default().join(STATE_FILE)
    }
}
