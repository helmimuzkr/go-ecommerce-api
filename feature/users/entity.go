package users

import (
	"e-commerce-api/feature/product/data"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	Username string `validate:"min=5,omitempty"`
	Fullname string
	Password string `validate:"min=5,omitempty"`
	Email    string `validate:"min=5,omitempty,email"`
	Address  string
	City     string
	Address  string
	Phone    string
	Product  []data.ProductNonGorm
}

type UserHandler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type UserService interface {
	Login(username, password string) (string, Core, error)
	Register(newUser Core) (Core, error)
	Profile(token interface{}) (interface{}, error)
	Update(token interface{}, updateData Core) (Core, error)
	Delete(token interface{}) error
}

type UserData interface {
	Login(username string) (Core, error)
	Register(newUser Core) (Core, error)
	Profile(id uint) (interface{}, error)
	Update(id uint, updateData Core) (Core, error)
	Delete(id uint) error
}
