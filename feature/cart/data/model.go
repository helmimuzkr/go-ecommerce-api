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

type GetAllCart struct {
	ID       uint
	UserID   uint
	Image    string
	Name     string
	Username string
	Price    uint
	Quantity uint
	Stock    uint
}

func ToCore(data GetAllCart) cart.Core {
	return cart.Core{
		ID:          0,
		UserID:      data.UserID,
		ProductID:   data.ID,
		Image:       data.Image,
		ProductName: data.Name,
		SellerName:  data.Username,
		Price:       int(data.Price),
		Quantity:    int(data.Quantity),
		Stock:       int(data.Stock),
		Subtotal:    0,
	}
}

func ToSliceCore(data []GetAllCart) []cart.Core {
	sliceCore := []cart.Core{}
	for _, v := range data {
		sliceCore = append(sliceCore, ToCore(v))
	}
	return sliceCore
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
