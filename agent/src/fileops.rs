/// File operations for agent-side file upload and download handling.
///
/// This module provides high-level APIs for uploading files to the operator
/// and downloading files from the operator using encrypted chunks.
use crate::filetransfer::{FileDecryptor, FileEncryptor};
use crate::protocol::{FileUploadRequest};
use crate::transport::https::HttpsTransport;
use anyhow::{anyhow, Context, Result};
use log::info;
use std::path::Path;

/// Uploads a file from the agent to the operator.
///
/// The file is split into encrypted chunks (256 KiB each) and sent sequentially.
/// Each chunk is encrypted with AES-256-GCM and authenticated with HMAC-SHA256.
pub async fn upload_file<P: AsRef<Path>>(
    transport: &HttpsTransport,
    agent_id: &str,
    token: &str,
    file_path: P,
    encryption_key: [u8; 32],
    hmac_key: [u8; 32],
) -> Result<()> {
    let file_path = file_path.as_ref();
    info!(
        "starting file upload: {} (agent: {})",
        file_path.display(),
        agent_id
    );

    let encryptor = FileEncryptor::new(encryption_key, hmac_key);
    let session_id = uuid::Uuid::new_v4().to_string();

    let chunks = encryptor
        .encrypt_file(file_path, session_id.clone())
        .context("failed to encrypt file")?;

    let total_chunks = chunks.len();
    info!(
        "file encrypted into {} chunks (session: {})",
        total_chunks, session_id
    );

    for (idx, chunk) in chunks.iter().enumerate() {
        info!(
            "uploading chunk {}/{} (session: {})",
            idx + 1,
            total_chunks,
            session_id
        );

        let request = FileUploadRequest {
            agent_id: agent_id.to_string(),
            token: token.to_string(),
            chunk: chunk.clone(),
        };

        let response = transport
            .upload_file_chunk(&request)
            .await
            .context(format!("failed to upload chunk {}", idx))?;

        if response.status != "ok" {
            return Err(anyhow!(
                "upload chunk {} failed: {}",
                idx,
                response.message.unwrap_or_default()
            ));
        }

        info!("chunk {} uploaded successfully", idx + 1);
    }

    info!(
        "file upload completed: {} (session: {})",
        file_path.display(),
        session_id
    );
    Ok(())
}

/// Downloads a file from the operator to the agent.
///
/// Fetches encrypted chunks sequentially from the operator and assembles them
/// into a complete file after decryption and verification.
pub async fn download_file<P: AsRef<Path>>(
    transport: &HttpsTransport,
    agent_id: &str,
    token: &str,
    file_id: &str,
    output_path: P,
    encryption_key: [u8; 32],
    hmac_key: [u8; 32],
) -> Result<()> {
    let output_path = output_path.as_ref();
    info!(
        "starting file download: {} (file_id: {}, agent: {})",
        output_path.display(),
        file_id,
        agent_id
    );

    let decryptor = FileDecryptor::new(encryption_key, hmac_key);
    let mut chunks = Vec::new();

    // First, fetch chunk 0 to determine total chunk count
    let first_chunk_response = transport
        .download_file_chunk(agent_id, token, file_id, 0)
        .await
        .context("failed to download first chunk")?;

    let total_chunks = first_chunk_response.chunk.total_chunks;
    chunks.push(first_chunk_response.chunk);
    info!(
        "file has {} chunks (file_id: {})",
        total_chunks, file_id
    );

    // Fetch remaining chunks
    for chunk_num in 1..total_chunks {
        info!("downloading chunk {}/{}", chunk_num + 1, total_chunks);

        let response = transport
            .download_file_chunk(agent_id, token, file_id, chunk_num)
            .await
            .context(format!("failed to download chunk {}", chunk_num))?;

        chunks.push(response.chunk);
    }

    info!(
        "all chunks downloaded, assembling file: {}",
        output_path.display()
    );

    decryptor
        .assemble_file(chunks, output_path)
        .context("failed to assemble downloaded file")?;

    info!(
        "file download completed: {} (file_id: {})",
        output_path.display(),
        file_id
    );
    Ok(())
}

/// Gets the default encryption keys for file transfer.
///
/// In production, these should be derived from the agent's token or a secure key exchange.
/// For now, this generates consistent keys based on the agent ID.
pub fn get_file_transfer_keys(agent_id: &str) -> ([u8; 32], [u8; 32]) {
    use sha2::{Digest, Sha256};

    // Derive keys from agent_id
    let mut hasher = Sha256::new();
    hasher.update(agent_id.as_bytes());
    hasher.update(b"encryption_key");
    let mut encryption_key = [0u8; 32];
    encryption_key.copy_from_slice(&hasher.finalize());

    let mut hasher = Sha256::new();
    hasher.update(agent_id.as_bytes());
    hasher.update(b"hmac_key");
    let mut hmac_key = [0u8; 32];
    hmac_key.copy_from_slice(&hasher.finalize());

    (encryption_key, hmac_key)
}
