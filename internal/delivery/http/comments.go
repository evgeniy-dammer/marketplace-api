package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/gin-gonic/gin"
)

// getComments is a get all comments delivery.
func (d *Delivery) getComments(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	results, err := d.ucComment.CommentGetAll(userID, organizationID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getComment is a get comment by id delivery.
func (d *Delivery) getComment(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	commentID := ctx.Param("id")

	if commentID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucComment.CommentGetOne(userID, organizationID, commentID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createComment register a comment in the system.
func (d *Delivery) createComment(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input comment.Comment
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	commentID, err := d.ucComment.CommentCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": commentID})
}

// updateComment is an update comment by id delivery.
func (d *Delivery) updateComment(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input comment.UpdateCommentInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucComment.CommentUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteComment is delete comment by id delivery.
func (d *Delivery) deleteComment(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")

	if organizationID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	commentID := ctx.Param("id")

	if commentID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucComment.CommentDelete(userID, organizationID, commentID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
