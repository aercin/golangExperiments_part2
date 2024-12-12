package services

import (
	"context"
	application_abstractions "go-poc/internal/application/abstractions"
	domain_abstractions "go-poc/internal/domain/abstractions"
	"go-poc/pkg/rabbitmq"

	"github.com/doug-martin/goqu/v9"
)

type eventDispatcher struct {
	uow      domain_abstractions.UnitOfWork
	producer rabbitmq.Producer
}

func NewEventDispatcher(uow domain_abstractions.UnitOfWork, producer rabbitmq.Producer) application_abstractions.EventDispatcher {
	return &eventDispatcher{
		uow:      uow,
		producer: producer,
	}
}

func (ed *eventDispatcher) DispatchEvents(ctx context.Context) error {

	query, _, err := goqu.Dialect("postgres").From("outbox_messages").Order(goqu.I("created_on").Asc()).Limit(100).ToSQL()

	if err != nil {
		return err
	}

	outboxMessages, err := ed.uow.GetOutboxRepo().Find(ctx, query)
	if err != nil {
		return err
	}

	for _, outboxMessage := range outboxMessages {

		err := ed.producer.PublishMessage(ctx, []byte(outboxMessage.Message))

		if err == nil {
			ed.uow.GetOutboxRepo().Delete(ctx, outboxMessage.Id)
		}
	}

	return nil
}
