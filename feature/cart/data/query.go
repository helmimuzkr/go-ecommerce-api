package data

import (
	"e-commerce-api/feature/cart"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type cartData struct {
	db *gorm.DB
}

func New(db *gorm.DB) cart.CartData {
	return &cartData{db: db}
}

func (cd *cartData) Add(userID uint, productID uint) error {

	var count int64
	err := cd.db.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	if err != nil {
		return err
	}

	// Check stock
	var stock int
	tx := cd.db.Raw("SELECT stock FROM products WHERE deleted_at IS NULL AND id = ? AND stock > 0", productID).First(&stock)
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("stock produk tidak tesedia")
	}

	if count == 0 {
		// Bikin row baru
		cart := &Cart{
			UserID:    userID,
			ProductID: productID,
			Quantity:  1,
		}

		tx := cd.db.Create(cart)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		// Update row yang ada
		tx := cd.db.Model(&Cart{}).Where("user_id = ? AND product_id = ?", userID, productID).Update("quantity", gorm.Expr("quantity + ?", 1))
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

// GetAll implements cart.CartData
func (cd *cartData) GetAll(userID uint) (interface{}, error) {
	carts := []GetAllCart{}

	err := cd.db.Raw("SELECT carts.id, carts.product_id, products.name, products.image, users.username, products.price, carts.user_id, carts.quantity, products.stock FROM carts JOIN products ON products.id = carts.product_id JOIN users ON users.id = products.seller_id WHERE carts.deleted_at IS NULL AND carts.user_id = ?", userID).Find(&carts).Error
	if err != nil {
		return nil, err
	}

	return ToSliceCore(carts), nil
}

// Delete implements cart.CartData
func (cd *cartData) Delete(userID uint, cartID uint) error {
	tx := cd.db.Where("id = ? AND user_id = ?", cartID, userID).Delete(&Cart{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected <= 0 {
		return errors.New("terjadi kesalahan pada server karena data user atau cart tidak ditemukan")
	}
	return nil
}

// Update implements cart.CartData
func (cd *cartData) Update(userID uint, cartID uint, quantity int) error {
	cart := &Cart{}
	tx := cd.db.Where("user_id = ? AND product_id = ?", userID, cartID).First(cart)
	if tx.Error != nil {
		return fmt.Errorf("Cart not found for user ID: %v and cart ID: %v", userID, cartID)
	}

	// Check if quantity is valid
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	// Update cart
	cart.Quantity = quantity
	tx = cd.db.Save(cart)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
