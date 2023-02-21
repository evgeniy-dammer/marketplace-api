package model

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var errStructHasNoValues = errors.New("update structure has no values")

// MyError my custom error.
type MyError struct {
	Message string `json:"message"`
}

// ResponseData is an alerts response data.
type ResponseData struct {
	User   User   `json:"user"`
	Tokens Tokens `json:"tokens"`
}

// StatusResponse is a status response data.
type StatusResponse struct {
	Status string `json:"status"`
}

// NewErrorResponse is a response with error.
func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, MyError{Message: message})
}
