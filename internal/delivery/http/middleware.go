package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	maxAge              = 300
)

var errUserIsNotFound = errors.New("user is not found")

// userIdentity validate access token.
func (d *Delivery) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)

	if header == "" {
		domain.NewErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")

		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 { //nolint:gomnd
		domain.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")

		return
	}

	// parse token and return user id
	userID, err := d.ucAuthentication.AuthenticationParseToken(headerParts[1])
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

		return
	}

	ctx.Set(userCtx, userID)
}

// getUserId returns user id from authorization context.
func (d *Delivery) getUserIDAndRole(ctx *gin.Context) (string, string, error) {
	userID, exists := ctx.Get(userCtx)

	if !exists {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, "user is not found")

		return "", "", errUserIsNotFound
	}

	idString, exists := userID.(string)

	if !exists {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, "user id is of invalid type")

		return "", "", errUserIsNotFound
	}

	role, err := d.ucAuthentication.AuthenticationGetUserRole(idString)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, "role is not found")

		return "", "", errUserIsNotFound
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
			domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

			return
		}

		if !enforced {
			domain.NewErrorResponse(ctx, http.StatusForbidden, "forbidden")

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
