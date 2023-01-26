package handler

import (
	"e-commerce-api/feature/cart"

	"github.com/labstack/echo/v4"
)

type cartHandler struct {
	srv cart.CartService
}

func New(c cart.CartService) cart.CartHandler {
	return &cartHandler{
		srv: c}
}

// Add implements cart.CartHandler
func (ch *cartHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic("unimplemented")
	}
}

// Delete implements cart.CartHandler
func (ch *cartHandler) Delete() echo.HandlerFunc {
	panic("unimplemented")
}

// GetAll implements cart.CartHandler
func (ch *cartHandler) GetAll() echo.HandlerFunc {
	panic("unimplemented")
}

// Update implements cart.CartHandler
func (ch *cartHandler) Update() echo.HandlerFunc {
	panic("unimplemented")
}
