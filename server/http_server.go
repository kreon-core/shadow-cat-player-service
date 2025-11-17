package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	tul "github.com/kreon-core/shadow-cat-common"
	"github.com/kreon-core/shadow-cat-common/logc"

	"sc-player-service/helper"
	"sc-player-service/infrastructure/config"
	"sc-player-service/middleware"
)

const (
	srvHost              = ""
	srvPost              = 8080
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

	useCORS(r)

	r.Use(chiMW.Recoverer)

	r.Use(chiMW.CleanPath)
	r.Use(chiMW.StripSlashes)

	r.Use(chiMW.RealIP)
	r.Use(chiMW.RequestID)
	r.Use(middleware.RequestLogger)

	r.Use(chiMW.Timeout(srvGatewayTimeout))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { helper.PlainText(w, http.StatusOK, "healthy") })

	r.Route("/api/v1", LoadRoutes(container))

	host := tul.OrElse(cfg.Host, srvHost)
	port := tul.OrElse(cfg.Port, srvPost)

	return &HTTPServer{
		Server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", host, port),
			Handler:           r,
			ReadTimeout:       tul.OrElse(cfg.ReadTimeout, srvReadTimeout),
			ReadHeaderTimeout: tul.OrElse(cfg.ReadHeaderTimeout, srvReadHeaderTimeout),
			WriteTimeout:      tul.OrElse(cfg.WriteTimeout, srvWriteTimeout),
			IdleTimeout:       tul.OrElse(cfg.IdleTimeout, srvIdleTimeout),
		},
	}
}

func useCORS(r chi.Router) {
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-Real-IP",
			"X-Request-ID",
		},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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
