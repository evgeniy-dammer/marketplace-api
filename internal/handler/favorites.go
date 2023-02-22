package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// createFavorite register an favorite in the system.
func (h *Handler) createFavorite(ctx *gin.Context) {
	var input model.Favorite

	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	err = h.services.Favorite.Create(userID, input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteFavorite is delete favorite by id handler.
func (h *Handler) deleteFavorite(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	itemID := ctx.Param("item_id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = h.services.Favorite.Delete(userID, itemID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
