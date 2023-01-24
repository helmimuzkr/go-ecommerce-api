package data

import (
	"e-commerce-api/feature/product"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	SellerID    uint
	Name        string
	Description string `gorm:"type:longtext"`
	Price       int
	Stock       int
	Image       string
}

type ProductNonGorm struct {
	ID    uint
	Image string
	Name  string
	Price int
	Stock int
}

type UserProduct struct {
	ID          uint
	Name        string
	Description string
	Price       int
	Stock       int
	Image       string
	Fullname    string
	City        string
	Avatar      string
}

func ToData(core product.Core) Product {
	return Product{
		Model:       gorm.Model{ID: core.ID},
		Name:        core.Name,
		Description: core.Description,
		Price:       core.Price,
		Stock:       core.Stock,
		Image:       core.Image,
	}
}

func ToCore(up UserProduct) product.Core {
	return product.Core{
		ID:          up.ID,
		Name:        up.Name,
		Description: up.Description,
		SellerName:  up.Fullname,
		City:        up.City,
		Avatar:      up.Avatar,
		Price:       up.Price,
		Stock:       up.Stock,
		Image:       up.Image,
	}
}

func ToSliceCore(up []UserProduct) []product.Core {
	temp := []product.Core{}
	for _, v := range up {
		c := ToCore(v)
		temp = append(temp, c)
	}

	return temp
}
