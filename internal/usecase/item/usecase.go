package item

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Item
	adapterCache   cache.Item
}

// New is a constructor for UseCase.
func New(storage storage.Item, cache cache.Item) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
