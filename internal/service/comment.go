package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// CommentService is an organization service.
type CommentService struct {
	repo repository.Comment
}

// NewCommentService is a constructor for CommentService.
func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

// GetAll returns all comments from the system.
func (s *CommentService) GetAll(userID string, organizationID string) ([]model.Comment, error) {
	comments, err := s.repo.GetAll(userID, organizationID)

	return comments, errors.Wrap(err, "comments select error")
}

// GetOne returns comment by id from the system.
func (s *CommentService) GetOne(userID string, organizationID string, commentID string) (model.Comment, error) {
	comment, err := s.repo.GetOne(userID, organizationID, commentID)

	return comment, errors.Wrap(err, "comment select error")
}

// Create inserts comment into system.
func (s *CommentService) Create(userID string, comment model.Comment) (string, error) {
	commentID, err := s.repo.Create(userID, comment)

	return commentID, errors.Wrap(err, "comment create error")
}

// Update updates comment by id in the system.
func (s *CommentService) Update(userID string, input model.UpdateCommentInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "comment update error")
}

// Delete deletes comment by id from the system.
func (s *CommentService) Delete(userID string, organizationID string, commentID string) error {
	err := s.repo.Delete(userID, organizationID, commentID)

	return errors.Wrap(err, "comment delete error")
}
