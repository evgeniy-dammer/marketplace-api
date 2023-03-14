package http

import (
	"net/http"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/image"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getImages
// @Summary Get all images method.
// @Description Get all images method.
// @Tags images
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Success 200		{array}  	image.Image		true  "Image List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/images/{org_id} [get].
func (d *Delivery) getImages(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getImages")
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

	results, err := d.ucImage.ImageGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getImage
// @Summary Get image by id method.
// @Description Get image by id method.
// @Tags images
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Image ID"
// @Success 200		{object}  	image.Image		true  "Image data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/images/{org_id}/{id} [get].
func (d *Delivery) getImage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getImage")
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

	imageID := ginCtx.Param("id")
	if imageID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucImage.ImageGetOne(ctx, userID, organizationID, imageID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createImage
// @Summary Create image method.
// @Description Create image method.
// @Tags images
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		image.CreateImageInput 	true  "Image data"
// @Success 200		{string}  	string					true  "Image ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/images/ [post].
func (d *Delivery) createImage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createImage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input image.CreateImageInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	imageID, err := d.ucImage.ImageCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": imageID})
}

// updateImage
// @Summary Update image method.
// @Description Update image method.
// @Tags images
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		image.UpdateImageInput 	true  "Image data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/images/ [patch].
func (d *Delivery) updateImage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateImage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input image.UpdateImageInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucImage.ImageUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteImage
// @Summary Delete image method.
// @Description Delete image method.
// @Tags images
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Image ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/images/{org_id}/{id} [delete].
func (d *Delivery) deleteImage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteImage")
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

	imageID := ginCtx.Param("id")
	if imageID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucImage.ImageDelete(ctx, userID, organizationID, imageID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
