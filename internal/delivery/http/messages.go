package http

import (
	"net/http"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/message"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getMessages
// @Summary Get all messages method.
// @Description Get all messages method.
// @Tags messages
// @Accept  json
// @Produce json
// @Security Bearer
// @Success 200		{array}  	message.Message		true  "Message List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/messages/ [get].
func (d *Delivery) getMessages(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getMessages")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	params := queryparameter.QueryParameter{}

	results, err := d.ucMessage.MessageGetAll(ctx, meta, params)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getMessage
// @Summary Get message by id method.
// @Description Get message by id method.
// @Tags messages
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "Message ID"
// @Success 200		{object}  	message.Message		true  "Message data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/messages/{id} [get].
func (d *Delivery) getMessage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getMessage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	messageID := ginCtx.Param("id")
	if messageID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucMessage.MessageGetOne(ctx, meta, messageID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createMessage
// @Summary Create message method.
// @Description Create message method.
// @Tags messages
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		message.CreateMessageInput 	true  "Message data"
// @Success 200		{string}  	string					true  "Message ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/messages/ [post].
func (d *Delivery) createMessage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createMessage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	var input message.CreateMessageInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	messageID, err := d.ucMessage.MessageCreate(ctx, meta, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": messageID})
}

// updateMessage
// @Summary Update message method.
// @Description Update message method.
// @Tags messages
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		message.UpdateMessageInput 	true  "Message data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/messages/ [patch].
func (d *Delivery) updateMessage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateMessage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	var input message.UpdateMessageInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucMessage.MessageUpdate(ctx, meta, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteMessage
// @Summary Delete message method.
// @Description Delete message method.
// @Tags messages
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "Message ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/messages/{id} [delete].
func (d *Delivery) deleteMessage(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteMessage")
		defer span.End()

		ctx = context.New(ctxt)
	}

	meta, err := d.parseMetadata(ginCtx)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	messageID := ginCtx.Param("id")
	if messageID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucMessage.MessageDelete(ctx, meta, messageID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
