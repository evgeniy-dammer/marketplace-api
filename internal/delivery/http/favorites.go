package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/favorite"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// createFavorite
// @Summary Create favorite item method.
// @Description Create favorite item method.
// @Tags favorites
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		favorite.Favorite 	true  "Favorite data"
// @Success 200		{object}  	StatusResponse		true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/favorites/ [post].
func (d *Delivery) createFavorite(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input favorite.Favorite
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	err = d.ucFavorite.FavoriteCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteFavorite
// @Summary Delete favorite item method.
// @Description Delete favorite item method.
// @Tags favorites
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   item_id path 		string 		   	true  "Item ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/favorites/{item_id} [delete].
func (d *Delivery) deleteFavorite(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	itemID := ginCtx.Param("item_id")
	if itemID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucFavorite.FavoriteDelete(ctx, userID, itemID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}