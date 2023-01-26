package order

import "github.com/labstack/echo/v4"

type Core struct {
	ID           uint
	Invoice      string
	Customer     string
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
	GetByID() echo.HandlerFunc
	Cancel() echo.HandlerFunc
}

type OrderService interface {
	Create(token interface{}, carts []Cart) (Core, error)
	GetAll(token interface{}) ([]Core, error)
	GetByID(token interface{}, orderID uint) (Core, error)
	Cancel(token interface{}, orderID uint) error
}

type OrderData interface {
	CreateOrder(userID uint, order Core, carts []Cart) (uint, error)
	GetItemById(userID uint, orderID uint) ([]OrderItem, error)
	GetAll(userID uint) ([]Core, error)
	GetByID(userID uint, orderID uint) (Core, error)
	Cancel(userID uint, orderID uint) error
}
