package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/organization"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// OrganizationGetAll selects all organizations from database.
func (r *Repository) OrganizationGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrganizationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var organizations []organization.Organization

	qry, args, err := r.organizationGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &organizations, qry, args...)

	return organizations, errors.Wrap(err, "organizations select query error")
}

// organizationGetAllQuery creates sql query.
func (r *Repository) organizationGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "name", "user_id", "address", "phone").From(organizationTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"is_deleted": false},
			squirrel.Or{
				squirrel.Like{"name": search},
				squirrel.Like{"address": search},
				squirrel.Like{"phone": search},
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

	if meta.RoleName != vendorRole {
		builder = builder.Where(squirrel.Eq{"user_id": meta.UserID})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortOrganization)...)
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

// OrganizationGetOne select organization by id from database.
func (r *Repository) OrganizationGetOne(ctxr context.Context, meta query.MetaData, organizationID string) (organization.Organization, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrganizationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var org organization.Organization

	builder := r.genSQL.Select("id", "name", "user_id", "address", "phone").From(organizationTable).
		Where(squirrel.Eq{"is_deleted": false, "id": organizationID})

	if meta.RoleName != vendorRole {
		builder = builder.Where(squirrel.Eq{"user_id": meta.UserID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return org, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &org, qry, args...)

	return org, errors.Wrap(err, "organization select query error")
}

// OrganizationCreate insert organization into database.
func (r *Repository) OrganizationCreate(ctxr context.Context, meta query.MetaData, input organization.CreateOrganizationInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrganizationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var organizationID string

	builder := r.genSQL.Insert(organizationTable).
		Columns("name", "user_id", "address", "phone", "user_created").
		Values(input.Name, meta.UserID, input.Address, input.Phone, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&organizationID)

	return organizationID, errors.Wrap(err, "organization create query error")
}

// OrganizationUpdate updates organization by id in database.
func (r *Repository) OrganizationUpdate(ctxr context.Context, meta query.MetaData, input organization.UpdateOrganizationInput) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrganizationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(organizationTable)

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
	}

	if input.Address != nil {
		builder = builder.Set("address", *input.Address)
	}

	if input.Phone != nil {
		builder = builder.Set("phone", *input.Phone)
	}

	builder = builder.Set("user_updated", meta.UserID).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": *input.ID})

	if meta.RoleName != vendorRole {
		builder = builder.Where(squirrel.Eq{"user_id": meta.UserID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "organization update query error")
}

// OrganizationDelete deletes organization by id from database.
func (r *Repository) OrganizationDelete(ctxr context.Context, meta query.MetaData, organizationID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.OrganizationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(organizationTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": organizationID})

	if meta.RoleName != vendorRole {
		builder = builder.Where(squirrel.Eq{"user_id": meta.UserID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "organization delete query error")
}
