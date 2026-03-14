import React from "react";

export function DemoSettings() {
  return (
    <div style={{ color: "rgba(235,235,245,0.72)", lineHeight: 1.4 }}>
      <div style={{ fontWeight: 900, letterSpacing: "0.06em", textTransform: "uppercase" }}>Demo Settings</div>
      <p>Visual-only settings mockups.</p>
      <div style={{ display: "grid", gap: 10, maxWidth: 520 }}>
        <label>
          <div style={{ fontSize: 12, marginBottom: 6, color: "rgba(235,235,245,0.72)" }}>Theme</div>
          <select disabled style={{ width: "100%", padding: 10, borderRadius: 12, background: "rgba(255,255,255,0.04)", color: "rgba(255,255,255,0.9)", border: "1px solid rgba(255,255,255,0.12)" }}>
            <option>Obsidian / Crimson</option>
          </select>
        </label>
        <label>
          <div style={{ fontSize: 12, marginBottom: 6, color: "rgba(235,235,245,0.72)" }}>Density</div>
          <input disabled value="High" style={{ width: "100%", padding: 10, borderRadius: 12, background: "rgba(255,255,255,0.04)", color: "rgba(255,255,255,0.9)", border: "1px solid rgba(255,255,255,0.12)" }} />
        </label>
      </div>
    </div>
  );
}

