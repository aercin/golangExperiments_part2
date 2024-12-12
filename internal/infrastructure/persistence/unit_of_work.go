package persistence

import (
	"go-poc/internal/domain/abstractions"

	"gorm.io/gorm"
)

type unitOfWork struct {
	Db         *gorm.DB
	StockRepo  abstractions.StockRepository
	InboxRepo  abstractions.InboxRepository
	OutboxRepo abstractions.OutboxRepository
}

func NewUow(db *gorm.DB,
	stock_repo abstractions.StockRepository,
	inbox_repo abstractions.InboxRepository,
	outbox_repo abstractions.OutboxRepository) abstractions.UnitOfWork {
	return &unitOfWork{
		Db:         db,
		StockRepo:  stock_repo,
		InboxRepo:  inbox_repo,
		OutboxRepo: outbox_repo,
	}
}

func (u *unitOfWork) GetInboxRepo() abstractions.InboxRepository {
	return u.InboxRepo
}

func (u *unitOfWork) GetOutboxRepo() abstractions.OutboxRepository {
	return u.OutboxRepo
}

func (u *unitOfWork) GetStockRepo() abstractions.StockRepository {
	return u.StockRepo
}

func (u *unitOfWork) Commit() error {
	return u.Db.Commit().Error
}

func (u *unitOfWork) Rollback() error {
	return u.Db.Rollback().Error
}
