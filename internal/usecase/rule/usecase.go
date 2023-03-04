package rule

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Rule
	adapterCache   cache.Rule
}

// New is a constructor for UseCase.
func New(storage storage.Rule, cache cache.Rule) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
