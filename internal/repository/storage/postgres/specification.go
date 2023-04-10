package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// SpecificationGetAll selects all specifications from database.
func (r *Repository) SpecificationGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var specifications []specification.Specification

	qry, args, err := r.specificationGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.SelectContext(ctx, &specifications, qry, args...)

	return specifications, errors.Wrap(err, "specifications select query error")
}

// specificationGetAllQuery creates sql query.
func (r *Repository) specificationGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"id", "item_id", "organization_id", "name_tm", "name_ru", "name_tr", "name_en",
		"description_tm", "description_ru", "description_tr", "description_en", "value").
		From(specificationTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.Or{
			squirrel.Like{"name_tm": search},
			squirrel.Like{"name_ru": search},
			squirrel.Like{"name_tr": search},
			squirrel.Like{"name_en": search},
			squirrel.Like{"description_tm": search},
			squirrel.Like{"description_ru": search},
			squirrel.Like{"description_tr": search},
			squirrel.Like{"description_en": search},
			squirrel.Like{"value": search},
		})
	}

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortSpecification)...)
	} else {
		builder = builder.OrderBy("name_tm ASC")
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

// SpecificationGetOne select specification by id from database.
func (r *Repository) SpecificationGetOne(ctxr context.Context, meta query.MetaData, specificationID string) (specification.Specification, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var spec specification.Specification

	builder := r.genSQL.Select("id", "item_id", "organization_id", "name_tm", "name_ru", "name_tr", "name_en",
		"description_tm", "description_ru", "description_tr", "description_en", "value").
		From(specificationTable).
		Where(squirrel.Eq{"id": specificationID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return spec, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.GetContext(ctx, &spec, qry, args...)

	return spec, errors.Wrap(err, "specification select query error")
}

// SpecificationCreate insert specification into database.
func (r *Repository) SpecificationCreate(ctxr context.Context, _ query.MetaData, input specification.CreateSpecificationInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var specificationID string

	builder := r.genSQL.Insert(specificationTable).
		Columns(
			"item_id", "organization_id", "name_tm", "name_ru", "name_tr", "name_en", "description_tm",
			"description_ru", "description_tr", "description_en", "value").
		Values(input.ItemID, input.OrganizationID, input.NameTm, input.NameRu, input.NameTr, input.NameEn,
			input.DescriptionTm, input.DescriptionRu, input.DescriptionTr, input.DescriptionEn, input.Value).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.databaseMaster.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&specificationID)

	return specificationID, errors.Wrap(err, "specification create query error")
}

// SpecificationUpdate updates specification by id in database.
func (r *Repository) SpecificationUpdate(ctxr context.Context, meta query.MetaData, input specification.UpdateSpecificationInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(specificationTable)

	if input.ItemID != nil {
		builder = builder.Set("item_id", *input.ItemID)
	}

	if input.OrganizationID != nil {
		builder = builder.Set("organization_id", *input.OrganizationID)
	}

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

	if input.DescriptionTm != nil {
		builder = builder.Set("description_tm", *input.DescriptionTm)
	}

	if input.DescriptionRu != nil {
		builder = builder.Set("description_ru", *input.DescriptionRu)
	}

	if input.DescriptionTr != nil {
		builder = builder.Set("description_tr", *input.DescriptionTr)
	}

	if input.DescriptionEn != nil {
		builder = builder.Set("description_en", *input.DescriptionEn)
	}

	if input.Value != nil {
		builder = builder.Set("value", *input.Value)
	}

	builder = builder.Where(squirrel.Eq{"id": *input.ID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "specification update query error")
}

// SpecificationDelete deletes specification by id from database.
func (r *Repository) SpecificationDelete(ctxr context.Context, meta query.MetaData, specificationID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Delete(specificationTable).Where(squirrel.Eq{"id": specificationID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "specification delete query error")
}
