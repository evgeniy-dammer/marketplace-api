package organization

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/organization"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/query"
	"github.com/evgeniy-dammer/marketplace-api/pkg/queryparameter"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// OrganizationGetAll returns all organizations from the system.
func (s *UseCase) OrganizationGetAll(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getAllWithCache(ctx, meta, params)
	}

	organizations, err := s.adapterStorage.OrganizationGetAll(ctx, meta, params)

	return organizations, errors.Wrap(err, "organization select error")
}

// getAllWithCache returns organizations from cache if exists.
func (s *UseCase) getAllWithCache(ctx context.Context, meta query.MetaData, params queryparameter.QueryParameter) ([]organization.Organization, error) { //nolint:lll
	organizations, err := s.adapterCache.OrganizationGetAll(ctx, meta, params)
	if err != nil {
		logger.Logger.Error("unable to get organizations from cache", zap.String("error", err.Error()))
	}

	if len(organizations) > 0 {
		return organizations, nil
	}

	organizations, err = s.adapterStorage.OrganizationGetAll(ctx, meta, params)

	if err != nil {
		return organizations, errors.Wrap(err, "organizations select failed")
	}

	if err = s.adapterCache.OrganizationSetAll(ctx, meta, params, organizations); err != nil {
		logger.Logger.Error("unable to add organizations into cache", zap.String("error", err.Error()))
	}

	return organizations, nil
}

// OrganizationGetOne returns organization by id from the system.
func (s *UseCase) OrganizationGetOne(ctx context.Context, meta query.MetaData, organizationID string) (organization.Organization, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return s.getOneWithCache(ctx, meta, organizationID)
	}

	org, err := s.adapterStorage.OrganizationGetOne(ctx, meta, organizationID)

	return org, errors.Wrap(err, "organization select error")
}

// getOneWithCache returns organization by id from cache if exists.
func (s *UseCase) getOneWithCache(ctx context.Context, meta query.MetaData, organizationID string) (organization.Organization, error) { //nolint:lll
	org, err := s.adapterCache.OrganizationGetOne(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get organization from cache", zap.String("error", err.Error()))
	}

	if org != (organization.Organization{}) {
		return org, nil
	}

	org, err = s.adapterStorage.OrganizationGetOne(ctx, meta, organizationID)

	if err != nil {
		return org, errors.Wrap(err, "organization select failed")
	}

	if err = s.adapterCache.OrganizationCreate(ctx, org); err != nil {
		logger.Logger.Error("unable to add organization into cache", zap.String("error", err.Error()))
	}

	return org, nil
}

// OrganizationCreate inserts organization into system.
func (s *UseCase) OrganizationCreate(ctx context.Context, meta query.MetaData, input organization.CreateOrganizationInput) (string, error) { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	organizationID, err := s.adapterStorage.OrganizationCreate(ctx, meta, input)
	if err != nil {
		return organizationID, errors.Wrap(err, "organization create error")
	}

	if s.isCacheOn {
		org, err := s.adapterStorage.OrganizationGetOne(ctx, meta, organizationID)
		if err != nil {
			return "", errors.Wrap(err, "organization select from database failed")
		}

		err = s.adapterCache.OrganizationCreate(ctx, org)
		if err != nil {
			return "", errors.Wrap(err, "organization create in cache failed")
		}

		err = s.adapterCache.OrganizationInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "organization invalidate users in cache failed")
		}
	}

	return organizationID, nil
}

// OrganizationUpdate updates organization by id in the system.
func (s *UseCase) OrganizationUpdate(ctx context.Context, meta query.MetaData, input organization.UpdateOrganizationInput) error { //nolint:lll
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrganizationUpdate(ctx, meta, input)
	if err != nil {
		return errors.Wrap(err, "organization update in database failed")
	}

	if s.isCacheOn {
		org, err := s.adapterStorage.OrganizationGetOne(ctx, meta, *input.ID)
		if err != nil {
			return errors.Wrap(err, "organization select from database failed")
		}

		err = s.adapterCache.OrganizationUpdate(ctx, org)
		if err != nil {
			return errors.Wrap(err, "organization update in cache failed")
		}

		err = s.adapterCache.OrganizationInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "organization invalidate users in cache failed")
		}
	}

	return nil
}

// OrganizationDelete deletes organization by id from the system.
func (s *UseCase) OrganizationDelete(ctx context.Context, meta query.MetaData, organizationID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.OrganizationDelete(ctx, meta, organizationID)
	if err != nil {
		return errors.Wrap(err, "organization delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.OrganizationDelete(ctx, organizationID)
		if err != nil {
			return errors.Wrap(err, "organization update in cache failed")
		}

		err = s.adapterCache.OrganizationInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate organizations in cache failed")
		}
	}

	return nil
}
