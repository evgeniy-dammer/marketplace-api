package order

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Order
	adapterCache   cache.Order
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Order, cache cache.Order, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
