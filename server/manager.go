package server

import (
	"context"
	"fmt"

	"sc-player-service/i12e/config"
)

type Manager struct {
	HTTPServer *HTTPServer
}

func New(ctx context.Context, cfg *config.Config) (*Manager, error) {
	httpServer, err := NewHTTPServer(cfg)
	if err != nil {
		return nil, fmt.Errorf("init_http_server -> %w", err)
	}

	return &Manager{
		HTTPServer: httpServer,
	}, nil
}

func (a *Manager) Start() {
	a.HTTPServer.Run()
}

func (a *Manager) Stop(ctx context.Context) error {
	err := a.HTTPServer.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop_http_server -> %w", err)
	}

	return nil
}
