package http

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/token"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
)

// AuthResponse is an auth response data.
type AuthResponse struct {
	// User
	User user.User `json:"user"`
	// Tokens
	Tokens token.Tokens `json:"tokens"`
}

// StatusResponse is a status response data.
type StatusResponse struct {
	// Status message
	Status string `json:"status"`
}
