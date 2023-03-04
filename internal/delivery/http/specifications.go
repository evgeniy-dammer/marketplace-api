package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/specification"
	"github.com/gin-gonic/gin"
)

// getSpecifications is a get all specifications delivery.
func (d *Delivery) getSpecifications(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := d.ucSpecification.SpecificationGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getSpecification is a get specification by id delivery.
func (d *Delivery) getSpecification(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	specificationID := ctx.Param("id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := d.ucSpecification.SpecificationGetOne(userID, organizationID, specificationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createSpecification register an specification in the system.
func (d *Delivery) createSpecification(ctx *gin.Context) {
	var input specification.Specification

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	specificationID, err := d.ucSpecification.SpecificationCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": specificationID})
}

// updateSpecification is an update specification by id delivery.
func (d *Delivery) updateSpecification(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input specification.UpdateSpecificationInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = d.ucSpecification.SpecificationUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteSpecification is delete specification by id delivery.
func (d *Delivery) deleteSpecification(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	specificationID := ctx.Param("id")

	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = d.ucSpecification.SpecificationDelete(userID, organizationID, specificationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
