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

type User struct {
	gorm.Model
	Name     string
	City     string
	Avatar   string
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
		ID:   p.ID,
		Name: p.Name,
		Seller: product.Seller{
			ID:     u.ID,
			Name:   u.Name,
			City:   u.City,
			Avatar: u.Avatar,
		},
		Description: p.Description,
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
