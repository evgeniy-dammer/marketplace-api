package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// getOrders is a get all orders handler.
func (h *Handler) getOrders(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := h.services.Order.GetAll(userID, organizationID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getOrder is a get order by id handler.
func (h *Handler) getOrder(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	orderID := ctx.Param("id")

	if organizationID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := h.services.Order.GetOne(userID, organizationID, orderID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createOrder register an order in the system.
func (h *Handler) createOrder(ctx *gin.Context) {
	var input model.Order

	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	orderID, err := h.services.Order.Create(userID, input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": orderID})
}

// updateOrder is an update order by id handler.
func (h *Handler) updateOrder(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input model.UpdateOrderInput
	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = h.services.Order.Update(userID, input); err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteOrder is delete order by id handler.
func (h *Handler) deleteOrder(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	orderID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = h.services.Order.Delete(userID, organizationID, orderID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
