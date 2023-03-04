package favorite

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Favorite
	adapterCache   cache.Favorite
}

// New is a constructor for UseCase.
func New(storage storage.Favorite, cache cache.Favorite) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
