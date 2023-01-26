package handler

import (
	"e-commerce-api/feature/cart"
	"e-commerce-api/helper"
	"net/http"
	"strconv"

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
		userID := c.Get("user").(uint)
		productID := c.Param("productID")
		id, err := strconv.Atoi(productID)

		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Add(userID, uint(id))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil update profil"))
	}
}

// Delete implements cart.CartHandler
func (ch *cartHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user").(uint)
		productID := c.Param("productID")
		id, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Delete(userID, uint(id))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil delete cart"))
	}
}

// GetAll implements cart.CartHandler
func (ch *cartHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user").(uint)

		carts, err := ch.srv.GetAll(userID)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil menunjukkan cart", carts))
	}
}

// Update implements cart.CartHandler
func (ch *cartHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("user").(uint)
		productID := c.Param("productID")
		quantity, err := strconv.Atoi(c.FormValue("quantity"))
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}
		id, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Update(userID, uint(id), quantity)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil update cart"))
	}
}
