package handler

import (
	"e-commerce-api/feature/product"
	"e-commerce-api/helper"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	srv product.ProductService
}

func New(s product.ProductService) product.ProductHandler {
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
		if fileHeader == nil {
			return c.JSON(helper.ErrorResponse("kesalahan input pada user karena tidak mengunggah gambar produk"))
		}
		file, _ := fileHeader.Open()

		pc := product.Core{}
		copier.Copy(&pc, &pr)

		if err := ph.srv.Add(token, pc, file); err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		return c.JSON(helper.SuccessResponse(201, "sukses menambah produk"))
	}
}

func (ph *productHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		str := c.QueryParam("page")
		page, _ := strconv.Atoi(str)

		p, res, err := ph.srv.GetAll(page)
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		paginate := pagination{
			Page:        p["page"].(int),
			Limit:       p["limit"].(int),
			Offset:      p["offset"].(int),
			TotalRecord: p["totalRecord"].(int),
			TotalPage:   p["totalPage"].(int),
		}

		productResp := ToListResponse(res)

		webResponse := productWithPagination{
			Pagination: paginate,
			Data:       productResp,
			Message:    "sukses menampilkan data",
		}

		return c.JSON(200, webResponse)
	}
}

func (ph *productHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		str := c.Param("product_id")
		productID, _ := strconv.Atoi(str)

		res, err := ph.srv.GetByID(uint(productID))
		if err != nil {
			return c.JSON(helper.ErrorResponse(err.Error()))
		}

		webResponse := ToResponse(res)

		return c.JSON(helper.SuccessResponse(200, "sukses menampilkan produk", webResponse))
	}
}

func (ph *productHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (ph *productHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
