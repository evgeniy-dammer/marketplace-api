package model

import (
	"net/http"
)

// DBConfig is a database config.
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig is a server config.
type ServerConfig struct {
	Handler        http.Handler
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	MaxHeaderBytes int
}
