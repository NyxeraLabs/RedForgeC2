package config

import "testing"

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("REDFORGE_JWT_SECRET", "test-secret")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/redforge")
	t.Setenv("REDFORGE_JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when REDFORGE_JWT_SECRET is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/redforge")
	t.Setenv("REDFORGE_JWT_SECRET", "test-secret")
	t.Setenv("REDFORGE_PORT", "")
	t.Setenv("REDFORGE_ADMIN_USER", "")
	t.Setenv("REDFORGE_ADMIN_PASS", "")
	t.Setenv("REDFORGE_TOKEN_EXPIRY_MIN", "")
	t.Setenv("REDFORGE_TLS_CERT_FILE", "")
	t.Setenv("REDFORGE_TLS_KEY_FILE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config load to succeed, got: %v", err)
	}

	if cfg.Port != "9080" {
		t.Fatalf("expected default port %q, got %q", "9080", cfg.Port)
	}
	if cfg.AdminUsername != "admin" {
		t.Fatalf("expected default admin username %q, got %q", "admin", cfg.AdminUsername)
	}
	if cfg.AdminPassword != "redforge-admin" {
		t.Fatalf("expected default admin password %q, got %q", "redforge-admin", cfg.AdminPassword)
	}
	if cfg.TokenExpiryMins != 60 {
		t.Fatalf("expected default token expiry %d, got %d", 60, cfg.TokenExpiryMins)
	}
}

func TestLoad_TokenExpiryParsing(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/redforge")
	t.Setenv("REDFORGE_JWT_SECRET", "test-secret")

	t.Setenv("REDFORGE_TOKEN_EXPIRY_MIN", "15")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config load to succeed, got: %v", err)
	}
	if cfg.TokenExpiryMins != 15 {
		t.Fatalf("expected token expiry %d, got %d", 15, cfg.TokenExpiryMins)
	}

	t.Setenv("REDFORGE_TOKEN_EXPIRY_MIN", "0")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("expected config load to succeed, got: %v", err)
	}
	if cfg.TokenExpiryMins != 60 {
		t.Fatalf("expected token expiry to remain default %d, got %d", 60, cfg.TokenExpiryMins)
	}

	t.Setenv("REDFORGE_TOKEN_EXPIRY_MIN", "-5")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("expected config load to succeed, got: %v", err)
	}
	if cfg.TokenExpiryMins != 60 {
		t.Fatalf("expected token expiry to remain default %d, got %d", 60, cfg.TokenExpiryMins)
	}

	t.Setenv("REDFORGE_TOKEN_EXPIRY_MIN", "not-a-number")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("expected config load to succeed, got: %v", err)
	}
	if cfg.TokenExpiryMins != 60 {
		t.Fatalf("expected token expiry to remain default %d, got %d", 60, cfg.TokenExpiryMins)
	}
}

