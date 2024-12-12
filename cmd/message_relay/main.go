package main

import (
	"context"
	"fmt"
	config_abstractions "go-poc/configs/abstractions"
	application_abstractions "go-poc/internal/application/abstractions"
	"go-poc/internal/interactor"
	"time"

	"github.com/google/uuid"
)

func main() {

	ioc := interactor.InitializeIoc().Scope(fmt.Sprintf("v", uuid.New()))

	interactor.RegisterScopeDependencies(ioc, false)

	var eventDispatcher application_abstractions.EventDispatcher
	var cfg config_abstractions.Config

	err := ioc.Invoke(func(ed application_abstractions.EventDispatcher, config config_abstractions.Config) {
		eventDispatcher = ed
		cfg = config
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Message relay service is started at %v \n", time.Now())

	for {

		if err := eventDispatcher.DispatchEvents(context.Background()); err != nil {
			fmt.Printf("An error occurred thats details : %v\n", err)
		}

		time.Sleep(time.Duration(cfg.GetValue("MessageRelay:CycleTime").(int)) * time.Second)
	}
}
