package http

import (
	"net/http"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/comment"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getComments
// @Summary Get all comments method.
// @Description Get all comments method.
// @Tags comments
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   		true  "Organization ID"
// @Success 200		{array}  	comment.Comment		true  "Comments List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/comments/{org_id} [get].
func (d *Delivery) getComments(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getComments")
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

	results, err := d.ucComment.CommentGetAll(ctx, userID, organizationID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getComment
// @Summary Get comment by id method.
// @Description Get comment by id method.
// @Tags comments
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   		true  "Organization ID"
// @Param   id	 	path 		string 		   		true  "Comment ID"
// @Success 200		{object}  	comment.Comment		true  "Comment data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/comments/{org_id}/{id} [get].
func (d *Delivery) getComment(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getComment")
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

	commentID := ginCtx.Param("id")

	if commentID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucComment.CommentGetOne(ctx, userID, organizationID, commentID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createComment
// @Summary Create comment method.
// @Description Create comment method.
// @Tags comments
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		comment.CreateCommentInput 	true  "Comment data"
// @Success 200		{string}  	string						true  "Comment ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/comments/ [post].
func (d *Delivery) createComment(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createComment")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input comment.CreateCommentInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	commentID, err := d.ucComment.CommentCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": commentID})
}

// updateComment
// @Summary Update comment method.
// @Description Update comment method.
// @Tags comments
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		comment.UpdateCommentInput 	true  "Comment data"
// @Success 200		{object}  	StatusResponse				true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/comments/ [patch].
func (d *Delivery) updateComment(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateComment")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input comment.UpdateCommentInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucComment.CommentUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteComment
// @Summary Delete comment method.
// @Description Delete comment method.
// @Tags comments
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   org_id 	path 		string 		   	true  "Organization ID"
// @Param   id	 	path 		string 		   	true  "Comment ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/comments/{org_id}/{id} [delete].
func (d *Delivery) deleteComment(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteComment")
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

	commentID := ginCtx.Param("id")

	if commentID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucComment.CommentDelete(ctx, userID, organizationID, commentID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
