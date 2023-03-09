package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/table"
	"github.com/gin-gonic/gin"
)

// getTables is a get all tables delivery.
func (d *Delivery) getTables(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := d.ucTable.TableGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getTable is a get table by id delivery.
func (d *Delivery) getTable(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	tableID := ctx.Param("id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucTable.TableGetOne(userID, organizationID, tableID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createTable register a table in the system.
func (d *Delivery) createTable(ctx *gin.Context) {
	var input table.Table

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	tableID, err := d.ucTable.TableCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": tableID})
}

// updateTable is an update table by id delivery.
func (d *Delivery) updateTable(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input table.UpdateTableInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucTable.TableUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteTable is delete table by id delivery.
func (d *Delivery) deleteTable(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	tableID := ctx.Param("id")

	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucTable.TableDelete(userID, organizationID, tableID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
