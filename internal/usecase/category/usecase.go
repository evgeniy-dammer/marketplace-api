package category

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is a category usecase.
type UseCase struct {
	adapterStorage storage.Category
	adapterCache   cache.Category
}

// New is a constructor for UseCase.
func New(storage storage.Category, cache cache.Category) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
