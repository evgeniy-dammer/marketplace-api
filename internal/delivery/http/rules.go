package http

import (
	"net/http"

	"github.com/evgeniy-dammer/emenu-api/internal/domain"
	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/gin-gonic/gin"
)

// getRules is a get all rules delivery.
func (d *Delivery) getRules(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	results, err := d.ucRule.RuleGetAll(userID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, results)
}

// getRule is a get rule by id delivery.
func (d *Delivery) getRule(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	ruleID := ctx.Param("id")

	list, err := d.ucRule.RuleGetOne(userID, ruleID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusCreated, list)
}

// createRule register an rule in the system.
func (d *Delivery) createRule(ctx *gin.Context) {
	var input rule.Rule

	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	ruleID, err := d.ucRule.RuleCreate(userID, input)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{"id": ruleID})
}

// updateRule is an update rule by id delivery.
func (d *Delivery) updateRule(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	var input rule.UpdateRuleInput
	if err = ctx.BindJSON(&input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())

		return
	}

	if err = d.ucRule.RuleUpdate(userID, input); err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// deleteRule is delete rule by id delivery.
func (d *Delivery) deleteRule(ctx *gin.Context) {
	userID, _, err := d.getUserIDAndRole(ctx)
	if err != nil {
		return
	}

	ruleID := ctx.Param("id")

	if ruleID == "" {
		domain.NewErrorResponse(ctx, http.StatusBadRequest, "empty id param")

		return
	}

	err = d.ucRule.RuleDelete(userID, ruleID)
	if err != nil {
		domain.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctx.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
