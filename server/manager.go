package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"sc-player-service/i12e/config"
)

type Manager struct {
	Container  *Container
	HTTPServer *HTTPServer
	PlayerDB   *pgxpool.Pool
}

func New(_ context.Context, cfg *config.Config) (*Manager, error) {
	// container := NewContainer()
	httpServer, err := NewHTTPServer(&cfg.HTTP)
	if err != nil {
		return nil, fmt.Errorf("init_http_server -> %w", err)
	}

	return &Manager{
		HTTPServer: httpServer,
	}, nil
}

func (a *Manager) Start() {
	var wg sync.WaitGroup
	wg.Go(func() { a.HTTPServer.Run() })
	wg.Wait()
}

func (a *Manager) Stop(ctx context.Context) error {
	err := a.HTTPServer.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop_http_server -> %w", err)
	}

	return nil
}
