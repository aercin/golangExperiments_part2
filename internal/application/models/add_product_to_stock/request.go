package add_product_to_stock

type Request struct {
	StockId   int64  `param:"id" validate:"required"`
	ProductId string `validate:"required"`
	Quantity  int    `validate:"required"`
}
