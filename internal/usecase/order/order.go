package order

import (
	"reflect"

	"github.com/evgeniy-dammer/marketplace-api/internal/domain/order"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
	"github.com/evgeniy-dammer/marketplace-api/pkg/logger"
	"github.com/evgeniy-dammer/marketplace-api/pkg/tracing"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// OrderGetAll returns all orders from the system.
func (s *UseCase) OrderGetAll(ctx context.Context, userID string, organizationID string) ([]order.Order, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrderGetAll")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getAllWithCache(ctx, s, userID, organizationID)
	}

	orders, err := s.adapterStorage.OrderGetAll(ctx, userID, organizationID)

	return orders, errors.Wrap(err, "orders select error")
}

// getAllWithCache returns orders from cache if exists.
func getAllWithCache(ctx context.Context, s *UseCase, userID string, organizationID string) ([]order.Order, error) {
	orders, err := s.adapterCache.OrderGetAll(ctx, organizationID)
	if err != nil {
		logger.Logger.Error("unable to get orders from cache", zap.String("error", err.Error()))
	}

	if len(orders) > 0 {
		return orders, nil
	}

	orders, err = s.adapterStorage.OrderGetAll(ctx, userID, organizationID)
	if err != nil {
		return orders, errors.Wrap(err, "orders select failed")
	}

	if err = s.adapterCache.OrderSetAll(ctx, organizationID, orders); err != nil {
		logger.Logger.Error("unable to add orders into cache", zap.String("error", err.Error()))
	}

	return orders, nil
}

// OrderGetOne returns order by id from the system.
func (s *UseCase) OrderGetOne(ctx context.Context, userID string, organizationID string, orderID string) (order.Order, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrderGetOne")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if s.isCacheOn {
		return getOneWithCache(ctx, s, userID, organizationID, orderID)
	}

	ordr, err := s.adapterStorage.OrderGetOne(ctx, userID, organizationID, orderID)

	return ordr, errors.Wrap(err, "order select error")
}

// getOneWithCache returns order by id from cache if exists.
func getOneWithCache(ctx context.Context, s *UseCase, userID string, organizationID string, orderID string) (order.Order, error) {
	ordr, err := s.adapterCache.OrderGetOne(ctx, orderID)
	if err != nil {
		logger.Logger.Error("unable to get order from cache", zap.String("error", err.Error()))
	}

	if !reflect.ValueOf(ordr).IsZero() {
		return ordr, nil
	}

	ordr, err = s.adapterStorage.OrderGetOne(ctx, userID, organizationID, orderID)
	if err != nil {
		return ordr, errors.Wrap(err, "order select failed")
	}

	if err = s.adapterCache.OrderCreate(ctx, ordr); err != nil {
		logger.Logger.Error("unable to add order into cache", zap.String("error", err.Error()))
	}

	return ordr, nil
}

// OrderCreate inserts order into system.
func (s *UseCase) OrderCreate(ctx context.Context, userID string, input order.CreateOrderInput) (string, error) {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrderCreate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	orderID, err := s.adapterStorage.OrderCreate(ctx, userID, input)
	if err != nil {
		return orderID, errors.Wrap(err, "order create error")
	}

	if s.isCacheOn {
		ordr, err := s.adapterStorage.OrderGetOne(ctx, userID, input.OrganizationID, orderID)
		if err != nil {
			return "", errors.Wrap(err, "order select from database failed")
		}

		err = s.adapterCache.OrderCreate(ctx, ordr)
		if err != nil {
			return "", errors.Wrap(err, "order create in cache failed")
		}

		err = s.adapterCache.OrderInvalidate(ctx)
		if err != nil {
			return "", errors.Wrap(err, "order invalidate users in cache failed")
		}
	}

	return orderID, nil
}

// OrderUpdate updates order by id in the system.
func (s *UseCase) OrderUpdate(ctx context.Context, userID string, input order.UpdateOrderInput) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrderUpdate")
		defer span.End()

		ctx = context.New(ctxt)
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrderUpdate(ctx, userID, input)
	if err != nil {
		return errors.Wrap(err, "order update in database failed")
	}

	if s.isCacheOn {
		ordr, err := s.adapterStorage.OrderGetOne(ctx, userID, *input.OrganizationID, *input.ID)
		if err != nil {
			return errors.Wrap(err, "order select from database failed")
		}

		err = s.adapterCache.OrderUpdate(ctx, ordr)
		if err != nil {
			return errors.Wrap(err, "order update in cache failed")
		}

		err = s.adapterCache.OrderInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "order invalidate users in cache failed")
		}
	}

	return nil
}

// OrderDelete deletes order by id from the system.
func (s *UseCase) OrderDelete(ctx context.Context, userID string, organizationID string, orderID string) error {
	if s.isTracingOn {
		ctxt, span := tracing.Tracer.Start(ctx, "Usecase.OrderDelete")
		defer span.End()

		ctx = context.New(ctxt)
	}

	err := s.adapterStorage.OrderDelete(ctx, userID, organizationID, orderID)
	if err != nil {
		return errors.Wrap(err, "order delete failed")
	}

	if s.isCacheOn {
		err = s.adapterCache.OrderDelete(ctx, orderID)
		if err != nil {
			return errors.Wrap(err, "order update in cache failed")
		}

		err = s.adapterCache.OrderInvalidate(ctx)
		if err != nil {
			return errors.Wrap(err, "invalidate orders in cache failed")
		}
	}

	return nil
}
