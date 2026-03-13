import { useState } from "react";
import { Button } from "../components/Button";
import { Card } from "../components/Card";
import { TextInput } from "../components/TextInput";
import { login } from "../lib/api";

export function AuthCard({
  token,
  onToken,
}: {
  token: string;
  onToken: (next: string) => void;
}) {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string>("");

  const isAuthenticated = token.length > 0;

  return (
    <Card title="Authentication">
      {!isAuthenticated ? (
        <div className="space-y-2">
          <div>
            <label className="block text-sm">Username</label>
            <TextInput value={username} onChange={(e) => setUsername(e.target.value)} />
          </div>
          <div>
            <label className="block text-sm">Password</label>
            <TextInput
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          {error ? <p className="text-sm text-red-300">{error}</p> : null}
          <Button
            variant="primary"
            className="mt-3"
            onClick={async () => {
              setError("");
              try {
                const nextToken = await login(username, password);
                onToken(nextToken);
              } catch (e) {
                setError(e instanceof Error ? e.message : "login failed");
              }
            }}
          >
            Login
          </Button>
        </div>
      ) : (
        <div className="space-y-2">
          <p className="text-sm">Authenticated. Token stored locally.</p>
          <Button
            variant="danger"
            className="mt-2"
            onClick={() => {
              onToken("");
            }}
          >
            Logout
          </Button>
        </div>
      )}
    </Card>
  );
}
