package repository

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CategoryPostgresql repository.
type CategoryPostgresql struct {
	db *sqlx.DB
}

// NewCategoryPostgresql is a constructor for CategoryPostgresql.
func NewCategoryPostgresql(db *sqlx.DB) *CategoryPostgresql {
	return &CategoryPostgresql{db: db}
}

// GetAll selects all categories from database.
func (r *CategoryPostgresql) GetAll(userID string, organizationID string) ([]model.Category, error) {
	var categories []model.Category

	query := fmt.Sprintf("SELECT id, name, parent_id, level, organisation_id FROM %s WHERE organisation_id = $1",
		categoryTable)

	err := r.db.Select(&categories, query, organizationID)

	return categories, errors.Wrap(err, "categories select query error")
}

// GetOne select category by id from database.
func (r *CategoryPostgresql) GetOne(userID string, organizationID string, categoryID string) (model.Category, error) {
	var user model.Category

	query := fmt.Sprintf(
		"SELECT id, name, parent_id, level, organisation_id FROM %s WHERE id = $1 AND organisation_id = $2",
		categoryTable,
	)

	err := r.db.Get(&user, query, categoryID, organizationID)

	return user, errors.Wrap(err, "category select query error")
}

// Create insert category into database.
func (r *CategoryPostgresql) Create(userID string, organizationID string, category model.Category) (string, error) {
	var categoryID string

	query := fmt.Sprintf(
		"INSERT INTO %s (name, parent_id, level, organisation_id) VALUES ($1, $2, $3, $4) RETURNING id",
		categoryTable,
	)

	row := r.db.QueryRow(query, category.Name, category.Parent, category.Level, organizationID)
	err := row.Scan(&categoryID)

	return categoryID, errors.Wrap(err, "category create query error")
}

// Update updates category by id in database.
func (r *CategoryPostgresql) Update(userID string, organizationID string, categoryID string, input model.UpdateCategoryInput) error { //nolint:lll
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.Parent != nil {
		setValues = append(setValues, fmt.Sprintf("parent_id=$%d", argID))
		args = append(args, *input.Parent)
		argID++
	}

	if input.Level != nil {
		setValues = append(setValues, fmt.Sprintf("level=$%d", argID))
		args = append(args, *input.Level)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s' AND organisation_id = '%s'",
		categoryTable, setQuery, categoryID, organizationID)
	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "category update query error")
}

// Delete deletes category by id from database.
func (r *CategoryPostgresql) Delete(userID string, organizationID string, categoryID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND organisation_id = $2", categoryTable)
	_, err := r.db.Exec(query, categoryID, organizationID)

	return errors.Wrap(err, "category delete query error")
}
