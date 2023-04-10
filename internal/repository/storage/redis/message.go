package redis

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/message"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

// MessageGetAll gets messages from cache.
func (r *Repository) MessageGetAll(ctxr context.Context, _ query.MetaData, _ queryparameter.QueryParameter) ([]message.Message, error) { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	messages := &message.ListMessage{}

	bytes, err := r.client.Get(ctx, messagesKey).Bytes()
	if err != nil {
		return *messages, errors.Wrap(err, "unable to get messages from cache")
	}

	if err = easyjson.Unmarshal(bytes, messages); err != nil {
		return *messages, errors.Wrap(err, "unable to unmarshal")
	}

	return *messages, nil
}

// MessageSetAll sets messages into cache.
func (r *Repository) MessageSetAll(ctxr context.Context, _ query.MetaData, _ queryparameter.QueryParameter, messages []message.Message) error { //nolint:lll
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageSetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	messageSlice := message.ListMessage(messages)

	bytes, err := easyjson.Marshal(messageSlice)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, messagesKey, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// MessageGetOne gets message by id from cache.
func (r *Repository) MessageGetOne(ctxr context.Context, messageID string) (message.Message, error) {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	var msg message.Message

	bytes, err := r.client.Get(ctx, messageKey+messageID).Bytes()
	if err != nil {
		return msg, errors.Wrap(err, "unable to get message from cache")
	}

	if err = easyjson.Unmarshal(bytes, &msg); err != nil {
		return msg, errors.Wrap(err, "unable to unmarshal")
	}

	return msg, nil
}

// MessageCreate sets message into cache.
func (r *Repository) MessageCreate(ctxr context.Context, msg message.Message) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	err = r.client.Set(ctx, messageKey+msg.ID, bytes, r.options.TTL).Err()

	return errors.Wrap(err, "setting key")
}

// MessageUpdate updates message by id in cache.
func (r *Repository) MessageUpdate(ctxr context.Context, msg message.Message) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	bytes, err := easyjson.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "unable to marshal json")
	}

	r.client.Set(ctx, messageKey+msg.ID, bytes, r.options.TTL)

	return nil
}

// MessageDelete deletes message by id from cache.
func (r *Repository) MessageDelete(ctxr context.Context, messageID string) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := r.client.Del(ctx, messageKey+messageID).Err()

	return errors.Wrap(err, "deleting key")
}

// MessageInvalidate invalidate messages cache.
func (r *Repository) MessageInvalidate(ctxr context.Context) error {
	ctx := ctxr.CopyWithTimeout(r.options.Timeout)
	defer ctx.Cancel()

	if r.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctxr, "Cache.MessageInvalidate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	iter := r.client.Scan(ctx, 0, messagesKey+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return errors.Wrap(err, "deleting key")
		}
	}

	return errors.Wrap(iter.Err(), "invalidate")
}
