package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// TableGetAll selects all tables from database.
func (r *Repository) TableGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]table.Table, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.TableGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var tables []table.Table

	qry, args, err := r.tableGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &tables, qry, args...)

	return tables, errors.Wrap(err, "tables select query error")
}

// tableGetAllQuery creates sql query.
func (r *Repository) tableGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "name", "organization_id").From(tableTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"is_deleted": false},
			squirrel.Like{"name": search},
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
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortTable)...)
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

// TableGetOne select table by id from database.
func (r *Repository) TableGetOne(ctxr context.Context, meta query.MetaData, tableID string) (table.Table, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.TableGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var tble table.Table

	builder := r.genSQL.Select("id", "name", "organization_id").From(tableTable).
		Where(squirrel.Eq{"is_deleted": false, "id": tableID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return tble, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &tble, qry, args...)

	return tble, errors.Wrap(err, "table select query error")
}

// TableCreate insert table into database.
func (r *Repository) TableCreate(ctxr context.Context, meta query.MetaData, input table.CreateTableInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.TableCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var tableID string

	builder := r.genSQL.Insert(tableTable).
		Columns("name", "organization_id", "user_created").
		Values(input.Name, input.OrganizationID, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)
	err = row.Scan(&tableID)

	return tableID, errors.Wrap(err, "table create query error")
}

// TableUpdate updates table by id in database.
func (r *Repository) TableUpdate(ctxr context.Context, meta query.MetaData, input table.UpdateTableInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.TableUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(tableTable)

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
	}

	if input.OrganizationID != nil {
		builder = builder.Set("organization_id", *input.OrganizationID)
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

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "table update query error")
}

// TableDelete deletes table by id from database.
func (r *Repository) TableDelete(ctxr context.Context, meta query.MetaData, tableID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.TableDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(tableTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": tableID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "table delete query error")
}
