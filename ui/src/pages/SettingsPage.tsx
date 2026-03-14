import React, { useEffect, useState } from "react";
import { getApiBase, setApiBase } from "../lib/storage";

export function SettingsPage() {
  const [darkMode, setDarkMode] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [serverUrl, setServerUrl] = useState("");
  const [savedUrl, setSavedUrl] = useState<string | null>(null);
  const [status, setStatus] = useState<string>("");

  useEffect(() => {
    const apiBase = getApiBase();
    setSavedUrl(apiBase);
    setServerUrl(apiBase || "");
  }, []);

  function saveServerUrl() {
    setApiBase(serverUrl.trim() || null);
    setSavedUrl(serverUrl.trim() || null);
    setStatus("Saved server URL (refresh page if needed)");
  }

  function resetServerUrl() {
    setServerUrl("");
    setApiBase(null);
    setSavedUrl(null);
    setStatus("Reset to default server URL");
  }

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">Settings</div>
          <div className="page-sub">Client preferences (mocked). Server-side settings coming soon.</div>
        </div>
      </div>

      <main className="main">
        <section className="panel">
          <div className="panel-header">
            <h2>UI Preferences</h2>
          </div>
          <div className="panel-body">
            <div className="form-row">
              <label>
                <input
                  type="checkbox"
                  checked={darkMode}
                  onChange={(e) => setDarkMode(e.target.checked)}
                />
                &nbsp;Enable dark theme (mock)
              </label>
            </div>
            <div className="form-row">
              <label>
                <input
                  type="checkbox"
                  checked={autoRefresh}
                  onChange={(e) => setAutoRefresh(e.target.checked)}
                />
                &nbsp;Auto-refresh dashboards (mock)
              </label>
            </div>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Connection</h2>
          </div>
          <div className="panel-body">
            <div className="form-row">
              <label>Teamserver URL</label>
              <input
                value={serverUrl}
                onChange={(e) => setServerUrl(e.target.value)}
                placeholder="https://localhost:9080"
              />
            </div>
            <div className="form-row">
              <button onClick={saveServerUrl}>Save</button>
              <button style={{ marginLeft: 8 }} onClick={resetServerUrl}>
                Reset
              </button>
            </div>
            <div style={{ marginTop: 14, color: "rgba(235, 235, 245, 0.7)" }}>
              {savedUrl ? (
                <span>
                  Current URL: <strong>{savedUrl}</strong>
                </span>
              ) : (
                <span>Using default server URL (based on current host).</span>
              )}
            </div>
            {status ? <div style={{ marginTop: 8 }}>{status}</div> : null}
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Notes</h2>
          </div>
          <div className="panel-body">
            <div style={{ color: "rgba(235, 235, 245, 0.7)" }}>
              These settings are currently stored locally in the browser. Server-backed configuration will be added in later milestones.
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
