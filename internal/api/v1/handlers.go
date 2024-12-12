package v1

import (
	"context"
	"fmt"
	configs_abstraction "go-poc/configs/abstractions"
	api_abstractions "go-poc/internal/api/abstractions"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/get_stock"
	"go-poc/internal/application/models/get_stock_product"
	"go-poc/internal/interactor"
	"net/http"
	"strconv"

	application_abstractions "go-poc/internal/application/abstractions"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type handlers struct {
	Configs configs_abstraction.Config
	Ioc     *dig.Container
}

func NewHandler(configs configs_abstraction.Config, ioc *dig.Container) api_abstractions.Handlers {
	return &handlers{
		Configs: configs,
		Ioc:     ioc,
	}
}

func (h *handlers) AddProductToStock(c echo.Context) error {

	request := new(add_product_to_stock.Request)

	c.Bind(&request)

	request.StockId, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	child_scope := h.Ioc.Scope(fmt.Sprintf("%v", uuid.New()))

	interactor.RegisterScopeDependencies(child_scope, true)

	var stock_service application_abstractions.StockService

	child_scope.Invoke(func(d_stock_svc application_abstractions.StockService) {
		stock_service = d_stock_svc
	})

	response := stock_service.AddProductToStock(ctx, *request)

	return c.JSON(http.StatusOK, response)
}

func (h *handlers) GetStock(c echo.Context) error {

	request := new(get_stock.Request)

	request.StockId, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	child_scope := h.Ioc.Scope(fmt.Sprintf("%v", uuid.New()))

	interactor.RegisterScopeDependencies(child_scope, false)

	var stock_service application_abstractions.StockService

	child_scope.Invoke(func(d_stock_svc application_abstractions.StockService) {
		stock_service = d_stock_svc
	})

	response := stock_service.GetStock(child_scope, ctx, *request)

	return c.JSON(http.StatusOK, response)
}

func (h *handlers) GetStockProduct(c echo.Context) error {
	request := new(get_stock_product.Request)

	request.StockId, _ = strconv.ParseInt(c.Param("stock_id"), 10, 64)
	request.ProductId = c.Param("product_id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	child_scope := h.Ioc.Scope(fmt.Sprintf("%v", uuid.New()))

	interactor.RegisterScopeDependencies(child_scope, false)

	var stock_service application_abstractions.StockService

	child_scope.Invoke(func(d_stock_svc application_abstractions.StockService) {
		stock_service = d_stock_svc
	})

	response := stock_service.GetStockProduct(ctx, *request)

	return c.JSON(http.StatusOK, response)
}
