package persistence

import (
	"context"
	"go-poc/internal/domain/abstractions"
	"go-poc/internal/domain/entities"

	"gorm.io/gorm"
)

type outboxRepo struct {
	Db *gorm.DB
}

func NewOutboxRepo(db *gorm.DB) abstractions.OutboxRepository {
	return &outboxRepo{
		Db: db,
	}
}

func (rep *outboxRepo) Find(ctx context.Context, query string) ([]entities.OutboxMessage, error) {

	var outboxMessages []entities.OutboxMessage

	if err := rep.Db.WithContext(ctx).Raw(query).Find(&outboxMessages).Error; err != nil {
		return nil, err
	}

	return outboxMessages, nil
}

func (rep *outboxRepo) Delete(ctx context.Context, id int64) error {
	if err := rep.Db.WithContext(ctx).Delete(&entities.OutboxMessage{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (rep *outboxRepo) Create(ctx context.Context, entity *entities.OutboxMessage) error {
	result := rep.Db.WithContext(ctx).Create(&entity)
	return result.Error
}
