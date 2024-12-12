package interactor

import (
	"go-poc/configs"
	"go-poc/internal/application/services"
	"go-poc/internal/infrastructure/persistence"

	"go-poc/pkg/rabbitmq"

	"go.uber.org/dig"
)

func InitializeIoc() *dig.Container {
	container := dig.New()

	container.Provide(configs.NewConfig)
	container.Provide(rabbitmq.InitRabbitMQ)
	container.Provide(rabbitmq.NewProducer)
	container.Provide(rabbitmq.NewConsumer)

	return container
}

func RegisterScopeDependencies(scope *dig.Scope, isInTransaction bool) {
	scope.Provide(func() bool {
		return isInTransaction
	})
	scope.Provide(persistence.NewDbConnection)
	scope.Provide(persistence.NewUow)
	scope.Provide(persistence.NewStockRepo)
	scope.Provide(persistence.NewInboxRepo)
	scope.Provide(persistence.NewOutboxRepo)
	scope.Provide(services.NewStockService)
	scope.Provide(services.NewEventDispatcher)
}
