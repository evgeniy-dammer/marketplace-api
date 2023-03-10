package order

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/evgeniy-dammer/emenu-api/pkg/context"
	"github.com/pkg/errors"
)

// OrderGetAll returns all orders from the system.
func (s *UseCase) OrderGetAll(ctx context.Context, userID string, organizationID string) ([]order.Order, error) {
	orders, err := s.adapterStorage.OrderGetAll(ctx, userID, organizationID)

	return orders, errors.Wrap(err, "orders select error")
}

// OrderGetOne returns order by id from the system.
func (s *UseCase) OrderGetOne(ctx context.Context, userID string, organizationID string, orderID string) (order.Order, error) {
	ordr, err := s.adapterStorage.OrderGetOne(ctx, userID, organizationID, orderID)

	return ordr, errors.Wrap(err, "order select error")
}

// OrderCreate inserts order into system.
func (s *UseCase) OrderCreate(ctx context.Context, userID string, input order.CreateOrderInput) (string, error) {
	orderID, err := s.adapterStorage.OrderCreate(ctx, userID, input)

	return orderID, errors.Wrap(err, "order create error")
}

// OrderUpdate updates order by id in the system.
func (s *UseCase) OrderUpdate(ctx context.Context, userID string, input order.UpdateOrderInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrderUpdate(ctx, userID, input)

	return errors.Wrap(err, "order update error")
}

// OrderDelete deletes order by id from the system.
func (s *UseCase) OrderDelete(ctx context.Context, userID string, organizationID string, orderID string) error {
	err := s.adapterStorage.OrderDelete(ctx, userID, organizationID, orderID)

	return errors.Wrap(err, "order delete error")
}
