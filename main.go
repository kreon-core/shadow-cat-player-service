package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kreon-core/shadow-cat-common/logc"

	"sc-player-service/i12e"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	defer exit()

	// Initialize global logger
	logc.InitializeLogger()

	// Load application configuration
	cfg, err := i12e.LoadConfigs()
	if err != nil {
		logc.Fatal("Failed to load configuration", err)
	}

	// Create server context
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize app
	app, err := NewApp(appCtx, cfg)
	if err != nil {
		logc.Fatal("Failed to create app", err)
	}

	// Handle termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go app.Start()

	// Wait for termination signal
	<-stop
	logc.Info("Shutting down application...")

	// Graceful shutdown
	shutdownCtx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	err = app.Stop(shutdownCtx)
	if err != nil {
		logc.Fatal("Forced shutdown HTTP server failed", err)
	}

	logc.Info("Exiting application gracefully")
}

func exit() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "application panicked: %+v\n", r)
		os.Exit(1)
	}
}
