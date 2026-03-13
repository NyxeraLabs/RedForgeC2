package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/config"
	"github.com/NyxeraLabs/RedForgeC2/teamserver/internal/server"
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

	httpServer := server.New(cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
