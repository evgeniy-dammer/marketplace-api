package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// OrganizationService is a organization service
type OrganizationService struct {
	repo repository.Organization
}

// NewOrganizationService is a constructor for OrganizationService
func NewOrganizationService(repo repository.Organization) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// GetAll returns all organizations from the system
func (s *OrganizationService) GetAll(userId string) ([]model.Organization, error) {
	return s.repo.GetAll(userId)
}

// GetOne returns organization by id from the system
func (s *OrganizationService) GetOne(userId string, organizationId string) (model.Organization, error) {
	return s.repo.GetOne(userId, organizationId)
}

// Create inserts organization into system
func (s *OrganizationService) Create(userId string, organization model.Organization) (string, error) {
	return s.repo.Create(userId, organization)
}

// Update updates organization by id in the system
func (s *OrganizationService) Update(userId string, organizationId string, input model.UpdateOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, organizationId, input)
}

// Delete deletes organization by id from the system
func (s *OrganizationService) Delete(userId string, organizationId string) error {
	return s.repo.Delete(userId, organizationId)
}
