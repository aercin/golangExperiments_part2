package main

import (
	"go-poc/configs"
	"go-poc/internal/api"
)

func main() {
	api.NewHttpServer(configs.NewConfig())
}
