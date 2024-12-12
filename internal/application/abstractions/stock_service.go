package abstractions

import (
	"context"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/decrease_stock"
	"go-poc/internal/application/models/get_stock"
	"go-poc/internal/application/models/get_stock_product"

	"go.uber.org/dig"
)

type StockService interface {
	AddProductToStock(ctx context.Context, request add_product_to_stock.Request) add_product_to_stock.Response

	GetStock(ioc *dig.Scope, ctx context.Context, request get_stock.Request) get_stock.Response

	GetStockProduct(ctx context.Context, request get_stock_product.Request) get_stock_product.Response

	DecreaseStock(ctx context.Context, request decrease_stock.Request) decrease_stock.Response
}
