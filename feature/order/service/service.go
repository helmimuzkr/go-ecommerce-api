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
	qry  order.OrderData
	snap snap.Client
}

func New(q order.OrderData, s snap.Client) order.OrderService {
	return &orderService{
		qry:  q,
		snap: s,
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
	items, err := os.qry.GetItemById(uint(userID), orderID)
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
	for _, v := range items {
		itemMidtrans := midtrans.ItemDetails{
			ID:           fmt.Sprintf("%d", v.ID),
			Name:         v.ProductName,
			MerchantName: v.ProductName,
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
	snapResp, _ := os.snap.CreateTransaction(req)

	updateOrder := order.Core{}
	updateOrder.PaymentURL = snapResp.RedirectURL
	updateOrder.PaymentToken = snapResp.Token
	// Update Query disini

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

func (os *orderService) GetByID(token interface{}, orderID uint) (order.Core, error) {
	return order.Core{}, nil
}

func (os *orderService) Cancel(token interface{}, orderID uint) error {
	return nil
}
