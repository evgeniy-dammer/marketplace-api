package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/gin-gonic/gin"
)

// getAllUsers is a get all users delivery.
func (d *Delivery) getAllUsers(ctx *gin.Context) {
	search := ctx.Query("search")
	roleID := ctx.Query("role_id")
	status := ctx.Query("status")

	results, err := d.ucUser.UserGetAll(search, status, roleID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getUser is a get user by id delivery.
func (d *Delivery) getUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucUser.UserGetOne(userID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// getAllRoles is a get all user roles delivery.
func (d *Delivery) getAllRoles(ctx *gin.Context) {
	results, err := d.ucUser.UserGetAllRoles()
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// createUser register a user in the system.
func (d *Delivery) createUser(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input user.User
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	insertID, err := d.ucUser.UserCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	// if all is ok
	ctx.JSON(http.StatusOK, map[string]interface{}{"id": insertID})
}

// updateUser is an update user by id delivery.
func (d *Delivery) updateUser(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input user.UpdateUserInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucUser.UserUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteUser is a delete user by id delivery.
func (d *Delivery) deleteUser(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	dUserID := ctx.Param("id")
	if dUserID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucUser.UserDelete(userID, dUserID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
