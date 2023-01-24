package service

import (
	"e-commerce-api/feature/product"
	"e-commerce-api/helper"
	"e-commerce-api/mocks"
	"errors"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testFile struct {
	Fieldname string
	Filename  string
	Content   []byte
}

func TestAddProduct(t *testing.T) {
	v := validator.New()
	data := mocks.NewProductData(t)
	srv := New(data, v)

	t.Run("Success add product", func(t *testing.T) {
		data.On("Add", uint(1), mock.Anything).Return(nil).Once()

		inSrv := product.Core{
			Name:        "Nike ardila",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		strToken, _ := helper.GenerateToken(uint(1))
		token := helper.ValidateToken(strToken)

		file, _ := os.Open("./test.png")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.Nil(t, err)
	})

	t.Run("Token invalid", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike ardila",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		token := jwt.New(jwt.SigningMethodHS256)

		file, _ := os.Open("./test.png")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "token tidak valid")
	})

	t.Run("Error in validation input", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Ni",    // Minimal 3 character
			Description: "Burua", // Minimal 5 character
			Price:       500000,  // Minimal 10k
		}

		strToken, _ := helper.GenerateToken(uint(1))
		token := helper.ValidateToken(strToken)

		file, _ := os.Open("./test.png")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Name value must be greater than 3 character")
	})

	t.Run("Wrong format upload image", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike ardila",
			Description: "Buruan dah beli",
			Price:       500000,
		}
		strToken, _ := helper.GenerateToken(uint(1))
		token := helper.ValidateToken(strToken)

		file, _ := os.Open("./false-test.pdf")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "kesalahan input user karena format gambar bukan jpg, jpeg, ataupun png")
	})

	t.Run("failed to upload image server error", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike ardila",
			Description: "Buruan dah beli",
			Price:       500000,
		}
		strToken, _ := helper.GenerateToken(uint(1))
		token := helper.ValidateToken(strToken)

		file, _ := os.Open("")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "gagal upload gambar karena kesalahan pada sistem server")
	})

	t.Run("Database error", func(t *testing.T) {
		data.On("Add", uint(1), mock.Anything).Return(errors.New("database error")).Once()

		inSrv := product.Core{
			Name:        "Nike ardila",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		strToken, _ := helper.GenerateToken(uint(1))
		token := helper.ValidateToken(strToken)

		file, _ := os.Open("./test.png")
		defer file.Close()

		err := srv.Add(token, inSrv, file)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "kesalahan pada sistem server")
	})
}
