/// File transfer module for encrypted chunked uploads/downloads over HTTPS.
///
/// This module provides functionality for splitting files into encrypted chunks
/// and reassembling them, using AES-256-GCM for encryption and HMAC-SHA256 for
/// authentication. All transfers are done over HTTPS with TLS 1.2+.
use crate::protocol::FileChunk;
use anyhow::{anyhow, Context, Result};
use aes_gcm::aead::{Aead, KeyInit};
use aes_gcm::{Aes256Gcm, Nonce};
use base64::{engine::general_purpose, Engine as _};
use hmac::{Hmac, Mac};
use rand::RngCore;
use sha2::Sha256;
use std::fs;
use std::path::Path;
use uuid::Uuid;

const CHUNK_SIZE: usize = 256 * 1024; // 256 KiB chunks
const NONCE_SIZE: usize = 12; // 96 bits for GCM

/// File encryption manager for splitting and encrypting files into chunks.
pub struct FileEncryptor {
    encryption_key: [u8; 32], // AES-256 key
    hmac_key: [u8; 32],       // HMAC key
}

impl FileEncryptor {
    /// Creates a new file encryptor with the given keys.
    /// Keys should be 32 bytes (256 bits) for AES-256.
    pub fn new(encryption_key: [u8; 32], hmac_key: [u8; 32]) -> Self {
        Self {
            encryption_key,
            hmac_key,
        }
    }

    /// Creates a new file encryptor with random keys (for session-specific encryption).
    pub fn new_with_random_keys() -> Self {
        let mut encryption_key = [0u8; 32];
        let mut hmac_key = [0u8; 32];
        let mut rng = rand::thread_rng();
        rng.fill_bytes(&mut encryption_key);
        rng.fill_bytes(&mut hmac_key);
        Self {
            encryption_key,
            hmac_key,
        }
    }

    /// Splits a file into encrypted chunks ready for transmission.
    pub fn encrypt_file<P: AsRef<Path>>(
        &self,
        file_path: P,
        _session_id: String,
    ) -> Result<Vec<FileChunk>> {
        let file_path = file_path.as_ref();
        let file_data = fs::read(file_path)
            .context(format!("failed to read file: {:?}", file_path))?;
        let filename = file_path
            .file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("file")
            .to_string();

        let total_size = file_data.len();
        let total_chunks = (total_size + CHUNK_SIZE - 1) / CHUNK_SIZE;

        let mut chunks = Vec::new();

        for (chunk_num, chunk_data) in file_data.chunks(CHUNK_SIZE).enumerate() {
            let encrypted = self.encrypt_chunk(chunk_data, chunk_num, total_chunks, &filename)?;
            chunks.push(encrypted);
        }

        Ok(chunks)
    }

    /// Encrypts a single chunk of data.
    fn encrypt_chunk(
        &self,
        data: &[u8],
        chunk_number: usize,
        total_chunks: usize,
        filename: &str,
    ) -> Result<FileChunk> {
        // Generate random nonce
        let mut nonce_bytes = [0u8; NONCE_SIZE];
        let mut rng = rand::thread_rng();
        rng.fill_bytes(&mut nonce_bytes);

        // Create cipher and encrypt
        let cipher = Aes256Gcm::new(&self.encryption_key.into());
        let nonce = Nonce::from_slice(&nonce_bytes);

        let ciphertext = cipher
            .encrypt(nonce, data)
            .map_err(|e| anyhow!("encryption failed: {}", e))?;

        // Calculate HMAC over chunk_number || total_chunks || nonce || ciphertext
        let mut hmac_input = Vec::new();
        hmac_input.extend_from_slice(&(chunk_number as u64).to_le_bytes());
        hmac_input.extend_from_slice(&(total_chunks as u64).to_le_bytes());
        hmac_input.extend_from_slice(&nonce_bytes);
        hmac_input.extend_from_slice(&ciphertext);

        let mut mac = <Hmac<Sha256> as KeyInit>::new_from_slice(&self.hmac_key)
            .map_err(|e| anyhow!("HMAC key error: {}", e))?;
        <_ as Mac>::update(&mut mac, &hmac_input);
        let hmac_result = mac.finalize();

        Ok(FileChunk {
            session_id: Uuid::new_v4().to_string(),
            chunk_number,
            total_chunks,
            filename: filename.to_string(),
            original_size: data.len(),
            nonce: general_purpose::STANDARD.encode(&nonce_bytes),
            data: general_purpose::STANDARD.encode(&ciphertext),
            hmac: general_purpose::STANDARD.encode(hmac_result.into_bytes()),
        })
    }
}

/// File decryption manager for assembling and decrypting file chunks.
pub struct FileDecryptor {
    encryption_key: [u8; 32],
    hmac_key: [u8; 32],
}

impl FileDecryptor {
    /// Creates a new file decryptor with the given keys.
    pub fn new(encryption_key: [u8; 32], hmac_key: [u8; 32]) -> Self {
        Self {
            encryption_key,
            hmac_key,
        }
    }

