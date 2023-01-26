package handler

import "e-commerce-api/feature/order"

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

type OrderDetailResponse struct {
	ID           uint           `json:"id"`
	Invoice      string         `json:"invoice"`
	CustomerID   uint           `json:"customer_id"`
	CustomerName string         `json:"customer_name"`
	Address      string         `json:"address"`
	Phone        string         `json:"phone"`
	OrderStatus  string         `json:"order_status"`
	OrderDate    string         `json:"order_date"`
	PaidDate     string         `json:"paid_date"`
	PaymentToken string         `json:"token"`
	PaymentURL   string         `json:"payment_link"`
	Total        int            `json:"total"`
	Items        []ItemResponse `json:"product_list"`
}

type ItemResponse struct {
	ID         uint   `json:"product_id"`
	SellerName string `json:"seller_name"`
	Image      string `json:"image"`
	Quantity   int    `json:"quantity"`
	Subtotal   int    `json:"subtotal"`
}

func ToOrderResponse(c order.Core) OrderDetailResponse {
	return OrderDetailResponse{
		ID:           c.ID,
		Invoice:      c.Invoice,
		CustomerID:   c.CustomerID,
		CustomerName: c.CustomerName,
		Address:      c.Address,
		Phone:        c.Phone,
		OrderStatus:  c.OrderStatus,
		OrderDate:    c.OrderDate,
		PaidDate:     c.PaidDate,
		Total:        c.Total,
		PaymentToken: c.PaymentToken,
		PaymentURL:   c.PaymentURL,
		Items:        ToItemResponse(c.Items),
	}
}

func ToItemResponse(c []order.OrderItem) []ItemResponse {
	listItem := []ItemResponse{}
	for _, v := range c {
		item := ItemResponse{}
		item.ID = v.ID
		item.SellerName = v.Seller
		item.Image = v.Image
		item.Quantity = v.Qty
		item.Subtotal = v.Subtotal

		listItem = append(listItem, item)
	}

	return listItem
}
