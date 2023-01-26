package data

import (
	_cart "e-commerce-api/feature/cart/data"
	_order "e-commerce-api/feature/order/data"
	"e-commerce-api/feature/product/data"
	"e-commerce-api/feature/users"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Fullname string
	Password string
	Email    string
	City     string
	Address  string
	Phone    string
	Product  []data.Product `gorm:"foreignKey:SellerID"`
	Orders   []_order.Order `gorm:"foreignKey:CustomerID"`
	Carts    []_cart.Cart `gorm:"foreignKey:UserID"`
}

type Product struct {
	ID           uint
	ProductImage string
	ProductName  string
	Price        int
	Stock        int
}

func ToCore(data User) users.Core {
	return users.Core{
		ID:       data.ID,
		Username: data.Username,
		Fullname: data.Fullname,
		Password: data.Password,
		Email:    data.Email,
		City:     data.City,
		Address:  data.Address,
		Phone:    data.Phone,
	}
}

func CoreToData(data users.Core) User {
	return User{
		Model:    gorm.Model{ID: data.ID},
		Username: data.Username,
		Fullname: data.Fullname,
		Password: data.Password,
		Email:    data.Email,
		City:     data.City,
		Address:  data.Address,
		Phone:    data.Phone,
	}
}
