package comment

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/comment"
	"github.com/pkg/errors"
)

// CommentGetAll returns all comments from the system.
func (s *UseCase) CommentGetAll(userID string, organizationID string) ([]comment.Comment, error) {
	comments, err := s.adapterStorage.CommentGetAll(userID, organizationID)

	return comments, errors.Wrap(err, "comments select error")
}

// CommentGetOne returns comment by id from the system.
func (s *UseCase) CommentGetOne(userID string, organizationID string, commentID string) (comment.Comment, error) {
	commentSingle, err := s.adapterStorage.CommentGetOne(userID, organizationID, commentID)

	return commentSingle, errors.Wrap(err, "comment select error")
}

// CommentCreate inserts comment into system.
func (s *UseCase) CommentCreate(userID string, comment comment.Comment) (string, error) {
	commentID, err := s.adapterStorage.CommentCreate(userID, comment)

	return commentID, errors.Wrap(err, "comment create error")
}

// CommentUpdate updates comment by id in the system.
func (s *UseCase) CommentUpdate(userID string, input comment.UpdateCommentInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CommentUpdate(userID, input)

	return errors.Wrap(err, "comment update error")
}

// CommentDelete deletes comment by id from the system.
func (s *UseCase) CommentDelete(userID string, organizationID string, commentID string) error {
	err := s.adapterStorage.CommentDelete(userID, organizationID, commentID)

	return errors.Wrap(err, "comment delete error")
}
