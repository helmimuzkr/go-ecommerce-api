package handler

import (
	"e-commerce-api/feature/order"
	"e-commerce-api/helper"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	srv order.OrderService
}

func New(s order.OrderService) order.OrderHandler {
	return &orderHandler{srv: s}
}

func (oh *orderHandler) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		var cartsReq []orderRequest
		if err := c.Bind(&cartsReq); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		carts := []order.Cart{}
		copier.Copy(&carts, &cartsReq)

		res, err := oh.srv.Create(token, carts)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := paymentResponse{PaymentToken: res.PaymentToken, PaymentURL: res.PaymentURL}

		return c.JSON(helper.SuccessResponse(200, "berhasil memproses order", response))
	}
}

func (oh *orderHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}

func (oh *orderHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error { return nil }
}

func (oh *orderHandler) Cancel() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
