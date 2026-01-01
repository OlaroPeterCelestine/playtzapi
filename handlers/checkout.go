package handlers

import (
	"database/sql"
	"encoding/json"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Order represents an order
type Order struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id,omitempty"`
	Items           []OrderItem     `json:"items"`
	Subtotal        float64         `json:"subtotal"`
	Shipping        float64         `json:"shipping"`
	Tax             float64         `json:"tax"`
	Total           float64         `json:"total"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"product_id"`
	MerchandiseID string  `json:"merchandise_id,omitempty"`
	Name          string  `json:"name"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
}

// ShippingAddress represents shipping information
type ShippingAddress struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CheckoutRequest represents the checkout request
type CheckoutRequest struct {
	CartID          string          `json:"cart_id"`
	UserID          string          `json:"user_id,omitempty"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
	PaymentMethod   string          `json:"payment_method"`
}

// CreateOrder creates a new order from checkout
func CreateOrder(c *gin.Context) {
	var checkout CheckoutRequest
	if err := c.ShouldBindJSON(&checkout); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if checkout.CartID == "" {
		c.JSON(400, gin.H{"error": "Cart ID is required"})
		return
	}

	// Get cart items
	rows, err := database.DB.Query(`
		SELECT ci.id, ci.merchandise_id, ci.quantity, m.name, m.price
		FROM cart_items ci
		LEFT JOIN merchandise m ON ci.merchandise_id = m.id
		WHERE ci.cart_id = $1
	`, checkout.CartID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cart items: " + err.Error()})
		return
	}
	defer rows.Close()

	var items []OrderItem
	var subtotal float64

	for rows.Next() {
		var item OrderItem
		var price float64
		err := rows.Scan(&item.ID, &item.MerchandiseID, &item.Quantity, &item.Name, &price)
		if err != nil {
			continue
		}
		item.ProductID = item.MerchandiseID
		item.Price = price
		subtotal += price * float64(item.Quantity)
		items = append(items, item)
	}

	if len(items) == 0 {
		c.JSON(400, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate totals
	shipping := 10.0 // Fixed shipping for now
	tax := subtotal * 0.1 // 10% tax
	total := subtotal + shipping + tax

	// Create order
	orderID := uuid.New().String()
	shippingAddrJSON, _ := json.Marshal(checkout.ShippingAddress)

	_, err = database.DB.Exec(
		"INSERT INTO orders (id, user_id, total, status, shipping_address, payment_method) VALUES ($1, $2, $3, $4, $5, $6)",
		orderID, checkout.UserID, total, "pending", string(shippingAddrJSON), checkout.PaymentMethod,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	// Create order items
	for _, item := range items {
		orderItemID := uuid.New().String()
		_, err = database.DB.Exec(
			"INSERT INTO order_items (id, order_id, merchandise_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			orderItemID, orderID, item.MerchandiseID, item.Quantity, item.Price,
		)
		if err != nil {
			// Log error but continue
			continue
		}
	}

	// Clear cart
	database.DB.Exec("DELETE FROM cart_items WHERE cart_id = $1", checkout.CartID)

	// Fetch created order
	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM orders WHERE id = $1",
		orderID,
	).Scan(&createdAt, &updatedAt)

	order := Order{
		ID:              orderID,
		UserID:          checkout.UserID,
		Items:           items,
		Subtotal:        subtotal,
		Shipping:        shipping,
		Tax:             tax,
		Total:           total,
		Status:          "pending",
		PaymentMethod:   checkout.PaymentMethod,
		ShippingAddress: checkout.ShippingAddress,
		CreatedAt:       createdAt.Format(time.RFC3339),
		UpdatedAt:       updatedAt.Format(time.RFC3339),
	}

	c.JSON(201, order)
}

// GetOrder returns a specific order
func GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order Order
	var shippingAddrJSON string
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, user_id, total, status, shipping_address, payment_method, created_at, updated_at FROM orders WHERE id = $1",
		id,
	).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &shippingAddrJSON, &order.PaymentMethod, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch order"})
		return
	}

	// Parse shipping address
	json.Unmarshal([]byte(shippingAddrJSON), &order.ShippingAddress)

	// Calculate subtotal, shipping, tax from total
	order.Subtotal = order.Total / 1.1 // Approximate (assuming 10% tax)
	order.Tax = order.Total - order.Subtotal - 10.0 // Assuming $10 shipping
	order.Shipping = 10.0
	order.CreatedAt = createdAt.Format(time.RFC3339)
	order.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Get order items
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.merchandise_id, oi.quantity, oi.price, m.name
		FROM order_items oi
		LEFT JOIN merchandise m ON oi.merchandise_id = m.id
		WHERE oi.order_id = $1
	`, id)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item OrderItem
			var price float64
			rows.Scan(&item.ID, &item.MerchandiseID, &item.Quantity, &price, &item.Name)
			item.ProductID = item.MerchandiseID
			item.Price = price
			order.Items = append(order.Items, item)
		}
	}

	c.JSON(200, order)
}

// GetOrders returns all orders (optionally filtered by user)
func GetOrders(c *gin.Context) {
	userID := c.Query("user_id")

	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = database.DB.Query(
			"SELECT id, user_id, total, status, shipping_address, payment_method, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC",
			userID,
		)
	} else {
		rows, err = database.DB.Query(
			"SELECT id, user_id, total, status, shipping_address, payment_method, created_at, updated_at FROM orders ORDER BY created_at DESC",
		)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var shippingAddrJSON string
		var createdAt, updatedAt time.Time
		err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &shippingAddrJSON, &order.PaymentMethod, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(shippingAddrJSON), &order.ShippingAddress)
		order.Subtotal = order.Total / 1.1
		order.Tax = order.Total - order.Subtotal - 10.0
		order.Shipping = 10.0
		order.CreatedAt = createdAt.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}

	c.JSON(200, orders)
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var update struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "processing": true, "shipped": true, "delivered": true, "cancelled": true,
	}
	if !validStatuses[update.Status] {
		c.JSON(400, gin.H{"error": "Invalid status"})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		update.Status, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update order: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Order status updated",
		"status":  update.Status,
	})
}
