package handler

import (
	"e-commerce-api/feature/users"
	"net/http"
	"strings"
)

type UserReponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ToResponse(data users.Core) UserReponse {
	return UserReponse{
		Username: data.Username,
		Email:    data.Email,
	}
}

type UpdateUserResp struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	City     string `json:"city"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

func UpdateUser(data users.Core) UpdateUserResp {
	return UpdateUserResp{
		Username: data.Username,
		Fullname: data.Fullname,
		Email:    data.Email,
		City:     data.City,
		Phone:    data.Phone,
	}
}

type AvatarResp struct {
	Avatar string `json:"avatar"`
}

func AvatarResponse(data users.Core) AvatarResp {
	return AvatarResp{
		Avatar: data.Avatar,
	}
}

func PrintSuccessReponse(code int, message string, data ...interface{}) (int, interface{}) {
	resp := map[string]interface{}{}
	if len(data) < 2 {
		resp["data"] = (data[0])
	} else {
		resp["data"] = (data[0])
		resp["token"] = data[1].(string)
	}

	if message != "" {
		resp["message"] = message
	}

	return code, resp
}

func PrintSuccessNoData(status int, message string, data interface{}) (int, map[string]interface{}) {
	result := make(map[string]interface{})
	result["message"] = message
	if data != nil {
		result["data"] = data
	}
	return status, result
}

func PrintErrorResponse(msg string) (int, interface{}) {
	resp := map[string]interface{}{}
	code := -1
	if msg != "" {
		resp["message"] = msg
	}

	if strings.Contains(msg, "server") {
		code = http.StatusInternalServerError
	} else if strings.Contains(msg, "format") {
		code = http.StatusBadRequest
	} else if strings.Contains(msg, "not found") {
		code = http.StatusNotFound
	}

	return code, resp
}
