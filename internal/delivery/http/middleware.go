package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
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

// getUserId returns user id from authorization context.
func (d *Delivery) getUserIDAndRole(ginCtx *gin.Context) (string, string, error) {
	ctx := context.New(ginCtx)

	userID, exists := ginCtx.Get(userCtx)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrUserIsNotFound)

		return "", "", ErrUserIsNotFound
	}

	idString, exists := userID.(string)
	if !exists {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrInvalidUserID)

		return "", "", ErrUserIsNotFound
	}

	role, err := d.ucAuthentication.AuthenticationGetUserRole(ctx, idString)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, ErrRoleIsNotFound)

		return "", "", ErrUserIsNotFound
	}

	return idString, role, nil
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
		_, userRole, err := d.getUserIDAndRole(ctx)
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
