package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// signIn login a user in the system
func (h *Handler) signIn(c *gin.Context) {
	var input model.SignInInput
	var tokens model.Tokens
	var user model.User

	// parsing JSON body
	if err := c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// check if user exists and create the token
	user, tokens, err := h.services.Authorization.GenerateToken("", input.Phone, input.Password)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// if all is ok
	c.JSON(http.StatusOK, model.ResponseData{User: user, Tokens: tokens})
}

// signUp register a user in the system
func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	// parsing JSON body
	if err := c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	statusId, err := h.services.User.GetActiveStatusId("active")

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// creating the user
	id, err := h.services.Authorization.CreateUser(input, statusId)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// if all is ok
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// refresh refreshes token
func (h *Handler) refresh(c *gin.Context) {
	var input model.RefreshToken
	var tokens model.Tokens
	var user model.User

	// parsing JSON body
	if err := c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	headerParts := strings.Split(input.Authorization, " ")

	if len(headerParts) != 2 {
		model.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		model.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// check if user exists and create the token
	user, tokens, err = h.services.Authorization.GenerateToken(userId, "", "")

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// if all is ok
	c.JSON(http.StatusOK, model.ResponseData{User: user, Tokens: tokens})
}
