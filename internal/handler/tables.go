package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// getTables is a get all tables handler.
func (h *Handler) getTables(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := h.services.Table.GetAll(userID, organizationID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getTable is a get table by id handler.
func (h *Handler) getTable(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	tableID := ctx.Param("id")

	if organizationID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := h.services.Table.GetOne(userID, organizationID, tableID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createTable register an table in the system.
func (h *Handler) createTable(ctx *gin.Context) {
	var input model.Table

	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	tableID, err := h.services.Table.Create(userID, input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": tableID})
}

// updateTable is an update table by id handler.
func (h *Handler) updateTable(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input model.UpdateTableInput
	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = h.services.Table.Update(userID, input); err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteTable is delete table by id handler.
func (h *Handler) deleteTable(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	tableID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = h.services.Table.Delete(userID, organizationID, tableID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
