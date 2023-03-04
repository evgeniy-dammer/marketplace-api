package table

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Table
	adapterCache   cache.Table
}

// New is a constructor for UseCase.
func New(storage storage.Table, cache cache.Table) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
