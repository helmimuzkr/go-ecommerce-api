package data

import (
	"e-commerce-api/feature/order"
	"fmt"

	"gorm.io/gorm"
)

type orderData struct {
	db *gorm.DB
}

func New(db *gorm.DB) order.OrderData {
	return &orderData{db: db}
}

func (od *orderData) CreateOrder(userID uint, newOrder order.Core, carts []order.Cart) (uint, error) {
	// Start transactions
	tx := od.db.Begin()

	// Create new order
	model := ToModel(newOrder)
	model.CustomerID = userID
	tx = tx.Create(&model)
	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}

	// Insert items to order
	tx = od.insertItem(tx, model.ID, carts)
	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}

	for _, cart := range carts {
		fmt.Println(cart)
	}
	// Delete cart after make an order
	// tx = od.deleteCart(tx, userID)
	// if tx.Error != nil {
	// 	tx.Rollback()
	// 	return 0, tx.Error
	// }

	tx.Commit()

	return model.ID, nil
}

func (od *orderData) insertItem(tx *gorm.DB, orderID uint, carts []order.Cart) *gorm.DB {
	for _, cart := range carts {
		oi := OrderItem{OrderID: orderID, ProductID: cart.ProductID, Quantity: cart.Quantity, Subtotal: cart.Subtotal}
		tx = tx.Exec("INSERT INTO order_items(order_id, product_id, quantity, subtotal) VALUES(?,?,?,?)", oi.OrderID, oi.ProductID, oi.Quantity, oi.Subtotal)
	}
	return tx
}

func (od *orderData) deleteCart(tx *gorm.DB, userID uint) *gorm.DB {
	tx = tx.Exec("UPDATE carts SET deleted_at=CURRENT_TIMESTAMP WHERE user_id=?", userID)
	return tx
}

func (od *orderData) GetItemById(userID uint, orderID uint) ([]order.OrderItem, error) {
	itemModels := []OrderItemModel{}
	query := "SELECT products.id, products.name, users.username, users.city, products.price, order_items.quantity, order_items.subtotal FROM order_items JOIN products ON products.id = order_items.product_id JOIN users ON users.id = products.seller_id WHERE order_items.deleted_at IS NULL AND order_items.order_id = ? AND users.id = ?"
	tx := od.db.Raw(query, userID, orderID).Find(&itemModels)
	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	return ToListCoreItem(itemModels), nil
}

func (od *orderData) GetAll(userID uint) ([]order.Core, error) {
	return nil, nil
}
func (od *orderData) GetByID(userID uint, orderID uint) (order.Core, error) {
	return order.Core{}, nil
}
func (od *orderData) Cancel(userID uint, orderID uint) error {
	return nil
}
