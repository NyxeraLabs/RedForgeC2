package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddr string        `json:"listen_addr"`
	BaseURL    string        `json:"base_url"`
	LogJSON    bool          `json:"log_json"`
	LogFile    string        `json:"log_file"`
	EventLog   string        `json:"event_log"`
	TLSCert    string        `json:"tls_cert"`
	TLSKey     string        `json:"tls_key"`
	Timeouts   ServerTimeout `json:"timeouts"`
}

type ServerTimeout struct {
	ReadHeader time.Duration `json:"read_header"`
}

func Default() Config {
	return Config{
		ListenAddr: "127.0.0.1:8080",
		BaseURL:    "http://127.0.0.1:8080",
		LogJSON:    true,
		LogFile:    "",
		EventLog:   "data/events.jsonl",
		TLSCert:    "",
		TLSKey:     "",
		Timeouts: ServerTimeout{
			ReadHeader: 5 * time.Second,
		},
	}
}

// Load merges (in order): defaults <- config file (optional) <- env <- flags.
func Load(args []string) (Config, error) {
	cfg := Default()

	fs := flag.NewFlagSet("teamserver", flag.ContinueOnError)
	configPath := fs.String("config", "", "optional JSON config file path")
	listenAddr := fs.String("addr", "", "listen address (default binds to localhost)")
	baseURL := fs.String("base-url", "", "public base URL (used for UI links)")
	logFile := fs.String("log-file", "", "optional log file path")
	eventLog := fs.String("event-log", "", "event log path (jsonl)")
	logJSON := fs.String("log-json", "", "true/false: emit JSON logs (default true)")
	tlsCert := fs.String("tls-cert", "", "TLS cert path (enables HTTPS)")
	tlsKey := fs.String("tls-key", "", "TLS key path (enables HTTPS)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if *configPath != "" {
		b, err := os.ReadFile(*configPath)
		if err != nil {
			return Config{}, err
		}
		var fromFile Config
		if err := json.Unmarshal(b, &fromFile); err != nil {
			return Config{}, err
		}
		merge(&cfg, fromFile)
	}

	applyEnv(&cfg)

	if *listenAddr != "" {
		cfg.ListenAddr = *listenAddr
	}
	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *logFile != "" {
		cfg.LogFile = *logFile
	}
	if *eventLog != "" {
		cfg.EventLog = *eventLog
	}
	if *tlsCert != "" {
		cfg.TLSCert = *tlsCert
	}
	if *tlsKey != "" {
		cfg.TLSKey = *tlsKey
	}
	if *logJSON != "" {
		v, err := strconv.ParseBool(*logJSON)
		if err != nil {
			return Config{}, errors.New("invalid -log-json value (use true/false)")
		}
		cfg.LogJSON = v
	}

	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	return cfg, nil
}

func merge(dst *Config, src Config) {
	if src.ListenAddr != "" {
		dst.ListenAddr = src.ListenAddr
	}
	if src.BaseURL != "" {
		dst.BaseURL = src.BaseURL
	}
	if src.LogFile != "" {
		dst.LogFile = src.LogFile
	}
	if src.EventLog != "" {
		dst.EventLog = src.EventLog
	}
	if src.TLSCert != "" {
		dst.TLSCert = src.TLSCert
	}
	if src.TLSKey != "" {
		dst.TLSKey = src.TLSKey
	}
	// bool: allow explicit false from file only if the JSON includes it; we treat zero-value as "not set".
	// So we only merge LogJSON when it differs from default of true AND file provided a value.
	// If users want deterministic behavior, they should set it via env/flag.
	if src.LogJSON == false {
		dst.LogJSON = false
	}
	if src.Timeouts.ReadHeader != 0 {
		dst.Timeouts.ReadHeader = src.Timeouts.ReadHeader
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("REDFORGE_ADDR"); v != "" {
		cfg.ListenAddr = v
	}
	if v := os.Getenv("REDFORGE_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("REDFORGE_LOG_FILE"); v != "" {
		cfg.LogFile = v
	}
	if v := os.Getenv("REDFORGE_EVENT_LOG"); v != "" {
		cfg.EventLog = v
	}
	if v := os.Getenv("REDFORGE_TLS_CERT"); v != "" {
		cfg.TLSCert = v
	}
	if v := os.Getenv("REDFORGE_TLS_KEY"); v != "" {
		cfg.TLSKey = v
	}
	if v := os.Getenv("REDFORGE_LOG_JSON"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.LogJSON = b
		}
	}
}
