package main

import (
	"context"
	"fmt"
	"go-poc/configs"
	"go-poc/internal/application/integration_events"
	"go-poc/internal/application/models/decrease_stock"
	"go-poc/internal/interactor"
	"go-poc/pkg/rabbitmq"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
)

func main() {

	fmt.Println("Consumer service is preparing now..")

	cfg := configs.NewConfig()

	ch, _ := rabbitmq.InitRabbitMQ(cfg)

	consumer, _ := rabbitmq.NewConsumer(ch, cfg)

	consumer.ConsumeMessages(context.Background(), func(msg []byte) bool {

		var consumedEvent integration_events.OrderPlacedEvent

		if err := jsoniter.Unmarshal(msg, &consumedEvent); err != nil {
			fmt.Println(err.Error()) //todo: loglama yapilabilir.
			return false
		}

		decrease_stock_req := decrease_stock.Request{
			MessageId: consumedEvent.MessageId,
			OrderNo:   consumedEvent.Message.OrderId,
		}
		copier.Copy(&decrease_stock_req.Items, &consumedEvent.Message.Items)

		uow := interactor.ResolveUow(cfg, true)
		stock_service := interactor.ResolveStockService(uow)

		decrease_stock_res := stock_service.DecreaseStock(context.Background(), decrease_stock_req)

		if decrease_stock_res.IsSuccess {
			uow.Commit()
		} else {
			uow.Rollback()
		}

		return decrease_stock_res.IsSuccess
	})
}
