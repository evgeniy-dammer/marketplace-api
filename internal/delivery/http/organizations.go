package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getOrganizations
// @Summary Get all organizations method.
// @Description Get all organizations method.
// @Tags organizations
// @Accept  json
// @Produce json
// @Security Bearer
// @Success 200		{array}  	organization.Organization	true  "Organization List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/organizations/ [get].
func (d *Delivery) getOrganizations(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	results, err := d.ucOrganization.OrganizationGetAll(ctx, userID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getOrganization
// @Summary Get organization by id method.
// @Description Get organization by id method.
// @Tags organizations
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   				true  "Organization ID"
// @Success 200		{object}  	organization.Organization	true  "Organization data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/organizations/{id} [get].
func (d *Delivery) getOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucOrganization.OrganizationGetOne(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createOrganization
// @Summary Create organization method.
// @Description Create organization method.
// @Tags organizations
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		organization.CreateOrganizationInput 	true  "Organization data"
// @Success 200		{string}  	string									true  "Organization ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/organizations/ [post].
func (d *Delivery) createOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input organization.CreateOrganizationInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	organizationID, err := d.ucOrganization.OrganizationCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": organizationID})
}

// updateOrganization
// @Summary Update organization method.
// @Description Update organization method.
// @Tags organizations
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		organization.UpdateOrganizationInput 	true  "Organization data"
// @Success 200		{object}  	StatusResponse							true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/organizations/ [patch].
func (d *Delivery) updateOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input organization.UpdateOrganizationInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if *input.ID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	if err = d.ucOrganization.OrganizationUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteOrganization
// @Summary Delete organization method.
// @Description Delete organization method.
// @Tags organizations
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "Organization ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/organizations/{id} [delete].
func (d *Delivery) deleteOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	organizationID := ginCtx.Param("id")
	if organizationID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucOrganization.OrganizationDelete(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
