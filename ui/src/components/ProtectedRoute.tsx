import React from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { getToken } from "../lib/storage";

export function ProtectedRoute() {
  const token = getToken();
  const loc = useLocation();
  if (!token) return <Navigate to="/login" replace state={{ from: loc.pathname }} />;
  return <Outlet />;
}

