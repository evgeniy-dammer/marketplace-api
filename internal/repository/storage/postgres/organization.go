package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/pkg/errors"
)

// OrganizationGetAll selects all organizations from database.
func (r *Repository) OrganizationGetAll(userID string) ([]organization.Organization, error) {
	var organizations []organization.Organization

	query := fmt.Sprintf("SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1",
		organizationTable)

	err := r.db.Select(&organizations, query, userID)

	return organizations, errors.Wrap(err, "organizations select query error")
}

// OrganizationGetOne select organization by id from database.
func (r *Repository) OrganizationGetOne(userID string, organizationID string) (organization.Organization, error) {
	var organization organization.Organization

	query := fmt.Sprintf(
		"SELECT id, name, user_id, address, phone FROM %s WHERE is_deleted = false AND user_id = $1 AND id = $2",
		organizationTable)
	err := r.db.Get(&organization, query, userID, organizationID)

	return organization, errors.Wrap(err, "organization select query error")
}

// OrganizationCreate insert organization into database.
func (r *Repository) OrganizationCreate(userID string, organization organization.Organization) (string, error) {
	var organizationID string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		organizationTable)

	row := r.db.QueryRow(createUserQuery, organization.Name, userID, organization.Address, organization.Phone, userID)

	err := row.Scan(&organizationID)

	return organizationID, errors.Wrap(err, "organization create query error")
}

// OrganizationUpdate updates organization by id in database.
func (r *Repository) OrganizationUpdate(userID string, input organization.UpdateOrganizationInput) error {
	setValues := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)
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

// OrganizationDelete deletes organization by id from database.
func (r *Repository) OrganizationDelete(userID string, organizationID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3",
		organizationTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, organizationID)

	return errors.Wrap(err, "organization delete query error")
}
