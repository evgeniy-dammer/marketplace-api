package http

import (
	"net/http"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getCategories
// @Summary Get all categories method.
// @Description Get all categories method.
// @Tags categories
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   		true  "Organization ID"
// @Success 200		{array}  	category.Category	true  "Category List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/categories/{org_id} [get].
func (d *Delivery) getCategories(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getCategories")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")

	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucCategory.CategoryGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getCategory
// @Summary Get category by id method.
// @Description Get category by id method.
// @Tags categories
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   		true  "Organization ID"
// @Param   id	 	path 		string 		   		true  "Category ID"
// @Success 200		{object}  	category.Category	true  "Category data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/categories/{org_id}/{id} [get].
func (d *Delivery) getCategory(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getCategory")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	categoryID := ginCtx.Param("id")
	if categoryID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucCategory.CategoryGetOne(ctx, userID, organizationID, categoryID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createCategory
// @Summary Create category method.
// @Description Create category method.
// @Tags categories
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		category.CreateCategoryInput 	true  "Category data"
// @Success 200		{string}  	string							true  "Category ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/categories/ [post].
func (d *Delivery) createCategory(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createCategory")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input category.CreateCategoryInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	categoryID, err := d.ucCategory.CategoryCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": categoryID})
}

// updateCategory
// @Summary Update category method.
// @Description Update category method.
// @Tags categories
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		category.UpdateCategoryInput 	true  "Category data"
// @Success 200		{object}  	StatusResponse					true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/categories/ [patch].
func (d *Delivery) updateCategory(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateCategory")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input category.UpdateCategoryInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucCategory.CategoryUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteCategory
// @Summary Delete category method.
// @Description Delete category method.
// @Tags categories
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Category ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/categories/{org_id}/{id} [delete].
func (d *Delivery) deleteCategory(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteCategory")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("org_id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	categoryID := ginCtx.Param("id")
	if categoryID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucCategory.CategoryDelete(ctx, userID, organizationID, categoryID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
