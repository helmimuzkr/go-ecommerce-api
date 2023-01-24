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

type User struct {
	gorm.Model
	Name     string
	City     string
	Products []Product `gorm:"foreignKey:SellerID"`
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

func ToCore(p Product, u User) product.Core {
	return product.Core{
		ID:          p.ID,
		Name:        p.Name,
		SellerName:  u.Name,
		Description: p.Description,
		City:        u.City,
		Price:       p.Price,
		Stock:       p.Stock,
		Image:       p.Image,
	}
}

func ToSliceCore(p []Product, u []User) []product.Core {
	temp := []product.Core{}
	for i := range p {
		c := ToCore(p[i], u[i])
		temp = append(temp, c)
	}

	return temp
}
