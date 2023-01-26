package handler

import (
	"e-commerce-api/feature/cart"
	"e-commerce-api/helper"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
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
		token := c.Get("user")
		productID := c.Param("product_id")
		id, err := strconv.Atoi(productID)

		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Add(token, uint(id))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil menambahkan produk ke keranjang"))
	}
}

// Delete implements cart.CartHandler
func (ch *cartHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		productID := c.Param("productID")
		id, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Delete(token, uint(id))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil menghapus produk di keranjang"))
	}
}

// GetAll implements cart.CartHandler
func (ch *cartHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		carts, err := ch.srv.GetAll(token)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		sliceResp := []GetAllCartResp{}
		copier.Copy(&sliceResp, carts.([]cart.Core))
		fmt.Println(sliceResp)
		fmt.Println(carts)
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil melihat keranjang", sliceResp))
	}
}

// Update implements cart.CartHandler
func (ch *cartHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		productID := c.Param("productID")
		quantity, err := strconv.Atoi(c.FormValue("quantity"))
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}
		id, err := strconv.Atoi(productID)
		if err != nil {
			return c.JSON(helper.ErrorResponse("Kesalahan pada input user"))
		}

		err = ch.srv.Update(token, uint(id), quantity)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil update keranjang"))
	}
}
