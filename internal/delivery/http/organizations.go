package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/gin-gonic/gin"
)

// getOrganizations is a get all organizations' delivery.
func (d *Delivery) getOrganizations(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	results, err := d.ucOrganization.OrganizationGetAll(userID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getOrganization is a get organization by id delivery.
func (d *Delivery) getOrganization(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("id")
	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucOrganization.OrganizationGetOne(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createOrganization register an organization in the system.
func (d *Delivery) createOrganization(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input organization.Organization
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	organizationID, err := d.ucOrganization.OrganizationCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": organizationID})
}

// updateOrganization is an update organization by id delivery.
func (d *Delivery) updateOrganization(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input organization.UpdateOrganizationInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if *input.ID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	if err = d.ucOrganization.OrganizationUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteOrganization is delete organization by id delivery.
func (d *Delivery) deleteOrganization(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("id")
	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucOrganization.OrganizationDelete(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
