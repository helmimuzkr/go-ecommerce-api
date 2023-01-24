package helper

import (
	"net/http"
	"strings"
)

func SuccessResponse(code int, message string, data ...any) (int, map[string]interface{}) {
	response := make(map[string]interface{})

	response["message"] = message

	switch len(data) {
	case 1:
		response["data"] = data[0]
	case 2:
		response["data"] = data[0]
		response["token"] = data[1]
	}

	return code, response
}

func ErrorResponse(msg string) (int, interface{}) {
	resp := map[string]interface{}{}
	code := http.StatusInternalServerError

	if msg != "" {
		resp["message"] = msg
	}

	switch true {
	case strings.Contains(msg, "server"):
		code = http.StatusInternalServerError
	case strings.Contains(msg, "format"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "tidak ditemukan"):
		code = http.StatusNotFound
	case strings.Contains(msg, "konflik"):
		code = http.StatusConflict
	case strings.Contains(msg, "duplikat"):
		code = http.StatusConflict
	case strings.Contains(msg, "kesalahan input"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "validation"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "unmarshal"):
		resp["message"] = "bad request"
		code = http.StatusBadRequest
	case strings.Contains(msg, "upload"):
		code = http.StatusInternalServerError
	}

	return code, resp
}
