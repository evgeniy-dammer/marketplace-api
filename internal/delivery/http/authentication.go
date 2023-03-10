package http

import (
	"net/http"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/token"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// signIn login a user in the system.
func (d *Delivery) signIn(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	var input user.SignInInput

	var tokens token.Tokens

	var usr user.User

	if err := ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	usr, tokens, err := d.ucAuthentication.AuthenticationGenerateToken(ctx, "", input.Phone, input.Password)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, AuthResponse{User: usr, Tokens: tokens})
}

// signUp register a user in the system.
func (d *Delivery) signUp(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	var input user.CreateUserInput
	if err := ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	userID, err := d.ucAuthentication.AuthenticationCreateUser(ctx, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": userID})
}

// refresh refreshes token.
func (d *Delivery) refresh(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	var input token.RefreshToken

	var tokens token.Tokens

	var usr user.User

	if err := ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	headerParts := strings.Split(input.Authorization, " ")
	if len(headerParts) != 2 { //nolint:gomnd
		NewErrorResponse(ginCtx, http.StatusUnauthorized, ErrInvalidAuthHeader)

		return
	}

	userID, err := d.ucAuthentication.AuthenticationParseToken(ctx, headerParts[1])
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusUnauthorized, err)

		return
	}

	usr, tokens, err = d.ucAuthentication.AuthenticationGenerateToken(ctx, userID, "", "")
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, AuthResponse{User: usr, Tokens: tokens})
}
