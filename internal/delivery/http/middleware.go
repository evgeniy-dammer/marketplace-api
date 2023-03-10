package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	log "github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

// userIdentity validate access token.
func (d *Delivery) userIdentity(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	header := ginCtx.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(ginCtx, http.StatusUnauthorized, ErrEmptyAuthHeader)

		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 { //nolint:gomnd
		NewErrorResponse(ginCtx, http.StatusUnauthorized, ErrInvalidAuthHeader)

		return
	}

	userID, err := d.ucAuthentication.AuthenticationParseToken(ctx, headerParts[1])
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusUnauthorized, err)

		return
	}

	ginCtx.Set(userCtx, userID)
}

// getUserID returns user id from authorization context.
func (d *Delivery) getUserID(ginCtx *gin.Context) (string, error) {
	userID, exists := ginCtx.Get(userCtx)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrUserIsNotFound)

		return "", ErrUserIsNotFound
	}

	idString, exists := userID.(string)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrInvalidUserID)

		return "", ErrUserIsNotFound
	}

	return idString, nil
}

// getUserRole returns users role from authorization context.
func (d *Delivery) getUserRole(ginCtx *gin.Context) (string, error) {
	ctx := context.New(ginCtx)

	userID, exists := ginCtx.Get(userCtx)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrUserIsNotFound)

		return "", ErrUserIsNotFound
	}

	idString, exists := userID.(string)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrInvalidUserID)

		return "", ErrUserIsNotFound
	}

	role, err := d.ucAuthentication.AuthenticationGetUserRole(ctx, idString)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrRoleIsNotFound)

		return "", ErrUserIsNotFound
	}

	return role, nil
}

// corsMiddleware middleware.
func (d *Delivery) corsMiddleware() gin.HandlerFunc {
	return cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS, UPDATE, PATCH",
		RequestHeaders:  "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token", //nolint:lll
		ExposedHeaders:  "",
		MaxAge:          maxAge * time.Second,
		Credentials:     false,
		ValidateHeaders: true,
	})
}

// Authorize determines if current subject has been authorized to take an action on an object.
func (d *Delivery) Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, err := d.getUserRole(ctx)
		if err != nil {
			return
		}

		enforced, err := enforce(userRole, obj, act, adapter)
		if err != nil {
			NewErrorResponse(ctx, http.StatusInternalServerError, err)

			return
		}

		if !enforced {
			NewErrorResponse(ctx, http.StatusUnauthorized, ErrAccessDenied)

			return
		}

		ctx.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to create enforcer: %w", err)
	}

	if err = enforcer.LoadPolicy(); err != nil {
		return false, fmt.Errorf("failed to load policy: %w", err)
	}

	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false, fmt.Errorf("failed enforcing: %w", err)
	}

	return ok, nil
}

func Tracer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := opentracing.SpanFromContext(ctx.Request.Context())

		if span == nil {
			span = StartSpanWithHeader(&ctx.Request.Header, "rest-request-"+ctx.Request.Method, ctx.Request.Method, ctx.Request.URL.Path) //nolint:lll
		}

		defer span.Finish()

		ctx.Request = ctx.Request.WithContext(opentracing.ContextWithSpan(ctx.Request.Context(), span))

		if traceID, ok := span.Context().(jaeger.SpanContext); ok {
			ctx.Header("uber-trace-id", traceID.TraceID().String())
		}

		ctx.Next()

		ext.HTTPStatusCode.Set(span, uint16(ctx.Writer.Status()))

		if len(ctx.Errors) == 0 {
			log.Info("", getContextFields(ctx)...)
		}
	}
}

func StartSpanWithHeader(header *http.Header, operationName, method, path string) opentracing.Span {
	var wireContext opentracing.SpanContext

	if header != nil {
		wireContext, _ = opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*header))
	}

	return StartSpanWithParent(wireContext, operationName, method, path)
}

// StartSpanWithParent will start a new span with a parent span.
func StartSpanWithParent(parent opentracing.SpanContext, operationName, method, path string) opentracing.Span {
	options := []opentracing.StartSpanOption{
		opentracing.Tag{Key: ext.SpanKindRPCServer.Key, Value: ext.SpanKindRPCServer.Value},
		opentracing.Tag{Key: string(ext.HTTPMethod), Value: method},
		opentracing.Tag{Key: string(ext.HTTPUrl), Value: path},
	}
	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	return opentracing.StartSpan(operationName, options...)
}
