package data

import (
	"e-commerce-api/feature/cart"
	"e-commerce-api/feature/product"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Quantity  int
}

type CartProductDetail struct {
	ID           uint
	ProductImage string
	ProductName  string
	Price        int
	Stock        int
	Quantity     int
}

func ToData(core cart.Core, quantity int) Cart {
	return Cart{
		UserID:    core.ID,
		ProductID: core.ID,
		Quantity:  quantity,
	}
}

func ToCartProduct(cd Cart, pd product.Core) CartProductDetail {
	return CartProductDetail{
		ID:           pd.ID,
		ProductImage: pd.Image,
		ProductName:  pd.Name,
		Price:        pd.Price,
		Stock:        pd.Stock,
		Quantity:     cd.Quantity,
	}
}
