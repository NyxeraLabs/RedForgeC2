package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load([]string{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ListenAddr == "" || cfg.BaseURL == "" {
		t.Fatalf("expected default listen/base url, got addr=%q base=%q", cfg.ListenAddr, cfg.BaseURL)
	}
	if cfg.EventLog == "" {
		t.Fatalf("expected default event log path")
	}
}

func TestLoadFromFileThenOverrideWithFlags(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "teamserver.json")
	if err := os.WriteFile(cfgPath, []byte(`{
  "listen_addr": "127.0.0.1:9999",
  "base_url": "http://127.0.0.1:9999",
  "log_json": false,
  "event_log": "`+filepath.Join(dir, "events.jsonl")+`"
}`), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := Load([]string{"-config", cfgPath, "-addr", "127.0.0.1:1234"})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ListenAddr != "127.0.0.1:1234" {
		t.Fatalf("expected flag override addr, got %q", cfg.ListenAddr)
	}
	if cfg.LogJSON != false {
		t.Fatalf("expected log_json false from file")
	}
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("REDFORGE_ADDR", "127.0.0.1:7777")
	t.Setenv("REDFORGE_LOG_JSON", "false")
	t.Setenv("REDFORGE_EVENT_LOG", "data/test-events.jsonl")

	cfg, err := Load([]string{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ListenAddr != "127.0.0.1:7777" {
		t.Fatalf("expected env addr override, got %q", cfg.ListenAddr)
	}
	if cfg.LogJSON != false {
		t.Fatalf("expected env log_json false")
	}
	if cfg.EventLog != "data/test-events.jsonl" {
		t.Fatalf("expected env event_log override, got %q", cfg.EventLog)
	}
}

func TestLoadRejectsInvalidLogJSONFlag(t *testing.T) {
	_, err := Load([]string{"-log-json", "notabool"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

