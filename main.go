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
	"sc-player-service/server"
)

const (
	shutdownTimeout = 30 * time.Second
)

func main() {
	defer exit()

	// Initialize global logger
	logc.InitializeLogger()

	// Load configuration
	i12e.LoadEnvs()
	cfg, err := i12e.LoadConfigs()
	if err != nil {
		logc.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Create server context
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize server
	server, err := server.New(appCtx, cfg)
	if err != nil {
		logc.Fatal().Err(err).Msg("Failed to create server")
	}

	// Handle termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	// Wait for termination signal
	<-stop
	logc.Info().Msg("Shutting down application...")

	// Graceful shutdown
	shutdownCtx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	err = server.Stop(shutdownCtx)
	if err != nil {
		logc.Fatal().Err(err).Msg("Forced shutdown HTTP server failed")
	}

	logc.Info().Msg("Exiting application gracefully")
}

func exit() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "Application panicked: %+v\n", r)
		os.Exit(1)
	}
}
