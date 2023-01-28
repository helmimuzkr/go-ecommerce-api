package handler

import (
	"e-commerce-api/feature/order"
	"e-commerce-api/helper"
	"strconv"

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
	return func(c echo.Context) error {
		token := c.Get("user")

		query := c.QueryParam("history")

		res, err := oh.srv.GetAll(token, query)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := []OrderHistoryResponse{}
		copier.Copy(&response, &res)

		return c.JSON(helper.SuccessResponse(200, "berhasil menampilkan list history order", response))
	}
}

func (oh *orderHandler) GetOrderBuy() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("order_id")
		orderID, _ := strconv.Atoi(str)

		res, err := oh.srv.GetOrderBuy(token, uint(orderID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := ToOrderResponse(res)

		return c.JSON(helper.SuccessResponse(200, "berhasil menampilkan detail order", response))
	}
}

func (oh *orderHandler) GetOrderSell() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		str := c.Param("order_id")
		orderID, _ := strconv.Atoi(str)

		res, err := oh.srv.GetOrderSell(token, uint(orderID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		response := ToOrderResponse(res)

		return c.JSON(helper.SuccessResponse(200, "berhasil menampilkan detail order", response))
	}
}

func (oh *orderHandler) Cancel() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		str := c.Param("order_id")
		orderID, _ := strconv.Atoi(str)

		if err := oh.srv.Cancel(token, uint(orderID)); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(200, "berhasil melakukan cancel order"))
	}
}

func (oh *orderHandler) Confirm() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		str := c.Param("order_id")
		orderID, _ := strconv.Atoi(str)

		if err := oh.srv.Confirm(token, uint(orderID)); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(200, "berhasil melakukan menerima order"))
	}
}

func (oh *orderHandler) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notificationPayload map[string]interface{}
		if err := c.Bind(&notificationPayload); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(200, "success menampilkan callback", notificationPayload))
	}
}
