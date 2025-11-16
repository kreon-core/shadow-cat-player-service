package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tul "github.com/kreon-core/shadow-cat-common"
	"github.com/kreon-core/shadow-cat-common/logc"

	"sc-player-service/i12e/config"
)

const (
	srvHost           = ""
	srvPost           = 8080
	srvReadTimeout    = 15 * time.Second
	srvWriteTimeout   = 15 * time.Second
	srvIdleTimeout    = 120 * time.Second
	srvGatewayTimeout = 60 * time.Second
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(cfg *config.HTTP) (*HTTPServer, error) {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(srvGatewayTimeout))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	host := tul.OrElse(cfg.Host, srvHost)
	port := tul.OrElse(cfg.Port, srvPost)

	return &HTTPServer{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      r,
			ReadTimeout:  tul.OrElse(cfg.ReadTimeout, srvReadTimeout),
			WriteTimeout: tul.OrElse(cfg.WriteTimeout, srvWriteTimeout),
			IdleTimeout:  tul.OrElse(cfg.IdleTimeout, srvIdleTimeout),
		},
	}, nil
}

func (s *HTTPServer) Run() {
	logc.Info("HTTP server is listening and serving", "address", s.Addr)

	err := s.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logc.Fatal("HTTP server failed to start", err)
	}
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		logc.Error("HTTP server shutdown failed", err)
		return err
	}

	logc.Info("HTTP server stopped gracefully")
	return nil
}
