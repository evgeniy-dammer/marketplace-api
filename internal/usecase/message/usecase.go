package message

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/cache"
	"github.com/evgeniy-dammer/marketplace-api/internal/usecase/adapters/storage"
)

// UseCase is a message usecase.
type UseCase struct {
	adapterStorage storage.Message
	adapterCache   cache.Message
	isTracingOn    bool
	isCacheOn      bool
}

// New is a constructor for UseCase.
func New(storage storage.Message, cache cache.Message, isTracingOn bool, isCacheOn bool) *UseCase {
	return &UseCase{adapterStorage: storage, adapterCache: cache, isTracingOn: isTracingOn, isCacheOn: isCacheOn}
}
