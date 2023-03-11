package specification

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/specification"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// SpecificationGetAll returns all specifications from the system.
func (s *UseCase) SpecificationGetAll(ctx context.Context, userID string, organizationID string) ([]specification.Specification, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.SpecificationGetAll")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	specifications, err := s.adapterStorage.SpecificationGetAll(ctx, userID, organizationID)

	return specifications, errors.Wrap(err, "specifications select error")
}

// SpecificationGetOne returns specification by id from the system.
func (s *UseCase) SpecificationGetOne(ctx context.Context, userID string, organizationID string, specificationID string) (specification.Specification, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.SpecificationGetOne")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	specification, err := s.adapterStorage.SpecificationGetOne(ctx, userID, organizationID, specificationID)

	return specification, errors.Wrap(err, "specification select error")
}

// SpecificationCreate inserts specification into system.
func (s *UseCase) SpecificationCreate(ctx context.Context, userID string, input specification.CreateSpecificationInput) (string, error) {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.SpecificationCreate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	specificationID, err := s.adapterStorage.SpecificationCreate(ctx, userID, input)

	return specificationID, errors.Wrap(err, "specification create error")
}

// SpecificationUpdate updates specification by id in the system.
func (s *UseCase) SpecificationUpdate(ctx context.Context, userID string, input specification.UpdateSpecificationInput) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.SpecificationUpdate")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.SpecificationUpdate(ctx, userID, input)

	return errors.Wrap(err, "specification update error")
}

// SpecificationDelete deletes specification by id from the system.
func (s *UseCase) SpecificationDelete(ctx context.Context, userID string, organizationID string, specificationID string) error {
	if s.isTracingOn {
		span, ctxt := opentracing.StartSpanFromContext(ctx, "Usecase.SpecificationDelete")
		defer span.Finish()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.SpecificationDelete(ctx, userID, organizationID, specificationID)

	return errors.Wrap(err, "specification delete error")
}
