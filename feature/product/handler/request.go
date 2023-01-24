package handler

type productRequest struct {
	Name        string `json:"product_name" form:"product_name"`
	Description string `json:"description" form:"description"`
	Price       int    `json:"price" form:"price"`
}
