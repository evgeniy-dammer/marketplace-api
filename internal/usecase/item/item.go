package item

import (
	"reflect"

	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/evgeniy-dammer/emenu-api/pkg/logger"
	"github.com/evgeniy-dammer/emenu-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ItemGetAll returns all items from the system.
func (s *UseCase) ItemGetAll(ctx context.Context, userID string, organizationID string) ([]item.Item, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ItemGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID, organizationID)
	}

	items, err := s.adapterStorage.ItemGetAll(ctx, userID, organizationID)

	return items, errors.Wrap(err, "items select error")
}

// getAllWithCache returns items from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) ([]item.Item, error) {
	items, err := s.adapterCache.ItemGetAll(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get items from cache", zap.String("error", err.Error()))
	}

	if len(items) > 0 {
		return items, nil
	}

	items, err = s.adapterStorage.ItemGetAll(ctx, userID, organizationID)
	if err != nil {
		return items, errors.Wrap(err, "items select failed")
	}

	if err = s.adapterCache.ItemSetAll(ctx, organizationID, items); err != nil {
		logger.Logger.Error("unable to add items into cache", zap.String("error", err.Error()))
	}

	return items, nil
}

// ItemGetOne returns item by id from the system.
func (s *UseCase) ItemGetOne(ctx context.Context, userID string, organizationID string, itemID string) (item.Item, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ItemGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID, itemID)
	}

	itemSingle, err := s.adapterStorage.ItemGetOne(ctx, userID, organizationID, itemID)

	return itemSingle, errors.Wrap(err, "item select error")
}

// getOneWithCache returns item by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string, itemID string) (item.Item, error) {
	itemSingle, err := s.adapterCache.ItemGetOne(ctx, itemID)
	if err != nil {
		logger.Logger.Error("unable to get item from cache", zap.String("error", err.Error()))
	}

	if !reflect.ValueOf(itemSingle).IsZero() {
		return itemSingle, nil
	}

	itemSingle, err = s.adapterStorage.ItemGetOne(ctx, userID, organizationID, itemID)
	if err != nil {
		return itemSingle, errors.Wrap(err, "item select failed")
	}

	if err = s.adapterCache.ItemCreate(ctx, itemSingle); err != nil {
		logger.Logger.Error("unable to add item into cache", zap.String("error", err.Error()))
	}

	return itemSingle, nil
}

// ItemCreate inserts item into system.
func (s *UseCase) ItemCreate(ctx context.Context, userID string, input item.CreateItemInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ItemCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	itemID, err := s.adapterStorage.ItemCreate(ctx, userID, input)
	if err != nil {
		return itemID, errors.Wrap(err, "item create error")
	}

	if s.isCacheOn {
		itm, err := s.adapterStorage.ItemGetOne(ctx, userID, input.OrganizationID, itemID)
		if err != nil {
			return "", errors.Wrap(err, "item select from database failed")
		}

		err = s.adapterCache.ItemCreate(ctx, itm)
		if err != nil {
			return "", errors.Wrap(err, "item create in cache failed")
		}

		err = s.adapterCache.ItemInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "item invalidate users in cache failed")
		}
	}

	return itemID, nil
}

// ItemUpdate updates item by id in the system.
func (s *UseCase) ItemUpdate(ctx context.Context, userID string, input item.UpdateItemInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ItemUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ItemUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "item update in database failed")
	}

	if s.isCacheOn {
		itm, err := s.adapterStorage.ItemGetOne(ctx, userID, *input.OrganizationID, *input.ID)
		if err != nil {
			return errors.Wrap(err, "item select from database failed")
		}

		err = s.adapterCache.ItemUpdate(ctx, itm)
		if err != nil {
			return errors.Wrap(err, "item update in cache failed")
		}

		err = s.adapterCache.ItemInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "item invalidate users in cache failed")
		}
	}

	return nil
}

// ItemDelete deletes item by id from the system.
func (s *UseCase) ItemDelete(ctx context.Context, userID string, organizationID string, itemID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.ItemDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.ItemDelete(ctx, userID, organizationID, itemID)
	if err != nil {
		return errors.Wrap(err, "item delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.ItemDelete(ctx, itemID)
		if err != nil {
			return errors.Wrap(err, "item update in cache failed")
		}

		err = s.adapterCache.ItemInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate items in cache failed")
		}
	}

	return nil
}
