package persistence

import (
	"go-poc/internal/domain/abstractions"

	"gorm.io/gorm"
)

type unitOfWork struct {
	Db *gorm.DB
}

func NewUow(db *gorm.DB, isInTransaction bool) abstractions.UnitOfWork {
	uow := &unitOfWork{}
	if isInTransaction {
		uow.Db = db.Begin()
	} else {
		uow.Db = db
	}
	return uow
}

func (u *unitOfWork) GetInboxRepo() abstractions.InboxRepository {
	return newInboxRepo(u.Db)
}

func (u *unitOfWork) GetOutboxRepo() abstractions.OutboxRepository {
	return newOutboxRepo(u.Db)
}

func (u *unitOfWork) GetStockRepo() abstractions.StockRepository {
	return newStockRepo(u.Db)
}

func (u *unitOfWork) Commit() error {
	return u.Db.Commit().Error
}

func (u *unitOfWork) Rollback() error {
	return u.Db.Rollback().Error
}
