package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
	"github.com/pkg/errors"
)

// ItemService is an organization service.
type ItemService struct {
	repo repository.Item
}

// NewItemService is a constructor for ItemService.
func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

// GetAll returns all items from the system.
func (s *ItemService) GetAll(userID string, organizationID string) ([]model.Item, error) {
	items, err := s.repo.GetAll(userID, organizationID)

	return items, errors.Wrap(err, "items select error")
}

// GetOne returns item by id from the system.
func (s *ItemService) GetOne(userID string, organizationID string, itemID string) (model.Item, error) {
	item, err := s.repo.GetOne(userID, organizationID, itemID)

	return item, errors.Wrap(err, "item select error")
}

// Create inserts item into system.
func (s *ItemService) Create(userID string, organizationID string, item model.Item) (string, error) {
	itemID, err := s.repo.Create(userID, organizationID, item)

	return itemID, errors.Wrap(err, "item create error")
}

// Update updates item by id in the system.
func (s *ItemService) Update(userID string, organizationID string, itemID string, input model.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.repo.Update(userID, organizationID, itemID, input)

	return errors.Wrap(err, "item update error")
}

// Delete deletes item by id from the system.
func (s *ItemService) Delete(userID string, organizationID string, itemID string) error {
	err := s.repo.Delete(userID, organizationID, itemID)

	return errors.Wrap(err, "item delete error")
}
