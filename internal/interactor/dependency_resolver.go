package interactor

import (
	config_abstraction "go-poc/configs/abstractions"
	application_abstractions "go-poc/internal/application/abstractions"
	"go-poc/internal/application/services"
	domain_abstraction "go-poc/internal/domain/abstractions"
	"go-poc/internal/infrastructure/persistence"
	"go-poc/pkg/rabbitmq"
)

func ResolveUow(cfg config_abstraction.Config, isInTransaction bool) domain_abstraction.UnitOfWork {
	db := persistence.NewDbConnection(cfg)
	uow := persistence.NewUow(db, isInTransaction)
	return uow
}

func ResolveStockService(uow domain_abstraction.UnitOfWork) application_abstractions.StockService {
	return services.NewStockService(uow)
}

func ResolveEventDispatcher(cfg config_abstraction.Config) application_abstractions.EventDispatcher {
	db := persistence.NewDbConnection(cfg)
	uow := persistence.NewUow(db, false)
	rabbitmq_chan, _ := rabbitmq.InitRabbitMQ(cfg)
	rabbitmq_producer, _ := rabbitmq.NewProducer(rabbitmq_chan, cfg)

	return services.NewEventDispatcher(uow, rabbitmq_producer)
}
