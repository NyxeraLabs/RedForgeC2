import { useState } from "react";

export default function App() {
  const [connected, setConnected] = useState(false);

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="p-4 border-b border-slate-800">
        <h1 className="text-2xl font-semibold">RedForgeC2 Operator Console</h1>
        <p className="text-sm text-slate-300">Alpha UI scaffold (connects to teamserver)</p>
      </header>

      <main className="p-6">
        <section className="rounded-lg border border-slate-800 bg-slate-900 p-6">
          <h2 className="text-xl font-semibold mb-2">Status</h2>
          <p>
            Teamserver connection: <span className="font-semibold">{connected ? "connected" : "disconnected"}</span>
          </p>
          <button
            onClick={() => setConnected((v) => !v)}
            className="mt-4 rounded bg-indigo-600 px-4 py-2 text-sm font-medium hover:bg-indigo-500"
          >
            {connected ? "Disconnect" : "Connect"}
          </button>
        </section>
      </main>
    </div>
  );
}
