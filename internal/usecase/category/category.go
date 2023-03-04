package category

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/pkg/errors"
)

// CategoryGetAll returns all categories from the system.
func (s *UseCase) CategoryGetAll(userID string, organizationID string) ([]category.Category, error) {
	categories, err := s.adapterStorage.CategoryGetAll(userID, organizationID)

	return categories, errors.Wrap(err, "categories select error")
}

// CategoryGetOne returns category by id from the system.
func (s *UseCase) CategoryGetOne(userID string, organizationID string, categoryID string) (category.Category, error) {
	category, err := s.adapterStorage.CategoryGetOne(userID, organizationID, categoryID)

	return category, errors.Wrap(err, "category select error")
}

// CategoryCreate inserts category into system.
func (s *UseCase) CategoryCreate(userID string, category category.Category) (string, error) {
	categoryID, err := s.adapterStorage.CategoryCreate(userID, category)

	return categoryID, errors.Wrap(err, "category create error")
}

// CategoryUpdate updates category by id in the system.
func (s *UseCase) CategoryUpdate(userID string, input category.UpdateCategoryInput) error { //nolint:lll
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CategoryUpdate(userID, input)

	return errors.Wrap(err, "category update error")
}

// CategoryDelete deletes category by id from the system.
func (s *UseCase) CategoryDelete(userID string, organizationID string, categoryID string) error {
	err := s.adapterStorage.CategoryDelete(userID, organizationID, categoryID)

	return errors.Wrap(err, "category delete error")
}
