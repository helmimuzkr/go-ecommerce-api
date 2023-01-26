package data

import (
	"e-commerce-api/feature/cart"

	"gorm.io/gorm"
)

type cartData struct {
	db *gorm.DB
}

func New(db *gorm.DB) cart.CartData {
	return &cartData{db: db}
}

func (cd *cartData) Add(userID uint, productID uint, quantity int) error {
	newCart := &Cart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	tx := cd.db.Create(newCart)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// GetAll implements cart.CartData
func (*cartData) GetAll(userID uint) ([]cart.Core, error) {
	panic("unimplemented")
}

// Delete implements cart.CartData
func (*cartData) Delete(userID uint, cartID uint) error {
	panic("unimplemented")
}

// Update implements cart.CartData
func (*cartData) Update(userID uint, cartID uint, quantity int) error {
	panic("unimplemented")
}
