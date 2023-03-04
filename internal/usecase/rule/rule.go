package rule

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/rule"
	"github.com/pkg/errors"
)

// RuleGetAll returns all rules from the system.
func (s *UseCase) RuleGetAll(userID string) ([]rule.Rule, error) {
	rules, err := s.adapterStorage.RuleGetAll(userID)

	return rules, errors.Wrap(err, "rules select error")
}

// RuleGetOne returns rule by id from the system.
func (s *UseCase) RuleGetOne(userID string, ruleID string) (rule.Rule, error) {
	rule, err := s.adapterStorage.RuleGetOne(userID, ruleID)

	return rule, errors.Wrap(err, "rule select error")
}

// RuleCreate inserts rule into system.
func (s *UseCase) RuleCreate(userID string, rule rule.Rule) (string, error) {
	ruleID, err := s.adapterStorage.RuleCreate(userID, rule)

	return ruleID, errors.Wrap(err, "rule create error")
}

// RuleUpdate updates rule by id in the system.
func (s *UseCase) RuleUpdate(userID string, input rule.UpdateRuleInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.RuleUpdate(userID, input)

	return errors.Wrap(err, "rule update error")
}

// RuleDelete deletes rule by id from the system.
func (s *UseCase) RuleDelete(userID string, ruleID string) error {
	err := s.adapterStorage.RuleDelete(userID, ruleID)

	return errors.Wrap(err, "rule delete error")
}
