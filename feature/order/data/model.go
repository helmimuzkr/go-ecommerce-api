package data

import (
	"e-commerce-api/feature/order"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID   uint
	Invoice      string
	TotalPrice   int
	OrderStatus  string
	OrderDate    string
	PaidDate     string
	PaymentURL   string
	PaymentToken string
	OrderItems   []OrderItem `foreignKey:"OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
	Subtotal  int
}

type OrderModel struct {
	ID           uint
	Invoice      string
	Fullname     string
	Address      string
	Phone        string
	OrderStatus  string
	OrderDate    string
	PaidDate     string
	Total        int
	PaymentURL   string
	PaymentToken string
	Items        []OrderItemModel `gorm:"-"`
}

type OrderItemModel struct {
	ID       uint
	Name     string
	Username string
	City     string
	Price    int
	Quantity int
	Subtotal int
}

func ToModel(oc order.Core) Order {
	return Order{
		Invoice:      oc.Invoice,
		TotalPrice:   oc.Total,
		OrderStatus:  oc.OrderStatus,
		OrderDate:    oc.OrderDate,
		PaymentToken: oc.PaymentToken,
		PaymentURL:   oc.PaymentURL,
	}
}

func ToCoreItem(oim OrderItemModel) order.OrderItem {
	return order.OrderItem{
		ID:          oim.ID,
		ProductName: oim.Name,
		Seller:      oim.Username,
		City:        oim.City,
		Price:       oim.Price,
		Qty:         oim.Quantity,
		Subtotal:    oim.Subtotal,
	}
}

func ToListCoreItem(items []OrderItemModel) []order.OrderItem {
	cores := []order.OrderItem{}
	for _, v := range items {
		cores = append(cores, ToCoreItem(v))
	}
	return cores
}

func ToCoreOrder(om OrderModel) order.Core {
	co := order.Core{
		ID:           om.ID,
		Invoice:      om.Invoice,
		Customer:     om.Fullname,
		Address:      om.Address,
		Phone:        om.Phone,
		OrderStatus:  om.OrderStatus,
		OrderDate:    om.OrderDate,
		PaidDate:     om.PaidDate,
		Total:        om.Total,
		PaymentToken: om.PaymentToken,
		PaymentURL:   om.PaymentURL,
		Items:        ToListCoreItem(om.Items),
	}
	return co
}

func ToListCoreOrder(models []OrderModel) []order.Core {
	cores := []order.Core{}
	for _, v := range models {
		cores = append(cores, ToCoreOrder(v))
	}
	return cores
}
