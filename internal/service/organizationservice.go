package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// OrganizationService is a organization service.
type OrganizationService struct {
	repo repository.Organization
}

// NewOrganizationService is a constructor for OrganizationService.
func NewOrganizationService(repo repository.Organization) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// GetAll returns all organizations from the system.
func (s *OrganizationService) GetAll(userID string) ([]model.Organization, error) {
	organizations, err := s.repo.GetAll(userID)

	return organizations, errors.Wrap(err, "organization select error")
}

// GetOne returns organization by id from the system.
func (s *OrganizationService) GetOne(userID string, organizationID string) (model.Organization, error) {
	organization, err := s.repo.GetOne(userID, organizationID)

	return organization, errors.Wrap(err, "organization select error")
}

// Create inserts organization into system.
func (s *OrganizationService) Create(userID string, organization model.Organization) (string, error) {
	organizationID, err := s.repo.Create(userID, organization)

	return organizationID, errors.Wrap(err, "organization create error")
}

// Update updates organization by id in the system.
func (s *OrganizationService) Update(userID string, input model.UpdateOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "organization update error")
}

// Delete deletes organization by id from the system.
func (s *OrganizationService) Delete(userID string, organizationID string) error {
	err := s.repo.Delete(userID, organizationID)

	return errors.Wrap(err, "organization delete error")
}
