package specification

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/specification"
	"github.com/pkg/errors"
)

// SpecificationGetAll returns all specifications from the system.
func (s *UseCase) SpecificationGetAll(userID string, organizationID string) ([]specification.Specification, error) {
	specifications, err := s.adapterStorage.SpecificationGetAll(userID, organizationID)

	return specifications, errors.Wrap(err, "specifications select error")
}

// SpecificationGetOne returns specification by id from the system.
func (s *UseCase) SpecificationGetOne(userID string, organizationID string, specificationID string) (specification.Specification, error) {
	specification, err := s.adapterStorage.SpecificationGetOne(userID, organizationID, specificationID)

	return specification, errors.Wrap(err, "specification select error")
}

// SpecificationCreate inserts specification into system.
func (s *UseCase) SpecificationCreate(userID string, specification specification.Specification) (string, error) {
	specificationID, err := s.adapterStorage.SpecificationCreate(userID, specification)

	return specificationID, errors.Wrap(err, "specification create error")
}

// SpecificationUpdate updates specification by id in the system.
func (s *UseCase) SpecificationUpdate(userID string, input specification.UpdateSpecificationInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.SpecificationUpdate(userID, input)

	return errors.Wrap(err, "specification update error")
}

// SpecificationDelete deletes specification by id from the system.
func (s *UseCase) SpecificationDelete(userID string, organizationID string, specificationID string) error {
	err := s.adapterStorage.SpecificationDelete(userID, organizationID, specificationID)

	return errors.Wrap(err, "specification delete error")
}
