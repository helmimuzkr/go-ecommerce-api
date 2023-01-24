package handler

import "e-commerce-api/feature/product"

type productResponse struct {
	ID          uint   `json:"product_id"`
	Name        string `json:"product_name"`
	Description string `json:"description"`
	SellerID    uint   `json:"seller_id"`
	SellerName  string `json:"seller_name"`
	City        string `json:"city"`
	Avatar      string `json:"avatar"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
}

type pagination struct {
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
	TotalRecord int `json:"total_record"`
	TotalPage   int `json:"total_page"`
}

type productWithPagination struct {
	Pagination pagination        `json:"pagination"`
	Data       []productResponse `json:"data"`
	Message    string            `json:"message"`
}

func ToResponse(core []product.Core) []productResponse {
	resp := []productResponse{}
	for _, v := range core {
		temp := productResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			SellerID:    v.Seller.ID,
			SellerName:  v.Seller.Name,
			City:        v.Seller.City,
			Avatar:      v.Seller.Avatar,
			Price:       v.Price,
			Stock:       v.Stock,
			Image:       v.Image,
		}

		resp = append(resp, temp)
	}

	return resp
}
