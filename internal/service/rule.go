package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// RuleService is an organization service.
type RuleService struct {
	repo repository.Rule
}

// NewRuleService is a constructor for RuleService.
func NewRuleService(repo repository.Rule) *RuleService {
	return &RuleService{repo: repo}
}

// GetAll returns all rules from the system.
func (s *RuleService) GetAll(userID string) ([]model.Rule, error) {
	rules, err := s.repo.GetAll(userID)

	return rules, errors.Wrap(err, "rules select error")
}

// GetOne returns rule by id from the system.
func (s *RuleService) GetOne(userID string, ruleID string) (model.Rule, error) {
	rule, err := s.repo.GetOne(userID, ruleID)

	return rule, errors.Wrap(err, "rule select error")
}

// Create inserts rule into system.
func (s *RuleService) Create(userID string, rule model.Rule) (string, error) {
	ruleID, err := s.repo.Create(userID, rule)

	return ruleID, errors.Wrap(err, "rule create error")
}

// Update updates rule by id in the system.
func (s *RuleService) Update(userID string, input model.UpdateRuleInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "rule update error")
}

// Delete deletes rule by id from the system.
func (s *RuleService) Delete(userID string, ruleID string) error {
	err := s.repo.Delete(userID, ruleID)

	return errors.Wrap(err, "rule delete error")
}
