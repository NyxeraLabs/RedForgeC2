use serde::{Deserialize, Serialize};
use std::{collections::BTreeMap, env, time::Duration};
use time::format_description::well_known::Rfc3339;
use time::OffsetDateTime;

#[derive(Debug, Serialize)]
struct RegisterRequest {
    #[serde(skip_serializing_if = "Option::is_none")]
    agent_id: Option<String>,
    os: String,
    arch: String,
    hostname: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    meta: Option<BTreeMap<String, String>>,
}

#[derive(Debug, Deserialize)]
struct RegisterResponse {
    agent_id: String,
    server_time: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
struct TelemetryEvent {
    timestamp: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    metrics: Option<BTreeMap<String, f64>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    labels: Option<BTreeMap<String, String>>,
}

#[derive(Debug, Serialize)]
struct TelemetryRequest {
    events: Vec<TelemetryEvent>,
}

#[derive(Debug, Deserialize)]
struct PollTasksResponse {
    tasks: Vec<Task>,
}

#[derive(Debug, Deserialize, Clone)]
#[allow(dead_code)]
struct Task {
    task_id: String,
    task_type: String,
    #[serde(default)]
    params: BTreeMap<String, serde_json::Value>,
    created_at: String,
    status: String,
}

#[derive(Debug, Serialize)]
struct SubmitResultRequest {
    #[serde(skip_serializing_if = "Option::is_none")]
    outcome: Option<BTreeMap<String, serde_json::Value>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    status: Option<String>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let server_url = env::var("REDFORGE_SERVER_URL").unwrap_or_else(|_| "http://127.0.0.1:8080".to_string());
    enforce_local_only(&server_url)?;

    let interval_ms: u64 = env::var("REDFORGE_INTERVAL_MS")
        .ok()
        .and_then(|v| v.parse().ok())
        .unwrap_or(1500);
    let once = env::var("REDFORGE_ONCE").ok().as_deref() == Some("1");

    let client = reqwest::Client::new();

    let req = RegisterRequest {
        agent_id: env::var("REDFORGE_AGENT_ID").ok(),
        os: env::consts::OS.to_string(),
        arch: env::consts::ARCH.to_string(),
        hostname: env::var("HOSTNAME").unwrap_or_else(|_| "lab-host".to_string()),
        meta: Some(BTreeMap::from([("mode".to_string(), "simulation".to_string())])),
    };

    let register: RegisterResponse = client
        .post(format!("{server_url}/api/v1/agents/register"))
        .json(&req)
        .send()
        .await?
        .error_for_status()?
        .json()
        .await?;

    eprintln!(
        "registered (simulation): agent_id={} server_time={}",
        register.agent_id, register.server_time
    );

    let mut counter: u64 = 0;
    loop {
        counter += 1;

        // Send mock telemetry
        let telemetry = TelemetryRequest {
            events: vec![TelemetryEvent {
                timestamp: now_rfc3339(),
                metrics: Some(BTreeMap::from([
                    ("cpu".to_string(), (counter % 100) as f64 / 100.0),
                    ("mem".to_string(), ((counter * 7) % 100) as f64 / 100.0),
                ])),
                labels: Some(BTreeMap::from([("source".to_string(), "rust_agent_sim".to_string())])),
            }],
        };

        client
            .post(format!("{server_url}/api/v1/agents/{}/telemetry", register.agent_id))
            .json(&telemetry)
            .send()
            .await?
            .error_for_status()?;

        // Poll tasks
        let poll: PollTasksResponse = client
            .get(format!("{server_url}/api/v1/agents/{}/tasks", register.agent_id))
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        for task in poll.tasks {
            let outcome = simulate_task(&task).await;
            let submit = SubmitResultRequest {
                outcome: Some(outcome),
                status: Some("simulated".to_string()),
            };

            client
                .post(format!(
                    "{server_url}/api/v1/agents/{}/tasks/{}/result",
                    register.agent_id, task.task_id
                ))
                .json(&submit)
                .send()
                .await?
                .error_for_status()?;
        }

        if once {
            break;
        }
        tokio::time::sleep(Duration::from_millis(interval_ms)).await;
    }

    Ok(())
}

async fn simulate_task(task: &Task) -> BTreeMap<String, serde_json::Value> {
    match task.task_type.as_str() {
        "echo" => {
            let msg = task
                .params
                .get("msg")
                .and_then(|v| v.as_str())
                .unwrap_or("");
            BTreeMap::from([("echo".to_string(), serde_json::Value::String(msg.to_string()))])
        }
        "sleep" => {
            let ms = task.params.get("ms").and_then(|v| v.as_u64()).unwrap_or(0);
            let clamped = ms.min(5_000);
            tokio::time::sleep(Duration::from_millis(clamped)).await;
            BTreeMap::from([
                ("slept_ms".to_string(), serde_json::Value::Number(clamped.into())),
                ("note".to_string(), serde_json::Value::String("clamped to 5000ms max".to_string())),
            ])
        }
        "collect_telemetry" => BTreeMap::from([(
            "note".to_string(),
            serde_json::Value::String("telemetry is sent periodically (mock)".to_string()),
        )]),
        other => BTreeMap::from([(
            "error".to_string(),
            serde_json::Value::String(format!("unknown task_type: {other}")),
        )]),
    }
}

fn enforce_local_only(server_url: &str) -> Result<(), Box<dyn std::error::Error>> {
    let allow_nonlocal = env::var("REDFORGE_UNSAFE_ALLOW_NONLOCAL").ok().as_deref() == Some("1");
    validate_server_url(server_url, allow_nonlocal).map_err(|e| e.into())
}

fn now_rfc3339() -> String {
    OffsetDateTime::now_utc()
        .format(&Rfc3339)
        .unwrap_or_else(|_| "1970-01-01T00:00:00Z".to_string())
}

fn validate_server_url(server_url: &str, allow_nonlocal: bool) -> Result<(), String> {
    if allow_nonlocal {
        return Ok(());
    }
    let url = reqwest::Url::parse(server_url).map_err(|e| format!("invalid server_url: {e}"))?;
    let host = url.host_str().unwrap_or_default();
    let is_local = host == "localhost" || host == "127.0.0.1" || host == "::1";
    if !is_local {
        return Err(format!(
            "refusing non-local server_url host={host}. Set REDFORGE_UNSAFE_ALLOW_NONLOCAL=1 for an isolated lab network."
        ));
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn validate_server_url_local_ok() {
        validate_server_url("http://127.0.0.1:8080", false).unwrap();
        validate_server_url("http://localhost:8080", false).unwrap();
    }

    #[test]
    fn validate_server_url_nonlocal_rejected() {
        let err = validate_server_url("http://example.com:8080", false).unwrap_err();
        assert!(err.contains("refusing non-local"));
    }

    #[test]
    fn validate_server_url_nonlocal_allowed_with_override() {
        validate_server_url("http://example.com:8080", true).unwrap();
    }
}
