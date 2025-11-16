package main

import (
	"context"
	"fmt"

	"sc-player-service/i12e/config"
)

type App struct {
	HTTPServer *HTTPServer
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	httpServer, err := NewHTTPServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("init_http_server -> %w", err)
	}

	return &App{
		HTTPServer: httpServer,
	}, nil
}

func (a *App) Start() {
	a.HTTPServer.Run()
}

func (a *App) Stop(ctx context.Context) error {
	err := a.HTTPServer.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop_http_server -> %w", err)
	}

	return nil
}
