package data

import (
	"e-commerce-api/feature/cart"
	"errors"

	"gorm.io/gorm"
)

type cartData struct {
	db *gorm.DB
}

func New(db *gorm.DB) cart.CartData {
	return &cartData{db: db}
}

func (cd *cartData) Add(userID uint, productID uint) error {

	quantity := 1
	cartLama := &Cart{}
	tx := cd.db.Where("user_id = ? AND product_id = ?", userID, productID).FirstOrCreate(cartLama)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected > 0 {

		cartLama.Quantity += 1
		tx := cd.db.Save(cartLama)
		if tx.Error != nil {
			return tx.Error
		}
	} else {

		cartBaru := &Cart{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		tx := cd.db.Create(cartBaru)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}

// GetAll implements cart.CartData
func (cd *cartData) GetAll(userID uint) ([]cart.Core, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}
