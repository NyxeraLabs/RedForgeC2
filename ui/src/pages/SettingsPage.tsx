import React, { useState } from "react";

export function SettingsPage() {
  const [darkMode, setDarkMode] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);

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
            <div style={{ marginTop: 14, color: "rgba(235, 235, 245, 0.7)" }}>
              These settings are currently stored locally in the browser. Server-backed configuration will be added in later milestones.
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}
