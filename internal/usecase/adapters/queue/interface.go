package queue

import (
	"github.com/evgeniy-dammer/marketplace-api/internal/domain/message"
	"github.com/evgeniy-dammer/marketplace-api/pkg/context"
)

// Message interface.
type Message interface {
	MessageSend(ctx context.Context, message message.Message) error
}
