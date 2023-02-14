package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

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
