package item

import (
	"github.com/evgeniy-dammer/emenu-api/internal/domain/item"
	"github.com/pkg/errors"
)

// ItemGetAll returns all items from the system.
func (s *UseCase) ItemGetAll(userID string, organizationID string) ([]item.Item, error) {
	items, err := s.adapterStorage.ItemGetAll(userID, organizationID)

	return items, errors.Wrap(err, "items select error")
}

// ItemGetOne returns item by id from the system.
func (s *UseCase) ItemGetOne(userID string, organizationID string, itemID string) (item.Item, error) {
	itemSingle, err := s.adapterStorage.ItemGetOne(userID, organizationID, itemID)

	return itemSingle, errors.Wrap(err, "item select error")
}

// ItemCreate inserts item into system.
func (s *UseCase) ItemCreate(userID string, item item.Item) (string, error) {
	itemID, err := s.adapterStorage.ItemCreate(userID, item)

	return itemID, errors.Wrap(err, "item create error")
}

// ItemUpdate updates item by id in the system.
func (s *UseCase) ItemUpdate(userID string, input item.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := s.adapterStorage.ItemUpdate(userID, input)

	return errors.Wrap(err, "item update error")
}

// ItemDelete deletes item by id from the system.
func (s *UseCase) ItemDelete(userID string, organizationID string, itemID string) error {
	err := s.adapterStorage.ItemDelete(userID, organizationID, itemID)

	return errors.Wrap(err, "item delete error")
}
