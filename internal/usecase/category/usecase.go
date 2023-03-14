package category

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/storage"
)

// UseCase is a category usecase.
type UseCase struct {
	adapterStorage storage.Category
	adapterCache   cache.Category
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Category, cache cache.Category, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
