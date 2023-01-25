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
	ID    uint   `json:"product_id"`
	Image string `json:"image"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type UserProduct struct {
	ID          uint
	Name        string
	Description string
	Price       int
	Stock       int
	Image       string
	Username    string
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
		SellerName:  up.Username,
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
