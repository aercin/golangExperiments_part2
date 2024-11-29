package abstractions

import (
	"context"
	"go-poc/internal/domain/entities"
)

type OutboxRepository interface {
	Find(ctx context.Context, query string) ([]entities.OutboxMessage, error)

	Delete(ctx context.Context, id int64) error

	Create(ctxt context.Context, entity *entities.OutboxMessage) error
}
