package message

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/message"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// MessageGetAll returns all messages from the system.
func (s *UseCase) MessageGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]message.Message, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.MessageGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	messages, err := s.adapterStorage.MessageGetAll(ctx, meta, params)

	return messages, errors.Wrap(err, "messages select error")
}

// getAllWithCache returns messages from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]message.Message, error) { //nolint:lll
	messages, err := s.adapterCache.MessageGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get messages from cache", zap.String("error", err.Error()))
	}

	if len(messages) > 0 {
		return messages, nil
	}

	messages, err = s.adapterStorage.MessageGetAll(ctx, meta, params)

	if err != nil {
		return messages, errors.Wrap(err, "messages select failed")
	}

	if err = s.adapterCache.MessageSetAll(ctx, meta, params, messages); err != nil {
		logger.Logger.Error("unable to add messages into cache", zap.String("error", err.Error()))
	}

	return messages, nil
}

// MessageGetOne returns message by id from the system.
func (s *UseCase) MessageGetOne(ctx context.Context, meta query.MetaData, messageID string) (message.Message, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.MessageGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, messageID)
	}

	rle, err := s.adapterStorage.MessageGetOne(ctx, meta, messageID)

	return rle, errors.Wrap(err, "message select error")
}

// getOneWithCache returns message by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, messageID string) (message.Message, error) {
	rle, err := s.adapterCache.MessageGetOne(ctx, messageID)
	if err != nil {
		logger.Logger.Error("unable to get message from cache", zap.String("error", err.Error()))
	}

	if rle != (message.Message{}) {
		return rle, nil
	}

	rle, err = s.adapterStorage.MessageGetOne(ctx, meta, messageID)

	if err != nil {
		return rle, errors.Wrap(err, "message select failed")
	}

	if err = s.adapterCache.MessageCreate(ctx, rle); err != nil {
		logger.Logger.Error("unable to add message into cache", zap.String("error", err.Error()))
	}

	return rle, nil
}

// MessageCreate inserts message into system.
func (s *UseCase) MessageCreate(ctx context.Context, meta query.MetaData, input message.CreateMessageInput) (string, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.MessageCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	messageID, err := s.adapterStorage.MessageCreate(ctx, meta, input)
	if err != nil {
		return messageID, errors.Wrap(err, "message create error")
	}

	if s.isCacheOn {
		rle, err := s.adapterStorage.MessageGetOne(ctx, meta, messageID)
		if err != nil {
			return "", errors.Wrap(err, "message select from database failed")
		}

		err = s.adapterCache.MessageCreate(ctx, rle)
		if err != nil {
			return "", errors.Wrap(err, "message create in cache failed")
		}

		err = s.adapterCache.MessageInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "message invalidate users in cache failed")
		}
	}

	return messageID, nil
}

// MessageUpdate updates message by id in the system.
func (s *UseCase) MessageUpdate(ctx context.Context, meta query.MetaData, input message.UpdateMessageInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.MessageUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.MessageUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "message update in database failed")
	}

	if s.isCacheOn {
		rle, err := s.adapterStorage.MessageGetOne(ctx, meta, *input.ID)
		if err != nil {
			return errors.Wrap(err, "message select from database failed")
		}

		err = s.adapterCache.MessageUpdate(ctx, rle)
		if err != nil {
			return errors.Wrap(err, "message update in cache failed")
		}

		err = s.adapterCache.MessageInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "message invalidate users in cache failed")
		}
	}

	return nil
}

// MessageDelete deletes message by id from the system.
func (s *UseCase) MessageDelete(ctx context.Context, meta query.MetaData, messageID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.MessageDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.MessageDelete(ctx, meta, messageID)
	if err != nil {
		return errors.Wrap(err, "message delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.MessageDelete(ctx, messageID)
		if err != nil {
			return errors.Wrap(err, "message update in cache failed")
		}

		err = s.adapterCache.MessageInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate messages in cache failed")
		}
	}

	return nil
}
