package handler

type productResponse struct {
	ID          uint
	Name        string
	Description string
	SellerName  string
	City        string
	Price       int
	Stock       int
	Image       string
}

type listPorductResponse []productResponse
