package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
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
func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)

	if header == "" {
		model.NewErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")

		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 { //nolint:gomnd
		model.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")

		return
	}

	// parse token and return user id
	userID, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

		return
	}

	ctx.Set(userCtx, userID)
}

// getUserId returns user id from authorization context.
func (h *Handler) getUserIDAndRole(ctx *gin.Context) (string, string, error) {
	userID, exists := ctx.Get(userCtx)

	if !exists {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, "user is not found")

		return "", "", errUserIsNotFound
	}

	idString, exists := userID.(string)

	if !exists {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, "user id is of invalid type")

		return "", "", errUserIsNotFound
	}

	role, err := h.services.GetUserRole(idString)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, "role is not found")

		return "", "", errUserIsNotFound
	}

	return idString, role, nil
}

// corsMiddleware middleware.
func (h *Handler) corsMiddleware() gin.HandlerFunc {
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
func (h *Handler) Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, userRole, err := h.getUserIDAndRole(ctx)
		if err != nil {
			return
		}

		enforced, err := enforce(userRole, obj, act, adapter)
		if err != nil {
			model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

			return
		}

		if !enforced {
			model.NewErrorResponse(ctx, http.StatusForbidden, "forbidden")

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
