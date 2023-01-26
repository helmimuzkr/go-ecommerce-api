package handler

type orderRequest struct {
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
	Subtotal  int  `json:"subtotal"`
}
