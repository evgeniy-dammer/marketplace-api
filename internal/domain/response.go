package domain

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// MyError my custom error.
type MyError struct {
	Message string `json:"message"`
}

// ResponseData is an alerts response data.
type ResponseData struct {
	User   user.User `json:"user"`
	Tokens Tokens    `json:"tokens"`
}

// StatusResponse is a status response data.
type StatusResponse struct {
	Status string `json:"status"`
}

// NewErrorResponse is a response with error.
func NewErrorResponse(c *gin.Context, statusCode int, err error) {
	logger.Logger.Error(err.Error())

	c.AbortWithStatusJSON(statusCode, MyError{Message: err.Error()})
}
