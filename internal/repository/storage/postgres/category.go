package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/pkg/errors"
)

// CategoryGetAll selects all categories from database.
func (r *Repository) CategoryGetAll(ctx context.Context, userID string, organizationID string) ([]category.Category, error) {
	ctx = ctx.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	var categories []category.Category

	query := fmt.Sprintf(
		"SELECT id, name_tm, name_ru, name_tr, name_en, parent_id, level, organization_id FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1",
		categoryTable)

	err := r.database.Select(&categories, query, organizationID)

	return categories, errors.Wrap(err, "categories select query error")
}

// CategoryGetOne select category by id from database.
func (r *Repository) CategoryGetOne(ctx context.Context, userID string, organizationID string, categoryID string) (category.Category, error) {
	ctx = ctx.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	var user category.Category

	query := fmt.Sprintf(
		"SELECT id, name_tm, name_ru, name_tr, name_en, parent_id, level, organization_id FROM %s "+
			"WHERE is_deleted = false AND id = $1 AND organization_id = $2",
		categoryTable,
	)

	err := r.database.Get(&user, query, categoryID, organizationID)

	return user, errors.Wrap(err, "category select query error")
}

// CategoryCreate insert category into database.
func (r *Repository) CategoryCreate(ctx context.Context, userID string, input category.CreateCategoryInput) (string, error) {
	ctx = ctx.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	var categoryID string

	query := fmt.Sprintf(
		"INSERT INTO %s (name_tm, name_ru, name_tr, name_en, parent_id, level, organization_id, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		categoryTable,
	)

	row := r.database.QueryRow(
		query,
		input.NameTm,
		input.NameRu,
		input.NameTr,
		input.NameEn,
		input.Parent,
		input.Level,
		input.OrganizationID,
		userID,
	)

	err := row.Scan(&categoryID)

	return categoryID, errors.Wrap(err, "category create query error")
}

// CategoryUpdate updates category by id in database.
func (r *Repository) CategoryUpdate(ctx context.Context, userID string, input category.UpdateCategoryInput) error {
	ctx = ctx.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	setValues := make([]string, 0, 8)
	args := make([]interface{}, 0, 8)
	argID := 1

	if input.NameTm != nil {
		setValues = append(setValues, fmt.Sprintf("name_tm=$%d", argID))
		args = append(args, *input.NameTm)
		argID++
	}

	if input.NameRu != nil {
		setValues = append(setValues, fmt.Sprintf("name_ru=$%d", argID))
		args = append(args, *input.NameRu)
		argID++
	}

	if input.NameTr != nil {
		setValues = append(setValues, fmt.Sprintf("name_tr=$%d", argID))
		args = append(args, *input.NameTr)
		argID++
	}

	if input.NameEn != nil {
		setValues = append(setValues, fmt.Sprintf("name_en=$%d", argID))
		args = append(args, *input.NameEn)
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
	_, err := r.database.Exec(query, args...)

	return errors.Wrap(err, "category update query error")
}

// CategoryDelete deletes category by id from database.
func (r *Repository) CategoryDelete(ctx context.Context, userID string, organizationID string, categoryID string) error {
	ctx = ctx.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		categoryTable,
	)

	_, err := r.database.Exec(query, time.Now().Format("2006-01-02 15:04:05"), userID, categoryID, organizationID)

	return errors.Wrap(err, "category delete query error")
}
