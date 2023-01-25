package data

import (
	"e-commerce-api/feature/users"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Name     string
	Password string
	Email    string
	City     string
	Phone    string
	Avatar   string

	//Product     []data.Product `gorm:"foreignKey:UserID"`
}

func ToCore(data User) users.Core {
	return users.Core{
		ID:       data.ID,
		Username: data.Username,
		Name:     data.Name,
		Password: data.Password,
		Email:    data.Email,
		City:     data.City,
		Phone:    data.Phone,
	}
}

func CoreToData(data users.Core) User {
	return User{
		Model:    gorm.Model{ID: data.ID},
		Username: data.Username,
		Name:     data.Name,
		Password: data.Password,
		Email:    data.Email,
		City:     data.City,
		Phone:    data.Phone,
	}
}
