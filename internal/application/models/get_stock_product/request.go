package get_stock_product

type Request struct {
	StockId   int64  `param:"stock_id" validate:"required"`
	ProductId string `param:"product_id" validate:"required"`
}
