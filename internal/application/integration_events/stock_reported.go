package integration_events

type StockReportedEvent struct {
	MessageId   string                    `json:"messageId"`
	MessageType []string                  `json:"messageType"`
	Message     StockReportedEventMessage `json:"message"`
}

type StockReportedEventMessage struct {
	OrderNo string `json:"orderId"`
}
