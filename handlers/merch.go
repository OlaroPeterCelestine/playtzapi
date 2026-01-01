package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Merchandise represents a merchandise item
type Merchandise struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Image       string  `json:"image,omitempty"`
	Stock       int     `json:"stock"`
	Active      bool    `json:"active"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

// GetMerch returns all merchandise items
func GetMerch(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description, price, image, stock, active, created_at, updated_at FROM merchandise ORDER BY created_at DESC")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch merchandise"})
		return
	}
	defer rows.Close()

	var items []Merchandise
	for rows.Next() {
		var item Merchandise
		var createdAt, updatedAt time.Time
		var price float64
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &price, &item.Image, &item.Stock, &item.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		item.Price = price
		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		items = append(items, item)
	}

	c.JSON(200, items)
}

// GetMerchByID returns a specific merchandise item
func GetMerchByID(c *gin.Context) {
	id := c.Param("id")

	var item Merchandise
	var createdAt, updatedAt time.Time
	var price float64
	err := database.DB.QueryRow(
		"SELECT id, name, description, price, image, stock, active, created_at, updated_at FROM merchandise WHERE id = $1",
		id,
	).Scan(&item.ID, &item.Name, &item.Description, &price, &item.Image, &item.Stock, &item.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Merchandise item not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch merchandise"})
		return
	}

	item.Price = price
	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, item)
}

// CreateMerch creates a new merchandise item
func CreateMerch(c *gin.Context) {
	var item Merchandise
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	item.ID = uuid.New().String()
	item.Active = true

	_, err := database.DB.Exec(
		"INSERT INTO merchandise (id, name, description, price, image, stock, active) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		item.ID, item.Name, item.Description, item.Price, item.Image, item.Stock, item.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create merchandise: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM merchandise WHERE id = $1",
		item.ID,
	).Scan(&createdAt, &updatedAt)

	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, item)
}

// UpdateMerch updates an existing merchandise item
func UpdateMerch(c *gin.Context) {
	id := c.Param("id")

	var item Merchandise
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	item.ID = id

	_, err := database.DB.Exec(
		"UPDATE merchandise SET name = $1, description = $2, price = $3, image = $4, stock = $5, active = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7",
		item.Name, item.Description, item.Price, item.Image, item.Stock, item.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update merchandise: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	err = database.DB.QueryRow(
		"SELECT created_at, updated_at FROM merchandise WHERE id = $1",
		id,
	).Scan(&createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Merchandise item not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated item"})
		return
	}

	item.CreatedAt = createdAt.Format(time.RFC3339)
	item.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, item)
}

// DeleteMerch deletes a merchandise item
func DeleteMerch(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM merchandise WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete merchandise: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Merchandise item not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Merchandise item deleted successfully"})
}
