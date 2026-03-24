use crate::config::AgentConfig;
use anyhow::{anyhow, Context};
use rand::Rng;
use reqwest::Client;
use std::time::Duration;

/// Result of a transport POST attempt.
pub struct PostResult {
    pub body: Vec<u8>,
    /// Number of retry attempts (0 = no retry).
    pub retries: usize,
    /// Backoff used before the final successful attempt.
    pub last_backoff_ms: Option<u64>,
}

/// HTTP(s) transport responsible for sending/receiving messages to the teamserver.
///
/// This is kept separate from the core agent bootstrap logic so we can easily
/// swap in alternative transports (TCP, DNS, etc.).
pub struct HttpsTransport {
    client: Client,
    base_url: String,
    max_retries: usize,
    initial_backoff: Duration,
}

impl HttpsTransport {
    /// Constructs a new HTTPS transport using the provided agent configuration.
    pub fn new(cfg: &AgentConfig) -> anyhow::Result<Self> {
        let mut builder = Client::builder().timeout(Duration::from_secs(20));

        let mut base_url = cfg.server_url.clone();
        if !base_url.to_lowercase().starts_with("https://") {
            base_url = format!("https://{}", base_url.trim_start_matches("http://"));
        }

        if let Some(ca_cert) = cfg.ca_cert_path.as_deref() {
            let cert = std::fs::read(ca_cert).context("reading CA certificate")?;
            let cert = reqwest::Certificate::from_pem(&cert).context("parse CA certificate")?;
            builder = builder.add_root_certificate(cert);
        }

        let client = builder
            .https_only(true)
            .danger_accept_invalid_certs(true)
            .build()
            .context("building https client")?;

        Ok(Self {
            client,
            base_url,
            max_retries: cfg.transport_max_retries,
            initial_backoff: Duration::from_millis(cfg.transport_backoff_ms),
        })
    }

    /// Posts JSON to the given path, returning the response body as bytes.
    ///
    /// Implements a basic retry/backoff strategy for transient errors.
    pub async fn post_json<T: serde::Serialize + ?Sized>(
        &self,
        path: &str,
        payload: &T,
    ) -> anyhow::Result<PostResult> {
        let url = format!(
            "{}/{}",
            self.base_url.trim_end_matches('/'),
            path.trim_start_matches('/')
        );

        for attempt in 0..=self.max_retries {
            let resp = self.client.post(&url).json(payload).send().await;

            match resp {
                Ok(r) if r.status().is_success() => {
                    let body = r.bytes().await.context("reading response bytes")?;
                    let backoff = if attempt == 0 {
                        None
                    } else {
                        Some(self.calculate_backoff(attempt - 1).as_millis() as u64)
                    };
                    return Ok(PostResult {
                        body: body.to_vec(),
                        retries: attempt,
                        last_backoff_ms: backoff,
                    });
                }
                Ok(r) if r.status().is_client_error() => {
                    return Err(anyhow!("client error: {}", r.status()));
                }
                Ok(r) => {
                    // Server error (5xx) - retry.
                    if attempt >= self.max_retries {
                        return Err(anyhow!("server error: {}", r.status()));
                    }
                }
                Err(e) => {
                    if attempt >= self.max_retries {
                        return Err(anyhow!("network error: {}", e));
                    }
                }
            }

            let backoff = self.calculate_backoff(attempt);
            tokio::time::sleep(backoff).await;
        }

        Err(anyhow!("failed to post after retries"))
    }

    fn calculate_backoff(&self, attempt: usize) -> Duration {
        // Exponential backoff with jitter.
        let base = self.initial_backoff.as_millis() as f64;
        let expo = base * (2f64.powi(attempt as i32));
        let jitter = rand::thread_rng().gen_range(0.0..base);
        Duration::from_millis((expo + jitter) as u64)
    }
}
