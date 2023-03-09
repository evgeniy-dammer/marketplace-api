package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/gin-gonic/gin"
)

// getOrders is a get all orders delivery.
func (d *Delivery) getOrders(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucOrder.OrderGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getOrder is a get order by id delivery.
func (d *Delivery) getOrder(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	orderID := ctx.Param("id")
	if orderID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucOrder.OrderGetOne(userID, organizationID, orderID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createOrder register an order in the system.
func (d *Delivery) createOrder(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input order.Order
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	orderID, err := d.ucOrder.OrderCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": orderID})
}

// updateOrder is an update order by id delivery.
func (d *Delivery) updateOrder(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input order.UpdateOrderInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucOrder.OrderUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteOrder is delete order by id delivery.
func (d *Delivery) deleteOrder(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	orderID := ctx.Param("id")
	if orderID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucOrder.OrderDelete(userID, organizationID, orderID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
