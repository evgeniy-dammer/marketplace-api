package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TablePostgresql repository.
type TablePostgresql struct {
	db *sqlx.DB
}

// NewTablePostgresql is a constructor for TablePostgresql.
func NewTablePostgresql(db *sqlx.DB) *TablePostgresql {
	return &TablePostgresql{db: db}
}

// GetAll selects all tables from database.
func (r *TablePostgresql) GetAll(userID string, organizationID string) ([]model.Table, error) {
	var tables []model.Table

	query := fmt.Sprintf("SELECT id, name, organisation_id FROM %s WHERE organisation_id = $1 ",
		tableTable)
	err := r.db.Select(&tables, query, organizationID)

	return tables, errors.Wrap(err, "tables select query error")
}

// GetOne select table by id from database.
func (r *TablePostgresql) GetOne(userID string, organizationID string, tableID string) (model.Table, error) {
	var table model.Table

	query := fmt.Sprintf(
		"SELECT id, name, organisation_id FROM %s WHERE organisation_id = $1 AND id = $2 ",
		tableTable,
	)
	err := r.db.Get(&table, query, organizationID, tableID)

	return table, errors.Wrap(err, "table select query error")
}

// Create insert table into database.
func (r *TablePostgresql) Create(userID string, organizationID string, table model.Table) (string, error) {
	var tableID string

	query := fmt.Sprintf("INSERT INTO %s (name, organisation_id) VALUES ($1, $2) RETURNING id",
		tableTable)
	row := r.db.QueryRow(query, table.Name, organizationID)
	err := row.Scan(&tableID)

	return tableID, errors.Wrap(err, "table create query error")
}

// Update updates table by id in database.
func (r *TablePostgresql) Update(userID string, organizationID string, tableID string, input model.UpdateTableInput) error { //nolint:lll
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.OrganisationID != nil {
		setValues = append(setValues, fmt.Sprintf("organisation_id=$%d", argID))
		args = append(args, *input.OrganisationID)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE organisation_id = '%s' AND id = '%s'",
		tableTable, setQuery, organizationID, tableID)
	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "table update query error")
}

// Delete deletes table by id from database.
func (r *TablePostgresql) Delete(userID string, organizationID string, tableID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE organisation_id = $1 AND id = $2", tableTable)
	_, err := r.db.Exec(query, organizationID, tableID)

	return errors.Wrap(err, "table delete query error")
}
