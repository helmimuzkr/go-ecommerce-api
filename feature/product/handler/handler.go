package handler

import (
	"e-commerce-api/feature/product"
	"e-commerce-api/helper"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	srv product.ProductService
}

func NewProductHandler(s product.ProductService) product.ProductHandler {
	return &productHandler{srv: s}
}

func (ph *productHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		pr := productRequest{}
		if err := c.Bind(&pr); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		fileHeader, _ := c.FormFile("image")

		pc := product.Core{}
		copier.Copy(&pc, &pr)

		if err := ph.srv.Add(token, pc, fileHeader); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return nil
	}
}

func (ph *productHandler) GetAll() echo.HandlerFunc
func (ph *productHandler) GetByID() echo.HandlerFunc
func (ph *productHandler) Update() echo.HandlerFunc
func (ph *productHandler) Delete() echo.HandlerFunc
