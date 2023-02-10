package service

import (
	"github.com/evgeniy-dammer/emenu-api/internal/model"
	"github.com/evgeniy-dammer/emenu-api/internal/repository"
)

// ItemService is a organization service
type ItemService struct {
	repo repository.Item
}

// NewItemService is a constructor for ItemService
func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

// GetAll returns all items from the system
func (s *ItemService) GetAll(userId string, organizationId string) ([]model.Item, error) {
	return s.repo.GetAll(userId, organizationId)
}

// GetOne returns item by id from the system
func (s *ItemService) GetOne(userId string, organizationId string, itemId string) (model.Item, error) {
	return s.repo.GetOne(userId, organizationId, itemId)
}

// Create inserts item into system
func (s *ItemService) Create(userId string, organizationId string, item model.Item) (string, error) {
	return s.repo.Create(userId, organizationId, item)
}

// Update updates item by id in the system
func (s *ItemService) Update(userId string, organizationId string, itemId string, input model.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, organizationId, itemId, input)
}

// Delete deletes item by id from the system
func (s *ItemService) Delete(userId string, organizationId string, itemId string) error {
	return s.repo.Delete(userId, organizationId, itemId)
}
