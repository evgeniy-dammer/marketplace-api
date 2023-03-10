package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getOrders
// @Summary Get all orders method.
// @Description Get all orders method.
// @Tags orders
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Success 200		{array}  	order.Order		true  "Order List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/orders/{org_id} [get].
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

// getOrder
// @Summary Get order by id method.
// @Description Get order by id method.
// @Tags orders
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Order ID"
// @Success 200		{object}  	order.Order		true  "Order data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/orders/{org_id}/{id} [get].
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

// createOrder
// @Summary Create order method.
// @Description Create order method.
// @Tags orders
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		order.CreateOrderInput 	true  "Order data"
// @Success 200		{string}  	string					true  "Order ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/orders/ [post].
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

// updateOrder
// @Summary Update order method.
// @Description Update order method.
// @Tags orders
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		order.UpdateOrderInput 	true  "Order data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/orders/ [patch].
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

// deleteOrder
// @Summary Delete order method.
// @Description Delete order method.
// @Tags orders
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Order ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/orders/{org_id}/{id} [delete].
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
