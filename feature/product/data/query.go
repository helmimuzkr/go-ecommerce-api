package data

import (
	"e-commerce-api/feature/product"

	"gorm.io/gorm"
)

type productData struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.ProductData {
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

func (pd *productData) GetAll(limit int, offset int) ([]product.Core, error) {
	dp := []Product{}
	du := []User{}

	queryProduct := "SELECT products.id, products.name, products.description, products.price, products.stock, products.image FROM products JOIN users ON users.id = products.seller_id ORDER BY products.id DESC LIMIT ? OFFSET ?"
	txProduct := pd.db.Raw(queryProduct, limit, offset).Find(&dp)
	if txProduct.Error != nil {
		return nil, txProduct.Error
	}

	// query := "SELECT products.id, products.name, products.description, products.price, products.stock, products.image, users.id, users.name, users.city, users.avatar FROM products JOIN users ON users.id = products.seller_id ORDER BY products.id DESC LIMIT ? OFFSET ?"
	queryUser := "SELECT users.id, users.name, users.city, users.avatar FROM products JOIN users ON users.id = products.seller_id ORDER BY products.id DESC LIMIT ? OFFSET ?"
	txUser := pd.db.Raw(queryUser, limit, offset).Find(&du)
	if txUser.Error != nil {
		return nil, txUser.Error
	}

	cores := ToSliceCore(dp, du)

	return cores, nil
}

func (pd *productData) CountProduct() (int, error) {
	var total int
	tx := pd.db.Raw("SELECT COUNT(id) FROM products WHERE deleted_at IS NULL").Find(&total)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return total, nil
}

func (pd *productData) GetByID(productID uint) (product.Core, error) {
	return product.Core{}, nil
}

func (pd *productData) Update(userID uint, productID uint, updateProduct product.Core) error {
	return nil
}
func (pd *productData) Delete(userID uint, productID uint) error {
	return nil
}
