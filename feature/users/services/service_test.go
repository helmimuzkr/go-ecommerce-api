package services

import (
	"e-commerce-api/feature/users"
	"e-commerce-api/helper"
	"e-commerce-api/mocks"
	"errors"

	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("sukses register", func(t *testing.T) {
		inputData := users.Core{Username: "audiz", Email: "audi@gmail.com", Password: "audi123"}
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com"}
		repo.On("Register", mock.Anything).Return(resData, nil).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		assert.Equal(t, resData.Username, res.Username)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		inputData := users.Core{Username: "audiz", Email: "audi@gmail.com", Password: "audi123"}
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com"}
		repo.On("Register", mock.Anything).Return(resData, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	t.Run("input tidak valid", func(t *testing.T) {
		inputData := users.Core{Username: "audi", Email: "audi@gmail.com", Password: "audi123"}
		// resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com"}
		// repo.On("Register", mock.Anything).Return(resData, errors.New("input tidak valid")).Once()
		srv := New(repo)
		res, err := srv.Register(inputData)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "greater")
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Berhasil login", func(t *testing.T) {
		inputUsername := "audiz"
		hashed, _ := helper.GeneratePassword("audi123")
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com", Password: hashed}

		repo.On("Login", inputUsername).Return(resData, nil).Once() // simulasi method login pada layer data

		srv := New(repo)
		token, res, err := srv.Login(inputUsername, "audi123")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Tidak ditemukan", func(t *testing.T) {
		inputUsername := "audiz"
		repo.On("Login", inputUsername).Return(users.Core{}, errors.New("data not found")).Once()

		srv := New(repo)
		token, res, err := srv.Login(inputUsername, "be1422")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Salah password", func(t *testing.T) {
		inputUsername := "audi@gmail.com"
		hashed, _ := helper.GeneratePassword("audi123")
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com", Password: hashed}
		repo.On("Login", inputUsername).Return(resData, nil).Once()

		srv := New(repo)
		token, res, err := srv.Login(inputUsername, "be1423")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password tidak sesuai")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		inputUsername := "audi@gmail.com"
		hashed, _ := helper.GeneratePassword("audi123")
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com", Password: hashed}
		repo.On("Login", inputUsername).Return(resData, errors.New("terdapat masalah pada server")).Once() // simulasi method login pada layer data

		srv := New(repo)
		token, res, err := srv.Login(inputUsername, "audi123")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}

func TestProfile(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Sukses lihat profile", func(t *testing.T) {
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com"}

		repo.On("Profile", uint(1)).Return(resData, nil).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Profile(pToken)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Profile(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Profile", uint(5)).Return(users.Core{}, errors.New("data not found")).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(5)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		repo.On("Profile", mock.Anything).Return(users.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("sukses update data", func(t *testing.T) {
		input := users.Core{Username: "audi", Email: "audi@gmail.com"}
		hashed, _ := helper.GeneratePassword("audi123")
		resData := users.Core{ID: uint(1), Username: "audi", Email: "audi@gmail.com", Password: hashed}
		repo.On("Update", uint(1), input).Return(resData, nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		assert.Equal(t, input.Username, res.Username)
		assert.Equal(t, input.Email, res.Email)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		input := users.Core{Username: "audi", Email: "audi@gmail.com"}
		srv := New(repo)
		_, token := helper.GenerateJWT(0)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		input := users.Core{Username: "audi", Email: "audi@gmail.com"}
		repo.On("Update", uint(5), input).Return(users.Core{}, errors.New("data not found")).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(5)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		input := users.Core{Username: "audi", Email: "audi@gmail.com"}
		repo.On("Update", uint(1), input).Return(users.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Sukses menghapus profile", func(t *testing.T) {
		repo.On("Delete", uint(1)).Return(nil).Once()

		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		err := srv.Delete(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Delete", uint(5)).Return(errors.New("data not found")).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(5)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})
}
