package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds teamserver runtime configuration.
type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AdminUsername   string
	AdminPassword   string
	TokenExpiryMins int
	TLSCertFile     string
	TLSKeyFile      string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	port := os.Getenv("REDFORGE_PORT")
	if port == "" {
		port = "9080"
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("REDFORGE_JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("REDFORGE_JWT_SECRET is required")
	}

	adminUser := os.Getenv("REDFORGE_ADMIN_USER")
	if adminUser == "" {
		adminUser = "admin"
	}

	adminPass := os.Getenv("REDFORGE_ADMIN_PASS")
	if adminPass == "" {
		adminPass = "redforge-admin"
	}

	expires := 60
	if v := os.Getenv("REDFORGE_TOKEN_EXPIRY_MIN"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			expires = parsed
		}
	}

	return &Config{
		Port:            port,
		DatabaseURL:     db,
		JWTSecret:       jwtSecret,
		AdminUsername:   adminUser,
		AdminPassword:   adminPass,
		TokenExpiryMins: expires,
		TLSCertFile:     os.Getenv("REDFORGE_TLS_CERT_FILE"),
		TLSKeyFile:      os.Getenv("REDFORGE_TLS_KEY_FILE"),
	}, nil
}
