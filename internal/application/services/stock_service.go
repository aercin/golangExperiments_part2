package services

import (
	"context"
	application_abstraction "go-poc/internal/application/abstractions"
	"go-poc/internal/application/integration_events"
	"go-poc/internal/application/models/add_product_to_stock"
	"go-poc/internal/application/models/decrease_stock"
	"go-poc/internal/application/models/get_stock"
	"go-poc/internal/application/models/get_stock_product"
	domain_abstraction "go-poc/internal/domain/abstractions"
	"go-poc/internal/domain/entities"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
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

func (s *stockService) GetStockProduct(ctx context.Context, request get_stock_product.Request) get_stock_product.Response {
	sql, _, _ := goqu.Dialect("postgres").
		Select("sp.*").
		From(goqu.T("stocks").As("s")).
		InnerJoin(goqu.T("stock_products").As("sp"), goqu.On(goqu.Ex{"s.id": goqu.I("sp.stock_id")})).
		Where(goqu.Ex{"s.id": request.StockId, "sp.product_id": request.ProductId}).ToSQL()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	stock, err := s.Uow.GetStockRepo().GetSpecificProduct(ctx, sql)

	if err != nil {
		return get_stock_product.Response{
			IsSuccess: false,
		}
	}

	res := get_stock_product.Response{
		IsSuccess: true,
	}

	copier.Copy(&res, &stock.StockProducts[0])

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

func (s *stockService) DecreaseStock(ctx context.Context, request decrease_stock.Request) decrease_stock.Response {

	inbox_repo := s.Uow.GetInboxRepo()

	if inbox_repo.Any(ctx, request.MessageId) {
		return decrease_stock.Response{
			IsSuccess: true,
		}
	}

	sql, _, _ := goqu.Dialect("postgre").From("stocks").Limit(1).ToSQL()

	stock_repo := s.Uow.GetStockRepo()

	stock, err := stock_repo.Get(ctx, sql)

	if err != nil || stock.Id == 0 {
		return decrease_stock.Response{
			IsSuccess: false,
		}
	}

	requestMap := make(map[string]int)
	for _, p := range request.Items {
		requestMap[p.ProductId] = p.Quantity
	}

	isProductAvailableInStock := true
	for _, sp := range stock.StockProducts {
		quantity, isExist := requestMap[sp.ProductId]
		if !isExist || quantity > sp.Quantity {
			isProductAvailableInStock = false
			break
		}
		sp.Quantity -= quantity
		sp.IsModified = true
	}

	inbox_repo.Create(ctx, &entities.InboxMessage{
		MessageId: request.MessageId,
		CreatedOn: time.Now(),
	})

	var serializedStockReportedEvent []byte

	if isProductAvailableInStock {

		stock_repo.Update(ctx, stock)

		serializedStockReportedEvent, _ = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(integration_events.StockReportedEvent{
			MessageId: uuid.NewString(),
			MessageType: []string{
				integration_events.StockDecreasedEventMessageType,
			},
			Message: integration_events.StockReportedEventMessage{
				OrderNo: request.OrderNo,
			},
		})
	} else {
		serializedStockReportedEvent, _ = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(integration_events.StockReportedEvent{
			MessageId: uuid.NewString(),
			MessageType: []string{
				integration_events.StockNotDecreasedEventMessageType,
			},
			Message: integration_events.StockReportedEventMessage{
				OrderNo: request.OrderNo,
			},
		})
	}

	s.Uow.GetOutboxRepo().Create(ctx, &entities.OutboxMessage{
		Message:   string(serializedStockReportedEvent),
		CreatedOn: time.Now(),
	})

	return decrease_stock.Response{
		IsSuccess: true,
	}
}
