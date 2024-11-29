package v1

import (
	"context"
	configs_abstraction "go-poc/configs/abstractions"
	api_abstractions "go-poc/internal/api/abstractions"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/get_stock"
	"go-poc/internal/interactor"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handlers struct {
	Configs configs_abstraction.Config
}

func NewHandler(configs configs_abstraction.Config) api_abstractions.Handlers {
	return &handlers{
		Configs: configs,
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

	uow := interactor.ResolveUow(h.Configs, true)

	stock_service := interactor.ResolveStockService(uow)

	response := stock_service.AddProductToStock(ctx, *request)

	if response.IsSuccess {
		uow.Commit()
	} else {
		uow.Rollback()
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handlers) GetStock(c echo.Context) error {

	request := new(get_stock.Request)

	request.StockId, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uow := interactor.ResolveUow(h.Configs, false)

	stock_service := interactor.ResolveStockService(uow)

	response := stock_service.GetStock(ctx, *request)

	return c.JSON(http.StatusOK, response)
}
