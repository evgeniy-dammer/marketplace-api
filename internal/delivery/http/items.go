package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/gin-gonic/gin"
)

// getItems is a get all items delivery.
func (d *Delivery) getItems(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := d.ucItem.ItemGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getItem is a get item by id delivery.
func (d *Delivery) getItem(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	itemID := ctx.Param("id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := d.ucItem.ItemGetOne(userID, organizationID, itemID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createItem register an item in the system.
func (d *Delivery) createItem(ctx *gin.Context) {
	var input item.Item

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	itemID, err := d.ucItem.ItemCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": itemID})
}

// updateItem is an update item by id delivery.
func (d *Delivery) updateItem(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input item.UpdateItemInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = d.ucItem.ItemUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteItem is delete item by id delivery.
func (d *Delivery) deleteItem(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	itemID := ctx.Param("id")

	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = d.ucItem.ItemDelete(userID, organizationID, itemID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
