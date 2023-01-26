package services

import (
	"e-commerce-api/feature/cart"
	"e-commerce-api/helper"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCartData struct {
	mock.Mock
}

func (m *mockCartData) Add(userID uint, productID uint) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *mockCartData) Delete(userID uint, cartID uint) error {
	args := m.Called(userID, cartID)
	return args.Error(0)
}

func (m *mockCartData) GetAll(userID uint) (interface{}, error) {
	args := m.Called(userID)
	return args.Get(0).([]cart.Core), args.Error(1)
}

func (m *mockCartData) GetAllServer(userID uint) (interface{}, error) {
	args := m.Called(userID)
	return nil, args.Error(1)
}

func (m *mockCartData) Update(userID uint, cartID uint, quantity int) error {
	args := m.Called(userID, cartID, quantity)
	return args.Error(0)
}

func TestAddCart(t *testing.T) {
	repo := new(mockCartData)
	srv := New(repo)

	t.Run("success", func(t *testing.T) {
		repo.On("Add", uint(1), uint(1)).Return(nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Add(pToken, 1)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		_, token := helper.GenerateJWT(0)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Add(pToken, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data not found")
		repo.AssertExpectations(t)
	})

	t.Run("failed to add to cart", func(t *testing.T) {
		repo.On("Add", uint(1), uint(1)).Return(errors.New("failed to add to cart")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Add(pToken, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to add to cart")
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mockCartData)
	srv := New(repo)
	t.Run("success delete item from cart", func(t *testing.T) {
		repo.On("Delete", uint(1), uint(1)).Return(nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 1)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error when cart id not found", func(t *testing.T) {
		repo.On("Delete", uint(1), uint(1)).Return(errors.New("cart not found")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		repo.AssertExpectations(t)
	})

	t.Run("error when server problem occured", func(t *testing.T) {
		repo.On("Delete", uint(1), uint(1)).Return(errors.New("server error")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	t.Run("error when trying to delete item that not exists", func(t *testing.T) {
		repo.On("Delete", uint(1), uint(1)).Return(errors.New("data not found")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, 1)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		repo.AssertExpectations(t)
	})
}
func TestGetAll(t *testing.T) {
	repo := new(mockCartData)
	srv := New(repo)
	t.Run("success get all item in cart", func(t *testing.T) {
		repo.On("GetAll", uint(1)).Return([]cart.Core{}, nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.GetAll(pToken)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("error when no item in cart", func(t *testing.T) {
		repo.On("GetAll", uint(1)).Return([]cart.Core{}, errors.New("data not found")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.GetAll(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Nil(t, res)
		repo.AssertExpectations(t)
	})
}

func TestUpdateCart(t *testing.T) {
	repo := new(mockCartData)
	srv := New(repo)
	t.Run("success update item quantity in cart", func(t *testing.T) {
		repo.On("Update", uint(1), uint(1), 2).Return(nil).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Update(pToken, 1, 2)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("error when server problem occured", func(t *testing.T) {
		repo.On("Update", uint(1), uint(1), 2).Return(errors.New("server error")).Once()
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Update(pToken, 1, 2)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	// t.Run("error when product not found", func(t *testing.T) {
	// 	repo.On("Update", uint(1), uint(1), 2).Return(nil, errors.New("product not found")).Once()
	// 	_, token := helper.GenerateJWT(1)
	// 	pToken := token.(*jwt.Token)
	// 	pToken.Valid = true
	// 	err := srv.Update(pToken, 1, 2)
	// 	assert.ErrorContains(t, err, "tidak ditemukan")
	// 	repo.AssertExpectations(t)
	// })

}
