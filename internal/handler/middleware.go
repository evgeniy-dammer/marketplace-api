package handler

import (
	"errors"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"net/http"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// userIdentity validate access token
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		model.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		model.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// parse token and return user id
	userId, err := h.services.ParseToken(headerParts[1])

	if err != nil {
		model.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

// getUserId returns user id from authorization context
func (h *Handler) getUserIdAndRole(c *gin.Context) (string, string, error) {
	id, ok := c.Get(userCtx)

	if !ok {
		model.NewErrorResponse(c, http.StatusInternalServerError, "user is not found")
		return "", "", errors.New("user is not found")
	}

	idString, ok := id.(string)

	if !ok {
		model.NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return "", "", errors.New("user is not found")
	}

	role, err := h.services.GetUserRole(idString)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, "role is not found")
		return "", "", errors.New("role is not found")
	}

	return idString, role, nil
}

// corsMiddleware middleware
func (h *Handler) corsMiddleware() gin.HandlerFunc {
	return cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS, UPDATE, PATCH",
		RequestHeaders:  "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token",
		ExposedHeaders:  "",
		MaxAge:          300 * time.Second,
		Credentials:     false,
		ValidateHeaders: true,
	})
}
