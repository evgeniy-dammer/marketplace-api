package authorization

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an authorization usecase.
type UseCase struct {
	adapterStorage storage.Authorization
	adapterCache   cache.Authorization
	isTracingOn    bool
	isCacheOn      bool
}

// New constructor for UseCase.
func New(storage storage.Authorization, cache cache.Authorization, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
