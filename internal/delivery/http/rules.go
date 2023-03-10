package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/gin-gonic/gin"
)

// getRules is a get all rules delivery.
func (d *Delivery) getRules(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// getRule is a get rule by id delivery.
func (d *Delivery) getRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// createRule register a rule in the system.
func (d *Delivery) createRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// updateRule is an update rule by id delivery.
func (d *Delivery) updateRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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

// deleteRule is delete rule by id delivery.
func (d *Delivery) deleteRule(ginCtx *gin.Context) {
	ctx := context.New(ginCtx)

	userID, _, err := d.getUserIDAndRole(ginCtx)
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
