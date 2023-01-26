package services

import (
	"e-commerce-api/feature/cart"
	"e-commerce-api/helper"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type cartUseCase struct {
	qry cart.CartData
	vld *validator.Validate
}

func New(cd cart.CartData) cart.CartService {
	return &cartUseCase{
		qry: cd,
		vld: validator.New(),
	}
}

// Add implements cart.CartService
func (cuc *cartUseCase) Add(token interface{}, productID uint) error {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return errors.New("data not found")
	}
	err := cuc.qry.Add(uint(id), productID)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements cart.CartService
func (cuc *cartUseCase) Delete(token interface{}, cartID uint) error {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return errors.New("data not found")
	}
	err := cuc.qry.Delete(uint(id), cartID)
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

// GetAll implements cart.CartService
func (cuc *cartUseCase) GetAll(token interface{}) ([]cart.Core, error) {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return nil, errors.New("data not found")
	}
	carts, err := cuc.qry.GetAll(uint(id))
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terdapat masalah pada server"
		}
		return nil, errors.New(msg)
	}
	return carts, nil
}

// Update implements cart.CartService
func (cuc *cartUseCase) Update(token interface{}, cartID uint, quantity int) error {
	id := helper.ExtractToken(token)
	if id <= 0 {
		return errors.New("data not found")
	}
	err := cuc.qry.Update(uint(id), cartID, quantity)
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
