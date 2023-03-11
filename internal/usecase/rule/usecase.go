package rule

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Rule
	adapterCache   cache.Rule
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Rule, cache cache.Rule, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
