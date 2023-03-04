package specification

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Specification
	adapterCache   cache.Specification
}

// New is a constructor for UseCase.
func New(storage storage.Specification, cache cache.Specification) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
