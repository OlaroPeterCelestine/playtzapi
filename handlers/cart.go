package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Name          string  `json:"name,omitempty"`
	Price         float64 `json:"price,omitempty"`
	Image         string  `json:"image,omitempty"`
	MerchandiseID string  `json:"merchandise_id,omitempty"`
}

// Cart represents a shopping cart
type Cart struct {
	ID       string     `json:"id"`
	UserID   string     `json:"user_id,omitempty"`
	Items    []CartItem `json:"items"`
	Subtotal float64    `json:"subtotal"`
	Total    float64    `json:"total"`
}

// GetCart returns the current cart
func GetCart(c *gin.Context) {
	cartID := c.Query("cart_id")
	userID := c.Query("user_id")

	if cartID == "" {
		// Return empty cart structure
		cart := Cart{
			ID:       uuid.New().String(),
			UserID:   userID,
			Items:    []CartItem{},
			Subtotal: 0,
			Total:    0,
		}
		c.JSON(200, cart)
		return
	}

	// Fetch cart from database
	var cart Cart
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, user_id, created_at, updated_at FROM cart WHERE id = $1",
		cartID,
	).Scan(&cart.ID, &cart.UserID, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		// Cart not found, return empty cart
		cart = Cart{
			ID:       cartID,
			UserID:   userID,
			Items:    []CartItem{},
			Subtotal: 0,
			Total:    0,
		}
		c.JSON(200, cart)
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch cart"})
		return
	}

	// Fetch cart items
	rows, err := database.DB.Query(`
		SELECT ci.id, ci.merchandise_id, ci.quantity, m.name, m.price, m.image
		FROM cart_items ci
		LEFT JOIN merchandise m ON ci.merchandise_id = m.id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at
	`, cartID)

	if err == nil {
		defer rows.Close()
		var subtotal float64
		for rows.Next() {
			var item CartItem
			var price float64
			rows.Scan(&item.ID, &item.MerchandiseID, &item.Quantity, &item.Name, &price, &item.Image)
			item.ProductID = item.MerchandiseID
			item.Price = price
			subtotal += price * float64(item.Quantity)
			cart.Items = append(cart.Items, item)
		}
		cart.Subtotal = subtotal
		cart.Total = subtotal // Can add shipping/tax later
	}

	c.JSON(200, cart)
}

// AddToCart adds an item to the cart
func AddToCart(c *gin.Context) {
	var req struct {
		CartID        string `json:"cart_id"`
		UserID        string `json:"user_id,omitempty"`
		MerchandiseID string `json:"merchandise_id"`
		Quantity      int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if req.MerchandiseID == "" {
		c.JSON(400, gin.H{"error": "Merchandise ID is required"})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// Get or create cart
	var cartID string
	if req.CartID == "" {
		cartID = uuid.New().String()
		_, err := database.DB.Exec(
			"INSERT INTO cart (id, user_id) VALUES ($1, $2)",
			cartID, req.UserID,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create cart: " + err.Error()})
			return
		}
	} else {
		cartID = req.CartID
		// Verify cart exists
		var exists bool
		database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM cart WHERE id = $1)", cartID).Scan(&exists)
		if !exists {
			_, err := database.DB.Exec(
				"INSERT INTO cart (id, user_id) VALUES ($1, $2)",
				cartID, req.UserID,
			)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to create cart: " + err.Error()})
				return
			}
		}
	}

	// Check if item already exists in cart
	var existingID string
	var existingQty int
	err := database.DB.QueryRow(
		"SELECT id, quantity FROM cart_items WHERE cart_id = $1 AND merchandise_id = $2",
		cartID, req.MerchandiseID,
	).Scan(&existingID, &existingQty)

	if err == nil {
		// Update quantity
		newQty := existingQty + req.Quantity
		_, err = database.DB.Exec(
			"UPDATE cart_items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
			newQty, existingID,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update cart item: " + err.Error()})
			return
		}
	} else {
		// Add new item
		itemID := uuid.New().String()
		_, err = database.DB.Exec(
			"INSERT INTO cart_items (id, cart_id, merchandise_id, quantity) VALUES ($1, $2, $3, $4)",
			itemID, cartID, req.MerchandiseID, req.Quantity,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to add item to cart: " + err.Error()})
			return
		}
	}

	c.JSON(201, gin.H{
		"message":       "Item added to cart",
		"cart_id":       cartID,
		"merchandise_id": req.MerchandiseID,
		"quantity":      req.Quantity,
	})
}

// UpdateCartItem updates a cart item
func UpdateCartItem(c *gin.Context) {
	var req struct {
		CartID   string `json:"cart_id"`
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(400, gin.H{"error": "Quantity must be greater than 0"})
		return
	}

	result, err := database.DB.Exec(
		"UPDATE cart_items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND cart_id = $3",
		req.Quantity, req.ItemID, req.CartID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update cart item: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Cart item updated",
		"item_id": req.ItemID,
	})
}

// RemoveFromCart removes an item from the cart
func RemoveFromCart(c *gin.Context) {
	cartID := c.Query("cart_id")
	itemID := c.Query("item_id")
	productID := c.Query("product_id")

	if cartID == "" {
		c.JSON(400, gin.H{"error": "Cart ID is required"})
		return
	}

	var err error
	if itemID != "" {
		_, err = database.DB.Exec("DELETE FROM cart_items WHERE id = $1 AND cart_id = $2", itemID, cartID)
	} else if productID != "" {
		_, err = database.DB.Exec("DELETE FROM cart_items WHERE cart_id = $1 AND merchandise_id = $2", cartID, productID)
	} else {
		c.JSON(400, gin.H{"error": "Item ID or Product ID is required"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to remove item: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Item removed from cart",
	})
}

// ClearCart clears all items from the cart
func ClearCart(c *gin.Context) {
	cartID := c.Query("cart_id")
	if cartID == "" {
		c.JSON(400, gin.H{"error": "Cart ID is required"})
		return
	}

	_, err := database.DB.Exec("DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear cart: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Cart cleared"})
}
