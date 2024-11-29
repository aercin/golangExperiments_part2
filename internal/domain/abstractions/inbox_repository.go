package abstractions

import (
	"context"
	"go-poc/internal/domain/entities"
)

type InboxRepository interface {
	Any(ctx context.Context, messageId string) bool

	Create(ctx context.Context, entity *entities.InboxMessage) error
}
