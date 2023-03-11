package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/image"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// ImageGetAll selects all images from database.
func (r *Repository) ImageGetAll(ctxr context.Context, userID string, organizationID string) ([]image.Image, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Database.ImageGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var images []image.Image

	query := fmt.Sprintf(
		"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 ",
		imageTable,
	)

	err := r.database.SelectContext(ctx, &images, query, organizationID)

	return images, errors.Wrap(err, "images select query error")
}

// ImageGetOne select image by id from database.
func (r *Repository) ImageGetOne(ctxr context.Context, userID string, organizationID string, imageID string) (image.Image, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Database.ImageGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var img image.Image

	query := fmt.Sprintf(
		"SELECT id, object_id, type, origin, middle, small, organization_id, is_main FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		imageTable,
	)
	err := r.database.GetContext(ctx, &img, query, organizationID, imageID)

	return img, errors.Wrap(err, "image select query error")
}

// ImageCreate insert image into database.
func (r *Repository) ImageCreate(ctxr context.Context, userID string, input image.CreateImageInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Database.ImageCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var imageID string

	query := fmt.Sprintf(
		"INSERT INTO %s (object_id, type, origin, middle, small, organization_id, is_main, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		imageTable,
	)

	row := r.database.QueryRowContext(
		ctx,
		query,
		input.ObjectID,
		input.ObjectType,
		input.Origin,
		input.Middle,
		input.Small,
		input.OrganizationID,
		input.IsMain,
		userID,
	)

	err := row.Scan(&imageID)

	return imageID, errors.Wrap(err, "image create query error")
}

// ImageUpdate updates image by id in database.
func (r *Repository) ImageUpdate(ctxr context.Context, userID string, input image.UpdateImageInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Database.ImageUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	setValues := make([]string, 0, 8)
	args := make([]interface{}, 0, 8)
	argID := 1

	if input.ObjectID != nil {
		setValues = append(setValues, fmt.Sprintf("object_id=$%d", argID))
		args = append(args, *input.ObjectID)
		argID++
	}

	if input.ObjectType != nil {
		setValues = append(setValues, fmt.Sprintf("type=$%d", argID))
		args = append(args, *input.ObjectType)
		argID++
	}

	if input.Origin != nil {
		setValues = append(setValues, fmt.Sprintf("origin=$%d", argID))
		args = append(args, *input.Origin)
		argID++
	}

	if input.Middle != nil {
		setValues = append(setValues, fmt.Sprintf("middle=$%d", argID))
		args = append(args, *input.Middle)
		argID++
	}

	if input.Small != nil {
		setValues = append(setValues, fmt.Sprintf("small=$%d", argID))
		args = append(args, *input.Small)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

	if input.IsMain != nil {
		setValues = append(setValues, fmt.Sprintf("is_main=$%d", argID))
		args = append(args, *input.IsMain)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		imageTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.database.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "image update query error")
}

// ImageDelete deletes image by id from database.
func (r *Repository) ImageDelete(ctxr context.Context, userID string, organizationID string, imageID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "Database.ImageDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		imageTable,
	)

	_, err := r.database.ExecContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"), userID, imageID, organizationID)

	return errors.Wrap(err, "image delete query error")
}
