package table

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/table"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// TableGetAll returns all tables from the system.
func (s *UseCase) TableGetAll(ctx context.Context, userID string, organizationID string) ([]table.Table, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.TableGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID, organizationID)
	}

	tables, err := s.adapterStorage.TableGetAll(ctx, userID, organizationID)

	return tables, errors.Wrap(err, "tables select error")
}

// getAllWithCache returns tables from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) ([]table.Table, error) {
	tables, err := s.adapterCache.TableGetAll(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get tables from cache", zap.String("error", err.Error()))
	}

	if len(tables) > 0 {
		return tables, nil
	}

	tables, err = s.adapterStorage.TableGetAll(ctx, userID, organizationID)
	if err != nil {
		return tables, errors.Wrap(err, "tables select failed")
	}

	if err = s.adapterCache.TableSetAll(ctx, organizationID, tables); err != nil {
		logger.Logger.Error("unable to add tables into cache", zap.String("error", err.Error()))
	}

	return tables, nil
}

// TableGetOne returns table by id from the system.
func (s *UseCase) TableGetOne(ctx context.Context, userID string, organizationID string, tableID string) (table.Table, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.TableGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID, tableID)
	}

	tbl, err := s.adapterStorage.TableGetOne(ctx, userID, organizationID, tableID)

	return tbl, errors.Wrap(err, "table select error")
}

// getOneWithCache returns table by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string, tableID string) (table.Table, error) {
	tble, err := s.adapterCache.TableGetOne(ctx, tableID)
	if err != nil {
		logger.Logger.Error("unable to get table from cache", zap.String("error", err.Error()))
	}

	if tble != (table.Table{}) {
		return tble, nil
	}

	tble, err = s.adapterStorage.TableGetOne(ctx, userID, organizationID, tableID)
	if err != nil {
		return tble, errors.Wrap(err, "table select failed")
	}

	if err = s.adapterCache.TableCreate(ctx, tble); err != nil {
		logger.Logger.Error("unable to add table into cache", zap.String("error", err.Error()))
	}

	return tble, nil
}

// TableCreate inserts table into system.
func (s *UseCase) TableCreate(ctx context.Context, userID string, input table.CreateTableInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.TableCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	tableID, err := s.adapterStorage.TableCreate(ctx, userID, input)
	if err != nil {
		return tableID, errors.Wrap(err, "table create error")
	}

	if s.isCacheOn {
		tble, err := s.adapterStorage.TableGetOne(ctx, userID, input.OrganizationID, tableID)
		if err != nil {
			return "", errors.Wrap(err, "table select from database failed")
		}

		err = s.adapterCache.TableCreate(ctx, tble)
		if err != nil {
			return "", errors.Wrap(err, "table create in cache failed")
		}

		err = s.adapterCache.TableInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "table invalidate users in cache failed")
		}
	}

	return tableID, nil
}

// TableUpdate updates table by id in the system.
func (s *UseCase) TableUpdate(ctx context.Context, userID string, input table.UpdateTableInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.TableUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.TableUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "table update in database failed")
	}

	if s.isCacheOn {
		tble, err := s.adapterStorage.TableGetOne(ctx, userID, *input.OrganizationID, *input.ID)
		if err != nil {
			return errors.Wrap(err, "table select from database failed")
		}

		err = s.adapterCache.TableUpdate(ctx, tble)
		if err != nil {
			return errors.Wrap(err, "table update in cache failed")
		}

		err = s.adapterCache.TableInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "table invalidate users in cache failed")
		}
	}

	return nil
}

// TableDelete deletes table by id from the system.
func (s *UseCase) TableDelete(ctx context.Context, userID string, organizationID string, tableID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.TableDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.TableDelete(ctx, userID, organizationID, tableID)
	if err != nil {
		return errors.Wrap(err, "table delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.TableDelete(ctx, tableID)
		if err != nil {
			return errors.Wrap(err, "table update in cache failed")
		}

		err = s.adapterCache.TableInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate tables in cache failed")
		}
	}

	return nil
}
