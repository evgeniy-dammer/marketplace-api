package domain

import (
	"net/http"
)

// ServerConfig entity.
type ServerConfig struct {
	Handler        http.Handler
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
}
