package handler

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// getOrganizations is a get all organizations handler
func (h *Handler) getOrganizations(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	results, err := h.services.Organization.GetAll(userId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, results)
}

// getOrganization is a get organization by id handler
func (h *Handler) getOrganization(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("id")

	if organizationId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.Organization.GetOne(userId, organizationId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, list)
}

// createOrganization register a organization in the system
func (h *Handler) createOrganization(c *gin.Context) {
	var input model.Organization

	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Organization.Create(userId, input)

	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// updateOrganization is an update organization by id handler
func (h *Handler) updateOrganization(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("id")

	if organizationId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var input model.UpdateOrganizationInput
	if err = c.BindJSON(&input); err != nil {
		model.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Organization.Update(userId, organizationId, input); err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteOrganization is delete organization by id handler
func (h *Handler) deleteOrganization(c *gin.Context) {
	userId, _, err := h.getUserIdAndRole(c)

	if err != nil {
		return
	}

	organizationId := c.Param("id")

	if userId == "" {
		model.NewErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	err = h.services.Organization.Delete(userId, organizationId)
	if err != nil {
		model.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
