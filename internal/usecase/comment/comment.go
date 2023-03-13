package comment

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CommentGetAll returns all comments from the system.
func (s *UseCase) CommentGetAll(ctx context.Context, userID string, organizationID string) ([]comment.Comment, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID, organizationID)
	}

	comments, err := s.adapterStorage.CommentGetAll(ctx, userID, organizationID)

	return comments, errors.Wrap(err, "comments select error")
}

// getAllWithCache returns comments from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) ([]comment.Comment, error) {
	comments, err := s.adapterCache.CommentGetAll(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get comments from cache", zap.String("error", err.Error()))
	}

	if len(comments) > 0 {
		return comments, nil
	}

	comments, err = s.adapterStorage.CommentGetAll(ctx, userID, organizationID)
	if err != nil {
		return comments, errors.Wrap(err, "comments select failed")
	}

	if err = s.adapterCache.CommentSetAll(ctx, organizationID, comments); err != nil {
		logger.Logger.Error("unable to add comments into cache", zap.String("error", err.Error()))
	}

	return comments, nil
}

// CommentGetOne returns comment by id from the system.
func (s *UseCase) CommentGetOne(ctx context.Context, userID string, organizationID string, commentID string) (comment.Comment, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID, commentID)
	}

	cmnt, err := s.adapterStorage.CommentGetOne(ctx, userID, organizationID, commentID)

	return cmnt, errors.Wrap(err, "comment select error")
}

// getOneWithCache returns comment by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string, itemID string) (comment.Comment, error) {
	cmnt, err := s.adapterCache.CommentGetOne(ctx, itemID)
	if err != nil {
		logger.Logger.Error("unable to get comment from cache", zap.String("error", err.Error()))
	}

	if cmnt != (comment.Comment{}) {
		return cmnt, nil
	}

	cmnt, err = s.adapterStorage.CommentGetOne(ctx, userID, organizationID, itemID)
	if err != nil {
		return cmnt, errors.Wrap(err, "comment select failed")
	}

	if err = s.adapterCache.CommentCreate(ctx, cmnt); err != nil {
		logger.Logger.Error("unable to add comment into cache", zap.String("error", err.Error()))
	}

	return cmnt, nil
}

// CommentCreate inserts comment into system.
func (s *UseCase) CommentCreate(ctx context.Context, userID string, input comment.CreateCommentInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	commentID, err := s.adapterStorage.CommentCreate(ctx, userID, input)
	if err != nil {
		return commentID, errors.Wrap(err, "comment create error")
	}

	if s.isCacheOn {
		cmnt, err := s.adapterStorage.CommentGetOne(ctx, userID, input.OrganizationID, commentID)
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
func (s *UseCase) CommentUpdate(ctx context.Context, userID string, input comment.UpdateCommentInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CommentUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "comment update in database failed")
	}

	if s.isCacheOn {
		cmnt, err := s.adapterStorage.CommentGetOne(ctx, userID, *input.OrganizationID, *input.ID)
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
func (s *UseCase) CommentDelete(ctx context.Context, userID string, organizationID string, commentID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CommentDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.CommentDelete(ctx, userID, organizationID, commentID)
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
