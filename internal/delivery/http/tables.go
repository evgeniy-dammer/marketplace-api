package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/table"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getTables
// @Summary Get all tables method.
// @Description Get all tables method.
// @Tags tables
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Success 200		{array}  	table.Table		true  "Table List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/tables/{org_id} [get].
func (d *Delivery) getTables(ginCtx *gin.Context) {
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

	results, err := d.ucTable.TableGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getTable
// @Summary Get table by id method.
// @Description Get table by id method.
// @Tags tables
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Table ID"
// @Success 200		{object}  	table.Table		true  "Table data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/tables/{org_id}/{id} [get].
func (d *Delivery) getTable(ginCtx *gin.Context) {
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

	tableID := ginCtx.Param("id")
	if tableID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucTable.TableGetOne(ctx, userID, organizationID, tableID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createTable
// @Summary Create table method.
// @Description Create table method.
// @Tags tables
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		table.CreateTableInput 	true  "Table data"
// @Success 200		{string}  	string					true  "Table ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/tables/ [post].
func (d *Delivery) createTable(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input table.CreateTableInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	tableID, err := d.ucTable.TableCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": tableID})
}

// updateTable
// @Summary Update table method.
// @Description Update table method.
// @Tags tables
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		table.UpdateTableInput 	true  "Table data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/tables/ [patch].
func (d *Delivery) updateTable(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input table.UpdateTableInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucTable.TableUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteTable
// @Summary Delete table method.
// @Description Delete table method.
// @Tags tables
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Table ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/tables/{org_id}/{id} [delete].
func (d *Delivery) deleteTable(ginCtx *gin.Context) {
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

	tableID := ginCtx.Param("id")
	if tableID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucTable.TableDelete(ctx, userID, organizationID, tableID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
