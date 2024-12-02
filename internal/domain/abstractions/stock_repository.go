package abstractions

import (
	"context"
	"go-poc/internal/domain/entities"
)

type StockRepository interface {
	Get(ctx context.Context, query string) (*entities.Stock, error)

	GetSpecificProduct(ctx context.Context, query string) (*entities.Stock, error)

	Create(ctx context.Context, entity *entities.Stock) error

	Update(ctx context.Context, entity *entities.Stock) error
}
