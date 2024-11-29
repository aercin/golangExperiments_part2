package services

import (
	"context"
	application_abstraction "go-poc/internal/application/abstractions"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/get_stock"
	domain_abstraction "go-poc/internal/domain/abstractions"
	"go-poc/internal/domain/entities"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jinzhu/copier"
)

type stockService struct {
	Uow domain_abstraction.UnitOfWork
}

func NewStockService(uow domain_abstraction.UnitOfWork) application_abstraction.StockService {
	return &stockService{
		Uow: uow,
	}
}

func (s *stockService) GetStock(ctx context.Context, request get_stock.Request) get_stock.Response {
	sql, _, _ := goqu.Dialect("postgres").From("stocks").Where(goqu.Ex{"id": request.StockId}).ToSQL()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stock, err := s.Uow.GetStockRepo().Get(ctx, sql)

	if err != nil {
		return get_stock.Response{
			IsSuccess: false,
		}
	}

	res := get_stock.Response{
		IsSuccess: true,
	}

	copier.Copy(&res.Products, &stock.StockProducts)

	return res
}

func (s *stockService) AddProductToStock(ctx context.Context, request add_product_to_stock.Request) add_product_to_stock.Response {
	sql, _, _ := goqu.Dialect("postgres").From("stocks").Where(goqu.Ex{"id": request.StockId}).ToSQL()

	stock_repo := s.Uow.GetStockRepo()

	stock, err := stock_repo.Get(ctx, sql)

	if err != nil {
		return add_product_to_stock.Response{
			IsSuccess: false,
		}
	}

	addOrUpdateStockProduct(stock, request.ProductId, request.Quantity)

	if stock.Id == 0 {
		stock.Id = request.StockId
		err = stock_repo.Create(ctx, stock)
	} else {
		err = stock_repo.Update(ctx, stock)
	}

	if err != nil {
		return add_product_to_stock.Response{
			IsSuccess: false,
		}
	}

	return add_product_to_stock.Response{
		IsSuccess: true,
	}
}

func addOrUpdateStockProduct(s *entities.Stock, productId string, quantity int) {
	for i, sp := range s.StockProducts {
		if sp.ProductId == productId {
			s.StockProducts[i].Quantity += quantity
			s.StockProducts[i].IsModified = true
			return
		}
	}

	s.StockProducts = append(s.StockProducts, entities.StockProduct{
		ProductId: productId,
		Quantity:  quantity,
	})
}
