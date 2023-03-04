package order

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/order"
	"github.com/pkg/errors"
)

// OrderGetAll returns all orders from the system.
func (s *UseCase) OrderGetAll(userID string, organizationID string) ([]order.Order, error) {
	orders, err := s.adapterStorage.OrderGetAll(userID, organizationID)

	return orders, errors.Wrap(err, "orders select error")
}

// OrderGetOne returns order by id from the system.
func (s *UseCase) OrderGetOne(userID string, organizationID string, orderID string) (order.Order, error) {
	order, err := s.adapterStorage.OrderGetOne(userID, organizationID, orderID)

	return order, errors.Wrap(err, "order select error")
}

// OrderCreate inserts order into system.
func (s *UseCase) OrderCreate(userID string, order order.Order) (string, error) {
	orderID, err := s.adapterStorage.OrderCreate(userID, order)

	return orderID, errors.Wrap(err, "order create error")
}

// OrderUpdate updates order by id in the system.
func (s *UseCase) OrderUpdate(userID string, input order.UpdateOrderInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.OrderUpdate(userID, input)

	return errors.Wrap(err, "order update error")
}

// OrderDelete deletes order by id from the system.
func (s *UseCase) OrderDelete(userID string, organizationID string, orderID string) error {
	err := s.adapterStorage.OrderDelete(userID, organizationID, orderID)

	return errors.Wrap(err, "order delete error")
}
