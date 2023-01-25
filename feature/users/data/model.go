package data

import (
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
	Phone    string
	Avatar   string
	Product  []data.Product `gorm:"foreignKey:SellerID"`
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
		Phone:    data.Phone,
	}
}
