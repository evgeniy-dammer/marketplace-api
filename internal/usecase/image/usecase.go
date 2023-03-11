package image

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an image usecase.
type UseCase struct {
	adapterStorage storage.Image
	adapterCache   cache.Image
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Image, cache cache.Image, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
