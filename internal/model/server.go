package model

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	maxHeaderBytes = 1 << 20 // 1MB
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

// Server is an http server.
type Server struct {
	httpServer *http.Server
}

// Run starts the http server.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	return errors.Wrap(s.httpServer.ListenAndServe(), "unable to start listening")
}

// Shutdown stops the http server.
func (s *Server) Shutdown(ctx context.Context) error {
	return errors.Wrap(s.httpServer.Shutdown(ctx), "unable to stop listening")
}
