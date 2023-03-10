package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CommentGetAll selects all comments from database.
func (r *Repository) CommentGetAll(ctxr context.Context, userID string, organizationID string) ([]comment.Comment, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CommentGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var comments []comment.Comment

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, content, status_id, rating, user_created, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 ",
		commentTable,
	)

	err := r.database.SelectContext(ctx, &comments, query, organizationID)

	return comments, errors.Wrap(err, "comments select query error")
}

// CommentGetOne select comment by id from database.
func (r *Repository) CommentGetOne(ctxr context.Context, userID string, organizationID string, commentID string) (comment.Comment, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CommentGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var commnt comment.Comment

	query := fmt.Sprintf(
		"SELECT id, item_id, organization_id, content, status_id, rating, user_created, created_at FROM %s "+
			"WHERE is_deleted = false AND organization_id = $1 AND id = $2 ",
		commentTable,
	)
	err := r.database.GetContext(ctx, &commnt, query, organizationID, commentID)

	return commnt, errors.Wrap(err, "comment select query error")
}

// CommentCreate insert comment into database.
func (r *Repository) CommentCreate(ctxr context.Context, userID string, input comment.CreateCommentInput) (string, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CommentCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	var commentID string

	query := fmt.Sprintf(
		"INSERT INTO %s (item_id, organization_id, content, status_id, rating, user_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		commentTable,
	)

	row := r.database.QueryRowContext(
		ctx,
		query,
		input.ItemID,
		input.OrganizationID,
		input.Content,
		input.Status,
		input.Rating,
		userID,
	)

	err := row.Scan(&commentID)

	return commentID, errors.Wrap(err, "comment create query error")
}

// CommentUpdate updates comment by id in database.
func (r *Repository) CommentUpdate(ctxr context.Context, userID string, input comment.UpdateCommentInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CommentUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	setValues := make([]string, 0, 7)
	args := make([]interface{}, 0, 7)
	argID := 1

	if input.ItemID != nil {
		setValues = append(setValues, fmt.Sprintf("item_id=$%d", argID))
		args = append(args, *input.ItemID)
		argID++
	}

	if input.Content != nil {
		setValues = append(setValues, fmt.Sprintf("content=$%d", argID))
		args = append(args, *input.Content)
		argID++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status_id=$%d", argID))
		args = append(args, *input.Status)
		argID++
	}

	if input.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argID))
		args = append(args, *input.Rating)
		argID++
	}

	if input.OrganizationID != nil {
		setValues = append(setValues, fmt.Sprintf("organization_id=$%d", argID))
		args = append(args, *input.OrganizationID)
		argID++
	}

	setValues = append(setValues, fmt.Sprintf("user_updated=$%d", argID))
	args = append(args, userID)
	argID++

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now().Format("2006-01-02 15:04:05"))

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE is_deleted = false AND organization_id = '%s' AND id = '%s'",
		commentTable, setQuery, *input.OrganizationID, *input.ID)

	_, err := r.database.ExecContext(ctx, query, args...)

	return errors.Wrap(err, "comment update query error")
}

// CommentDelete deletes comment by id from database.
func (r *Repository) CommentDelete(ctxr context.Context, userID string, organizationID string, commentID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctxr, "RepositoryDatabase.CommentDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET is_deleted = true, deleted_at = $1, user_deleted = $2 "+
			"WHERE is_deleted = false AND id = $3 AND organization_id = $4",
		commentTable,
	)

	_, err := r.database.ExecContext(ctx, query, time.Now().Format("2006-01-02 15:04:05"), userID, commentID, organizationID)

	return errors.Wrap(err, "comment delete query error")
}
