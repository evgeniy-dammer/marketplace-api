package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// CategoryService is a category service.
type CategoryService struct {
	repo repository.Category
}

// NewCategoryService is a constructor for CategoryService.
func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll returns all categories from the system.
func (s *CategoryService) GetAll(userID string, organizationID string) ([]model.Category, error) {
	categories, err := s.repo.GetAll(userID, organizationID)

	return categories, errors.Wrap(err, "categories select error")
}

// GetOne returns category by id from the system.
func (s *CategoryService) GetOne(userID string, organizationID string, categoryID string) (model.Category, error) {
	category, err := s.repo.GetOne(userID, organizationID, categoryID)

	return category, errors.Wrap(err, "category select error")
}

// Create inserts category into system.
func (s *CategoryService) Create(userID string, category model.Category) (string, error) {
	categoryID, err := s.repo.Create(userID, category)

	return categoryID, errors.Wrap(err, "category create error")
}

// Update updates category by id in the system.
func (s *CategoryService) Update(userID string, input model.UpdateCategoryInput) error { //nolint:lll
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "category update error")
}

// Delete deletes category by id from the system.
func (s *CategoryService) Delete(userID string, organizationID string, categoryID string) error {
	err := s.repo.Delete(userID, organizationID, categoryID)

	return errors.Wrap(err, "category delete error")
}
