package repository

import (
	"fmt"
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

// CategoryPostgresql repository
type CategoryPostgresql struct {
	db *sqlx.DB
}

// NewCategoryPostgresql is a constructor for CategoryPostgresql
func NewCategoryPostgresql(db *sqlx.DB) *CategoryPostgresql {
	return &CategoryPostgresql{db: db}
}

// GetAll selects all categories from database
func (r *CategoryPostgresql) GetAll(userId string, organizationId string) ([]model.Category, error) {
	var categories []model.Category

	query := fmt.Sprintf("SELECT id, name, parent_id, level, organisation_id FROM %s WHERE organisation_id = $1",
		categoryTable)

	err := r.db.Select(&categories, query, organizationId)

	return categories, err
}

// GetOne select category by id from database
func (r *CategoryPostgresql) GetOne(userId string, organizationId string, categoryId string) (model.Category, error) {
	var user model.Category

	query := fmt.Sprintf("SELECT id, name, parent_id, level, organisation_id FROM %s WHERE id = $1 AND organisation_id = $2",
		categoryTable)
	err := r.db.Get(&user, query, categoryId, organizationId)

	return user, err
}

// Create insert category into database
func (r *CategoryPostgresql) Create(userId string, organizationId string, category model.Category) (string, error) {
	var id string

	query := fmt.Sprintf(
		"INSERT INTO %s (name, parent_id, level, organisation_id) VALUES ($1, $2, $3, $4) RETURNING id",
		categoryTable,
	)

	row := r.db.QueryRow(query, category.Name, category.Parent, category.Level, organizationId)

	err := row.Scan(&id)

	return id, err
}

// Update updates category by id in database
func (r *CategoryPostgresql) Update(userId string, organizationId string, categoryId string, input model.UpdateCategoryInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Parent != nil {
		setValues = append(setValues, fmt.Sprintf("parent_id=$%d", argId))
		args = append(args, *input.Parent)
		argId++
	}

	if input.Level != nil {
		setValues = append(setValues, fmt.Sprintf("level=$%d", argId))
		args = append(args, *input.Level)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = '%s' AND organisation_id = '%s'",
		categoryTable, setQuery, categoryId, organizationId)
	_, err := r.db.Exec(query, args...)

	return err
}

// Delete deletes category by id from database
func (r *CategoryPostgresql) Delete(userId string, organizationId string, categoryId string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND organisation_id = $2", categoryTable)
	_, err := r.db.Exec(query, categoryId, organizationId)

	return err
}
