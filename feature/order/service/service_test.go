package service

import (
	"e-commerce-api/feature/order"
	"e-commerce-api/helper"
	"e-commerce-api/mocks"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	carts := []order.Cart{
		{
			ProductID: 1,
			Quantity:  1,
			Subtotal:  10000,
		},
	}
	// resGetItem := []order.OrderItem{
	// 	{
	// 		ID:          1,
	// 		ProductName: "Nike",
	// 		Seller:      "John",
	// 		City:        "Depok",
	// 		Image:       "www.cloudinary.com/product.jpg",
	// 		Price:       10000,
	// 		Qty:         1,
	// 		Subtotal:    10000,
	// 	},
	// }

	inCreate := order.Core{
		Invoice:     fmt.Sprintf("INV/%d/%s", 1, time.Now().Format("20060102/150405")),
		OrderStatus: "Pending",
		OrderDate:   time.Now().Format("02-01-2006"),
		Total:       10000,
	}

	// t.Run("Success create order", func(t *testing.T) {
	// 	repo.On("CreateOrder", uint(1), inCreate, carts).Return(resCreate, nil).Once()
	// 	repo.On("GetItemBuy", uint(1), uint(resCreate)).Return(resGetItem, nil).Once()
	// })

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		actual, err := srv.Create(token, carts)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
		assert.Empty(t, actual)
	})

	t.Run("Create order data not found", func(t *testing.T) {
		repo.On("CreateOrder", uint(1), inCreate, carts).Return(uint(0), errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, carts)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Empty(t, actual)
	})
	t.Run("Create order server error", func(t *testing.T) {
		repo.On("CreateOrder", uint(1), inCreate, carts).Return(uint(0), errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, carts)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
		assert.Empty(t, actual)
	})
	t.Run("Get item not found", func(t *testing.T) {
		resCreate := uint(1)
		repo.On("CreateOrder", uint(1), inCreate, carts).Return(resCreate, nil).Once()
		repo.On("GetItemBuy", uint(1), resCreate).Return(order.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, carts)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Empty(t, actual)
	})
	t.Run("Get item server error", func(t *testing.T) {
		resCreate := uint(1)
		repo.On("CreateOrder", uint(1), inCreate, carts).Return(resCreate, nil).Once()
		repo.On("GetItemBuy", uint(1), resCreate).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.Create(token, carts)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
		assert.Empty(t, actual)
	})
}

func TestGetAllOrder(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	resRepo := []order.Core{
		{
			ID:          1,
			Invoice:     fmt.Sprintf("INV/%d/%s", 1, time.Now().Format("20060102/150405")),
			OrderStatus: "Pending",
			OrderDate:   time.Now().Format("20060102/150405"),
		},
	}

	t.Run("Success get all buy order", func(t *testing.T) {
		repo.On("GetListOrderBuy", uint(1)).Return(resRepo, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "buy")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, resRepo[0].ID, actual[0].ID)
	})

	t.Run("Success get all sell order", func(t *testing.T) {
		repo.On("GetListOrderSell", uint(1)).Return(resRepo, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "sell")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, resRepo[0].ID, actual[0].ID)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		actual, err := srv.GetAll(token, "buy")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
		assert.Empty(t, actual)
	})

	t.Run("Get all buy not found", func(t *testing.T) {
		repo.On("GetListOrderBuy", uint(1)).Return(nil, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "buy")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Nil(t, actual)
	})

	t.Run("Get all sell not found", func(t *testing.T) {
		repo.On("GetListOrderSell", uint(1)).Return(nil, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "sell")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Nil(t, actual)
	})
	t.Run("Get all buy server error", func(t *testing.T) {
		repo.On("GetListOrderBuy", uint(1)).Return(nil, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "buy")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "kesalahan pada sistem server")
		assert.Nil(t, actual)
	})

	t.Run("Get all sell server error", func(t *testing.T) {
		repo.On("GetListOrderSell", uint(1)).Return(nil, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "sell")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "kesalahan pada sistem server")
		assert.Nil(t, actual)
	})

	t.Run("Wrong input query from user", func(t *testing.T) {
		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetAll(token, "false")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Nil(t, actual)
	})
}

