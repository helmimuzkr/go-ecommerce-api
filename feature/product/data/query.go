package data

import (
	"e-commerce-api/feature/product"

	"gorm.io/gorm"
)

type productData struct {
	db *gorm.DB
}

func NewProductData(db *gorm.DB) product.ProductData {
	return &productData{db: db}
}

func (pd *productData) Add(userID uint, newProduct product.Core) error {
	d := ToData(newProduct)
	d.SellerID = userID

	tx := pd.db.Create(&d)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (pd *productData) GetAll(limit int, offset int) ([]product.Core, error)
func (pd *productData) GetByID(productID uint) (product.Core, error)
func (pd *productData) Update(userID uint, productID uint, updateProduct product.Core) error
func (pd *productData) Delete(userID uint, productID uint) error
