import React, { useEffect, useState } from "react";
import { adminCreateUser, adminListUsers, AdminUser } from "../lib/api";
import { useToasts } from "../components/ToastProvider";

export function AdminUsersPage() {
  const toasts = useToasts();
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [username, setUsername] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"admin" | "operator" | "observer">("operator");

  async function refresh() {
    try {
      const list = await adminListUsers();
      setUsers(list);
    } catch (e) {
      toasts.push({ kind: "bad", title: "User list failed", message: String(e) });
    }
  }

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function create() {
    try {
      await adminCreateUser({ username, password, role, display_name: displayName });
      toasts.push({ kind: "good", title: "User created", message: username });
      setUsername("");
      setDisplayName("");
      setPassword("");
      setRole("operator");
      await refresh();
    } catch (e) {
      toasts.push({ kind: "bad", title: "Create failed", message: String(e) });
    }
  }

  return (
    <div className="page">
      <div className="page-head">
        <div>
          <div className="page-title">User Management</div>
          <div className="page-sub">Admin-only: create users and assign roles.</div>
        </div>
        <button className="btn btn-ghost" onClick={refresh}>
          Refresh
        </button>
      </div>

      <div className="grid-2">
        <section className="panel">
          <div className="panel-header">
            <h2>Create User</h2>
          </div>
          <div className="panel-body">
            <div className="form-row">
              <label>Username</label>
              <input value={username} onChange={(e) => setUsername(e.target.value)} />
            </div>
            <div className="form-row">
              <label>Display Name</label>
              <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} />
            </div>
            <div className="form-row">
              <label>Password (min 10 chars)</label>
              <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
            </div>
            <div className="form-row">
              <label>Role</label>
              <select value={role} onChange={(e) => setRole(e.target.value as any)}>
                <option value="admin">Admin</option>
                <option value="operator">Operator</option>
                <option value="observer">Observer</option>
              </select>
            </div>
            <button className="btn btn-primary" onClick={create} disabled={!username || password.length < 10}>
              Create
            </button>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Users</h2>
          </div>
          <div className="panel-body">
            <table className="table">
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Role</th>
                  <th>Display</th>
                  <th>Updated</th>
                </tr>
              </thead>
              <tbody>
                {users.map((u) => (
                  <tr key={u.username}>
                    <td>{u.username}</td>
                    <td>{u.role}</td>
                    <td>{u.display_name || "-"}</td>
                    <td>{new Date(u.updated_at).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </div>
  );
}

