package handler

import "e-commerce-api/feature/cart"

type CartResponse struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	Image       string `json:"image"`
	ProductName string `json:"product_name"`
	SellerName  string `json:"seller_name"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Stock       int    `json:"stock"`
}

type GetAllCartResp struct {
	ID          uint   `json:"id"`
	ProductID   uint   `json:"product_id"`
	Image       string `json:"image"`
	ProductName string `json:"product_name"`
	SellerName  string `json:"seller_name"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Stock       int    `json:"stock"`
}

func GetAllResponse(data cart.Core) GetAllCartResp {
	return GetAllCartResp{
		ID:          data.ID,
		ProductID:   data.ProductID,
		Image:       data.Image,
		ProductName: data.ProductName,
		SellerName:  data.SellerName,
		Price:       data.Price,
		Quantity:    data.Quantity,
		Stock:       data.Stock,
	}
}
