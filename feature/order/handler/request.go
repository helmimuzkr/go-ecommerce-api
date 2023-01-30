package handler

type orderRequest struct {
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Subtotal  int  `json:"subtotal"`
}

type cartRequest struct {
	CartID []int `json:"cart_id"`
}

type webHookRequest struct {
	OrderID string `json:"order_id"`
}
