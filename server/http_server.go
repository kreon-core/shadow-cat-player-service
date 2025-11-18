package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/kreon-core/shadow-cat-common/logc"
	"github.com/kreon-core/shadow-cat-common/mwc"
	"github.com/kreon-core/shadow-cat-common/resc"
	"github.com/kreon-core/shadow-cat-common/utlc"

	"sc-player-service/infrastructure/config"
	"sc-player-service/middleware"
)

const (
	srvReadTimeout       = 15 * time.Second
	srvReadHeaderTimeout = 15 * time.Second
	srvWriteTimeout      = 15 * time.Second
	srvIdleTimeout       = 120 * time.Second
	srvGatewayTimeout    = 60 * time.Second
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(cfg *config.HTTP, container *Container) *HTTPServer {
	r := chi.NewRouter()

	r.Use(mwc.CORS(nil))

	r.Use(chiMW.Recoverer)

	r.Use(chiMW.CleanPath)
	r.Use(chiMW.StripSlashes)

	r.Use(chiMW.RealIP)
	r.Use(chiMW.RequestID)

	r.Use(middleware.RequestLogger)

	r.Use(chiMW.Timeout(srvGatewayTimeout))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		resc.PlainText(w, http.StatusOK, "healthy")
	})

	r.Route("/api/v1", LoadRoutes(container))

	return &HTTPServer{
		Server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           r,
			ReadTimeout:       utlc.OrElse(cfg.ReadTimeout, srvReadTimeout),
			ReadHeaderTimeout: utlc.OrElse(cfg.ReadHeaderTimeout, srvReadHeaderTimeout),
			WriteTimeout:      utlc.OrElse(cfg.WriteTimeout, srvWriteTimeout),
			IdleTimeout:       utlc.OrElse(cfg.IdleTimeout, srvIdleTimeout),
		},
	}
}

func (s *HTTPServer) Run() {
	logc.Info().Str("address", s.Addr).Msg("HTTP server is listening and serving")
	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logc.Fatal().Err(err).Msg("HTTP server failed to start")
	}
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		logc.Error().Err(err).Msg("HTTP server shutdown failed")
		return err
	}

	logc.Info().Msg("HTTP server stopped gracefully")
	return nil
}
