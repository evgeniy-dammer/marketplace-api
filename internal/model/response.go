package model

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MyError
type MyError struct {
	Message string `json:"message"`
}

// ResponseData
type ResponseData struct {
	User   User   `json:"user"`
	Tokens Tokens `json:"tokens"`
}

// StatusResponse
type StatusResponse struct {
	Status string `json:"status"`
}

// NewErrorResponse is a response with error
func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, MyError{Message: message})
}
