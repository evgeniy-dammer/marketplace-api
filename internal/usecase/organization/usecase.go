package organization

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is a organization usecase.
type UseCase struct {
	adapterStorage storage.Organization
	adapterCache   cache.Organization
}

// New is a constructor for UseCase.
func New(storage storage.Organization, cache cache.Organization) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
