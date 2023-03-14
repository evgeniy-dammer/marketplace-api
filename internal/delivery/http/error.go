package http

import (
	"errors"

	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

var (
	ErrEmptyIDParam      = errors.New("empty id param")
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrEmptyAuthHeader   = errors.New("empty auth header")
	ErrUserIsNotFound    = errors.New("user is not found")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrRoleIsNotFound    = errors.New("role is not found")
	ErrAccessDenied      = errors.New("access denied")
)

// ErrorResponse my custom error.
type ErrorResponse struct {
	// Error message
	Message string `json:"message"`
}

// NewErrorResponse is a response with error.
func NewErrorResponse(c *gin.Context, statusCode int, err error) {
	logger.Logger.Error(err.Error())

	c.AbortWithStatusJSON(statusCode, ErrorResponse{Message: err.Error()})
}
