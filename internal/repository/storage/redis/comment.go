package redis

import (
	"encoding/json"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
)

// CommentGetAll gets comments from cache.
func (r *Repository) CommentGetAll(ctxr context.Context, organizationID string) ([]comment.Comment, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var comments []comment.Comment

	bytes, err := r.client.Get(ctx, commentsKey+"o."+organizationID).Bytes()
	if err != nil {
		return comments, errors.Wrap(err, "unable to get comments from cache")
	}

	if err = json.Unmarshal(bytes, &comments); err != nil {
		return comments, errors.Wrap(err, "unable to unmarshal")
	}

	return comments, nil
}

// CommentSetAll sets comments into cache.
func (r *Repository) CommentSetAll(ctxr context.Context, organizationID string, comments []comment.Comment) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(comments)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, commentsKey+"o."+organizationID, bytes, r.options.Ttl).Err()

	return err
}

// CommentGetOne gets comment by id from cache.
func (r *Repository) CommentGetOne(ctxr context.Context, commentID string) (comment.Comment, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var usr comment.Comment

	bytes, err := r.client.Get(ctx, commentKey+commentID).Bytes()
	if err != nil {
		return usr, errors.Wrap(err, "unable to get comment from cache")
	}

	if err = json.Unmarshal(bytes, &usr); err != nil {
		return usr, errors.Wrap(err, "unable to unmarshal")
	}

	return usr, nil
}

// CommentCreate sets comment into cache.
func (r *Repository) CommentCreate(ctxr context.Context, usr comment.Comment) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, commentKey+usr.ID, bytes, r.options.Ttl).Err()

	return err
}

// CommentUpdate updates comment by id in cache.
func (r *Repository) CommentUpdate(ctxr context.Context, usr comment.Comment) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := json.Marshal(usr)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, commentKey+usr.ID, bytes, r.options.Ttl)

	return nil
}

// CommentDelete deletes comment by id from cache.
func (r *Repository) CommentDelete(ctxr context.Context, commentID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, commentKey+commentID).Err()

	return err
}

// CommentInvalidate invalidate comments cache.
func (r *Repository) CommentInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.CommentInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, commentsKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}

	return iter.Err()
}
