package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

// OrganizationPostgresql repository
type OrganizationPostgresql struct {
	db *sqlx.DB
}

// NewOrganizationPostgresql is a constructor for OrganizationPostgresql
func NewOrganizationPostgresql(db *sqlx.DB) *OrganizationPostgresql {
	return &OrganizationPostgresql{db: db}
}

// GetAll selects all organizations from database
func (r *OrganizationPostgresql) GetAll(userId string) ([]model.Organization, error) {
	var organizations []model.Organization

	query := fmt.Sprintf("SELECT id, name, user_id, address, phone FROM %s WHERE user_id = $1",
		organisationTable)

	err := r.db.Select(&organizations, query, userId)

	return organizations, err
}

// GetOne select organization by id from database
func (r *OrganizationPostgresql) GetOne(userId string, organizationId string) (model.Organization, error) {
	var organization model.Organization

	query := fmt.Sprintf(
		"SELECT id, name, user_id, address, phone FROM %s WHERE user_id = $1 AND id = $2",
		organisationTable)
	err := r.db.Get(&organization, query, userId, organizationId)

	return organization, err
}

// Create insert organization into database
func (r *OrganizationPostgresql) Create(userId string, organization model.Organization) (string, error) {
	var id string

	createUserQuery := fmt.Sprintf(
		"INSERT INTO %s (name, user_id, address, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		organisationTable)

	row := r.db.QueryRow(createUserQuery, organization.Name, userId, organization.Address, organization.Phone)

	err := row.Scan(&id)

	return id, err
}

// Update updates organization by id in database
func (r *OrganizationPostgresql) Update(userId string, organizationId string, input model.UpdateOrganizationInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *input.Address)
		argId++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *input.Phone)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = '%s' AND user_id = '%s'",
		organisationTable, setQuery, organizationId, userId)

	_, err := r.db.Exec(query, args...)

	return err
}

// Delete deletes organization by id from database
func (r *OrganizationPostgresql) Delete(userId string, organizationId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", organisationTable)
	_, err := r.db.Exec(query, organizationId, userId)

	return err
}
