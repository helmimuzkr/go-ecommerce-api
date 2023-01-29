package order

import "github.com/labstack/echo/v4"

type Core struct {
	ID           uint
	CustomerID   uint
	CustomerName string
	Invoice      string
	Address      string
	Phone        string
	OrderStatus  string
	OrderDate    string
	PaidDate     string
	Total        int
	PaymentToken string
	PaymentURL   string
	Items        []OrderItem
}

type OrderItem struct {
	ID          uint
	ProductName string
	Seller      string
	City        string
	Image       string
	Price       int
	Qty         int
	Subtotal    int
}

type OrderHandler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetOrderBuy() echo.HandlerFunc
	GetOrderSell() echo.HandlerFunc
	Cancel() echo.HandlerFunc
	Confirm() echo.HandlerFunc
	Callback() echo.HandlerFunc
}

type OrderService interface {
	Create(token interface{}, carts []int) (Core, error)
	GetAll(token interface{}, history string) ([]Core, error)
	GetOrderBuy(token interface{}, orderID uint) (Core, error)
	GetOrderSell(token interface{}, orderID uint) (Core, error)
	Cancel(token interface{}, orderID uint) error
	Confirm(token interface{}, orderID uint) error
	UpdateStatus(invoice string, status string, paidAt string) error
}

type OrderData interface {
	CreateOrder(userID uint, order Core) (uint, error)
	CreateOrderItem(userID uint, orderID uint, cart uint) error
	GetOrderBuy(userID uint, orderID uint) (Core, error)
	GetOrderSell(userID uint, orderID uint) (Core, error)
	ListOrderBuy(userID uint) ([]Core, error)
	ListOrderSell(userID uint) ([]Core, error)
	Confirm(orderID uint, updateOrder Core) error
	UpdateStatus(invoice string, updateOrder Core) error
}
