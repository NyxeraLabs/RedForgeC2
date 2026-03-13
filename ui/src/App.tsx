import React from "react";
import { BrowserRouter, Link, Route, Routes, useNavigate } from "react-router-dom";
import { ToastProvider } from "./components/ToastProvider";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { LoginPage } from "./pages/LoginPage";
import { DashboardPage } from "./pages/DashboardPage";
import { AdminUsersPage } from "./pages/AdminUsersPage";
import { ProfilePage } from "./pages/ProfilePage";
import { getMe } from "./lib/api";
import { setToken } from "./lib/storage";

function Shell() {
  const nav = useNavigate();
  const [me, setMe] = React.useState<{ username: string; role: string } | null>(null);

  React.useEffect(() => {
    getMe()
      .then((m) => setMe({ username: m.username, role: m.role }))
      .catch(() => setMe(null));
  }, []);

  function logout() {
    setToken(null);
    nav("/login", { replace: true });
  }

  return (
    <div className="shell">
      <div className="shell-nav">
        <div className="shell-brand">REDFORGE-C2</div>
        {me ? <div className="shell-meta">{me.username} ({me.role})</div> : null}
        <Link className="shell-link" to="/">
          Dashboard
        </Link>
        <Link className="shell-link" to="/profile">
          Profile
        </Link>
        {me?.role === "admin" ? (
          <Link className="shell-link" to="/admin/users">
            Users
          </Link>
        ) : null}
        <button className="shell-link" onClick={logout}>
          Logout
        </button>
      </div>
      <div className="shell-main">
        <Routes>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/admin/users" element={<AdminUsersPage />} />
        </Routes>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <ToastProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedRoute />}>
            <Route path="/*" element={<Shell />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ToastProvider>
  );
}