func TestGetOrderSell(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	resRepo := order.Core{
		ID: 1,
		Items: []order.OrderItem{
			{
				ID:          1,
				ProductName: "Nike",
			},
		},
	}

	t.Run("Success get order selling", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(resRepo, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderSell(token, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, resRepo.ID, actual.ID)
		assert.NotNil(t, actual.Items)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		actual, err := srv.GetOrderSell(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
		assert.Empty(t, actual)
	})

	t.Run("not found", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderSell(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Empty(t, actual)
	})

	t.Run("erver error", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderSell(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "kesalahan pada sistem server")
		assert.Empty(t, actual)
	})

	t.Run("order item not found", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderSell(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Empty(t, actual)
	})
}

func TestGetOrderBuy(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	resRepo := order.Core{
		ID: 1,
		Items: []order.OrderItem{
			{
				ID:          1,
				ProductName: "Nike",
			},
		},
	}

	t.Run("Success get order selling", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(resRepo, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderBuy(token, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, resRepo.ID, actual.ID)
		assert.NotNil(t, actual.Items)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		actual, err := srv.GetOrderBuy(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
		assert.Empty(t, actual)
	})

	t.Run("not found", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderBuy(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Empty(t, actual)
	})

	t.Run("erver error", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderBuy(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "kesalahan pada sistem server")
		assert.Empty(t, actual)
	})

	t.Run("order item not found", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		actual, err := srv.GetOrderBuy(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
		assert.Empty(t, actual)
	})
}

func TestCancelOrder(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	t.Run("Success cancel order", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemBuy", uint(1), uint(1)).Return(resRepo, nil).Once()

		resRepo.OrderStatus = "CANCELED"
		repo.On("Update", uint(1), uint(1), resRepo).Return(nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.Nil(t, err)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
	})

	t.Run("Get item not found", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
	})

	t.Run("Get item server error", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Order status cant be accepted", func(t *testing.T) {
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Cancel order server error", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		resRepo.OrderStatus = "CANCELED"
		repo.On("Update", uint(1), uint(1), resRepo).Return(errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Cancel, order id not found", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		resRepo.OrderStatus = "CANCELED"
		repo.On("Update", uint(1), uint(1), resRepo).Return(errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
	})
}
func TestConfirm(t *testing.T) {
	repo := mocks.NewOrderData(t)
	srv := New(repo)

	t.Run("Success confirm order", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemSell", uint(1), uint(1)).Return(resRepo, nil).Once()

		resRepo.OrderStatus = "ACCEPTED"
		repo.On("Confirm", uint(1), resRepo).Return(nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Confirm(token, uint(1))

		assert.Nil(t, err)
	})

	t.Run("Token invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		err := srv.Confirm(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token tidak valid")
	})

	t.Run("Get item not found", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Confirm(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
	})

	t.Run("Get item server error", func(t *testing.T) {
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Confirm(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Order status cant be accepted", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "CANCELED"}
		repo.On("GetItemSell", uint(1), uint(1)).Return(resRepo, nil).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Confirm(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak bisa diterima")
	})

	t.Run("Cancel order server error", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemBuy", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		resRepo.OrderStatus = "CANCELED"
		repo.On("Update", uint(1), uint(1), resRepo).Return(errors.New("database error")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Cancel(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terjadi kesalahan pada sistem server")
	})

	t.Run("Confirm, order id not found", func(t *testing.T) {
		resRepo := order.Core{OrderStatus: "Pending"}
		repo.On("GetItemSell", uint(1), uint(1)).Return(order.Core{}, nil).Once()

		resRepo.OrderStatus = "ACCEPTED"
		repo.On("Confirm", uint(1), resRepo).Return(errors.New("not found")).Once()

		_, raw := helper.GenerateJWT(1)
		token := raw.(*jwt.Token)
		token.Valid = true

		err := srv.Confirm(token, uint(1))

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "data tidak ditemukan")
	})
}
