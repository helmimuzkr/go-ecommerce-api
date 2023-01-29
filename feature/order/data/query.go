package data

import (
	"e-commerce-api/feature/order"

	"gorm.io/gorm"
)

type orderData struct {
	db *gorm.DB
}

func New(db *gorm.DB) order.OrderData {
	return &orderData{db: db}
}

func (od *orderData) CreateOrder(userID uint, newOrder order.Core) (uint, error) {
	// Create new order
	model := ToModel(newOrder)
	model.CustomerID = userID
	tx := od.db.Create(&model)
	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}

	tx = tx.Exec("UPDATE carts SET deleted_at=CURRENT_TIMESTAMP WHERE user_id = ?", userID)
	if tx.Error != nil {
		tx.Rollback()
		return 0, tx.Error
	}

	return model.ID, nil
}

func (od *orderData) CreateOrderItem(orderID uint, cartID uint) error {
	tx := od.db.Begin()

	// Convert cart items to order items
	item := OrderItem{}
	tx = tx.Raw("SELECT product_id, quantity FROM carts WHERE deleted_at IS NULL AND id = ?", cartID).Find(&item)
	if tx.Error != nil {
		return tx.Error
	}

	// Get the price product for counting the subtotal
	var productPrice int
	tx = tx.Raw("SELECT price FROM products WHERE deleted_at IS NULL AND id = ? AND stock > 0", item.ProductID).First(&productPrice)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	item.OrderID = orderID
	item.Subtotal = item.Quantity * productPrice // Calculate subtotal

	// Insert item to table order_items
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()

	return nil
}

func (od *orderData) GetOrderSell(userID uint, orderID uint) (order.Core, error) {
	// Tampilkan data detail order dengan kondisi
	// order item = order id
	// seller produk = penjual
	o := OrderModel{}
	qOrder := "SELECT orders.id, orders.invoice, orders.customer_id, users.fullname, users.address, users.city, users.phone, orders.order_status, orders.order_date, orders.paid_date, orders.total_price, orders.payment_token, orders.payment_url FROM orders JOIN users ON users.id = orders.customer_id JOIN order_items ON order_items.order_id = orders.id JOIN products ON products.id = order_items.product_id WHERE orders.id = ? AND products.seller_id = ?"
	tx := od.db.Raw(qOrder, orderID, userID).Find(&o)
	if tx.Error != nil {
		return order.Core{}, tx.Error
	}

	// Tampilkan seluruh data item dengan kondisi
	// order item = order id
	// seller produk = penjual
	itemModels := []OrderItemModel{}
	qItems := "SELECT products.id, products.name, users.username, users.city, products.price, products.image, order_items.quantity, order_items.subtotal FROM order_items JOIN products ON products.id = order_items.product_id JOIN users ON users.id = products.seller_id WHERE order_items.deleted_at IS NULL AND order_items.order_id = ? AND products.seller_id = ?"
	tx = od.db.Raw(qItems, orderID, userID).Find(&itemModels)
	if tx.Error != nil {
		tx.Rollback()
		return order.Core{}, tx.Error
	}

	o.Items = itemModels

	return ToCoreOrder(o), nil
}

func (od *orderData) GetOrderBuy(userID uint, orderID uint) (order.Core, error) {
	o := OrderModel{}
	qOrder := "SELECT orders.id, orders.invoice,orders.customer_id, users.fullname, users.address, users.city, users.phone, orders.order_status, orders.order_date, orders.paid_date, orders.total_price, orders.payment_token, orders.payment_url FROM orders JOIN users ON users.id = orders.customer_id Where orders.id = ? AND orders.customer_id = ?"
	tx := od.db.Raw(qOrder, orderID, userID).Find(&o)
	if tx.Error != nil {
		return order.Core{}, tx.Error
	}

	itemModels := []OrderItemModel{}
	qItems := "SELECT products.id, products.name, users.username, users.city, products.price, products.image, order_items.quantity, order_items.subtotal FROM order_items JOIN products ON products.id = order_items.product_id JOIN users ON users.id = products.seller_id WHERE order_items.deleted_at IS NULL AND order_items.order_id = ?"
	tx = od.db.Raw(qItems, orderID).Find(&itemModels)
	if tx.Error != nil {
		tx.Rollback()
		return order.Core{}, tx.Error
	}

	o.Items = itemModels

	return ToCoreOrder(o), nil
}

func (od *orderData) ListOrderBuy(userID uint) ([]order.Core, error) {
	o := []OrderModel{}
	query := "SELECT id, invoice, order_status, order_date FROM orders WHERE customer_id = ?"
	tx := od.db.Raw(query, userID).Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ToListCoreOrder(o), nil
}

func (od *orderData) ListOrderSell(userID uint) ([]order.Core, error) {
	o := []OrderModel{}
	query := "SELECT orders.id, orders.invoice, orders.order_status, orders.order_date FROM orders JOIN order_items ON order_items.order_id = orders.id JOIN products ON products.id = order_items.product_id WHERE products.seller_id = ?"
	tx := od.db.Raw(query, userID).Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ToListCoreOrder(o), nil
}

func (od *orderData) Confirm(orderID uint, updateOrder order.Core) error {
	data := ToModel(updateOrder)
	tx := od.db.Where("id = ?", orderID).Updates(&data)
	if tx.Error != nil {
		return tx.Error
	}

	om := []OrderItem{}
	tx = od.db.Where("order_id = ?", orderID).Find(&om)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	for _, v := range om {
		tx = tx.Exec("UPDATE products SET stock=stock-? WHERE id=?", v.Quantity, v.ProductID)
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	return nil
}

func (od *orderData) UpdateStatus(invoice string, updateOrder order.Core) error {
	data := ToModel(updateOrder)
	tx := od.db.Where("invoice = ?", invoice).Updates(&data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
