package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// TableService is an organization service.
type TableService struct {
	repo repository.Table
}

// NewTableService is a constructor for TableService.
func NewTableService(repo repository.Table) *TableService {
	return &TableService{repo: repo}
}

// GetAll returns all tables from the system.
func (s *TableService) GetAll(userID string, organizationID string) ([]model.Table, error) {
	tables, err := s.repo.GetAll(userID, organizationID)

	return tables, errors.Wrap(err, "tables select error")
}

// GetOne returns table by id from the system.
func (s *TableService) GetOne(userID string, organizationID string, tableID string) (model.Table, error) {
	table, err := s.repo.GetOne(userID, organizationID, tableID)

	return table, errors.Wrap(err, "table select error")
}

// Create inserts table into system.
func (s *TableService) Create(userID string, table model.Table) (string, error) {
	tableID, err := s.repo.Create(userID, table)

	return tableID, errors.Wrap(err, "table create error")
}

// Update updates table by id in the system.
func (s *TableService) Update(userID string, input model.UpdateTableInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "table update error")
}

// Delete deletes table by id from the system.
func (s *TableService) Delete(userID string, organizationID string, tableID string) error {
	err := s.repo.Delete(userID, organizationID, tableID)

	return errors.Wrap(err, "table delete error")
}
