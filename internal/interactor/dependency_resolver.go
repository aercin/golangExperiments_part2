package interactor

import (
	config_abstraction "go-poc/configs/abstractions"
	application_abstractions "go-poc/internal/application/abstractions"
	"go-poc/internal/application/services"
	domain_abstraction "go-poc/internal/domain/abstractions"
	"go-poc/internal/infrastructure/persistence"
)

func ResolveUow(config config_abstraction.Config, isInTransaction bool) domain_abstraction.UnitOfWork {
	db := persistence.NewDbConnection(config)
	uow := persistence.NewUow(db, isInTransaction)
	return uow
}

func ResolveStockService(uow domain_abstraction.UnitOfWork) application_abstractions.StockService {
	return services.NewStockService(uow)
}
