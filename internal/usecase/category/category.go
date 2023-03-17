package category

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/category"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CategoryGetAll returns all categories from the system.
func (s *UseCase) CategoryGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]category.Category, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CategoryGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	categories, err := s.adapterStorage.CategoryGetAll(ctx, meta, params)

	return categories, errors.Wrap(err, "categories select error")
}

// getAllWithCache returns categories from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]category.Category, error) {
	categories, err := s.adapterCache.CategoryGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get categories from cache", zap.String("error", err.Error()))
	}

	if len(categories) > 0 {
		return categories, nil
	}

	categories, err = s.adapterStorage.CategoryGetAll(ctx, meta, params)
	if err != nil {
		return categories, errors.Wrap(err, "categories select failed")
	}

	if err = s.adapterCache.CategorySetAll(ctx, meta, params, categories); err != nil {
		logger.Logger.Error("unable to add categories into cache", zap.String("error", err.Error()))
	}

	return categories, nil
}

// CategoryGetOne returns category by id from the system.
func (s *UseCase) CategoryGetOne(ctx context.Context, meta query.MetaData, categoryID string) (category.Category, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CategoryGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, categoryID)
	}

	ctgry, err := s.adapterStorage.CategoryGetOne(ctx, meta, categoryID)

	return ctgry, errors.Wrap(err, "category select error")
}

// getOneWithCache returns category by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, categoryID string) (category.Category, error) {
	ctgry, err := s.adapterCache.CategoryGetOne(ctx, categoryID)
	if err != nil {
		logger.Logger.Error("unable to get category from cache", zap.String("error", err.Error()))
	}

	if ctgry != (category.Category{}) {
		return ctgry, nil
	}

	ctgry, err = s.adapterStorage.CategoryGetOne(ctx, meta, categoryID)
	if err != nil {
		return ctgry, errors.Wrap(err, "category select failed")
	}

	if err = s.adapterCache.CategoryCreate(ctx, ctgry); err != nil {
		logger.Logger.Error("unable to add category into cache", zap.String("error", err.Error()))
	}

	return ctgry, nil
}

// CategoryCreate inserts category into system.
func (s *UseCase) CategoryCreate(ctx context.Context, meta query.MetaData, input category.CreateCategoryInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CategoryCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	categoryID, err := s.adapterStorage.CategoryCreate(ctx, meta, input)
	if err != nil {
		return categoryID, errors.Wrap(err, "category create error")
	}

	if s.isCacheOn {
		ctgry, err := s.adapterStorage.CategoryGetOne(ctx, meta, categoryID)
		if err != nil {
			return "", errors.Wrap(err, "category select from database failed")
		}

		err = s.adapterCache.CategoryCreate(ctx, ctgry)
		if err != nil {
			return "", errors.Wrap(err, "category create in cache failed")
		}

		err = s.adapterCache.CategoryInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "category invalidate users in cache failed")
		}
	}

	return categoryID, nil
}

// CategoryUpdate updates category by id in the system.
func (s *UseCase) CategoryUpdate(ctx context.Context, meta query.MetaData, input category.UpdateCategoryInput) error { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CategoryUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.CategoryUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "category update in database failed")
	}

	if s.isCacheOn {
		ctgry, err := s.adapterStorage.CategoryGetOne(ctx, meta, *input.ID)
		if err != nil {
			return errors.Wrap(err, "category select from database failed")
		}

		err = s.adapterCache.CategoryUpdate(ctx, ctgry)
		if err != nil {
			return errors.Wrap(err, "category update in cache failed")
		}

		err = s.adapterCache.CategoryInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "category invalidate users in cache failed")
		}
	}

	return nil
}

// CategoryDelete deletes category by id from the system.
func (s *UseCase) CategoryDelete(ctx context.Context, meta query.MetaData, categoryID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.CategoryDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.CategoryDelete(ctx, meta, categoryID)
	if err != nil {
		return errors.Wrap(err, "category delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.CategoryDelete(ctx, categoryID)
		if err != nil {
			return errors.Wrap(err, "category update in cache failed")
		}

		err = s.adapterCache.CategoryInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate categories in cache failed")
		}
	}

	return nil
}
