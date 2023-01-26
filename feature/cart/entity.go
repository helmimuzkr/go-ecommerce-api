package cart

import (
	prd "e-commerce-api/feature/product/data"
	"e-commerce-api/feature/users"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	UserID   []users.Core
	Product  []prd.ProductNonGorm
	Quantity int
	Subtotal int
}

type CartHandler interface {
	Add() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type CartService interface {
	Add(token interface{}, productID uint, quantity int) error
	GetAll(token interface{}) ([]Core, error)
	Update(token interface{}, cartID uint, quantity int) error
	Delete(token interface{}, cartID uint) error
}

type CartData interface {
	Add(userID uint, productID uint, quantity int) error
	GetAll(userID uint) ([]Core, error)
	Update(userID uint, cartID uint, quantity int) error
	Delete(userID uint, cartID uint) error
}
