package authentication

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an authentication usecase.
type UseCase struct {
	adapterStorage storage.Authentication
	adapterCache   cache.Authentication
	isTracingOn    bool
	isCacheOn      bool
}

// New constructor for UseCase.
func New(storage storage.Authentication, cache cache.Authentication, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
