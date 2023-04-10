package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// CategoryGetAll selects all categories from database.
func (r *Repository) CategoryGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]category.Category, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CategoryGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var categories []category.Category

	qry, args, err := r.categoryGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.SelectContext(ctx, &categories, qry, args...)

	return categories, errors.Wrap(err, "categories select query error")
}

// categoryGetAllQuery creates sql query.
func (r *Repository) categoryGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"id", "name_tm", "name_ru", "name_tr", "name_en", "parent_id", "level", "organization_id").
		From(categoryTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"is_deleted": false},
			squirrel.Or{
				squirrel.Like{"name_tm": search},
				squirrel.Like{"name_ru": search},
				squirrel.Like{"name_tr": search},
				squirrel.Like{"name_en": search},
			},
		})
	} else {
		builder = builder.Where(squirrel.Eq{"is_deleted": false})
	}

	switch {
	case !params.StartDate.IsZero() && params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": time.Now().Format("2006-01-02 15:04:05")},
		})
	case params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": time.Now().Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	case !params.StartDate.IsZero() && !params.EndDate.IsZero():
		builder = builder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": params.StartDate.Format("2006-01-02 15:04:05")},
			squirrel.LtOrEq{"created_at": params.EndDate.Format("2006-01-02 15:04:05")},
		})
	}

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortCategory)...)
	} else {
		builder = builder.OrderBy("created_at DESC")
	}

	if params.Pagination.Limit > 0 {
		builder = builder.Limit(params.Pagination.Limit)
	}

	if params.Pagination.Offset > 0 {
		builder = builder.Offset(params.Pagination.Offset)
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", nil, errors.Wrap(err, "unable to build a query string")
	}

	return qry, args, nil
}

// CategoryGetOne select category by id from database.
func (r *Repository) CategoryGetOne(ctxr context.Context, meta query.MetaData, categoryID string) (category.Category, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CategoryGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var ctgry category.Category

	builder := r.genSQL.Select(
		"id", "name_tm", "name_ru", "name_tr", "name_en", "parent_id", "level", "organization_id").
		From(categoryTable).
		Where(squirrel.Eq{"is_deleted": false, "id": categoryID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return ctgry, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.GetContext(ctx, &ctgry, qry, args...)

	return ctgry, errors.Wrap(err, "category select query error")
}

// CategoryCreate insert category into database.
func (r *Repository) CategoryCreate(ctxr context.Context, meta query.MetaData, input category.CreateCategoryInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CategoryCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var categoryID string

	builder := r.genSQL.Insert(categoryTable).
		Columns(
			"name_tm", "name_ru", "name_tr", "name_en", "parent_id", "level", "organization_id", "user_created").
		Values(
			input.NameTm, input.NameRu, input.NameTr, input.NameEn, input.Parent, input.Level,
			input.OrganizationID, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.databaseMaster.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&categoryID)

	return categoryID, errors.Wrap(err, "category create query error")
}

// CategoryUpdate updates category by id in database.
func (r *Repository) CategoryUpdate(ctxr context.Context, meta query.MetaData, input category.UpdateCategoryInput) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CategoryUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(categoryTable)

	if input.NameTm != nil {
		builder = builder.Set("name_tm", *input.NameTm)
	}

	if input.NameRu != nil {
		builder = builder.Set("name_ru", *input.NameRu)
	}

	if input.NameTr != nil {
		builder = builder.Set("name_tr", *input.NameTr)
	}

	if input.NameEn != nil {
		builder = builder.Set("name_en", *input.NameEn)
	}

	if input.Parent != nil {
		builder = builder.Set("parent_id", *input.Parent)
	}

	if input.Level != nil {
		builder = builder.Set("level", *input.Level)
	}

	builder = builder.Set("user_updated", meta.UserID).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": *input.ID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "category update query error")
}

// CategoryDelete deletes category by id from database.
func (r *Repository) CategoryDelete(ctxr context.Context, meta query.MetaData, categoryID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CategoryDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(categoryTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": categoryID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "category delete query error")
}
