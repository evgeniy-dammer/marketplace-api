package http

import (
	"net/http"

	_ "github.com/evgeniy-dammer/marketplace-api/internal/domain/role"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/user"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getAllUsers
// @Summary Get all users method.
// @Description Get all users method.
// @Tags users
// @Accept  json
// @Produce json
// @Security Bearer
// @Success 200		{array}  	user.User		true  "User List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/ [get].
func (d *Delivery) getAllUsers(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getAllUsers")
		defer span.End()

		ctx = context.New(ctxt)
	}

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

// getUser
// @Summary Get user by id method.
// @Description Get user by id method.
// @Tags users
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "User ID"
// @Success 200		{object}  	user.User		true  "User data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/{id} [get].
func (d *Delivery) getUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

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

// getAllRoles
// @Summary Get all roles method.
// @Description Get all roles method.
// @Tags roles
// @Accept  json
// @Produce json
// @Security Bearer
// @Success 200		{array}  	role.Role		true  "Role List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/roles/ [get].
func (d *Delivery) getAllRoles(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getAllRoles")
		defer span.End()

		ctx = context.New(ctxt)
	}

	results, err := d.ucUser.UserGetAllRoles(ctx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// createUser
// @Summary Create user method.
// @Description Create user method.
// @Tags users
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		user.CreateUserInput 	true  "User data"
// @Success 200		{string}  	string					true  "User ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/ [post].
func (d *Delivery) createUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
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

// updateUser
// @Summary Update user method.
// @Description Update user method.
// @Tags users
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		user.UpdateUserInput 	true  "User data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/ [patch].
func (d *Delivery) updateUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
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

// deleteUser
// @Summary Delete user method.
// @Description Delete user method.
// @Tags users
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "User ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/users/{id} [delete].
func (d *Delivery) deleteUser(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteUser")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
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
