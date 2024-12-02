package persistence

import (
	"context"
	"go-poc/internal/domain/abstractions"
	"go-poc/internal/domain/entities"

	"gorm.io/gorm"
)

type stockRepo struct {
	Db *gorm.DB
}

func newStockRepo(db *gorm.DB) abstractions.StockRepository {
	return &stockRepo{
		Db: db,
	}
}

func (rep *stockRepo) Get(ctx context.Context, query string) (*entities.Stock, error) {
	var stock *entities.Stock
	if err := rep.Db.WithContext(ctx).Preload("StockProducts").Raw(query).Find(&stock).Error; err != nil {
		return nil, err
	}
	return stock, nil
}

func (rep *stockRepo) GetSpecificProduct(ctx context.Context, query string) (*entities.Stock, error) {
	var stockProduct *entities.StockProduct
	if err := rep.Db.WithContext(ctx).Raw(query).Find(&stockProduct).Error; err != nil {
		return nil, err
	}
	stock := &entities.Stock{
		Id: stockProduct.StockId,
		StockProducts: []entities.StockProduct{
			*stockProduct,
		},
	}
	return stock, nil
}

func (rep *stockRepo) Create(ctx context.Context, entity *entities.Stock) error {
	result := rep.Db.WithContext(ctx).Create(&entity)
	return result.Error
}

func (rep *stockRepo) Update(ctx context.Context, entity *entities.Stock) error {
	result := rep.Db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return result.Error
	}
	for _, product := range entity.StockProducts {
		if product.Id != 0 && product.IsModified {
			result = rep.Db.WithContext(ctx).Save(&product)
			if result.Error != nil {
				return result.Error
			}
		}
	}
	return result.Error
}
