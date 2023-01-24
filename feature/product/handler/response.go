package handler

import "e-commerce-api/feature/product"

type productResponse struct {
	ID          uint   `json:"product_id"`
	Name        string `json:"product_name"`
	Description string `json:"description"`
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

func ToResponse(core product.Core) productResponse {
	return productResponse{
		ID:          core.ID,
		Name:        core.Name,
		Description: core.Description,
		SellerName:  core.SellerName,
		City:        core.City,
		Avatar:      core.Avatar,
		Price:       core.Price,
		Stock:       core.Stock,
		Image:       core.Image,
	}
}

func ToListResponse(cores []product.Core) []productResponse {
	resp := []productResponse{}
	for _, v := range cores {
		temp := productResponse{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			SellerName:  v.SellerName,
			City:        v.City,
			Avatar:      v.Avatar,
			Price:       v.Price,
			Stock:       v.Stock,
			Image:       v.Image,
		}

		resp = append(resp, temp)
	}

	return resp
}
