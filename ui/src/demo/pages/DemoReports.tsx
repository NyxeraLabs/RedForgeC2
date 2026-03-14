import React from "react";

export function DemoReports() {
  return (
    <div style={{ color: "rgba(235,235,245,0.72)", lineHeight: 1.4 }}>
      <div style={{ fontWeight: 900, letterSpacing: "0.06em", textTransform: "uppercase" }}>Demo Reports</div>
      <p>This is a placeholder page for visuals (tables/cards/charts). No data is fetched in demo mode.</p>
      <ul>
        <li>Engagement summary card</li>
        <li>Timeline chart</li>
        <li>Export buttons (disabled)</li>
      </ul>
    </div>
  );
}

