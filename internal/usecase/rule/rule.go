package rule

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// RuleGetAll returns all rules from the system.
func (s *UseCase) RuleGetAll(ctx context.Context, userID string) ([]rule.Rule, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.RuleGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID)
	}

	rules, err := s.adapterStorage.RuleGetAll(ctx, userID)

	return rules, errors.Wrap(err, "rules select error")
}

// getAllWithCache returns rules from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string) ([]rule.Rule, error) {
	rules, err := s.adapterCache.RuleGetAll(ctx)
	if err != nil {
		logger.Logger.Error("unable to get rules from cache", zap.String("error", err.Error()))
	}

	if len(rules) > 0 {
		return rules, nil
	}

	rules, err = s.adapterStorage.RuleGetAll(ctx, userID)

	if err != nil {
		return rules, errors.Wrap(err, "rules select failed")
	}

	if err = s.adapterCache.RuleSetAll(ctx, rules); err != nil {
		logger.Logger.Error("unable to add rules into cache", zap.String("error", err.Error()))
	}

	return rules, nil
}

// RuleGetOne returns rule by id from the system.
func (s *UseCase) RuleGetOne(ctx context.Context, userID string, ruleID string) (rule.Rule, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.RuleGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, ruleID)
	}

	rle, err := s.adapterStorage.RuleGetOne(ctx, userID, ruleID)

	return rle, errors.Wrap(err, "rule select error")
}

// getOneWithCache returns rule by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, ruleID string) (rule.Rule, error) {
	rle, err := s.adapterCache.RuleGetOne(ctx, ruleID)
	if err != nil {
		logger.Logger.Error("unable to get rule from cache", zap.String("error", err.Error()))
	}

	if rle != (rule.Rule{}) {
		return rle, nil
	}

	rle, err = s.adapterStorage.RuleGetOne(ctx, userID, ruleID)

	if err != nil {
		return rle, errors.Wrap(err, "rule select failed")
	}

	if err = s.adapterCache.RuleCreate(ctx, rle); err != nil {
		logger.Logger.Error("unable to add rule into cache", zap.String("error", err.Error()))
	}

	return rle, nil
}

// RuleCreate inserts rule into system.
func (s *UseCase) RuleCreate(ctx context.Context, userID string, input rule.CreateRuleInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.RuleCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	ruleID, err := s.adapterStorage.RuleCreate(ctx, userID, input)
	if err != nil {
		return ruleID, errors.Wrap(err, "rule create error")
	}

	if s.isCacheOn {
		rle, err := s.adapterStorage.RuleGetOne(ctx, userID, ruleID)
		if err != nil {
			return "", errors.Wrap(err, "rule select from database failed")
		}

		err = s.adapterCache.RuleCreate(ctx, rle)
		if err != nil {
			return "", errors.Wrap(err, "rule create in cache failed")
		}

		err = s.adapterCache.RuleInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "rule invalidate users in cache failed")
		}
	}

	return ruleID, nil
}

// RuleUpdate updates rule by id in the system.
func (s *UseCase) RuleUpdate(ctx context.Context, userID string, input rule.UpdateRuleInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.RuleUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.RuleUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "rule update in database failed")
	}

	if s.isCacheOn {
		rle, err := s.adapterStorage.RuleGetOne(ctx, userID, *input.ID)
		if err != nil {
			return errors.Wrap(err, "rule select from database failed")
		}

		err = s.adapterCache.RuleUpdate(ctx, rle)
		if err != nil {
			return errors.Wrap(err, "rule update in cache failed")
		}

		err = s.adapterCache.RuleInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "rule invalidate users in cache failed")
		}
	}

	return nil
}

// RuleDelete deletes rule by id from the system.
func (s *UseCase) RuleDelete(ctx context.Context, userID string, ruleID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.RuleDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.RuleDelete(ctx, userID, ruleID)
	if err != nil {
		return errors.Wrap(err, "rule delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.RuleDelete(ctx, ruleID)
		if err != nil {
			return errors.Wrap(err, "rule update in cache failed")
		}

		err = s.adapterCache.RuleInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate rules in cache failed")
		}
	}

	return nil
}
