package handler

type addProductReq struct {
	Name        string `json:"product_name" form:"product_name"`
	Description string `json:"description" form:"description"`
	Price       int    `json:"price" form:"price"`
}
type updateProductReq struct {
	Name        string `json:"product_name" form:"product_name"`
	Description string `json:"description" form:"description"`
	Price       int    `json:"price" form:"price"`
	Stock       int    `json:"stock" form:"stock"`
}
