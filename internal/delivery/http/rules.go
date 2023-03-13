package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/gin-gonic/gin"
)

// getRules
// @Summary Get all rules method.
// @Description Get all rules method.
// @Tags rules
// @Accept  json
// @Produce json
// @Security Bearer
// @Success 200		{array}  	rule.Rule		true  "Rule List"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/rules/ [get].
func (d *Delivery) getRules(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getRules")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	results, err := d.ucRule.RuleGetAll(ctx, userID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, results)
}

// getRule
// @Summary Get rule by id method.
// @Description Get rule by id method.
// @Tags rules
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "Rule ID"
// @Success 200		{object}  	rule.Rule		true  "Rule data"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/rules/{id} [get].
func (d *Delivery) getRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.getRule")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	ruleID := ginCtx.Param("id")
	if ruleID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	list, err := d.ucRule.RuleGetOne(ctx, userID, ruleID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, list)
}

// createRule
// @Summary Create rule method.
// @Description Create rule method.
// @Tags rules
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		rule.CreateRuleInput 	true  "Rule data"
// @Success 200		{string}  	string					true  "Rule ID"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/rules/ [post].
func (d *Delivery) createRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.createRule")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input rule.CreateRuleInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	ruleID, err := d.ucRule.RuleCreate(ctx, userID, input)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, map[string]interface{}{"id": ruleID})
}

// updateRule
// @Summary Update rule method.
// @Description Update rule method.
// @Tags rules
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   input 	body 		rule.UpdateRuleInput 	true  "Rule data"
// @Success 200		{object}  	StatusResponse			true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/rules/ [patch].
func (d *Delivery) updateRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.updateRule")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	var input rule.UpdateRuleInput
	if err = ginCtx.BindJSON(&input); err != nil {
		NewErrorResponse(ginCtx, http.StatusBadRequest, err)

		return
	}

	if err = d.ucRule.RuleUpdate(ctx, userID, input); err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}

// deleteRule
// @Summary Delete rule method.
// @Description Delete rule method.
// @Tags rules
// @Accept  json
// @Produce json
// @Security Bearer
// @Param   id	 	path 		string 		   	true  "Rule ID"
// @Success 200		{object}  	StatusResponse	true  "OK"
// @Failure 400 	{object}    ErrorResponse
// @Failure 401	 	{object}	ErrorResponse
// @Failure 404 	{object} 	ErrorResponse
// @Failure 500 	{object} 	ErrorResponse
// @Router /api/v1/rules/{id} [delete].
func (d *Delivery) deleteRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	if d.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ginCtx.Request.Context(), "Delivery.deleteRule")
		defer span.End()

		ctx = context.New(ctxt)
	}

	userID, err := d.getUserID(ginCtx)
	if err != nil {
		return
	}

	ruleID := ginCtx.Param("id")
	if ruleID == "" {
		NewErrorResponse(ginCtx, http.StatusBadRequest, ErrEmptyIDParam)

		return
	}

	err = d.ucRule.RuleDelete(ctx, userID, ruleID)
	if err != nil {
		NewErrorResponse(ginCtx, http.StatusInternalServerError, err)

		return
	}

	ginCtx.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
