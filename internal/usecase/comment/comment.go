package comment

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/comment"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CommentGetAll returns all comments from the system.
func (s *UseCase) CommentGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]comment.Comment, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	comments, err := s.adapterStorage.CommentGetAll(ctx, meta, params)

	return comments, errors.Wrap(err, "comments select error")
}

// getAllWithCache returns comments from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]comment.Comment, error) { //nolint:lll
	comments, err := s.adapterCache.CommentGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get comments from cache", zap.String("error", err.Error()))
	}

	if len(comments) > 0 {
		return comments, nil
	}

	comments, err = s.adapterStorage.CommentGetAll(ctx, meta, params)
	if err != nil {
		return comments, errors.Wrap(err, "comments select failed")
	}

	if err = s.adapterCache.CommentSetAll(ctx, meta, params, comments); err != nil {
		logger.Logger.Error("unable to add comments into cache", zap.String("error", err.Error()))
	}

	return comments, nil
}

// CommentGetOne returns comment by id from the system.
func (s *UseCase) CommentGetOne(ctx context.Context, meta query.MetaData, commentID string) (comment.Comment, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, commentID)
	}

	cmnt, err := s.adapterStorage.CommentGetOne(ctx, meta, commentID)

	return cmnt, errors.Wrap(err, "comment select error")
}

// getOneWithCache returns comment by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, itemID string) (comment.Comment, error) {
	cmnt, err := s.adapterCache.CommentGetOne(ctx, itemID)
	if err != nil {
		logger.Logger.Error("unable to get comment from cache", zap.String("error", err.Error()))
	}

	if cmnt != (comment.Comment{}) {
		return cmnt, nil
	}

	cmnt, err = s.adapterStorage.CommentGetOne(ctx, meta, itemID)
	if err != nil {
		return cmnt, errors.Wrap(err, "comment select failed")
	}

	if err = s.adapterCache.CommentCreate(ctx, cmnt); err != nil {
		logger.Logger.Error("unable to add comment into cache", zap.String("error", err.Error()))
	}

	return cmnt, nil
}

// CommentCreate inserts comment into system.
func (s *UseCase) CommentCreate(ctx context.Context, meta query.MetaData, input comment.CreateCommentInput) (string, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	commentID, err := s.adapterStorage.CommentCreate(ctx, meta, input)
	if err != nil {
		return commentID, errors.Wrap(err, "comment create error")
	}

	if s.isCacheOn {
		cmnt, err := s.adapterStorage.CommentGetOne(ctx, meta, commentID)
		if err != nil {
			return "", errors.Wrap(err, "comment select from database failed")
		}

		err = s.adapterCache.CommentCreate(ctx, cmnt)
		if err != nil {
			return "", errors.Wrap(err, "comment create in cache failed")
		}

		err = s.adapterCache.CommentInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "comment invalidate users in cache failed")
		}
	}

	return commentID, nil
}

// CommentUpdate updates comment by id in the system.
func (s *UseCase) CommentUpdate(ctx context.Context, meta query.MetaData, input comment.UpdateCommentInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CommentUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "comment update in database failed")
	}

	if s.isCacheOn {
		cmnt, err := s.adapterStorage.CommentGetOne(ctx, meta, *input.ID)
		if err != nil {
			return errors.Wrap(err, "comment select from database failed")
		}

		err = s.adapterCache.CommentUpdate(ctx, cmnt)
		if err != nil {
			return errors.Wrap(err, "comment update in cache failed")
		}

		err = s.adapterCache.CommentInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "comment invalidate users in cache failed")
		}
	}

	return nil
}

// CommentDelete deletes comment by id from the system.
func (s *UseCase) CommentDelete(ctx context.Context, meta query.MetaData, commentID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.CommentDelete(ctx, meta, commentID)
	if err != nil {
		return errors.Wrap(err, "comment delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.CommentDelete(ctx, commentID)
		if err != nil {
			return errors.Wrap(err, "comment update in cache failed")
		}

		err = s.adapterCache.CommentInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate comments in cache failed")
		}
	}

	return nil
}
