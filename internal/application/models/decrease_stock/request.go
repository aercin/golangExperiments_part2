package decrease_stock

type Request struct {
	MessageId string
	OrderNo   string
	Items     []OrderItem
}

type OrderItem struct {
	ProductId string
	Quantity  int
}
