import { useMemo, useState } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import { TextInput } from "../components/TextInput";
import type { TaskResult } from "../lib/api";
import { createTask } from "../lib/api";

function decodeBase64ToBytes(b64: string): Uint8Array {
  const binary = atob(b64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);
  return bytes;
}

async function fileToBase64(file: File): Promise<string> {
  const reader = new FileReader();
  const promise = new Promise<string>((resolve, reject) => {
    reader.onerror = () => reject(new Error("failed to read file"));
    reader.onload = () => resolve(String(reader.result));
  });
  reader.readAsDataURL(file);
  const dataUrl = await promise;
  const idx = dataUrl.indexOf("base64,");
  if (idx === -1) throw new Error("unexpected file encoding");
  return dataUrl.slice(idx + "base64,".length);
}

export function FilesCard({
  isAuthenticated,
  token,
  selectedAgentId,
  results,
  onRefreshResults,
}: {
  isAuthenticated: boolean;
  token: string;
  selectedAgentId: string;
  results: TaskResult[];
  onRefreshResults: () => void;
}) {
  const [browsePath, setBrowsePath] = useState("/tmp");
  const [uploadRemotePath, setUploadRemotePath] = useState("/tmp/upload.bin");
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [downloadRemotePath, setDownloadRemotePath] = useState("/tmp/file.bin");
  const [lastLsTask, setLastLsTask] = useState<string>("");
  const [lastUploadTask, setLastUploadTask] = useState<string>("");
  const [lastDownloadTask, setLastDownloadTask] = useState<string>("");
  const [error, setError] = useState<string>("");

  const findResult = (taskId: string) => results.find((r) => r.task_id === taskId);

  const lsResult = useMemo(() => (lastLsTask ? findResult(lastLsTask) : undefined), [results, lastLsTask]);
  const uploadResult = useMemo(
    () => (lastUploadTask ? findResult(lastUploadTask) : undefined),
    [results, lastUploadTask]
  );
  const downloadResult = useMemo(
    () => (lastDownloadTask ? findResult(lastDownloadTask) : undefined),
    [results, lastDownloadTask]
  );

  return (
    <Card title="Files (v0)">
      {!isAuthenticated ? (
        <p className="text-sm text-slate-400">Login to manage files.</p>
      ) : !selectedAgentId ? (
        <p className="text-sm text-slate-400">Select an agent to manage files.</p>
      ) : (
        <div className="space-y-6">
          {error ? <p className="text-sm text-red-300">{error}</p> : null}

          <div className="rounded border border-slate-800 bg-slate-950 p-4 space-y-2">
            <div className="text-sm font-medium">Browse</div>
            <div className="flex gap-2">
              <TextInput value={browsePath} onChange={(e) => setBrowsePath(e.target.value)} />
              <Button
                onClick={async () => {
                  setError("");
                  try {
                    const resp = await createTask(token, selectedAgentId, "ls", [browsePath], 60);
                    setLastLsTask(resp.task_id);
                  } catch (e) {
                    setError(e instanceof Error ? e.message : "ls failed");
                  }
                }}
              >
                List
              </Button>
              <Button onClick={onRefreshResults}>Refresh</Button>
            </div>
            <pre className="max-h-56 overflow-y-auto rounded bg-slate-900 p-3 text-xs">
              {lsResult
                ? `${lsResult.error ? `ERROR: ${lsResult.error}\n\n` : ""}${lsResult.output}`
                : lastLsTask
                  ? "Waiting for result (refresh)..."
                  : "Run ls to see directory listing."}
            </pre>
          </div>

          <div className="rounded border border-slate-800 bg-slate-950 p-4 space-y-2">
            <div className="text-sm font-medium">Upload</div>
            <div className="grid gap-2 md:grid-cols-2">
              <div>
                <label className="block text-sm">Remote path</label>
                <TextInput
                  value={uploadRemotePath}
                  onChange={(e) => setUploadRemotePath(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm">Local file</label>
                <input
                  type="file"
                  className="w-full text-sm"
                  onChange={(e) => setUploadFile(e.target.files?.[0] || null)}
                />
              </div>
            </div>
            <div className="flex gap-2">
              <Button
                variant="primary"
                disabled={!uploadFile}
                onClick={async () => {
                  if (!uploadFile) return;
                  setError("");
                  try {
                    const b64 = await fileToBase64(uploadFile);
                    const resp = await createTask(
                      token,
                      selectedAgentId,
                      "upload",
                      [uploadRemotePath, b64],
                      120
                    );
                    setLastUploadTask(resp.task_id);
                  } catch (e) {
                    setError(e instanceof Error ? e.message : "upload failed");
                  }
                }}
              >
                Upload
              </Button>
              <Button onClick={onRefreshResults}>Refresh</Button>
            </div>
            <div className="text-xs text-slate-400">
              Upload uses agent command: <span className="font-mono">upload &lt;path&gt; &lt;base64&gt;</span>
            </div>
            <pre className="max-h-40 overflow-y-auto rounded bg-slate-900 p-3 text-xs">
              {uploadResult
                ? `${uploadResult.error ? `ERROR: ${uploadResult.error}\n\n` : ""}${uploadResult.output}`
                : lastUploadTask
                  ? "Waiting for result (refresh)..."
                  : "Select a file and upload."}
            </pre>
          </div>

          <div className="rounded border border-slate-800 bg-slate-950 p-4 space-y-2">
            <div className="text-sm font-medium">Download</div>
            <div className="flex gap-2">
              <TextInput
                value={downloadRemotePath}
                onChange={(e) => setDownloadRemotePath(e.target.value)}
              />
              <Button
                variant="primary"
                onClick={async () => {
                  setError("");
                  try {
                    const resp = await createTask(token, selectedAgentId, "download", [downloadRemotePath], 120);
                    setLastDownloadTask(resp.task_id);
                  } catch (e) {
                    setError(e instanceof Error ? e.message : "download failed");
                  }
                }}
              >
                Fetch
              </Button>
              <Button onClick={onRefreshResults}>Refresh</Button>
            </div>
            <div className="flex flex-wrap gap-2">
              <Button
                disabled={!downloadResult || downloadResult.status !== "success" || !downloadResult.output}
                onClick={() => {
                  if (!downloadResult?.output) return;
                  const bytes = decodeBase64ToBytes(downloadResult.output.trim());
                  const blob = new Blob([bytes], { type: "application/octet-stream" });
                  const url = URL.createObjectURL(blob);
                  const a = document.createElement("a");
                  a.href = url;
                  a.download = downloadRemotePath.split("/").pop() || "download.bin";
                  a.click();
                  URL.revokeObjectURL(url);
                }}
              >
                Save File
              </Button>
            </div>
            <pre className="max-h-40 overflow-y-auto rounded bg-slate-900 p-3 text-xs">
              {downloadResult
                ? `${downloadResult.error ? `ERROR: ${downloadResult.error}\n\n` : ""}${downloadResult.output.slice(0, 5000)}`
                : lastDownloadTask
                  ? "Waiting for result (refresh)..."
                  : "Fetch a file to get base64 output."}
            </pre>
          </div>
        </div>
      )}
    </Card>
  );
}

