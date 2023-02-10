package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getCategories is a get all categories handler
func (h *Handler) getCategories(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")

	results, err := h.services.Category.GetAll(userId, organizationId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// getCategory is a get category by id handler
func (h *Handler) getCategory(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	categoryId := c.Param("id")

	if categoryId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.Category.GetOne(userId, organizationId, categoryId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, list)
}

// createCategory register a category in the system
func (h *Handler) createCategory(c *gin.Context) {
	var input model.Category

	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")

	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.Create(userId, organizationId, input)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// updateCategory is an update category by id handler
func (h *Handler) updateCategory(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	categoryId := c.Param("id")

	if categoryId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var input model.UpdateCategoryInput
	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Category.Update(userId, organizationId, categoryId, input); err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteCategory is a delete category by id handler
func (h *Handler) deleteCategory(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	categoryId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	err = h.services.Category.Delete(userId, organizationId, categoryId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
