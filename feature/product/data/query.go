package data

import (
	"e-commerce-api/feature/product"
	"errors"

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
	up := []UserProduct{}
	query := "SELECT products.id, products.name, products.description, products.price, products.stock, products.image, products.created_at, users.username, users.city FROM products JOIN users ON users.id = products.seller_id WHERE products.deleted_at IS NULL AND products.stock > 0 ORDER BY products.id DESC LIMIT ? OFFSET ?"
	txProduct := pd.db.Raw(query, limit, offset).Find(&up)
	if txProduct.Error != nil {
		return nil, txProduct.Error
	}

	return ToSliceCore(up), nil
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
	up := UserProduct{}
	query := "SELECT products.id, products.name, products.description, products.price, products.stock, products.image, products.created_at, users.username, users.city FROM products JOIN users ON users.id = products.seller_id WHERE products.deleted_at IS NULL AND products.id = ?"
	tx := pd.db.Raw(query, productID).First(&up)
	if tx.Error != nil {
		return product.Core{}, tx.Error
	}

	return ToCore(up), nil
}

func (pd *productData) Update(userID uint, productID uint, updateProduct product.Core) error {
	up := ToData(updateProduct)
	tx := pd.db.Model(&Product{}).Where("id = ? AND seller_id = ? AND deleted_at IS NULL", productID, userID).Updates(&up)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected <= 0 {
		return errors.New("terjadi kesalahan pada server karena data user atau product tidak ditemukan")
	}

	return nil
}

func (pd *productData) Delete(userID uint, productID uint) error {
	tx := pd.db.Exec("UPDATE products SET deleted_at=CURRENT_TIMESTAMP, image='' WHERE id=? AND seller_id=?", productID, userID)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected <= 0 {
		return errors.New("terjadi kesalahan pada server karena data user atau product tidak ditemukan")
	}

	return nil
}
