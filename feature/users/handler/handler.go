package handler

import (
	"e-commerce-api/feature/users"
	"e-commerce-api/helper"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userControl struct {
	srv users.UserService
}

func New(srv users.UserService) users.UserHandler {
	return &userControl{
		srv: srv,
	}
}

func (uc *userControl) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := LoginRequest{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "format inputan salah")
		}

		token, res, err := uc.srv.Login(input.Username, input.Password)
		if err != nil {
			return c.JSON(PrintErrorResponse(err.Error()))
		}

		return c.JSON(PrintSuccessReponse(http.StatusOK, "berhasil login", ToResponse(res), token))
	}
}
func (uc *userControl) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := RegisterRequest{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "format inputan salah")
		}

		res, err := uc.srv.Register(*ReqToCore(input))
		if err != nil {
			return c.JSON(PrintErrorResponse(err.Error()))
		}

		return c.JSON(PrintSuccessReponse(http.StatusCreated, "berhasil mendaftar", ToResponse(res)))
	}
}
func (uc *userControl) Profile() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		res, err := uc.srv.Profile(token)
		if err != nil {
			return c.JSON(PrintErrorResponse(err.Error()))
		}

		return c.JSON(PrintSuccessReponse(http.StatusOK, "berhasil lihat profil", res))
	}
}

func (uc *userControl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {

		// var formHeader *multipart.FileHeader
		// var err error

		// if c.Request().MultipartForm != nil && c.Request().MultipartForm.File["file"] != nil {
		// 	formHeader, err = c.FormFile("file")
		// 	if err != nil {
		// 		return c.JSON(http.StatusBadRequest, "File is required")
		// 	}
		// } else {
		// 	formHeader = nil
		// }

		token := c.Get("user")
		input := UpdateRequest{}

		//cek input json dengan format yang benar
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "format inputan salah")
		}

		//validasi input data json
		// if err := c.Validate(input); err != nil {
		// 	return c.JSON(http.StatusBadRequest, "Invalid input data")
		// }

		_, err := uc.srv.Update(token, *ReqToCore(input))
		if err != nil {
			return c.JSON(PrintErrorResponse(err.Error()))
		}
		return c.JSON(helper.SuccessResponse(http.StatusOK, "berhasil update profil"))
	}
}

func (uc *userControl) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		err := uc.srv.Delete(token)
		if err != nil {
			return c.JSON(PrintErrorResponse(err.Error()))
		}

		return c.JSON(PrintSuccessNoData(http.StatusOK, "berhasil delete profil", err))
	}
}
