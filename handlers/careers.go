package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Career represents a career listing
type Career struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	Department   string `json:"department,omitempty"`
	Location     string `json:"location,omitempty"`
	Type         string `json:"type,omitempty"`
	Requirements string `json:"requirements,omitempty"`
	Active       bool   `json:"active"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// GetCareers returns all career listings
func GetCareers(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, description, department, location, type, active, created_at, updated_at FROM careers ORDER BY created_at DESC")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch careers"})
		return
	}
	defer rows.Close()

	var careers []Career
	for rows.Next() {
		var career Career
		var createdAt, updatedAt time.Time
		err := rows.Scan(&career.ID, &career.Title, &career.Description, &career.Department, &career.Location, &career.Type, &career.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		career.CreatedAt = createdAt.Format(time.RFC3339)
		career.UpdatedAt = updatedAt.Format(time.RFC3339)
		careers = append(careers, career)
	}

	c.JSON(200, careers)
}

// GetCareerByID returns a specific career listing
func GetCareerByID(c *gin.Context) {
	id := c.Param("id")

	var career Career
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, title, description, department, location, type, active, created_at, updated_at FROM careers WHERE id = $1",
		id,
	).Scan(&career.ID, &career.Title, &career.Description, &career.Department, &career.Location, &career.Type, &career.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Career listing not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch career"})
		return
	}

	career.CreatedAt = createdAt.Format(time.RFC3339)
	career.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, career)
}

// CreateCareer creates a new career listing
func CreateCareer(c *gin.Context) {
	var career Career
	if err := c.ShouldBindJSON(&career); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	career.ID = uuid.New().String()
	career.Active = true

	_, err := database.DB.Exec(
		"INSERT INTO careers (id, title, description, department, location, type, active) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		career.ID, career.Title, career.Description, career.Department, career.Location, career.Type, career.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create career: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM careers WHERE id = $1",
		career.ID,
	).Scan(&createdAt, &updatedAt)

	career.CreatedAt = createdAt.Format(time.RFC3339)
	career.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, career)
}

// UpdateCareer updates an existing career listing
func UpdateCareer(c *gin.Context) {
	id := c.Param("id")

	var career Career
	if err := c.ShouldBindJSON(&career); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	career.ID = id

	_, err := database.DB.Exec(
		"UPDATE careers SET title = $1, description = $2, department = $3, location = $4, type = $5, active = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7",
		career.Title, career.Description, career.Department, career.Location, career.Type, career.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update career: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	err = database.DB.QueryRow(
		"SELECT created_at, updated_at FROM careers WHERE id = $1",
		id,
	).Scan(&createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Career listing not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated career"})
		return
	}

	career.CreatedAt = createdAt.Format(time.RFC3339)
	career.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, career)
}

// DeleteCareer deletes a career listing
func DeleteCareer(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM careers WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete career: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Career listing not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Career listing deleted successfully"})
}
