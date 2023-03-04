package image

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an image usecase.
type UseCase struct {
	adapterStorage storage.Image
	adapterCache   cache.Image
}

// New is a constructor for UseCase.
func New(storage storage.Image, cache cache.Image) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
