package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getAllUsers is a get all users handler
func (h *Handler) getAllUsers(c *gin.Context) {
	search := c.Query("search")
	roleId := c.Query("role_id")
	status := c.Query("status")

	results, err := h.services.User.GetAll(search, status, roleId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// getUser is a get user by id handler
func (h *Handler) getUser(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.User.GetOne(userId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, list)
}

// getAllRoles is a get all user roles handler
func (h *Handler) getAllRoles(c *gin.Context) {
	results, err := h.services.User.GetAllRoles()
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// createUser register a user in the system
func (h *Handler) createUser(c *gin.Context) {
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
	id, err := h.services.User.Create(input, statusId)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// if all is ok
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// updateUser is an update user by id handler
func (h *Handler) updateUser(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var input model.UpdateUserInput
	if err := c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.Update(userId, input); err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteUser is a delete user by id handler
func (h *Handler) deleteUser(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	err := h.services.User.Delete(userId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
