package repository

import (
	"fmt"
	"strings"
	"time"

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

	query := fmt.Sprintf(
		"SELECT id, name, parent_id, level, organization_id FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1",
		categoryTable)

	err := r.db.Select(&categories, query, organizationID)

	return categories, errors.Wrap(err, "categories select query error")
}

// GetOne select category by id from database.
func (r *CategoryPostgresql) GetOne(userID string, organizationID string, categoryID string) (model.Category, error) {
	var user model.Category

	query := fmt.Sprintf(
		"SELECT id, name, parent_id, level, organization_id FROM %s "+
			"WHERE is_deleted = false AND id = $1 AND organization_id = $2",
		categoryTable,
	)

	err := r.db.Get(&user, query, categoryID, organizationID)

	return user, errors.Wrap(err, "category select query error")
}

// Create insert category into database.
func (r *CategoryPostgresql) Create(userID string, category model.Category) (string, error) {
	var categoryID string

	query := fmt.Sprintf(
		"INSERT INTO %s (name, parent_id, level, organization_id, user_created) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		categoryTable,
	)

	row := r.db.QueryRow(query, category.Name, category.Parent, category.Level, category.OrganizationID, userID)
	err := row.Scan(&categoryID)

	return categoryID, errors.Wrap(err, "category create query error")
}

// Update updates category by id in database.
func (r *CategoryPostgresql) Update(userID string, input model.UpdateCategoryInput) error { //nolint:lll
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
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND id = '%s' AND organization_id = '%s'",
		categoryTable, setQuery, *input.ID, *input.OrganizationID)
	_, err := r.db.Exec(query, args...)

	return errors.Wrap(err, "category update query error")
}

// Delete deletes category by id from database.
func (r *CategoryPostgresql) Delete(userID string, organizationID string, categoryID string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		categoryTable,
	)

	_, err := r.db.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, categoryID, organizationID)

	return errors.Wrap(err, "category delete query error")
}
