package handler

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/gin-gonic/gin"
)

// getComments is a get all comments handler.
func (h *Handler) getComments(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)

	organizationID := ctx.Param("org_id")

	if err != nil {
		return
	}

	results, err := h.services.Comment.GetAll(userID, organizationID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getComment is a get comment by id handler.
func (h *Handler) getComment(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	commentID := ctx.Param("id")

	if organizationID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "invalid id param")

		return
	}

	list, err := h.services.Comment.GetOne(userID, organizationID, commentID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createComment register an comment in the system.
func (h *Handler) createComment(ctx *gin.Context) {
	var input model.Comment

	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	commentID, err := h.services.Comment.Create(userID, input)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": commentID})
}

// updateComment is an update comment by id handler.
func (h *Handler) updateComment(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input model.UpdateCommentInput
	if err = ctx.BindJSON(&input); err != nil {
		model.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = h.services.Comment.Update(userID, input); err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}

// deleteComment is delete comment by id handler.
func (h *Handler) deleteComment(ctx *gin.Context) {
	userID, _, err := h.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	organizationID := ctx.Param("org_id")
	commentID := ctx.Param("id")

	if userID == "" {
		model.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = h.services.Comment.Delete(userID, organizationID, commentID)
	if err != nil {
		model.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, model.StatusResponse{Status: "ok"})
}
