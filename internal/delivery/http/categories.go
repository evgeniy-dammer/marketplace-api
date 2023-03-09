package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/gin-gonic/gin"
)

// getCategories is a get all categories delivery.
func (d *Delivery) getCategories(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucCategory.CategoryGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getCategory is a get category by id delivery.
func (d *Delivery) getCategory(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	categoryID := ctx.Param("id")

	if categoryID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucCategory.CategoryGetOne(userID, organizationID, categoryID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createCategory register a category in the system.
func (d *Delivery) createCategory(ctx *gin.Context) {
	var input category.Category

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	categoryID, err := d.ucCategory.CategoryCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": categoryID})
}

// updateCategory is an update category by id delivery.
func (d *Delivery) updateCategory(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input category.UpdateCategoryInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucCategory.CategoryUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteCategory is a delete category by id delivery.
func (d *Delivery) deleteCategory(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	categoryID := ctx.Param("id")

	if categoryID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucCategory.CategoryDelete(userID, organizationID, categoryID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
