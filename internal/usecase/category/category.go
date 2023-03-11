package category

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/category"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// CategoryGetAll returns all categories from the system.
func (s *UseCase) CategoryGetAll(ctx context.Context, userID string, organizationID string) ([]category.Category, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CategoryGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	categories, err := s.adapterStorage.CategoryGetAll(ctx, userID, organizationID)

	return categories, errors.Wrap(err, "categories select error")
}

// CategoryGetOne returns category by id from the system.
func (s *UseCase) CategoryGetOne(ctx context.Context, userID string, organizationID string, categoryID string) (category.Category, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CategoryGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	ctgry, err := s.adapterStorage.CategoryGetOne(ctx, userID, organizationID, categoryID)

	return ctgry, errors.Wrap(err, "category select error")
}

// CategoryCreate inserts category into system.
func (s *UseCase) CategoryCreate(ctx context.Context, userID string, input category.CreateCategoryInput) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CategoryCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	categoryID, err := s.adapterStorage.CategoryCreate(ctx, userID, input)

	return categoryID, errors.Wrap(err, "category create error")
}

// CategoryUpdate updates category by id in the system.
func (s *UseCase) CategoryUpdate(ctx context.Context, userID string, input category.UpdateCategoryInput) error { //nolint:lll
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CategoryUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CategoryUpdate(ctx, userID, input)

	return errors.Wrap(err, "category update error")
}

// CategoryDelete deletes category by id from the system.
func (s *UseCase) CategoryDelete(ctx context.Context, userID string, organizationID string, categoryID string) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.CategoryDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.CategoryDelete(ctx, userID, organizationID, categoryID)

	return errors.Wrap(err, "category delete error")
}
