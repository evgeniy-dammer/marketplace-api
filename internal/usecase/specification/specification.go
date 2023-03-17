package specification

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/specification"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SpecificationGetAll returns all specifications from the system.
func (s *UseCase) SpecificationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.SpecificationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	specifications, err := s.adapterStorage.SpecificationGetAll(ctx, meta, params)

	return specifications, errors.Wrap(err, "specifications select error")
}

// getAllWithCache returns specifications from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]specification.Specification, error) {
	specifications, err := s.adapterCache.SpecificationGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get specifications from cache", zap.String("error", err.Error()))
	}

	if len(specifications) > 0 {
		return specifications, nil
	}

	specifications, err = s.adapterStorage.SpecificationGetAll(ctx, meta, params)
	if err != nil {
		return specifications, errors.Wrap(err, "specifications select failed")
	}

	if err = s.adapterCache.SpecificationSetAll(ctx, meta, params, specifications); err != nil {
		logger.Logger.Error("unable to add specifications into cache", zap.String("error", err.Error()))
	}

	return specifications, nil
}

// SpecificationGetOne returns specification by id from the system.
func (s *UseCase) SpecificationGetOne(ctx context.Context, meta query.MetaData, specificationID string) (specification.Specification, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.SpecificationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, specificationID)
	}

	spec, err := s.adapterStorage.SpecificationGetOne(ctx, meta, specificationID)

	return spec, errors.Wrap(err, "specification select error")
}

// getOneWithCache returns specification by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, specificationID string) (specification.Specification, error) {
	spec, err := s.adapterCache.SpecificationGetOne(ctx, specificationID)
	if err != nil {
		logger.Logger.Error("unable to get specification from cache", zap.String("error", err.Error()))
	}

	if spec != (specification.Specification{}) {
		return spec, nil
	}

	spec, err = s.adapterStorage.SpecificationGetOne(ctx, meta, specificationID)
	if err != nil {
		return spec, errors.Wrap(err, "specification select failed")
	}

	if err = s.adapterCache.SpecificationCreate(ctx, spec); err != nil {
		logger.Logger.Error("unable to add specification into cache", zap.String("error", err.Error()))
	}

	return spec, nil
}

// SpecificationCreate inserts specification into system.
func (s *UseCase) SpecificationCreate(ctx context.Context, meta query.MetaData, input specification.CreateSpecificationInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.SpecificationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	specificationID, err := s.adapterStorage.SpecificationCreate(ctx, meta, input)
	if err != nil {
		return specificationID, errors.Wrap(err, "specification create error")
	}

	if s.isCacheOn {
		spec, err := s.adapterStorage.SpecificationGetOne(ctx, meta, specificationID)
		if err != nil {
			return "", errors.Wrap(err, "specification select from database failed")
		}

		err = s.adapterCache.SpecificationCreate(ctx, spec)
		if err != nil {
			return "", errors.Wrap(err, "specification create in cache failed")
		}

		err = s.adapterCache.SpecificationInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "specification invalidate users in cache failed")
		}
	}

	return specificationID, nil
}

// SpecificationUpdate updates specification by id in the system.
func (s *UseCase) SpecificationUpdate(ctx context.Context, meta query.MetaData, input specification.UpdateSpecificationInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.SpecificationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.SpecificationUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "specification update in database failed")
	}

	if s.isCacheOn {
		spec, err := s.adapterStorage.SpecificationGetOne(ctx, meta, *input.ID)
		if err != nil {
			return errors.Wrap(err, "specification select from database failed")
		}

		err = s.adapterCache.SpecificationUpdate(ctx, spec)
		if err != nil {
			return errors.Wrap(err, "specification update in cache failed")
		}

		err = s.adapterCache.SpecificationInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "specification invalidate users in cache failed")
		}
	}

	return nil
}

// SpecificationDelete deletes specification by id from the system.
func (s *UseCase) SpecificationDelete(ctx context.Context, meta query.MetaData, specificationID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.SpecificationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.SpecificationDelete(ctx, meta, specificationID)
	if err != nil {
		return errors.Wrap(err, "specification delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.SpecificationDelete(ctx, specificationID)
		if err != nil {
			return errors.Wrap(err, "specification update in cache failed")
		}

		err = s.adapterCache.SpecificationInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate specifications in cache failed")
		}
	}

	return nil
}
