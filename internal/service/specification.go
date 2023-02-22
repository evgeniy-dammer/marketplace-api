package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// SpecificationService is an organization service.
type SpecificationService struct {
	repo repository.Specification
}

// NewSpecificationService is a constructor for SpecificationService.
func NewSpecificationService(repo repository.Specification) *SpecificationService {
	return &SpecificationService{repo: repo}
}

// GetAll returns all specifications from the system.
func (s *SpecificationService) GetAll(userID string, organizationID string) ([]model.Specification, error) {
	specifications, err := s.repo.GetAll(userID, organizationID)

	return specifications, errors.Wrap(err, "specifications select error")
}

// GetOne returns specification by id from the system.
func (s *SpecificationService) GetOne(userID string, organizationID string, specificationID string) (model.Specification, error) {
	specification, err := s.repo.GetOne(userID, organizationID, specificationID)

	return specification, errors.Wrap(err, "specification select error")
}

// Create inserts specification into system.
func (s *SpecificationService) Create(userID string, specification model.Specification) (string, error) {
	specificationID, err := s.repo.Create(userID, specification)

	return specificationID, errors.Wrap(err, "specification create error")
}

// Update updates specification by id in the system.
func (s *SpecificationService) Update(userID string, input model.UpdateSpecificationInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "specification update error")
}

// Delete deletes specification by id from the system.
func (s *SpecificationService) Delete(userID string, organizationID string, specificationID string) error {
	err := s.repo.Delete(userID, organizationID, specificationID)

	return errors.Wrap(err, "specification delete error")
}
