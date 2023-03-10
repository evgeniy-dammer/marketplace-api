package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getOrders is a get all orders delivery.
func (d *Delivery) getOrders(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucOrder.OrderGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getOrder is a get order by id delivery.
func (d *Delivery) getOrder(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	orderID := ginCtx.Param("id")
	if orderID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucOrder.OrderGetOne(ctx, userID, organizationID, orderID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createOrder register an order in the system.
func (d *Delivery) createOrder(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input order.CreateOrderInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	orderID, err := d.ucOrder.OrderCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": orderID})
}

// updateOrder is an update order by id delivery.
func (d *Delivery) updateOrder(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input order.UpdateOrderInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucOrder.OrderUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteOrder is delete order by id delivery.
func (d *Delivery) deleteOrder(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if userID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	orderID := ginCtx.Param("id")
	if orderID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucOrder.OrderDelete(ctx, userID, organizationID, orderID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
