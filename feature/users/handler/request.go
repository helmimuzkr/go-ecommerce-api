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
	Username string `json:"username" form:"username"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Address  string `json:"address" form:"address"`
	Avatar   string `json:"avatar" form:"avatar"`
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
		res.Email = cnv.Email
		res.Password = cnv.Password
		res.Address = cnv.Address
		res.Avatar = cnv.Avatar
	default:
		return nil
	}

	return &res
}
