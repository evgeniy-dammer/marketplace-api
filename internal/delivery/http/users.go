package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/user"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getAllUsers is a get all users delivery.
func (d *Delivery) getAllUsers(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	search := ginCtx.Query("search")
	roleID := ginCtx.Query("role_id")
	status := ginCtx.Query("status")

	results, err := d.ucUser.UserGetAll(ctx, search, status, roleID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getUser is a get user by id delivery.
func (d *Delivery) getUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID := ginCtx.Param("id")
	if userID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucUser.UserGetOne(ctx, userID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// getAllRoles is a get all user roles delivery.
func (d *Delivery) getAllRoles(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	results, err := d.ucUser.UserGetAllRoles(ctx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// createUser register a user in the system.
func (d *Delivery) createUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input user.CreateUserInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	insertID, err := d.ucUser.UserCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	// if all is ok
	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": insertID})
}

// updateUser is an update user by id delivery.
func (d *Delivery) updateUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	var input user.UpdateUserInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucUser.UserUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteUser is a delete user by id delivery.
func (d *Delivery) deleteUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
	if err != nil {
		return
	}

	dUserID := ginCtx.Param("id")
	if dUserID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucUser.UserDelete(ctx, userID, dUserID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
