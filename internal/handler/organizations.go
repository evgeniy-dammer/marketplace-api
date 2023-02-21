package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// getOrganizations is a get all organizations' handler.
func (h *Handler) getOrganizations(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	results, err := h.services.Organization.GetAll(userID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getOrganization is a get organization by id handler.
func (h *Handler) getOrganization(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("id")

	if organizationID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := h.services.Organization.GetOne(userID, organizationID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createOrganization register an organization in the system.
func (h *Handler) createOrganization(ctx *gin.Context) {
	var input model.Organization

	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	organizationID, err := h.services.Organization.Create(userID, input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": organizationID})
}

// updateOrganization is an update organization by id handler.
func (h *Handler) updateOrganization(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input model.UpdateOrganizationInput
	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if *input.ID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	if err = h.services.Organization.Update(userID, input); err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteOrganization is delete organization by id handler.
func (h *Handler) deleteOrganization(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = h.services.Organization.Delete(userID, organizationID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
