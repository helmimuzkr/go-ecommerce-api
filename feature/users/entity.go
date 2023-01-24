package users

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	Username string
	Name     string
	Email    string
	Password string
	Address  string
	Avatar   string
}

type UserHandler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Profile() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	UpdatePwd() echo.HandlerFunc
}

type UserService interface {
	Login(username, password string) (string, Core, error)
	Register(newUser Core) (Core, error)
	Profile(token interface{}) (Core, error)
	Update(token interface{}, file multipart.FileHeader, updateData Core) (Core, error)
	Delete(token interface{}) error
	UpdatePwd(token interface{}) error
}

type UserData interface {
	Login(username string) (Core, error)
	Register(newUser Core) (Core, error)
	Profile(id uint) (Core, error)
	Update(id uint, updateData Core) (Core, error)
	Delete(id uint) error
	UpdatePwd(id uint, newPwd string) error
}
