package cart

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID          uint
	UserID      uint
	ProductID   uint
	Image       string
	ProductName string
	SellerName  string
	Price       int
	Quantity    int
	Stock       int
	Subtotal    int
}

type CartHandler interface {
	Add() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type CartService interface {
	Add(token interface{}, productID uint) error
	GetAll(token interface{}) (interface{}, error)
	Update(token interface{}, cartID uint, quantity int) error
	Delete(token interface{}, cartID uint) error
}

type CartData interface {
	Add(userID uint, productID uint) error
	GetAll(userID uint) (interface{}, error)
	Update(userID uint, cartID uint, quantity int) error
	Delete(userID uint, cartID uint) error
}
