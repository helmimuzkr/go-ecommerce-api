package data

import (
	"e-commerce-api/features/users"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string
	Name        string
	Email       string
	Password    string
	Address     string
	PhoneNumber string
	Avatar      string
	//Product     []data.Product `gorm:"foreignKey:UserID"`
}

func ToCore(data User) users.Core {
	return users.Core{
		ID:       data.ID,
		Username: data.Username,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Address:  data.Address,
		Avatar:   data.Avatar,
	}
}

func CoreToData(data users.Core) User {
	return User{
		Model:    gorm.Model{ID: data.ID},
		Username: data.Username,
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Address:  data.Address,
		Avatar:   data.Avatar,
	}
}
