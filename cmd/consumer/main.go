package main

import (
	"context"
	"fmt"
	"go-poc/internal/application/integration_events"
	"go-poc/internal/application/models/decrease_stock"
	"go-poc/internal/interactor"

	application_abstractions "go-poc/internal/application/abstractions"
	"go-poc/pkg/rabbitmq"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
)

func main() {

	fmt.Println("Consumer service is preparing now..")

	ioc := interactor.InitializeIoc().Scope(fmt.Sprintf("%v", uuid.New()))
	interactor.RegisterScopeDependencies(ioc, true)

	var consumer rabbitmq.Consumer
	var stock_service application_abstractions.StockService

	ioc.Invoke(func(c rabbitmq.Consumer, stock_svc application_abstractions.StockService) {
		consumer = c
		stock_service = stock_svc
	})

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

		decrease_stock_res := stock_service.DecreaseStock(context.Background(), decrease_stock_req)

		return decrease_stock_res.IsSuccess
	})
}
