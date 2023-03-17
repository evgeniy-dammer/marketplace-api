package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// userIdentity validate access token.
func (d *Delivery) userIdentity(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.userIdentity")
		defer span.End()

		ctx = context.New(ctxt)
	}

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

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getUserRole")
		defer span.End()

		ctx = context.New(ctxt)
	}

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

	role, err := d.ucAuthorization.AuthorizationGetUserRole(ctx, idString)
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
	return func(ginCtx *gin.Context) {
		userRole, err := d.getUserRole(ginCtx)
		if err != nil {
			return
		}

		enforced, err := enforce(userRole, obj, act, adapter)
		if err != nil {
			NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

			return
		}

		if !enforced {
			NewErrorResponse(ginCtx, http.StatusUnauthorized, ErrAccessDenied)

			return
		}

		ginCtx.Next()
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

func (d *Delivery) parseMetadata(ginCtx *gin.Context) (query.MetaData, error) {
	metaUserID, err := d.getUserID(ginCtx)
	if err != nil {
		return query.MetaData{}, err
	}

	return query.MetaData{
		UserID:         metaUserID,
		OrganizationID: ginCtx.Query(organizationQueryKey),
	}, nil
}
