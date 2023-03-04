package comment

import (
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/emenu-api/internal/usecase/adapters/storage"
)

// UseCase is an organization usecase.
type UseCase struct {
	adapterStorage storage.Comment
	adapterCache   cache.Comment
}

// New is a constructor for UseCase.
func New(storage storage.Comment, cache cache.Comment) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache}
}
