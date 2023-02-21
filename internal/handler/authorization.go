package handler

import (
	"net/http"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// signIn login a user in the system.
func (h *Handler) signIn(ctx *gin.Context) {
	var input model.SignInInput

	var tokens model.Tokens

	var user model.User

	if err := ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	user, tokens, err := h.services.Authorization.GenerateToken("", input.Phone, input.Password)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, model.ResponseData{User: user, Tokens: tokens})
}

// signUp register a user in the system.
func (h *Handler) signUp(ctx *gin.Context) {
	var input model.User

	if err := ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, map[string]interface{}{"id": userID})
}

// refresh refreshes token.
func (h *Handler) refresh(ctx *gin.Context) {
	var input model.RefreshToken

	var tokens model.Tokens

	var user model.User

	// parsing JSON body
	if err := ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	headerParts := strings.Split(input.Authorization, " ")

	if len(headerParts) != 2 { //nolint:gomnd
		model.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")

		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

		return
	}

	// check if user exists and create the token
	user, tokens, err = h.services.Authorization.GenerateToken(userID, "", "")

	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, model.ResponseData{User: user, Tokens: tokens})
}
