package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/message"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
)

// MessageGetAll selects all messages from database.
func (r *Repository) MessageGetAll(ctxr context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]message.Message, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.MessageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var messages []message.Message

	qry, args, err := r.messageGetAllQuery(meta, params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.SelectContext(ctx, &messages, qry, args...)

	return messages, errors.Wrap(err, "messages select query error")
}

// messageGetAllQuery creates sql query.
func (r *Repository) messageGetAllQuery(_ query.MetaData, params queryparameter.QueryParameter) (string, []interface{}, error) { //nolint:lll
	builder := r.genSQL.Select("id", "title", "body", "user_id", "is_public", "is_published").From(messageTable)

	if params.Search != "" {
		search := "%" + params.Search + "%"

		builder = builder.Where(squirrel.Or{
			squirrel.Like{"title": search},
			squirrel.Like{"body": search},
		})
	}

	if len(params.Sorts) > 0 {
		builder = builder.OrderBy(params.Sorts.Parsing(mappingSortMessage)...)
	} else {
		builder = builder.OrderBy("id ASC")
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

// MessageGetOne select message by id from database.
func (r *Repository) MessageGetOne(ctxr context.Context, _ query.MetaData, messageID string) (message.Message, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.MessageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var msg message.Message

	builder := r.genSQL.Select("id", "title", "body", "user_id", "is_public", "is_published").From(messageTable).
		Where(squirrel.Eq{"id": messageID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return msg, errors.Wrap(err, "unable to build a query string")
	}

	err = r.databaseSlave.GetContext(ctx, &msg, qry, args...)

	return msg, errors.Wrap(err, "message select query error")
}

// MessageCreate insert message into database.
func (r *Repository) MessageCreate(ctxr context.Context, _ query.MetaData, input message.CreateMessageInput) (string, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.MessageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var messageID string

	builder := r.genSQL.Insert(messageTable).
		Columns("title", "body", "user_id", "is_public").
		Values(input.Title, input.Body, input.UserID, input.IsPublic).
		Suffix("RETURNING \"id\"")

	qry, args, err := builder.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "unable to build a query string")
	}

	row := r.databaseMaster.QueryRowContext(ctx, qry, args...)
	err = row.Scan(&messageID)

	return messageID, errors.Wrap(err, "message create query error")
}

// MessageUpdate updates message by id in database.
func (r *Repository) MessageUpdate(ctxr context.Context, _ query.MetaData, input message.UpdateMessageInput) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.MessageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Update(tableTable)

	if input.Title != nil {
		builder = builder.Set("title", *input.Title)
	}

	if input.Body != nil {
		builder = builder.Set("body", *input.Body)
	}

	if input.UserID != nil {
		builder = builder.Set("user_id", *input.UserID)
	}

	if input.IsPublic != nil {
		builder = builder.Set("is_public", *input.IsPublic)
	}

	if input.IsPublished != nil {
		builder = builder.Set("is_published", *input.IsPublished)
	}

	builder = builder.Where(squirrel.Eq{"id": *input.ID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "message update query error")
}

// MessageDelete deletes message by id from database.
func (r *Repository) MessageDelete(ctxr context.Context, _ query.MetaData, messageID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Database.MessageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	builder := r.genSQL.Delete(messageTable).Where(squirrel.Eq{"id": messageID})

	qry, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "unable to build a query string")
	}

	_, err = r.databaseMaster.ExecContext(ctx, qry, args...)

	return errors.Wrap(err, "message delete query error")
}
