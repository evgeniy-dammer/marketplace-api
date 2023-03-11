package favorite

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Favorite
	adapterCache   cache.Favorite
	isTracingOn    bool
}

// New is a constructor for UseCase.
func New(storage storage.Favorite, cache cache.Favorite, isTracingOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn}
}
