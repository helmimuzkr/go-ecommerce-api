package handler

import "e-commerce-api/feature/users"

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UpdateRequest struct {
	Username string `json:"username" form:"username" validate:"omitempty"`
	Name     string `json:"name" form:"name" validate:"omitempty"`
	Password string `json:"password" form:"password" validate:"omitempty"`
	Email    string `json:"email" form:"email" validate:"omitempty"`
	City     string `json:"city" form:"city" validate:"omitempty"`
	Phone    string `json:"phone" form:"phone" validate:"omitempty"`
	Avatar   string `json:"avatar" form:"avatar" validate:"omitempty"`
}

func ReqToCore(data interface{}) *users.Core {
	res := users.Core{}

	switch data.(type) {
	case LoginRequest:
		cnv := data.(LoginRequest)
		res.Username = cnv.Username
		res.Password = cnv.Password
	case RegisterRequest:
		cnv := data.(RegisterRequest)
		res.Email = cnv.Email
		res.Username = cnv.Username
		res.Password = cnv.Password
	case UpdateRequest:
		cnv := data.(UpdateRequest)
		res.Username = cnv.Username
		res.Name = cnv.Name
		res.Password = cnv.Password
		res.Email = cnv.Email
		res.City = cnv.City
		res.Phone = cnv.Phone
	default:
		return nil
	}

	return &res
}
