package rabbitmq

import (
	"context"
	"time"

	"go-poc/configs/abstractions"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msg []byte) error
}

type producer struct {
	ch  *amqp.Channel
	cfg abstractions.Config
}

func NewProducer(ch *amqp.Channel, cfg abstractions.Config) (Producer, error) {
	if err := initQueue(ch, cfg.GetValue("RabbitMQ:ProduceQueue").(string)); err != nil { //idempotent
		return nil, err
	}

	if err := initExchange(ch, cfg.GetValue("RabbitMQ:ProduceQueue").(string)); err != nil { //idempotent
		return nil, err
	}

	return &producer{
		ch:  ch,
		cfg: cfg,
	}, nil
}

func (p *producer) PublishMessage(ctx context.Context, msg []byte) error {

	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.cfg.GetValue("RabbitMQ:ProduceTimeout").(int))*time.Second)

	defer cancel()

	if err := p.ch.PublishWithContext(ctx,
		p.cfg.GetValue("RabbitMQ:ProduceQueue").(string), // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		}); err != nil {
		return err
	}

	return nil
}
