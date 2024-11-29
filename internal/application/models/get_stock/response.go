package get_stock

type Response struct {
	IsSuccess bool
	Products  []stockProduct
}

type stockProduct struct {
	ProductId string
	Quantity  int
}
