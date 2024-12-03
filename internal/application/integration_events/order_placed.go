package integration_events

type OrderPlacedEvent struct {
	MessageId string  `json:"messageId"`
	Message   Message `json:"message"`
}

type Message struct {
	OrderId string      `json:"orderId"`
	Items   []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
