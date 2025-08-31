package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wdsjk/avito-shop/internal/config"
)

type Server struct {
	addr string
}

func NewServer(cfg *config.Config) *Server {
	return &Server{addr: cfg.HTTPServer.Address}
}

func (s *Server) Start(cfg *config.Config, log *slog.Logger, handler http.Handler) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigint

		log.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.IdleTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("http server shutdown error", "error", err)
		}

		close(idleConnsClosed)
	}()

	log.Info("starting server", "address", s.addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	log.Info("server gracefully stopped")
	return nil
}
