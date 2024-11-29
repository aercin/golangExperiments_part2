package get_stock

type Request struct {
	StockId int64 `param:"id" validate:"required"`
}
