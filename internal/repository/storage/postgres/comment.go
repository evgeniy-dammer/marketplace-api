package postgres

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/comment"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// CommentGetAll selects all comments from database.
func (r *Repository) CommentGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]comment.Comment, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CommentGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var comments []comment.Comment

	qry, args, err := r.commentGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.SelectContext(ctx, &comments, qry, args...)

	return comments, errors.Wrap(err, "comments select query error")
}

// commentGetAllQuery creates sql query.
func (r *Repository) commentGetAllQuery(meta query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select(
		"id", "item_id", "organization_id", "content", "status_id", "rating", "user_created", "created_at").
		From(commentTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.And{
			squirrel.Eq{"is_deleted": false},
			squirrel.Like{"content": search},
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
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortComment)...)
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

// CommentGetOne select comment by id from database.
func (r *Repository) CommentGetOne(ctxr context.Context, meta query.MetaData, commentID string) (comment.Comment, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CommentGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var commnt comment.Comment

	builder := r.genSQL.Select(
		"id", "item_id", "organization_id", "content", "status_id", "rating", "user_created", "created_at").
		From(commentTable).
		Where(squirrel.Eq{"is_deleted": false, "id": commentID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return commnt, errors.Wrap(err, "unable to build a query string")
	}

	err = r.database.GetContext(ctx, &commnt, qry, args...)

	return commnt, errors.Wrap(err, "comment select query error")
}

// CommentCreate insert comment into database.
func (r *Repository) CommentCreate(ctxr context.Context, meta query.MetaData, input comment.CreateCommentInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CommentCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var commentID string

	builder := r.genSQL.Insert(commentTable).
		Columns("item_id", "organization_id", "content", "status_id", "rating", "user_created").
		Values(input.ItemID, input.OrganizationID, input.Content, input.Status, input.Rating, meta.UserID).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.database.QueryRowContext(ctx, qry, args...)

	err = row.Scan(&commentID)

	return commentID, errors.Wrap(err, "comment create query error")
}

// CommentUpdate updates comment by id in database.
func (r *Repository) CommentUpdate(ctxr context.Context, meta query.MetaData, input comment.UpdateCommentInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CommentUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(commentTable)

	if input.ItemID != nil {
		builder = builder.Set("item_id", *input.ItemID)
	}

	if input.Content != nil {
		builder = builder.Set("content", *input.Content)
	}

	if input.Status != nil {
		builder = builder.Set("status_id", *input.Status)
	}

	if input.Rating != nil {
		builder = builder.Set("rating", *input.Rating)
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

	return errors.Wrap(err, "comment update query error")
}

// CommentDelete deletes comment by id from database.
func (r *Repository) CommentDelete(ctxr context.Context, meta query.MetaData, commentID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.CommentDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(commentTable).
		Set("is_deleted", true).
		Set("user_deleted", meta.UserID).
		Set("deleted_at", time.Now().UTC()).
		Where(squirrel.Eq{"is_deleted": false, "id": commentID})

	if meta.OrganizationID != "" {
		builder = builder.Where(squirrel.Eq{"organization_id": meta.OrganizationID})
	}

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.database.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "comment delete query error")
}
