package http

import (
	"errors"

	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
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

func getContextFields(ctx *gin.Context) []zap.Field {
	fields := []zap.Field{
		zap.Int("status", ctx.Writer.Status()),
		zap.String("method", ctx.Request.Method),
		zap.String("path", ctx.Request.URL.Path),
		zap.String("query", ctx.Request.URL.RawQuery),
		zap.String("ip", ctx.ClientIP()),
		zap.String("user-agent", ctx.Request.UserAgent()),
	}

	if span := opentracing.SpanFromContext(ctx.Request.Context()); span != nil {
		if jaegerSpan, ok := span.Context().(jaeger.SpanContext); ok {
			fields = append(fields, zap.Stringer("traceID", jaegerSpan.TraceID()))
		}
	}

	return fields
}
