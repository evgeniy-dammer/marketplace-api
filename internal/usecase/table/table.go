package table

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/table"
	"github.com/pkg/errors"
)

// TableGetAll returns all tables from the system.
func (s *UseCase) TableGetAll(userID string, organizationID string) ([]table.Table, error) {
	tables, err := s.adapterStorage.TableGetAll(userID, organizationID)

	return tables, errors.Wrap(err, "tables select error")
}

// TableGetOne returns table by id from the system.
func (s *UseCase) TableGetOne(userID string, organizationID string, tableID string) (table.Table, error) {
	tbl, err := s.adapterStorage.TableGetOne(userID, organizationID, tableID)

	return tbl, errors.Wrap(err, "table select error")
}

// TableCreate inserts table into system.
func (s *UseCase) TableCreate(userID string, table table.Table) (string, error) {
	tableID, err := s.adapterStorage.TableCreate(userID, table)

	return tableID, errors.Wrap(err, "table create error")
}

// TableUpdate updates table by id in the system.
func (s *UseCase) TableUpdate(userID string, input table.UpdateTableInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.TableUpdate(userID, input)

	return errors.Wrap(err, "table update error")
}

// TableDelete deletes table by id from the system.
func (s *UseCase) TableDelete(userID string, organizationID string, tableID string) error {
	err := s.adapterStorage.TableDelete(userID, organizationID, tableID)

	return errors.Wrap(err, "table delete error")
}
