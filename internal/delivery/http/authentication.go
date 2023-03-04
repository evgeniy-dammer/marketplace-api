package http

import (
	"net/http"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// signIn login a user in the system.
func (d *Delivery) signIn(ctx *gin.Context) {
	var input user.SignInInput

	var tokens domain.Tokens

	var user user.User

	if err := ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	user, tokens, err := d.ucAuthentication.AuthenticationGenerateToken("", input.Phone, input.Password)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, domain.ResponseData{User: user, Tokens: tokens})
}

// signUp register a user in the system.
func (d *Delivery) signUp(ctx *gin.Context) {
	var input user.User

	if err := ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	userID, err := d.ucAuthentication.AuthenticationCreateUser(input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, map[string]interface{}{"id": userID})
}

// refresh refreshes token.
func (d *Delivery) refresh(ctx *gin.Context) {
	var input domain.RefreshToken

	var tokens domain.Tokens

	var user user.User

	// parsing JSON body
	if err := ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	headerParts := strings.Split(input.Authorization, " ")

	if len(headerParts) != 2 { //nolint:gomnd
		domain.NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")

		return
	}

	userID, err := d.ucAuthentication.AuthenticationParseToken(headerParts[1])
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())

		return
	}

	// check if user exists and create the token
	user, tokens, err = d.ucAuthentication.AuthenticationGenerateToken(userID, "", "")

	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, domain.ResponseData{User: user, Tokens: tokens})
}
