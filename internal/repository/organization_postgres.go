package repository

import (
	"fmt"
	"strings"
	"time"

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

	query := fmt.Sprintf("SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1",
		organizationTable)

	err := r.db.Select(&organizations, query, userID)

	return organizations, errors.Wrap(err, "organizations select query error")
}

// GetOne select organization by id from database.
func (r *OrganizationPostgresql) GetOne(userID string, organizationID string) (model.Organization, error) {
	var organization model.Organization

	query := fmt.Sprintf(
		"SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1 AND id = $2",
		organizationTable)
	err := r.db.Get(&organization, query, userID, organizationID)

	return organization, errors.Wrap(err, "organization select query error")
}

// Create insert organization into database.
func (r *OrganizationPostgresql) Create(userID string, organization model.Organization) (string, error) {
	var organizationID string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		organizationTable)

	row := r.db.QueryRow(createUserQuery, organization.Name, userID, organization.Address, organization.Phone, userID)

	err := row.Scan(&organizationID)

	return organizationID, errors.Wrap(err, "organization create query error")
}

// Update updates organization by id in database.
func (r *OrganizationPostgresql) Update(userID string, input model.UpdateOrganizationInput) error {
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
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE is_deleted = false AND id = '%s' AND user_id = '%s'",
		organizationTable, setQuery, *input.ID, userID)

	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "organization update query error")
}

// Delete deletes organization by id from database.
func (r *OrganizationPostgresql) Delete(userID string, organizationID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3",
		organizationTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, organizationID)

	return errors.Wrap(err, "organization delete query error")
}
