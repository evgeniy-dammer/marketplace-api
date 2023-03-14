package comment

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Comment
	adapterCache   cache.Comment
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Comment, cache cache.Comment, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
