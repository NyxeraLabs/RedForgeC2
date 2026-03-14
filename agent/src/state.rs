use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use zeroize::Zeroizing;

const STATE_FILE: &str = ".redforge_state.json";

/// Agent persisted state (saved to disk).
///
/// The token is kept in memory using `Zeroizing` to ensure it is wiped on drop.
#[derive(Debug)]
pub struct AgentState {
    pub agent_id: String,
    pub token: Option<Zeroizing<String>>,
}

#[derive(Serialize, Deserialize)]
struct AgentStateSerde {
    agent_id: String,
    token: Option<String>,
}

impl AgentState {
    pub fn load() -> Option<Self> {
        let path = Self::path();
        let data = fs::read_to_string(path).ok()?;
        let deserialized: AgentStateSerde = serde_json::from_str(&data).ok()?;
        Some(Self {
            agent_id: deserialized.agent_id,
            token: deserialized.token.map(Zeroizing::new),
        })
    }

    pub fn save(&self) -> std::io::Result<()> {
        let path = Self::path();
        let data = serde_json::to_string_pretty(&AgentStateSerde {
            agent_id: self.agent_id.clone(),
            token: self.token.as_ref().map(|t| t.to_string()),
        })?;
        fs::write(path, data)
    }

    fn path() -> PathBuf {
        std::env::current_dir().unwrap_or_default().join(STATE_FILE)
    }
}
