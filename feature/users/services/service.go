package services

import (
	"e-commerce-api/config"
	"e-commerce-api/feature/users"
	"e-commerce-api/helper"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type userUseCase struct {
	qry users.UserData
	vld *validator.Validate
}

func New(ud users.UserData) users.UserService {
	return &userUseCase{
		qry: ud,
		vld: validator.New(),
	}
}

func (uuc *userUseCase) Login(username, password string) (string, users.Core, error) {
	res, err := uuc.qry.Login(username)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terdapat masalah pada server"
		}
		return "", users.Core{}, errors.New(msg)
	}

	if err := helper.CheckPassword(res.Password, password); err != nil {
		log.Println("login compare", err.Error())
		return "", users.Core{}, errors.New("password tidak sesuai")
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = res.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	useToken, _ := token.SignedString([]byte(config.JWT_KEY))

	return useToken, res, nil

}
func (uuc *userUseCase) Register(newUser users.Core) (users.Core, error) {
	hashed, err := helper.GeneratePassword(newUser.Password)
	if err != nil {
		log.Println("bcrypt error ", err.Error())
		return users.Core{}, errors.New("password process error")
	}

	err = uuc.vld.Struct(&newUser)
	if err != nil {
		msg := helper.ValidationErrorHandle(err)
		fmt.Println("msg", msg)
		return users.Core{}, errors.New(msg)
	}

	newUser.Password = string(hashed)
	res, err := uuc.qry.Register(newUser)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "duplicated") {
			msg = "data sudah terdaftar"
		} else {
			msg = "terdapat masalah pada server"
		}
		return users.Core{}, errors.New(msg)
	}

	return res, nil
}
func (uuc *userUseCase) Profile(token interface{}) (users.Core, error) {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return users.Core{}, errors.New("data not found")
	}
	res, err := uuc.qry.Profile(uint(id))
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terdapat masalah pada server"
		}
		return users.Core{}, errors.New(msg)
	}

	return res, nil
}

func (uuc *userUseCase) Update(token interface{}, updateData users.Core) (users.Core, error) {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return users.Core{}, errors.New("data not found")
	}
	//ambil user data
	// res, err := uuc.qry.Profile(uint(id))
	// if err != nil {
	// 	return users.Core{}, err
	// }

	//ambil avatar ID
	// oldAvatarPublicID := helper.GetPublicID(user.Avatar)

	// if formHeader != nil {
	// 	// cek file yang di upload
	// 	formFile, err := formHeader.Open()
	// 	if err != nil {
	// 		return users.Core{}, errors.New("input tidak sesuai")
	// 	}

	// 	uploadUrl, err := helper.UploadFile(formFile)
	// 	if err != nil {
	// 		return users.Core{}, errors.New("input tidak sesuai")
	// 	}

	// 	//ambil avatar ID baru
	// 	newAvatarPublicID := helper.GetPublicID(uploadUrl)

	// 	if newAvatarPublicID != oldAvatarPublicID {
	// 		// hapus avatar lama
	// 		err = helper.DestroyFile(oldAvatarPublicID)
	// 		if err != nil {
	// 			return users.Core{}, err
	// 		}
	// 	}
	// 	updateData.Avatar = uploadUrl
	// } else {
	// 	updateData.Avatar = user.Avatar
	// }

	//update user data
	res, err := uuc.qry.Update(uint(id), updateData)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terdapat masalah pada server"
		}
		return users.Core{}, errors.New(msg)
	}

	return res, nil
}

func (uuc *userUseCase) Delete(token interface{}) error {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return errors.New("data not found")
	}

	err := uuc.qry.Delete(uint(id))
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terdapat masalah pada server"
		}
		return errors.New(msg)
	}

	return nil
}
