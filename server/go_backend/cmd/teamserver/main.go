package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"redforgec2/server/go_backend/internal/auth"
	"redforgec2/server/go_backend/internal/config"
	"redforgec2/server/go_backend/internal/events"
	"redforgec2/server/go_backend/internal/httpapi"
	"redforgec2/server/go_backend/internal/logging"
	"redforgec2/server/go_backend/internal/store"
)

func main() {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	logger, closeLog, err := logging.NewLogger(logging.Options{JSON: cfg.LogJSON, LogFile: cfg.LogFile})
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}
	defer func() { _ = closeLog() }()

	ev, err := events.New(cfg.EventLog)
	if err != nil {
		logger.Error("event_log_init_failed", "err", err, "path", cfg.EventLog)
		os.Exit(2)
	}
	defer func() { _ = ev.Close() }()

	var authSvc *auth.Service
	if cfg.Auth.Enabled {
		adminUser := os.Getenv("REDFORGE_ADMIN_USER")
		adminPass := os.Getenv("REDFORGE_ADMIN_PASS")
		if cfg.Auth.AllowDevAdminFromEnv && (adminUser == "" || adminPass == "") {
			logger.Error("auth_enabled_missing_admin_env", "need", "REDFORGE_ADMIN_USER and REDFORGE_ADMIN_PASS")
			os.Exit(2)
		}
		var users []auth.User
		if cfg.Auth.AllowDevAdminFromEnv {
			hash, err := auth.HashPassword(adminPass)
			if err != nil {
				logger.Error("auth_hash_failed", "err", err)
				os.Exit(2)
			}
			users = append(users, auth.User{Username: adminUser, PasswordHash: hash, Role: auth.RoleAdmin})
		}
		svc, err := auth.New(true, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL, users)
		if err != nil {
			logger.Error("auth_init_failed", "err", err)
			os.Exit(2)
		}
		authSvc = svc
		logger.Info("auth_enabled", "token_ttl", cfg.Auth.TokenTTL.String(), "dev_admin_from_env", cfg.Auth.AllowDevAdminFromEnv)
	}

	st := store.New()
	router := httpapi.NewRouter(httpapi.Deps{Store: st, Logger: logger, Events: ev, Auth: authSvc})

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           router,
		ReadHeaderTimeout: cfg.Timeouts.ReadHeader,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},
	}

	logger.Info("teamserver_listen", "addr", cfg.ListenAddr, "base_url", cfg.BaseURL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if cfg.TLSCert != "" && cfg.TLSKey != "" {
			errCh <- srv.ListenAndServeTLS(cfg.TLSCert, cfg.TLSKey)
			return
		}
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown_signal")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server_error", "err", err)
			os.Exit(1)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("shutdown_failed", "err", err)
	} else {
		logger.Info("shutdown_complete")
	}

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server_error", "err", err)
			os.Exit(1)
		}
	default:
	}
}
