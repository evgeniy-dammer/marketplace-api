package rule

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/pkg/errors"
)

// RuleGetAll returns all rules from the system.
func (s *UseCase) RuleGetAll(ctx context.Context, userID string) ([]rule.Rule, error) {
	rules, err := s.adapterStorage.RuleGetAll(ctx, userID)

	return rules, errors.Wrap(err, "rules select error")
}

// RuleGetOne returns rule by id from the system.
func (s *UseCase) RuleGetOne(ctx context.Context, userID string, ruleID string) (rule.Rule, error) {
	rule, err := s.adapterStorage.RuleGetOne(ctx, userID, ruleID)

	return rule, errors.Wrap(err, "rule select error")
}

// RuleCreate inserts rule into system.
func (s *UseCase) RuleCreate(ctx context.Context, userID string, input rule.CreateRuleInput) (string, error) {
	ruleID, err := s.adapterStorage.RuleCreate(ctx, userID, input)

	return ruleID, errors.Wrap(err, "rule create error")
}

// RuleUpdate updates rule by id in the system.
func (s *UseCase) RuleUpdate(ctx context.Context, userID string, input rule.UpdateRuleInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.RuleUpdate(ctx, userID, input)

	return errors.Wrap(err, "rule update error")
}

// RuleDelete deletes rule by id from the system.
func (s *UseCase) RuleDelete(ctx context.Context, userID string, ruleID string) error {
	err := s.adapterStorage.RuleDelete(ctx, userID, ruleID)

	return errors.Wrap(err, "rule delete error")
}
