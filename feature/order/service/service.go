package service

import (
	"e-commerce-api/feature/order"
	"e-commerce-api/helper"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type orderService struct {
	qry order.OrderData
}

func New(q order.OrderData) order.OrderService {
	return &orderService{
		qry: q,
	}
}

func (os *orderService) Create(token interface{}, carts []order.Cart) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, errors.New("token tidak valid")
	}

	// Buat new order
	newOrder := order.Core{}
	newOrder.Invoice = fmt.Sprintf("INV/%d/%s", userID, time.Now().Format("20060102/150405"))
	newOrder.OrderStatus = "Pending"
	newOrder.OrderDate = time.Now().Format("02-01-2006")
	// Hitung total belanjaan
	for _, v := range carts {
		newOrder.Total += v.Subtotal
	}

	// Buat order
	orderID, err := os.qry.CreateOrder(uint(userID), newOrder, carts)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return order.Core{}, errors.New(msg)
	}

	// Ambil order item untuk dimasukkan ke midtrans
	res, err := os.qry.GetItemBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return order.Core{}, errors.New(msg)
	}

	// Konversi list item tadi ke list item midtrans
	listItemMidtrans := []midtrans.ItemDetails{}
	for _, v := range res.Items {
		itemMidtrans := midtrans.ItemDetails{
			ID:           fmt.Sprintf("%d", v.ID),
			Name:         v.ProductName,
			MerchantName: v.Seller,
			Price:        int64(v.Price),
			Qty:          int32(v.Qty),
		}
		listItemMidtrans = append(listItemMidtrans, itemMidtrans)
	}

	// Buat request untuk midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newOrder.Invoice,
			GrossAmt: int64(newOrder.Total),
		},
		Items: &listItemMidtrans,
	}
	// Buat transaksi
	s := helper.NewSnapMidtrans()
	snapResp, _ := s.CreateTransaction(req)

	updateOrder := order.Core{}
	updateOrder.PaymentURL = snapResp.RedirectURL
	updateOrder.PaymentToken = snapResp.Token

	// Update payment url dan token
	err = os.qry.Update(uint(userID), orderID, updateOrder)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return order.Core{}, errors.New(msg)
	}

	return updateOrder, nil
}

func (os *orderService) GetAll(token interface{}, history string) ([]order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return nil, errors.New("token tidak valid")
	}

	var listOrder []order.Core
	switch history {
	case "buy":
		res, err := os.qry.GetListOrderBuy(uint(userID))
		if err != nil {
			log.Println(err)
			msg := ""
			if strings.Contains(err.Error(), "not found") {
				msg = "data tidak ditemukan"
			} else {
				msg = "terjadi kesalahan pada sistem server"
			}
			return nil, errors.New(msg)
		}
		listOrder = res
	case "sell":
		res, err := os.qry.GetListOrderSell(uint(userID))
		if err != nil {
			log.Println(err)
			msg := ""
			if strings.Contains(err.Error(), "not found") {
				msg = "data tidak ditemukan"
			} else {
				msg = "terjadi kesalahan pada sistem server"
			}
			return nil, errors.New(msg)
		}
		listOrder = res
	default:
		return nil, errors.New("query parameter pada url tidak ditemukan")
	}

	return listOrder, nil
}

func (os *orderService) GetOrderBuy(token interface{}, orderID uint) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, errors.New("token tidak valid")
	}

	res, err := os.qry.GetItemBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return order.Core{}, errors.New(msg)
	}

	if res.ID == 0 {
		return order.Core{}, errors.New("data tidak ditemukan")
	}

	return res, nil
}

func (os *orderService) GetOrderSell(token interface{}, orderID uint) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, errors.New("token tidak valid")
	}

	res, err := os.qry.GetItemSell(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return order.Core{}, errors.New(msg)
	}

	if len(res.Items) < 1 {
		return order.Core{}, errors.New("data tidak ditemukan")
	}

	return res, nil
}

func (os *orderService) Cancel(token interface{}, orderID uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	res, err := os.qry.GetItemBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	if res.OrderStatus == "ACCEPTED" {
		return errors.New("terjadi kesalahan input pada user. order status yang sudah diterima tidak bisa dibatalkan")
	}

	updateOrder := order.Core{OrderStatus: "CANCELED"}
	if err := os.qry.Update(uint(userID), orderID, updateOrder); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	return nil
}

func (os *orderService) Confirm(token interface{}, orderID uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	res, err := os.qry.GetItemSell(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	if res.OrderStatus == "CANCELED" {
		return errors.New("terjadi kesalahan input pada user. order status yang sudah dibatalkan tidak bisa diterima kembali")
	}

	updateOrder := order.Core{OrderStatus: "ACCEPTED"}
	if err := os.qry.Confirm(orderID, updateOrder); err != nil {
		log.Println(err)
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data tidak ditemukan"
		} else {
			msg = "terjadi kesalahan pada sistem server"
		}
		return errors.New(msg)
	}

	return nil
}
