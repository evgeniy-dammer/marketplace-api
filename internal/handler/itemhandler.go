package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getItems is a get all items handler
func (h *Handler) getItems(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	organizationId := c.Param("org_id")

	if err != nil {
		return
	}

	results, err := h.services.Item.GetAll(userId, organizationId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// getItem is a get item by id handler
func (h *Handler) getItem(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	itemId := c.Param("id")

	if organizationId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.Item.GetOne(userId, organizationId, itemId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, list)
}

// createItem register a item in the system
func (h *Handler) createItem(c *gin.Context) {
	var input model.Item

	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")

	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Item.Create(userId, organizationId, input)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// updateItem is an update item by id handler
func (h *Handler) updateItem(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	itemId := c.Param("id")

	if organizationId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var input model.UpdateItemInput
	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Item.Update(userId, organizationId, itemId, input); err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteItem is delete item by id handler
func (h *Handler) deleteItem(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("org_id")
	itemId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	err = h.services.Item.Delete(userId, organizationId, itemId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
