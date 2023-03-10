package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getOrganizations is a get all organizations' delivery.
func (d *Delivery) getOrganizations(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// getOrganization is a get organization by id delivery.
func (d *Delivery) getOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// createOrganization register an organization in the system.
func (d *Delivery) createOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// updateOrganization is an update organization by id delivery.
func (d *Delivery) updateOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// deleteOrganization is delete organization by id delivery.
func (d *Delivery) deleteOrganization(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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
