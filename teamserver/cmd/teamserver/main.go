package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/config"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/db"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	logger := log.New(os.Stdout, "[teamserver] ", log.LstdFlags|log.Lmsgprefix)

	logger.Printf("starting teamserver (version=%s)", version)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := connectWithRetry(ctx, cfg.DatabaseURL, 30*time.Second)
	if err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}
	defer pool.Close()

	httpServer := server.New(cfg, logger, pool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if os.Getenv("REDFORGE_CLI_DASHBOARD") == "1" {
		go runDashboard(ctx, pool)
	}

	go func() {
		if err := httpServer.Listen(ctx); err != nil && err != context.Canceled {
			logger.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	logger.Println("shutdown signal received, shutting down...")

	cancel()
	// Give the server time to shut down cleanly.
	time.Sleep(500 * time.Millisecond)

	logger.Println("stopped")
}

func connectWithRetry(ctx context.Context, databaseURL string, timeout time.Duration) (*pgxpool.Pool, error) {
	deadline := time.Now().Add(timeout)
	for {
		pool, err := db.Init(ctx, databaseURL)
		if err == nil {
			return pool, nil
		}
		if time.Now().After(deadline) {
			return nil, err
		}
		// Wait a bit and retry (common case: postgres is still starting)
		time.Sleep(1 * time.Second)
	}
}
