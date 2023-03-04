package organization

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/pkg/errors"
)

// OrganizationGetAll returns all organizations from the system.
func (s *UseCase) OrganizationGetAll(userID string) ([]organization.Organization, error) {
	organizations, err := s.adapterStorage.OrganizationGetAll(userID)

	return organizations, errors.Wrap(err, "organization select error")
}

// OrganizationGetOne returns organization by id from the system.
func (s *UseCase) OrganizationGetOne(userID string, organizationID string) (organization.Organization, error) {
	org, err := s.adapterStorage.OrganizationGetOne(userID, organizationID)

	return org, errors.Wrap(err, "organization select error")
}

// OrganizationCreate inserts organization into system.
func (s *UseCase) OrganizationCreate(userID string, organization organization.Organization) (string, error) {
	organizationID, err := s.adapterStorage.OrganizationCreate(userID, organization)

	return organizationID, errors.Wrap(err, "organization create error")
}

// OrganizationUpdate updates organization by id in the system.
func (s *UseCase) OrganizationUpdate(userID string, input organization.UpdateOrganizationInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrganizationUpdate(userID, input)

	return errors.Wrap(err, "organization update error")
}

// OrganizationDelete deletes organization by id from the system.
func (s *UseCase) OrganizationDelete(userID string, organizationID string) error {
	err := s.adapterStorage.OrganizationDelete(userID, organizationID)

	return errors.Wrap(err, "organization delete error")
}
