package api

import (
	"go-poc/internal/api/abstractions"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/get_stock"

	"go-poc/pkg/data_validation"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Group, h abstractions.Handlers) {
	e.POST("/:id/product", h.AddProductToStock, data_validation.ValidationMiddleware(&add_product_to_stock.Request{}))
	e.GET("/:id", h.GetStock, data_validation.ValidationMiddleware(&get_stock.Request{}))
}
