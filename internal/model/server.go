package model

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

// Run starts the http server.
func (s *Server) Run(srvConfig ServerConfig) error {
	s.httpServer = &http.Server{
		Addr:           ":" + srvConfig.Port,
		Handler:        srvConfig.Handler,
		MaxHeaderBytes: srvConfig.MaxHeaderBytes,
		ReadTimeout:    time.Duration(srvConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(srvConfig.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(srvConfig.IdleTimeout) * time.Second,
	}

	return errors.Wrap(s.httpServer.ListenAndServe(), "unable to start listening")
}

// Shutdown stops the http server.
func (s *Server) Shutdown(ctx context.Context) error {
	return errors.Wrap(s.httpServer.Shutdown(ctx), "unable to stop listening")
}
