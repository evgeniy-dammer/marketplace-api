package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// OrderService is an organization service.
type OrderService struct {
	repo repository.Order
}

// NewOrderService is a constructor for OrderService.
func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

// GetAll returns all orders from the system.
func (s *OrderService) GetAll(userID string, organizationID string) ([]model.Order, error) {
	orders, err := s.repo.GetAll(userID, organizationID)

	return orders, errors.Wrap(err, "orders select error")
}

// GetOne returns order by id from the system.
func (s *OrderService) GetOne(userID string, organizationID string, orderID string) (model.Order, error) {
	order, err := s.repo.GetOne(userID, organizationID, orderID)

	return order, errors.Wrap(err, "order select error")
}

// Create inserts order into system.
func (s *OrderService) Create(userID string, order model.Order) (string, error) {
	orderID, err := s.repo.Create(userID, order)

	return orderID, errors.Wrap(err, "order create error")
}

// Update updates order by id in the system.
func (s *OrderService) Update(userID string, input model.UpdateOrderInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, input)

	return errors.Wrap(err, "order update error")
}

// Delete deletes order by id from the system.
func (s *OrderService) Delete(userID string, organizationID string, orderID string) error {
	err := s.repo.Delete(userID, organizationID, orderID)

	return errors.Wrap(err, "order delete error")
}
