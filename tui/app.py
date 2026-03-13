from __future__ import annotations

import json
import os
import time
from dataclasses import dataclass
from typing import Any

import requests
from textual import on
from textual.app import App, ComposeResult
from textual.containers import Container, Horizontal, Vertical
from textual.widgets import Button, Footer, Header, Input, Label, ListItem, ListView, Static, TextLog


@dataclass(frozen=True)
class Agent:
    agent_id: str
    hostname: str
    os: str
    arch: str
    version: str
    last_seen: str


class AgentItem(ListItem):
    def __init__(self, agent: Agent) -> None:
        super().__init__()
        self.agent = agent

    def compose(self) -> ComposeResult:
        short_id = self.agent.agent_id[:8]
        yield Label(f"{self.agent.hostname}  ({short_id})")


class RedForgeTUI(App):
    CSS = """
    Screen { background: #0b1220; }
    #status { padding: 1 2; color: #cbd5e1; }
    #left { width: 40%; border: tall #334155; }
    #right { width: 60%; border: tall #334155; }
    #controls { height: 3; }
    #cmd { width: 40%; }
    #args { width: 60%; }
    """

    BINDINGS = [
        ("r", "refresh", "Refresh"),
        ("q", "quit", "Quit"),
    ]

    def __init__(self) -> None:
        super().__init__()
        self.base_url = os.environ.get("REDFORGE_TEAMSERVER_URL", "http://localhost:9080").rstrip("/")
        self.username = os.environ.get("REDFORGE_USERNAME", "admin")
        self.password = os.environ.get("REDFORGE_PASSWORD", "redforge")
        self.token = os.environ.get("REDFORGE_TOKEN", "")
        self.selected_agent_id: str = ""

    def compose(self) -> ComposeResult:
        yield Header()
        yield Static("", id="status")
        with Horizontal():
            with Vertical(id="left"):
                yield Label("Agents", classes="title")
                yield ListView(id="agents")
                with Horizontal(id="controls"):
                    yield Button("Login", id="login", variant="primary")
                    yield Button("Refresh", id="refresh")
            with Vertical(id="right"):
                yield Label("Results", classes="title")
                yield TextLog(id="results", highlight=False, wrap=True)
                with Container():
                    yield Label("Command")
                    yield Input(value="ls", placeholder="command", id="cmd")
                    yield Label("Args (space separated)")
                    yield Input(value="/tmp", placeholder="args", id="args")
                with Horizontal():
                    yield Button("Send Task", id="send", variant="success")
                    yield Button("Refresh Results", id="refresh_results")
        yield Footer()

    def _set_status(self, msg: str) -> None:
        self.query_one("#status", Static).update(msg)

    def _auth_headers(self) -> dict[str, str]:
        return {"Authorization": f"Bearer {self.token}"}

    def _login(self) -> None:
        resp = requests.post(
            f"{self.base_url}/api/login",
            headers={"Content-Type": "application/json"},
            data=json.dumps({"username": self.username, "password": self.password}),
            timeout=10,
        )
        resp.raise_for_status()
        body = resp.json()
        self.token = body.get("token", "")
        if not self.token:
            raise RuntimeError("login response missing token")

    def _list_agents(self) -> list[Agent]:
        resp = requests.get(
            f"{self.base_url}/api/operator/agents",
            headers=self._auth_headers(),
            timeout=10,
        )
        resp.raise_for_status()
        items: list[dict[str, Any]] = resp.json()
        out: list[Agent] = []
        for a in items:
            out.append(
                Agent(
                    agent_id=a.get("agent_id", ""),
                    hostname=a.get("hostname", ""),
                    os=a.get("os", ""),
                    arch=a.get("arch", ""),
                    version=a.get("version", ""),
                    last_seen=a.get("last_seen", ""),
                )
            )
        return out

    def _create_task(self, agent_id: str, command: str, args: list[str]) -> str:
        resp = requests.post(
            f"{self.base_url}/api/operator/task",
            headers={"Content-Type": "application/json", **self._auth_headers()},
            data=json.dumps(
                {
                    "agent_id": agent_id,
                    "command": command,
                    "args": args,
                    "timeout_seconds": 60,
                }
            ),
            timeout=10,
        )
        resp.raise_for_status()
        body = resp.json()
        return body.get("task_id", "")

    def _list_results(self, agent_id: str) -> list[dict[str, Any]]:
        resp = requests.get(
            f"{self.base_url}/api/operator/results",
            params={"agent_id": agent_id},
            headers=self._auth_headers(),
            timeout=10,
        )
        resp.raise_for_status()
        return resp.json()

    def on_mount(self) -> None:
        self._set_status(f"teamserver={self.base_url}  token={'set' if self.token else 'unset'}")
        self.set_interval(5.0, self._tick, pause=False)

    def _tick(self) -> None:
        if not self.token:
            return
        try:
            self.refresh_agents()
            if self.selected_agent_id:
                self.refresh_results()
        except Exception:
            pass

    def refresh_agents(self) -> None:
        agents_view = self.query_one("#agents", ListView)
        agents_view.clear()
        for a in self._list_agents():
            agents_view.append(AgentItem(a))
        self._set_status(
            f"teamserver={self.base_url}  agents={len(agents_view.children)}  selected={self.selected_agent_id[:8] if self.selected_agent_id else '-'}"
        )

    def refresh_results(self) -> None:
        if not self.selected_agent_id:
            return
        results = self._list_results(self.selected_agent_id)
        log = self.query_one("#results", TextLog)
        log.clear()
        log.write(f"agent={self.selected_agent_id}\nrefreshed={time.strftime('%Y-%m-%d %H:%M:%S')}\n")
        log.write(json.dumps(results, indent=2))

    def action_refresh(self) -> None:
        if not self.token:
            self._set_status("not authenticated; press Login")
            return
        try:
            self.refresh_agents()
            self.refresh_results()
        except Exception as e:
            self._set_status(f"refresh failed: {e}")

    @on(Button.Pressed, "#login")
    def on_login_pressed(self) -> None:
        try:
            self._login()
            self._set_status("authenticated; refreshing agents...")
            self.refresh_agents()
        except Exception as e:
            self._set_status(f"login failed: {e}")

    @on(Button.Pressed, "#refresh")
    def on_refresh_pressed(self) -> None:
        self.action_refresh()

    @on(Button.Pressed, "#refresh_results")
    def on_refresh_results_pressed(self) -> None:
        if not self.token:
            self._set_status("not authenticated; press Login")
            return
        try:
            self.refresh_results()
        except Exception as e:
            self._set_status(f"results failed: {e}")

    @on(ListView.Selected, "#agents")
    def on_agent_selected(self, event: ListView.Selected) -> None:
        if isinstance(event.item, AgentItem):
            self.selected_agent_id = event.item.agent.agent_id
            self._set_status(f"selected agent {self.selected_agent_id[:8]}")
            if self.token:
                try:
                    self.refresh_results()
                except Exception:
                    pass

    @on(Button.Pressed, "#send")
    def on_send_pressed(self) -> None:
        if not self.token:
            self._set_status("not authenticated; press Login")
            return
        if not self.selected_agent_id:
            self._set_status("select an agent first")
            return

        command = self.query_one("#cmd", Input).value.strip()
        args_raw = self.query_one("#args", Input).value.strip()
        args = [a for a in args_raw.split(" ") if a]
        if not command:
            self._set_status("command required")
            return

        try:
            task_id = self._create_task(self.selected_agent_id, command, args)
            self._set_status(f"task queued: {task_id}")
        except Exception as e:
            self._set_status(f"task failed: {e}")


if __name__ == "__main__":
    RedForgeTUI().run()

