package postgres

import (
	"fmt"
	"strings"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// SpecificationGetAll selects all specifications from database.
func (r *Repository) SpecificationGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var specifications []specification.Specification

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, "+
			"description_tr, description_en, value FROM %s WHERE organization_id = $1 ",
		specificationTable,
	)

	err := r.database.SelectContext(ctx, &specifications, query, meta.OrganizationID)

	return specifications, errors.Wrap(err, "specifications select query error")
}

// SpecificationGetOne select specification by id from database.
func (r *Repository) SpecificationGetOne(ctxr context.Context, meta query.MetaData, specificationID string) (specification.Specification, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var spec specification.Specification

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, name_tm, name_ru, name_tr, name_en, description_tm, description_ru, "+
			"description_tr, description_en, value FROM %s WHERE organization_id = $1 AND id = $2 ",
		specificationTable,
	)

	err := r.database.GetContext(ctx, &spec, query, meta.OrganizationID, specificationID)

	return spec, errors.Wrap(err, "specification select query error")
}

// SpecificationCreate insert specification into database.
func (r *Repository) SpecificationCreate(ctxr context.Context, meta query.MetaData, input specification.CreateSpecificationInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.SpecificationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var specificationID string

	query := fmt.Sprintf(
		"INSERT INTO %s (item_id, organization_id, name_tm, name_ru, name_tr, name_en,  description_tm, description_ru, "+
			"description_tr, description_en, value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		specificationTable,
	)

	row := r.database.QueryRowContext(
		ctx,
		query,
		input.ItemID,
		input.OrganizationID,
		input.NameTm,
		input.NameRu,
		input.NameTr,
		input.NameEn,
		input.DescriptionTm,
		input.DescriptionRu,
		input.DescriptionTr,
		input.DescriptionEn,
		input.Value,
	)

	err := row.Scan(&specificationID)

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

	setValues := make([]string, 0, 12)
	args := make([]interface{}, 0, 12)
	argID := 1

	if input.ItemID != nil {
		setValues = append(setValues, fmt.Sprintf("item_id=$%d", argID))
		args = append(args, *input.ItemID)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

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

	if input.DescriptionTm != nil {
		setValues = append(setValues, fmt.Sprintf("description_tm=$%d", argID))
		args = append(args, *input.DescriptionTm)
		argID++
	}

	if input.DescriptionRu != nil {
		setValues = append(setValues, fmt.Sprintf("description_ru=$%d", argID))
		args = append(args, *input.DescriptionRu)
		argID++
	}

	if input.DescriptionTr != nil {
		setValues = append(setValues, fmt.Sprintf("description_tr=$%d", argID))
		args = append(args, *input.DescriptionTr)
		argID++
	}

	if input.DescriptionEn != nil {
		setValues = append(setValues, fmt.Sprintf("description_en=$%d", argID))
		args = append(args, *input.DescriptionEn)
		argID++
	}

	if input.Value != nil {
		setValues = append(setValues, fmt.Sprintf("value=$%d", argID))
		args = append(args, *input.Value)
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE organization_id = '%s' AND id = '%s'",
		specificationTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.database.ExecContext(ctx, query, args...)

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

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND organization_id = $2", specificationTable)

	_, err := r.database.ExecContext(ctx, query, specificationID, meta.OrganizationID)

	return errors.Wrap(err, "specification delete query error")
}
