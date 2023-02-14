package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// getAllUsers is a get all users handler.
func (h *Handler) getAllUsers(ctx *gin.Context) {
	search := ctx.Query("search")
	roleID := ctx.Query("role_id")
	status := ctx.Query("status")

	results, err := h.services.User.GetAll(search, status, roleID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getUser is a get user by id handler.
func (h *Handler) getUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := h.services.User.GetOne(userID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// getAllRoles is a get all user roles handler.
func (h *Handler) getAllRoles(ctx *gin.Context) {
	results, err := h.services.User.GetAllRoles()
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// createUser register a user in the system.
func (h *Handler) createUser(ctx *gin.Context) {
	var input model.User

	// parsing JSON body
	if err := ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	statusID, err := h.services.User.GetActiveStatusID("active")
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	userID, err := h.services.User.Create(input, statusID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, map[string]interface{}{"id": userID})
}

// updateUser is an update user by id handler.
func (h *Handler) updateUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	var input model.UpdateUserInput
	if err := ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.User.Update(userID, input); err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteUser is a delete user by id handler.
func (h *Handler) deleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err := h.services.User.Delete(userID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
