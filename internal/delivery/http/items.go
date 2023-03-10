package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getItems
// @Summary Get all items method.
// @Description Get all items method.
// @Tags items
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Success 200		{array}  	item.Item		true  "Item List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/items/{org_id} [get].
func (d *Delivery) getItems(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucItem.ItemGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getItem
// @Summary Get item by id method.
// @Description Get item by id method.
// @Tags items
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "item ID"
// @Success 200		{object}  	item.Item		true  "item data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/items/{org_id}/{id} [get].
func (d *Delivery) getItem(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	itemID := ginCtx.Param("id")
	if itemID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucItem.ItemGetOne(ctx, userID, organizationID, itemID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createItem
// @Summary Create item method.
// @Description Create item method.
// @Tags items
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		item.CreateItemInput 	true  "Item data"
// @Success 200		{string}  	string					true  "Item ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/items/ [post].
func (d *Delivery) createItem(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input item.CreateItemInput

	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	itemID, err := d.ucItem.ItemCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": itemID})
}

// updateItem
// @Summary Update item method.
// @Description Update item method.
// @Tags items
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		item.UpdateItemInput 	true  "Item data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/items/ [patch].
func (d *Delivery) updateItem(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input item.UpdateItemInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucItem.ItemUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteItem
// @Summary Delete item method.
// @Description Delete item method.
// @Tags items
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Item ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/items/{org_id}/{id} [delete].
func (d *Delivery) deleteItem(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if userID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	itemID := ginCtx.Param("id")
	if itemID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucItem.ItemDelete(ctx, userID, organizationID, itemID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
