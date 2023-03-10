package category

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/pkg/errors"
)

// CategoryGetAll returns all categories from the system.
func (s *UseCase) CategoryGetAll(ctx context.Context, userID string, organizationID string) ([]category.Category, error) {
	categories, err := s.adapterStorage.CategoryGetAll(ctx, userID, organizationID)

	return categories, errors.Wrap(err, "categories select error")
}

// CategoryGetOne returns category by id from the system.
func (s *UseCase) CategoryGetOne(ctx context.Context, userID string, organizationID string, categoryID string) (category.Category, error) {
	ctgry, err := s.adapterStorage.CategoryGetOne(ctx, userID, organizationID, categoryID)

	return ctgry, errors.Wrap(err, "category select error")
}

// CategoryCreate inserts category into system.
func (s *UseCase) CategoryCreate(ctx context.Context, userID string, input category.CreateCategoryInput) (string, error) {
	categoryID, err := s.adapterStorage.CategoryCreate(ctx, userID, input)

	return categoryID, errors.Wrap(err, "category create error")
}

// CategoryUpdate updates category by id in the system.
func (s *UseCase) CategoryUpdate(ctx context.Context, userID string, input category.UpdateCategoryInput) error { //nolint:lll
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CategoryUpdate(ctx, userID, input)

	return errors.Wrap(err, "category update error")
}

// CategoryDelete deletes category by id from the system.
func (s *UseCase) CategoryDelete(ctx context.Context, userID string, organizationID string, categoryID string) error {
	err := s.adapterStorage.CategoryDelete(ctx, userID, organizationID, categoryID)

	return errors.Wrap(err, "category delete error")
}
