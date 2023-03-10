package http

import (
	"net/http"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/token"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// signIn
// @Summary SignIn user method.
// @Description SignIn user method.
// @Tags authentication
// @Accept  json
// @Produce json
// @Param   input 	body 		user.SignInInput 	true  "Username and Password"
// @Success 200		{object}  	AuthResponse		true  "User data and Tokens"
// @Failure 400 	{object}    ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /signin/ [post].
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

// signUp
// @Summary SignUp user method.
// @Description SignUp user method.
// @Tags authentication
// @Accept  json
// @Produce json
// @Param   input 	body 		user.CreateUserInput 	true  "User data"
// @Success 200		{string}  	string					true  "User ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /signup/ [post].
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

// refresh
// @Summary Refresh token method.
// @Description Refresh token method.
// @Tags authentication
// @Accept  json
// @Produce json
// @Param   input 	body 		token.RefreshToken		true  "Refresh token"
// @Success 200		{object}  	AuthResponse			true  "User data and Tokens"
// @Failure 400 	{object}    ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /refresh/ [post].
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
