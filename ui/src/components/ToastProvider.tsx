import React, { createContext, useCallback, useContext, useMemo, useState } from "react";

export type Toast = {
  id: string;
  title: string;
  message?: string;
  kind: "info" | "good" | "warn" | "bad";
  createdAt: number;
};

type ToastCtx = {
  push: (t: Omit<Toast, "id" | "createdAt"> & { ttlMs?: number }) => void;
};

const Ctx = createContext<ToastCtx | null>(null);

function uid() {
  return Math.random().toString(16).slice(2) + Date.now().toString(16);
}

export function ToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const remove = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const push = useCallback(
    (t: Omit<Toast, "id" | "createdAt"> & { ttlMs?: number }) => {
      const id = uid();
      const toast: Toast = { id, title: t.title, message: t.message, kind: t.kind, createdAt: Date.now() };
      setToasts((prev) => [toast, ...prev].slice(0, 6));
      const ttl = t.ttlMs ?? 4500;
      window.setTimeout(() => remove(id), ttl);
    },
    [remove]
  );

  const value = useMemo(() => ({ push }), [push]);

  return (
    <Ctx.Provider value={value}>
      {children}
      <div className="toasts" aria-live="polite" aria-relevant="additions">
        {toasts.map((t) => (
          <div key={t.id} className={`toast ${t.kind}`}>
            <div className="toast-title">{t.title}</div>
            {t.message ? <div className="toast-msg">{t.message}</div> : null}
          </div>
        ))}
      </div>
    </Ctx.Provider>
  );
}

export function useToasts() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error("useToasts must be used within ToastProvider");
  return ctx;
}

