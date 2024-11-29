package abstractions

type UnitOfWork interface {
	Commit() error
	Rollback() error
	GetInboxRepo() InboxRepository
	GetOutboxRepo() OutboxRepository
	GetStockRepo() StockRepository
}
