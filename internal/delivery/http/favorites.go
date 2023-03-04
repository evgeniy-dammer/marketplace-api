package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/gin-gonic/gin"
)

// createFavorite register an favorite in the system.
func (d *Delivery) createFavorite(ctx *gin.Context) {
	var input favorite.Favorite

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	err = d.ucFavorite.FavoriteCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteFavorite is delete favorite by id delivery.
func (d *Delivery) deleteFavorite(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	itemID := ctx.Param("item_id")

	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = d.ucFavorite.FavoriteDelete(userID, itemID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
