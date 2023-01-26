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

type Cart struct {
	ID        uint
	OrderID   uint
	ProductID uint
	Quantity  int
	Subtotal  int
}

type OrderHandler interface {
	Create() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetOrderBuy() echo.HandlerFunc
	GetOrderSell() echo.HandlerFunc
	Cancel() echo.HandlerFunc
}

type OrderService interface {
	Create(token interface{}, carts []Cart) (Core, error)
	GetAll(token interface{}, history string) ([]Core, error)
	GetOrderBuy(token interface{}, orderID uint) (Core, error)
	GetOrderSell(token interface{}, orderID uint) (Core, error)
	Cancel(token interface{}, orderID uint) error
}

type OrderData interface {
	CreateOrder(userID uint, order Core, carts []Cart) (uint, error)
	GetItemBuy(userID uint, orderID uint) (Core, error)
	GetItemSell(userID uint, orderID uint) (Core, error)
	GetListOrderBuy(userID uint) ([]Core, error)
	GetListOrderSell(userID uint) ([]Core, error)
	GetByID(userID uint, orderID uint) (Core, error)
	Cancel(userID uint, orderID uint) error
}
