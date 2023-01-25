package service

import (
	"e-commerce-api/config"
	"e-commerce-api/feature/product"
	"e-commerce-api/helper"
	"e-commerce-api/mocks"
	"errors"
	"mime/multipart"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	v := validator.New()
	cld := config.NewCloudinary(*config.InitConfig())
	data := mocks.NewProductData(t)
	srv := New(data, v, cld)

	// t.Run("Success add product", func(t *testing.T) {
	// 	// data.On("Add", uint(1), mock.Anything).Return( nil).Once()

	// 	inSrv := product.Core{
	// 		Name:        "Nike",
	// 		Description: "Buruan dah beli",
	// 		Price:       500000,
	// 	}

	// 	_, raw := helper.GenerateJWT(1)
	// 	token := raw.(*jwt.Token)
	// 	token.Valid = true

	// 	file, err := os.Open("./file-test/test.png")
	// 	if err != nil {
	// 		fmt.Println("Error saat membuka file", err)
	// 	}

	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: file.Name(),
	// 	}

	// 	err = srv.Add(token, inSrv, fileHeader)

	// 	assert.Nil(t, err)
	// })

	t.Run("Token invalid", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		token := jwt.New(jwt.SigningMethodHS256)

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Add(token, inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "token tidak valid")
	})

	t.Run("Error in validation input", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Ni",    // Minimal 3 character
			Description: "Burua", // Minimal 5 character
			Price:       500000,  // Minimal 10k
		}

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Add(token, inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Name value must be greater than 3 character")
	})

	// t.Run("Wrong format upload image", func(t *testing.T) {
	// 	inSrv := product.Core{
	// 		Name:        "Nike",
	// 		Description: "Buruan dah beli",
	// 		Price:       500000,
	// 	}
	// 	_, raw := helper.GenerateJWT(1)
	// 	token := raw.(*jwt.Token)
	// 	token.Valid = true

	// 	file, _ := os.Open("./file-test/false-test.pdf")
	// 	defer file.Close()
	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: file.Name(),
	// 	}

	// 	err := srv.Add(token, inSrv, fileHeader)

	// 	assert.NotNil(t, err)
	// 	assert.EqualError(t, err, "kesalahan input user karena format gambar bukan jpg, jpeg, ataupun png")
	// })

	t.Run("failed to upload image server error", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}
		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Add(token, inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "gagal upload gambar karena kesalahan pada sistem server")
	})

	// t.Run("Database error", func(t *testing.T) {
	// 	data.On("Add", uint(1), mock.Anything).Return(errors.New("database error")).Once()

	// 	inSrv := product.Core{
	// 		Name:        "Nike",
	// 		Description: "Buruan dah beli",
	// 		Price:       500000,
	// 	}

	// 	_, raw := helper.GenerateJWT(1)
	// 	token := raw.(*jwt.Token)
	// 	token.Valid = true

	// 	file, _ := os.Open("./file-test/test.png")
	// 	defer file.Close()
	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: file.Name(),
	// 	}

	// 	err := srv.Add(token, inSrv, fileHeader)

	// 	assert.NotNil(t, err)
	// 	assert.EqualError(t, err, "kesalahan pada sistem server")
	// })

}

