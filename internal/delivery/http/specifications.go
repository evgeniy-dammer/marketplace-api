package http

import (
	"net/http"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getSpecifications
// @Summary Get all specifications method.
// @Description Get all specifications method.
// @Tags specifications
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   					true  "Specification ID"
// @Success 200		{array}  	specification.Specification		true  "Specification List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/specifications/{org_id} [get].
func (d *Delivery) getSpecifications(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getSpecifications")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		return
	}

	params := queryparameter.QueryParameter{}

	results, err := d.ucSpecification.SpecificationGetAll(ctx, meta, params)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getSpecification
// @Summary Get specification by id method.
// @Description Get specification by id method.
// @Tags specifications
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   					true  "Organization ID"
// @Param   id	 	path 		string 		   					true  "Specification ID"
// @Success 200		{object}  	specification.Specification		true  "Specification data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/specifications/{org_id}/{id} [get].
func (d *Delivery) getSpecification(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getSpecification")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		return
	}

	specificationID := ginCtx.Param("id")
	if specificationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucSpecification.SpecificationGetOne(ctx, meta, specificationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createSpecification
// @Summary Create specification method.
// @Description Create specification method.
// @Tags specifications
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		specification.CreateSpecificationInput 	true  "Specification data"
// @Success 200		{string}  	string									true  "Specification ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/specifications/ [post].
func (d *Delivery) createSpecification(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createSpecification")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		return
	}

	var input specification.CreateSpecificationInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	specificationID, err := d.ucSpecification.SpecificationCreate(ctx, meta, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": specificationID})
}

// updateSpecification
// @Summary Update specification method.
// @Description Update specification method.
// @Tags specifications
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		specification.UpdateSpecificationInput 	true  "Specification data"
// @Success 200		{object}  	StatusResponse							true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/specifications/ [patch].
func (d *Delivery) updateSpecification(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateSpecification")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		return
	}

	var input specification.UpdateSpecificationInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucSpecification.SpecificationUpdate(ctx, meta, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteSpecification
// @Summary Delete specification method.
// @Description Delete specification method.
// @Tags specifications
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Specification ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/specifications/{org_id}/{id} [delete].
func (d *Delivery) deleteSpecification(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteSpecification")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		return
	}

	specificationID := ginCtx.Param("id")
	if specificationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucSpecification.SpecificationDelete(ctx, meta, specificationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
