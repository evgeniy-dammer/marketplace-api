package model

import (
	"context"
	"net/http"
	"time"
)

// Server is an http server
type Server struct {
	httpServer *http.Server
}

// Run starts the http server
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	
	// start listening
	return s.httpServer.ListenAndServe()
}

// Shutdown stops the http server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
