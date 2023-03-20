package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/image"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// ImageGetAll selects all images from database.
func (r *Repository) ImageGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]image.Image, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ImageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var images []image.Image

	qry, args, err := r.imageGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &images, qry, args...)

	return images, errors.Wrap(err, "images select query error")
}

// tableGetAllQuery creates sql query.
func (r *Repository) imageGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"id", "object_id", "type", "origin", "middle", "small", "organization_id", "is_main").
		From(imageTable).
		Where(squirrel.Eq{"is_deleted": false})

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
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortImage)...)
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

// ImageGetOne select image by id from database.
func (r *Repository) ImageGetOne(ctxr context.Context, meta query.MetaData, imageID string) (image.Image, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ImageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var img image.Image

	builder := r.genSQL.Select(
		"id", "object_id", "type", "origin", "middle", "small", "organization_id", "is_main").
		From(imageTable).
		Where(squirrel.Eq{"is_deleted": false, "id": imageID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return img, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &img, qry, args...)

	return img, errors.Wrap(err, "image select query error")
}

// ImageCreate insert image into database.
func (r *Repository) ImageCreate(ctxr context.Context, meta query.MetaData, input image.CreateImageInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ImageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var imageID string

	builder := r.genSQL.Insert(imageTable).
		Columns("object_id", "type", "origin", "middle", "small", "organization_id", "is_main", "user_created").
		Values(input.ObjectID, input.ObjectType, input.Origin, input.Middle, input.Small, input.OrganizationID,
			input.IsMain, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&imageID)

	return imageID, errors.Wrap(err, "image create query error")
}

// ImageUpdate updates image by id in database.
func (r *Repository) ImageUpdate(ctxr context.Context, meta query.MetaData, input image.UpdateImageInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ImageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(imageTable)

	if input.ObjectID != nil {
		builder = builder.Set("object_id", *input.ObjectID)
	}

	if input.ObjectType != nil {
		builder = builder.Set("type", *input.ObjectType)
	}

	if input.Origin != nil {
		builder = builder.Set("origin", *input.Origin)
	}

	if input.Middle != nil {
		builder = builder.Set("middle", *input.Middle)
	}

	if input.Small != nil {
		builder = builder.Set("small", *input.Small)
	}

	if input.OrganizationID != nil {
		builder = builder.Set("organization_id", *input.OrganizationID)
	}

	if input.IsMain != nil {
		builder = builder.Set("is_main", *input.IsMain)
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

	return errors.Wrap(err, "image update query error")
}

// ImageDelete deletes image by id from database.
func (r *Repository) ImageDelete(ctxr context.Context, meta query.MetaData, imageID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.ImageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(imageTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": imageID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "image delete query error")
}
