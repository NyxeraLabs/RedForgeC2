import React, { useEffect, useState } from "react";
import { changePassword, getMe, updateProfile } from "../lib/api";
import { useToasts } from "../components/ToastProvider";
import styles from "./ProfilePage.module.css";

export function ProfilePage() {
  const toasts = useToasts();
  const [me, setMe] = useState<{ username: string; role: string; display_name?: string } | null>(null);
  const [displayName, setDisplayName] = useState("");
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");

  useEffect(() => {
    getMe()
      .then((m) => {
        setMe(m);
        setDisplayName(m.display_name || "");
      })
      .catch((e) => toasts.push({ kind: "bad", title: "Profile load failed", message: String(e) }));
  }, [toasts]);

  async function saveProfile() {
    try {
      await updateProfile(displayName);
      toasts.push({ kind: "good", title: "Profile updated" });
    } catch (e) {
      toasts.push({ kind: "bad", title: "Update failed", message: String(e) });
    }
  }

  async function savePassword() {
    try {
      await changePassword(oldPassword, newPassword);
      setOldPassword("");
      setNewPassword("");
      toasts.push({ kind: "good", title: "Password changed" });
    } catch (e) {
      toasts.push({ kind: "bad", title: "Password change failed", message: String(e) });
    }
  }

  return (
    <div className={`page ${styles.scope}`}>
      <div className="page-head">
        <div>
          <div className="page-title">Profile</div>
          <div className="page-sub">{me ? `${me.username} (${me.role})` : "…"}</div>
        </div>
      </div>

      <div className="grid-2">
        <section className="panel">
          <div className="panel-header">
            <h2>Profile Info</h2>
          </div>
          <div className="panel-body">
            <div className="form-row">
              <label>Display Name</label>
              <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} placeholder="Operator" />
            </div>
            <button className="btn btn-primary" onClick={saveProfile}>
              Save
            </button>
          </div>
        </section>

        <section className="panel">
          <div className="panel-header">
            <h2>Change Password</h2>
          </div>
          <div className="panel-body">
            <div className="form-row">
              <label>Old Password</label>
              <input type="password" value={oldPassword} onChange={(e) => setOldPassword(e.target.value)} />
            </div>
            <div className="form-row">
              <label>New Password</label>
              <input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} />
            </div>
            <button className="btn btn-primary" onClick={savePassword} disabled={!oldPassword || !newPassword}>
              Change
            </button>
          </div>
        </section>
      </div>
    </div>
  );
}
