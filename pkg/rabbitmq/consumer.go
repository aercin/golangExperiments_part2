package rabbitmq

import (
	"context"
	"fmt"

	"go-poc/configs/abstractions"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ConsumeMessages(ctx context.Context, cb func([]byte) bool) error
}

type consumer struct {
	ch  *amqp.Channel
	cfg abstractions.Config
}

func NewConsumer(ch *amqp.Channel, cfg abstractions.Config) (Consumer, error) {
	if err := initQueue(ch, cfg.GetValue("RabbitMQ:ConsumeQueue").(string)); err != nil {
		return nil, err
	}

	if err := initExchange(ch, cfg.GetValue("RabbitMQ:ConsumeQueue").(string)); err != nil {
		return nil, err
	}

	if err := bindExchangeToQueue(ch, cfg.GetValue("RabbitMQ:ConsumeQueue").(string), cfg.GetValue("RabbitMQ:ConsumeQueue").(string)); err != nil {
		return nil, err
	}

	return &consumer{
		ch:  ch,
		cfg: cfg,
	}, nil
}

func (c *consumer) ConsumeMessages(ctx context.Context, cb func([]byte) bool) error {

	msgs, err := c.ch.Consume(
		c.cfg.GetValue("RabbitMQ:ConsumeQueue").(string), // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return err
	}

	fmt.Println("Incoming messages are listening...")

	for msg := range msgs {
		isSuccessed := cb(msg.Body)
		if isSuccessed {
			msg.Ack(false)
		} else {
			msg.Nack(false, true)
		}
	}

	return nil
}
