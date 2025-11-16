package app

import (
	"context"

	"sc-player-service/i12e/config"
)

type HTTPServer struct {
	Config *config.Config
}

func NewHTTPServer(cfg *config.Config) (*HTTPServer, error) {
	return &HTTPServer{
		Config: cfg,
	}, nil
}

func (s *HTTPServer) Run() {
	// Placeholder for HTTP server logic
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	// Placeholder for HTTP server shutdown logic
	return nil
}
