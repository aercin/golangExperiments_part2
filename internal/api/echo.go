package api

import (
	"fmt"

	"go-poc/configs/abstractions"

	v1 "go-poc/internal/api/v1"

	"go-poc/pkg/data_validation"

	"github.com/labstack/echo/v4"

	"go-poc/internal/interactor"
)

func NewHttpServer(c abstractions.Config) {
	e := echo.New()
	e.Validator = data_validation.NewEchoDataValidator()

	eg_v1 := e.Group("api/v1.0/stock")

	NewRouter(eg_v1, v1.NewHandler(c, interactor.InitializeIoc()))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", c.GetValue("HttpServer:Port"))))
}
