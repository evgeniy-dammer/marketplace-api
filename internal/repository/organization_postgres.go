package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// OrganizationPostgresql repository.
type OrganizationPostgresql struct {
	db *sqlx.DB
}

// NewOrganizationPostgresql is a constructor for OrganizationPostgresql.
func NewOrganizationPostgresql(db *sqlx.DB) *OrganizationPostgresql {
	return &OrganizationPostgresql{db: db}
}

// GetAll selects all organizations from database.
func (r *OrganizationPostgresql) GetAll(userID string) ([]model.Organization, error) {
	var organizations []model.Organization

	query := fmt.Sprintf("SELECT id, name, user_id, address, phone FROM %s WHERE user_id = $1",
		organisationTable)

	err := r.db.Select(&organizations, query, userID)

	return organizations, errors.Wrap(err, "organizations select query error")
}

// GetOne select organization by id from database.
func (r *OrganizationPostgresql) GetOne(userID string, organizationID string) (model.Organization, error) {
	var organization model.Organization

	query := fmt.Sprintf(
		"SELECT id, name, user_id, address, phone FROM %s WHERE user_id = $1 AND id = $2",
		organisationTable)
	err := r.db.Get(&organization, query, userID, organizationID)

	return organization, errors.Wrap(err, "organization select query error")
}

// Create insert organization into database.
func (r *OrganizationPostgresql) Create(userID string, organization model.Organization) (string, error) {
	var organizationID string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		organisationTable)

	row := r.db.QueryRow(createUserQuery, organization.Name, userID, organization.Address, organization.Phone)

	err := row.Scan(&organizationID)

	return organizationID, errors.Wrap(err, "organization create query error")
}

// Update updates organization by id in database.
func (r *OrganizationPostgresql) Update(userID string, organizationID string, input model.UpdateOrganizationInput) error { //nolint:lll
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argID))
		args = append(args, *input.Address)
		argID++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argID))
		args = append(args, *input.Phone)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = '%s' AND user_id = '%s'",
		organisationTable, setQuery, organizationID, userID)

	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "organization update query error")
}

// Delete deletes organization by id from database.
func (r *OrganizationPostgresql) Delete(userID string, organizationID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", organisationTable)
	_, err := r.db.Exec(query, organizationID, userID)

	return errors.Wrap(err, "organization delete query error")
}
