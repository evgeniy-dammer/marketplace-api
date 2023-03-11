package user

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is a user usecase.
type UseCase struct {
	adapterStorage storage.User
	adapterCache   cache.User
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.User, cache cache.User, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