    /// Verifies and decrypts a chunk.
    pub fn decrypt_chunk(&self, chunk: &FileChunk) -> Result<Vec<u8>> {
        // Decode inputs
        let nonce_bytes = general_purpose::STANDARD
            .decode(&chunk.nonce)
            .context("failed to decode nonce")?;
        let ciphertext = general_purpose::STANDARD
            .decode(&chunk.data)
            .context("failed to decode ciphertext")?;
        let received_hmac = general_purpose::STANDARD
            .decode(&chunk.hmac)
            .context("failed to decode HMAC")?;

        // Verify HMAC
        let mut hmac_input = Vec::new();
        hmac_input.extend_from_slice(&(chunk.chunk_number as u64).to_le_bytes());
        hmac_input.extend_from_slice(&(chunk.total_chunks as u64).to_le_bytes());
        hmac_input.extend_from_slice(&nonce_bytes);
        hmac_input.extend_from_slice(&ciphertext);

        let mut mac = <Hmac<Sha256> as KeyInit>::new_from_slice(&self.hmac_key)
            .map_err(|e| anyhow!("HMAC key error: {}", e))?;
        <_ as Mac>::update(&mut mac, &hmac_input);

        mac.verify_slice(&received_hmac)
            .map_err(|_| anyhow!("HMAC verification failed - chunk may be corrupted or tampered"))?;

        // Decrypt
        let cipher = Aes256Gcm::new(&self.encryption_key.into());
        let nonce = Nonce::from_slice(&nonce_bytes);

        let plaintext = cipher
            .decrypt(nonce, &ciphertext[..])
            .map_err(|e| anyhow!("decryption failed: {}", e))?;

        Ok(plaintext)
    }

    /// Assembles chunks into a complete file.
    pub fn assemble_file<P: AsRef<Path>>(
        &self,
        chunks: Vec<FileChunk>,
        output_path: P,
    ) -> Result<()> {
        let output_path = output_path.as_ref();

        if chunks.is_empty() {
            return Err(anyhow!("no chunks provided"));
        }

        // Verify chunk integrity
        let expected_total = chunks[0].total_chunks;
        if chunks.len() != expected_total {
            return Err(anyhow!(
                "expected {} chunks, got {}",
                expected_total,
                chunks.len()
            ));
        }

        // Sort chunks by chunk_number
        let mut sorted_chunks = chunks;
        sorted_chunks.sort_by_key(|c| c.chunk_number);

        // Verify sequential numbering
        for (i, chunk) in sorted_chunks.iter().enumerate() {
            if chunk.chunk_number != i {
                return Err(anyhow!(
                    "chunk numbering is not sequential at position {}",
                    i
                ));
            }
        }

        // Decrypt and assemble
        let mut output = Vec::new();
        for chunk in sorted_chunks {
            let decrypted = self.decrypt_chunk(&chunk)?;
            output.extend_from_slice(&decrypted);
        }

        fs::write(output_path, &output)
            .context(format!("failed to write file: {:?}", output_path))?;

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use std::io::Write;
    use tempfile::TempDir;

    #[test]
    fn test_encrypt_decrypt_chunk() {
        let key = [42u8; 32];
        let hmac_key = [43u8; 32];
        let encryptor = FileEncryptor::new(key, hmac_key);
        let decryptor = FileDecryptor::new(key, hmac_key);

        let data = b"Hello, World! This is a test message.";
        let chunk = encryptor
            .encrypt_chunk(data, 0, 1, "test.txt")
            .expect("encryption failed");

        let decrypted = decryptor
            .decrypt_chunk(&chunk)
            .expect("decryption failed");
        assert_eq!(decrypted, data);
    }

    #[test]
    fn test_file_encrypt_decrypt() {
        let dir = TempDir::new().expect("failed to create temp dir");
        let file_path = dir.path().join("test.txt");
        let mut file = fs::File::create(&file_path).expect("failed to create file");
        file.write_all(b"Test file content for encryption")
            .expect("write failed");
        drop(file);

        let key = [42u8; 32];
        let hmac_key = [43u8; 32];
        let encryptor = FileEncryptor::new(key, hmac_key);
        let decryptor = FileDecryptor::new(key, hmac_key);

        let session_id = Uuid::new_v4().to_string();
        let chunks = encryptor
            .encrypt_file(&file_path, session_id)
            .expect("encryption failed");
        assert!(!chunks.is_empty());

        let output_path = dir.path().join("decrypted.txt");
        decryptor
            .assemble_file(chunks, &output_path)
            .expect("assembly failed");

        let original = fs::read(&file_path).expect("read original failed");
        let decrypted = fs::read(&output_path).expect("read decrypted failed");
        assert_eq!(original, decrypted);
    }

    #[test]
    fn test_hmac_verification() {
        let key = [42u8; 32];
        let hmac_key = [43u8; 32];
        let encryptor = FileEncryptor::new(key, hmac_key);
        let decryptor = FileDecryptor::new(key, hmac_key);

        let data = b"Important data";
        let mut chunk = encryptor
            .encrypt_chunk(data, 0, 1, "test.txt")
            .expect("encryption failed");

        // Tamper with the chunk
        chunk.data = general_purpose::STANDARD.encode(b"tampered");

        let result = decryptor.decrypt_chunk(&chunk);
        assert!(result.is_err());
    }
}
