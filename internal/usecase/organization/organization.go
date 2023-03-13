package organization

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/organization"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// OrganizationGetAll returns all organizations from the system.
func (s *UseCase) OrganizationGetAll(ctx context.Context, userID string) ([]organization.Organization, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID)
	}

	organizations, err := s.adapterStorage.OrganizationGetAll(ctx, userID)

	return organizations, errors.Wrap(err, "organization select error")
}

// getAllWithCache returns organizations from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string) ([]organization.Organization, error) {
	organizations, err := s.adapterCache.OrganizationGetAll(ctx)
	if err != nil {
		logger.Logger.Error("unable to get organizations from cache", zap.String("error", err.Error()))
	}

	if len(organizations) > 0 {
		return organizations, nil
	}

	organizations, err = s.adapterStorage.OrganizationGetAll(ctx, userID)

	if err != nil {
		return organizations, errors.Wrap(err, "organizations select failed")
	}

	if err = s.adapterCache.OrganizationSetAll(ctx, organizations); err != nil {
		logger.Logger.Error("unable to add organizations into cache", zap.String("error", err.Error()))
	}

	return organizations, nil
}

// OrganizationGetOne returns organization by id from the system.
func (s *UseCase) OrganizationGetOne(ctx context.Context, userID string, organizationID string) (organization.Organization, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID)
	}

	org, err := s.adapterStorage.OrganizationGetOne(ctx, userID, organizationID)

	return org, errors.Wrap(err, "organization select error")
}

// getOneWithCache returns organization by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) (organization.Organization, error) {
	org, err := s.adapterCache.OrganizationGetOne(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get organization from cache", zap.String("error", err.Error()))
	}

	if org != (organization.Organization{}) {
		return org, nil
	}

	org, err = s.adapterStorage.OrganizationGetOne(ctx, userID, organizationID)

	if err != nil {
		return org, errors.Wrap(err, "organization select failed")
	}

	if err = s.adapterCache.OrganizationCreate(ctx, org); err != nil {
		logger.Logger.Error("unable to add organization into cache", zap.String("error", err.Error()))
	}

	return org, nil
}

// OrganizationCreate inserts organization into system.
func (s *UseCase) OrganizationCreate(ctx context.Context, userID string, input organization.CreateOrganizationInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	organizationID, err := s.adapterStorage.OrganizationCreate(ctx, userID, input)
	if err != nil {
		return organizationID, errors.Wrap(err, "organization create error")
	}

	if s.isCacheOn {
		org, err := s.adapterStorage.OrganizationGetOne(ctx, userID, organizationID)
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
func (s *UseCase) OrganizationUpdate(ctx context.Context, userID string, input organization.UpdateOrganizationInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrganizationUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "organization update in database failed")
	}

	if s.isCacheOn {
		org, err := s.adapterStorage.OrganizationGetOne(ctx, userID, *input.ID)
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
func (s *UseCase) OrganizationDelete(ctx context.Context, userID string, organizationID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrganizationDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.OrganizationDelete(ctx, userID, organizationID)
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
