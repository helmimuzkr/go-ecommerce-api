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

func (os *orderService) Create(token interface{}, carts []int) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, helper.NewCustErr(401, "token tidak valid")
	}

	if len(carts) == 0 {
		return order.Core{}, helper.NewCustErr(404, "tidak bisa memproses order dikarenakan cart kosong")
	}

	// Buat new order
	newOrder := order.Core{}
	newOrder.Invoice = fmt.Sprintf("INV-%d-%s", userID, time.Now().Format("20060102-150405"))
	newOrder.OrderStatus = "Pending"
	newOrder.OrderDate = time.Now().Format("02-01-2006")
	newOrder.PaidDate = "waiting payment"

	// Buat order
	orderID, err := os.qry.CreateOrder(uint(userID), newOrder)
	if err != nil {
		log.Println(err)
		return order.Core{}, helper.NewCustErr(500, "terjadi kesalahan pada sistem server")
	}

	// convert dari cart item ke order item
	for _, cartID := range carts {
		err := os.qry.CreateOrderItem(orderID, uint(cartID))
		if err != nil {
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal membuat order karena item tidak ditemukan"
			}
			return order.Core{}, helper.NewCustErr(code, msg)
		}
	}

	// Ambil item yang sudah dibuat tadi untuk request midtrans
	orders, err := os.qry.GetOrderBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal membuat order karena item tidak ditemukan"
		}
		return order.Core{}, helper.NewCustErr(code, msg)
	}

	// Konversi list item tadi ke list item midtrans
	listItemMidtrans := []midtrans.ItemDetails{}
	for _, v := range orders.Items {
		itemMidtrans := midtrans.ItemDetails{
			ID:           fmt.Sprintf("%d", v.ID),
			Name:         v.ProductName,
			MerchantName: v.Seller,
			Price:        int64(v.Price),
			Qty:          int32(v.Qty),
		}
		listItemMidtrans = append(listItemMidtrans, itemMidtrans)
	}

	// Menghitung total harga
	var totalPrice int
	for _, v := range orders.Items {
		totalPrice += v.Subtotal
	}

	// Buat request untuk midtrans
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newOrder.Invoice,
			GrossAmt: int64(totalPrice),
		},
		Items: &listItemMidtrans,
	}
	// Buat transaksi
	s := helper.NewSnapMidtrans()
	snapResp, _ := s.CreateTransaction(req)

	updateOrder := order.Core{}
	updateOrder.PaymentURL = snapResp.RedirectURL
	updateOrder.PaymentToken = snapResp.Token
	updateOrder.Total = totalPrice

	// Update order
	err = os.qry.UpdateStatus(newOrder.Invoice, updateOrder)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal membuat order"
		}
		return order.Core{}, helper.NewCustErr(code, msg)
	}

	return updateOrder, nil
}

func (os *orderService) GetAll(token interface{}, history string) ([]order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return nil, helper.NewCustErr(401, "token tidak valid")
	}

	var listOrder []order.Core
	switch history {
	case "buy":
		res, err := os.qry.ListOrderBuy(uint(userID))
		if err != nil {
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal menampilkan list order pembelian"
			}
			return nil, helper.NewCustErr(code, msg)
		}
		listOrder = res // Assign res to listOrder
	case "sell":
		res, err := os.qry.ListOrderSell(uint(userID))
		if err != nil {
			log.Println(err)
			msg := ""
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal menampilkan list order penjualan"
			}
			return nil, helper.NewCustErr(code, msg)
		}
		listOrder = res // Assign res to listOrder
	default:
		return nil, helper.NewCustErr(400, "query parameter pada url tidak ditemukan")
	}

	return listOrder, nil
}

func (os *orderService) GetOrderBuy(token interface{}, orderID uint) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, helper.NewCustErr(401, "token tidak valid")
	}

	res, err := os.qry.GetOrderBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal menampilkan detail order pembelian"
		}
		return order.Core{}, helper.NewCustErr(code, msg)
	}

	if res.CustomerID != uint(userID) {
		return order.Core{}, helper.NewCustErr(404, "gagal menampilkan detail order pembelian dikarenakan data tidak ditemukan")
	}

	return res, nil
}

func (os *orderService) GetOrderSell(token interface{}, orderID uint) (order.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return order.Core{}, errors.New("token tidak valid")
	}

	res, err := os.qry.GetOrderSell(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal menampilkan detail order penjualan"
		}
		return order.Core{}, helper.NewCustErr(code, msg)
	}

	if len(res.Items) < 1 {
		return order.Core{}, helper.NewCustErr(404, "gagal menampilkan detail order penjualan dikarenakan data tidak ditemukan")
	}

	return res, nil
}

func (os *orderService) Confirm(token interface{}, orderID uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("token tidak valid")
	}

	res, err := os.qry.GetOrderSell(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal menerima order"
		}
		return helper.NewCustErr(code, msg)
	}

	if res.OrderStatus == "CANCELED" {
		return helper.NewCustErr(400, "gagal menerima order dikarenkan order sudah dibatalkan")
	}

	updateOrder := order.Core{
		OrderStatus:  "ACCEPTED",
		PaymentToken: "",
		PaymentURL:   "",
	}
	if err := os.qry.Confirm(orderID, updateOrder); err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal menerima order"
		}
		return helper.NewCustErr(code, msg)
	}

	return nil
}

func (os *orderService) Cancel(token interface{}, orderID uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return helper.NewCustErr(401, "token tidak valid")
	}

	res, err := os.qry.GetOrderBuy(uint(userID), orderID)
	if err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal melakukan pembatalan order"
		}
		return helper.NewCustErr(code, msg)
	}

	if res.OrderStatus == "ACCEPTED" {
		return helper.NewCustErr(400, "gagal melakukan pembatalan order dikarenakan pembayaran sudah diterima")
	}

	updateOrder := order.Core{
		OrderStatus:  "CANCELED",
		PaymentToken: "",
		PaymentURL:   "",
	}
	if err := os.qry.UpdateStatus(res.Invoice, updateOrder); err != nil {
		log.Println(err)
		code, msg := 500, "terjadi kesalahan pada sistem server"
		if strings.Contains(err.Error(), "not found") {
			code, msg = 404, "gagal melakukan pembatalan order"
		}
		return helper.NewCustErr(code, msg)
	}

	return nil
}

func (os *orderService) UpdateStatus(invoice string, status string, paidAt string) error {
	updateOrder := order.Core{}
	switch status {
	case "settlement":
		updateOrder.OrderStatus = "ACCEPTED"
		updateOrder.PaymentToken = ""
		updateOrder.PaymentURL = ""
		updateOrder.PaidDate = paidAt
		if err := os.qry.UpdateStatus(invoice, updateOrder); err != nil {
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal melakukan udpate status order"
			}
			return helper.NewCustErr(code, msg)
		}
	case "cancel":
		updateOrder.OrderStatus = "CANCEL"
		updateOrder.PaymentToken = ""
		updateOrder.PaymentURL = ""
		if err := os.qry.UpdateStatus(invoice, updateOrder); err != nil {
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal melakukan udpate status order"
			}
			return helper.NewCustErr(code, msg)
		}
	default:
		updateOrder.OrderStatus = "PENDING"
		updateOrder.PaidDate = paidAt
		if err := os.qry.UpdateStatus(invoice, updateOrder); err != nil {
			log.Println(err)
			code, msg := 500, "terjadi kesalahan pada sistem server"
			if strings.Contains(err.Error(), "not found") {
				code, msg = 404, "gagal melakukan udpate status order"
			}
			return helper.NewCustErr(code, msg)
		}
	}

	return nil
}
