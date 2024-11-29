package persistence

import (
	"context"
	"go-poc/internal/domain/abstractions"
	"go-poc/internal/domain/entities"

	"gorm.io/gorm"
)

type inboxRepo struct {
	Db *gorm.DB
}

func newInboxRepo(db *gorm.DB) abstractions.InboxRepository {
	return &inboxRepo{
		Db: db,
	}
}

func (rep *inboxRepo) Any(ctx context.Context, messageId string) bool {
	var inboxMessage entities.InboxMessage
	if result := rep.Db.WithContext(ctx).Where("message_id = ?", messageId).First(&inboxMessage); result.Error != nil {
		return false
	}
	return true
}

func (rep *inboxRepo) Create(ctx context.Context, entity *entities.InboxMessage) error {
	result := rep.Db.WithContext(ctx).Create(entity)
	return result.Error
}
