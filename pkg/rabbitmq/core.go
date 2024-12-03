package rabbitmq

import (
	"go-poc/configs/abstractions"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(cfg abstractions.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(cfg.GetValue("RabbitMQ:BrokerAddress").(string))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func initQueue(ch *amqp.Channel, queue_name string) error {
	if _, err := ch.QueueDeclare(
		queue_name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	return nil
}

func initExchange(ch *amqp.Channel, exchange_name string) error {
	if err := ch.ExchangeDeclare(
		exchange_name, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		return err
	}
	return nil
}

func bindExchangeToQueue(ch *amqp.Channel, exchange_name, queue_name string) error {
	if err := ch.QueueBind(
		queue_name,    // queue name
		"",            // routing key
		exchange_name, // exchange
		false,
		nil); err != nil {
		return err
	}
	return nil
}
