package server

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Server is an http server.
type Server struct {
	httpServer *http.Server
}

// Config entity.
type Config struct {
	Handler        http.Handler
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
}

// Run starts the http server.
func (s *Server) Run(cfg Config) error {
	s.httpServer = &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        cfg.Handler,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(cfg.IdleTimeout) * time.Second,
	}

	return errors.Wrap(s.httpServer.ListenAndServe(), "unable to start listening")
}

// Shutdown stops the http server.
func (s *Server) Shutdown(ctx context.Context) error {
	return errors.Wrap(s.httpServer.Shutdown(ctx), "unable to stop listening")
}
