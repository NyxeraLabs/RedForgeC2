import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { DemoLayout } from "./DemoLayout";
import { DemoDashboard } from "./pages/DemoDashboard";
import { DemoTelemetry } from "./pages/DemoTelemetry";
import { DemoMap } from "./pages/DemoMap";
import { DemoReports } from "./pages/DemoReports";
import { DemoSettings } from "./pages/DemoSettings";

export function DemoApp() {
  return (
    <Routes>
      <Route element={<DemoLayout />}>
        <Route index element={<DemoDashboard />} />
        <Route path="telemetry" element={<DemoTelemetry />} />
        <Route path="map" element={<DemoMap />} />
        <Route path="reports" element={<DemoReports />} />
        <Route path="settings" element={<DemoSettings />} />
        <Route path="*" element={<Navigate to="/demo" replace />} />
      </Route>
    </Routes>
  );
}

