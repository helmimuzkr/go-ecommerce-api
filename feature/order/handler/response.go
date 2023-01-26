package handler

type paymentResponse struct {
	PaymentToken string `json:"token"`
	PaymentURL   string `json:"payment_link"`
}

type OrderHistoryResponse struct {
	ID          uint   `json:"order_id"`
	Invoice     string `json:"invoice"`
	OrderDate   string `json:"order_date"`
	OrderStatus string `json:"order_status"`
}
