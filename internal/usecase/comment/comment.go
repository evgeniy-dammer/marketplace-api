package comment

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CommentGetAll returns all comments from the system.
func (s *UseCase) CommentGetAll(ctx context.Context, userID string, organizationID string) ([]comment.Comment, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CommentGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	comments, err := s.adapterStorage.CommentGetAll(ctx, userID, organizationID)

	return comments, errors.Wrap(err, "comments select error")
}

// CommentGetOne returns comment by id from the system.
func (s *UseCase) CommentGetOne(ctx context.Context, userID string, organizationID string, commentID string) (comment.Comment, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CommentGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	commentSingle, err := s.adapterStorage.CommentGetOne(ctx, userID, organizationID, commentID)

	return commentSingle, errors.Wrap(err, "comment select error")
}

// CommentCreate inserts comment into system.
func (s *UseCase) CommentCreate(ctx context.Context, userID string, input comment.CreateCommentInput) (string, error) {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CommentCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	commentID, err := s.adapterStorage.CommentCreate(ctx, userID, input)

	return commentID, errors.Wrap(err, "comment create error")
}

// CommentUpdate updates comment by id in the system.
func (s *UseCase) CommentUpdate(ctx context.Context, userID string, input comment.UpdateCommentInput) error {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CommentUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CommentUpdate(ctx, userID, input)

	return errors.Wrap(err, "comment update error")
}

// CommentDelete deletes comment by id from the system.
func (s *UseCase) CommentDelete(ctx context.Context, userID string, organizationID string, commentID string) error {
	if viper.GetBool("service.tracing") {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CommentDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.CommentDelete(ctx, userID, organizationID, commentID)

	return errors.Wrap(err, "comment delete error")
}