func TestGetAll(t *testing.T) {
	v := validator.New()
	cld := config.NewCloudinary(*config.InitConfig())
	data := mocks.NewProductData(t)
	srv := New(data, v, cld)

	t.Run("Success get all product", func(t *testing.T) {
		resData := []product.Core{
			{
				ID:          1,
				Name:        "Nike",
				SellerName:  "John",
				City:        "Jakarta",
				Avatar:      "www.cloudinary.com/file/avatar1.jpg",
				Description: "PDIP Solid solid solid",
				Price:       1000000,
				Stock:       1,
				Image:       "www.cloudinary.com/file/image-product2.jpg",
			},
			{
				ID:          2,
				Name:        "Adidas",
				SellerName:  "Doe",
				City:        "Palembang",
				Avatar:      "www.cloudinary.com/file/avatar3.jpg",
				Description: "PDIP Solid solid solid",
				Price:       1000000,
				Stock:       10,
				Image:       "www.cloudinary.com/file/image-product1.jpg",
			},
		}
		data.On("CountProduct").Return(2, nil).Once()
		data.On("GetAll", 10, 0).Return(resData, nil).Once()

		paginate, actual, err := srv.GetAll(1)

		assert.Nil(t, err)
		for i := range actual {
			assert.Equal(t, resData[i].ID, actual[i].ID)
		}
		assert.NotNil(t, paginate)
		data.AssertExpectations(t)
	})

	t.Run("Error while count record", func(t *testing.T) {
		data.On("CountProduct").Return(0, errors.New("database error")).Once()

		paginate, actual, err := srv.GetAll(1)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
		assert.Nil(t, actual)
		assert.Nil(t, paginate)
		data.AssertExpectations(t)
	})

	t.Run("No record data", func(t *testing.T) {
		data.On("CountProduct").Return(0, nil).Once()

		paginate, actual, err := srv.GetAll(1)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data tidak ditemukan")
		assert.Nil(t, actual)
		assert.Nil(t, paginate)
		data.AssertExpectations(t)
	})

	t.Run("Database error while query get all  product", func(t *testing.T) {
		data.On("CountProduct").Return(2, nil).Once()
		data.On("GetAll", 10, 0).Return(nil, errors.New("database error"))

		paginate, actual, err := srv.GetAll(1)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
		assert.Nil(t, actual)
		assert.Nil(t, paginate)
		data.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	v := validator.New()
	cld := config.NewCloudinary(*config.InitConfig())
	data := mocks.NewProductData(t)
	srv := New(data, v, cld)

	t.Run("Success get product", func(t *testing.T) {
		resData := product.Core{
			ID:          1,
			Name:        "Nike",
			SellerName:  "John",
			City:        "Jakarta",
			Avatar:      "www.cloudinary.com/file/avatar1.jpg",
			Description: "PDIP Solid solid solid",
			Price:       1000000,
			Stock:       1,
			Image:       "www.cloudinary.com/file/image-product2.jpg",
		}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()

		actual, err := srv.GetByID(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, resData.ID, actual.ID)
		data.AssertExpectations(t)
	})

	t.Run("Product not found", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("not found")).Once()

		actual, err := srv.GetByID(uint(1))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data tidak ditemukan")
		assert.Empty(t, actual)
		data.AssertExpectations(t)
	})

	t.Run("Error in database", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("database error")).Once()

		actual, err := srv.GetByID(uint(1))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
		assert.Empty(t, actual)
		data.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	v := validator.New()
	cld := config.NewCloudinary(*config.InitConfig())
	data := mocks.NewProductData(t)
	srv := New(data, v, cld)

	t.Run("Success update product wihtout image", func(t *testing.T) {
		input := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		data.On("Update", uint(1), uint(1), input).Return(nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		var fileHeader *multipart.FileHeader

		err := srv.Update(token, uint(1), input, fileHeader)

		assert.Nil(t, err)
	})

	t.Run("error update wihtout image because not found", func(t *testing.T) {
		input := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		data.On("Update", uint(1), uint(1), input).Return(errors.New("terjadi kesalahan pada server karena data user atau product tidak ditemukan")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		var fileHeader *multipart.FileHeader

		err := srv.Update(token, uint(1), input, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada server karena data user atau product tidak ditemukan")
	})

	t.Run("error update wihtout image because server errro", func(t *testing.T) {
		input := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		data.On("Update", uint(1), uint(1), input).Return(errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		var fileHeader *multipart.FileHeader

		err := srv.Update(token, uint(1), input, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Token invalid", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		token := jwt.New(jwt.SigningMethodHS256)

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Update(token, uint(1), inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "token tidak valid")
	})

	t.Run("Error in validation input", func(t *testing.T) {
		inSrv := product.Core{
			Name:        "Ni",    // Minimal 3 character
			Description: "Burua", // Minimal 5 character
			Price:       500000,  // Minimal 10k
		}

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Update(token, uint(1), inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "Name value must be greater than 3 character")
	})

	t.Run("Failed to get product by id not found", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("not found")).Once()

		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Update(token, uint(1), inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data tidak ditemukan")
	})

	t.Run("Failed to get product by id server error", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("database error")).Once()

		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Update(token, uint(1), inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
	})

	// t.Run("Wrong format upload image", func(t *testing.T) {
	// 	inSrv := product.Core{
	// 		Name:        "Nike",
	// 		Description: "Buruan dah beli",
	// 		Price:       500000,
	// 	}
	// 	_, raw := helper.GenerateJWT(1)
	// 	token := raw.(*jwt.Token)
	// 	token.Valid = true

	// 	file, _ := os.Open("./file-test/false-test.pdf")
	// 	defer file.Close()
	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: file.Name(),
	// 	}

	// 	err := srv.Add(token, inSrv, fileHeader)

	// 	assert.NotNil(t, err)
	// 	assert.EqualError(t, err, "kesalahan input user karena format gambar bukan jpg, jpeg, ataupun png")
	// })

	t.Run("failed to upload image server error", func(t *testing.T) {
		resData := product.Core{Avatar: "www.cloudinary.com/file/avatar1.jpg"}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()

		inSrv := product.Core{
			Name:        "Nike",
			Description: "Buruan dah beli",
			Price:       500000,
		}
		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		file, _ := os.Open("./file-test/test.png")
		defer file.Close()
		fileHeader := &multipart.FileHeader{
			Filename: file.Name(),
		}

		err := srv.Update(token, uint(1), inSrv, fileHeader)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "gagal upload gambar karena kesalahan pada sistem server")
	})

	// t.Run("Database error", func(t *testing.T) {
	// 	data.On("Add", uint(1), mock.Anything).Return(errors.New("database error")).Once()

	// 	inSrv := product.Core{
	// 		Name:        "Nike",
	// 		Description: "Buruan dah beli",
	// 		Price:       500000,
	// 	}

	// 	_, raw := helper.GenerateJWT(1)
	// 	token := raw.(*jwt.Token)
	// 	token.Valid = true

	// 	file, _ := os.Open("./file-test/test.png")
	// 	defer file.Close()
	// 	fileHeader := &multipart.FileHeader{
	// 		Filename: file.Name(),
	// 	}

	// 	err := srv.Add(token, inSrv, fileHeader)

	// 	assert.NotNil(t, err)
	// 	assert.EqualError(t, err, "kesalahan pada sistem server")
	// })
}

func TestDelete(t *testing.T) {
	v := validator.New()
	cld := config.NewCloudinary(*config.InitConfig())
	data := mocks.NewProductData(t)
	srv := New(data, v, cld)

	t.Run("Success delete product", func(t *testing.T) {
		resData := product.Core{}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.Nil(t, err)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "token tidak valid")
	})

	t.Run("Failed to get product by id not found", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data tidak ditemukan")
	})

	t.Run("Failed to get product by id server error", func(t *testing.T) {
		data.On("GetByID", uint(1)).Return(product.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Failed to delete product because not found", func(t *testing.T) {
		resData := product.Core{}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()
		data.On("Delete", uint(1), uint(1)).Return(errors.New("tidak ditemukan")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
	})

	t.Run("Failed to delete product because database error", func(t *testing.T) {
		resData := product.Core{}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()
		data.On("Delete", uint(1), uint(1)).Return(errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Failed to delete image on cloud", func(t *testing.T) {
		resData := product.Core{Avatar: "www.cloudinary.com/file/avatar1.jpg"}
		data.On("GetByID", uint(1)).Return(resData, nil).Once()
		data.On("Delete", uint(1), uint(1)).Return(nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Delete(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to destroy image")
	})
}
