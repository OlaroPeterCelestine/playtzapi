package handlers

import (
	"database/sql"
	"playtz-api/database"
	"playtz-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRooms returns all rooms
func GetRooms(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, genre, description, gradient, text_color, image, active, created_at, updated_at FROM rooms ORDER BY name")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch rooms"})
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var createdAt, updatedAt time.Time
		err := rows.Scan(&room.ID, &room.Name, &room.Genre, &room.Description, &room.Gradient, &room.TextColor, &room.Image, &room.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		room.CreatedAt = createdAt.Format(time.RFC3339)
		room.UpdatedAt = updatedAt.Format(time.RFC3339)
		rooms = append(rooms, room)
	}

	c.JSON(200, rooms)
}

// GetRoomByID returns a specific room
func GetRoomByID(c *gin.Context) {
	id := c.Param("id")

	var room models.Room
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, name, genre, description, gradient, text_color, image, active, created_at, updated_at FROM rooms WHERE id = $1",
		id,
	).Scan(&room.ID, &room.Name, &room.Genre, &room.Description, &room.Gradient, &room.TextColor, &room.Image, &room.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Room not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch room"})
		return
	}

	room.CreatedAt = createdAt.Format(time.RFC3339)
	room.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, room)
}

// CreateRoom creates a new room
func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	room.ID = uuid.New().String()
	room.Active = true

	_, err := database.DB.Exec(
		"INSERT INTO rooms (id, name, genre, description, gradient, text_color, image, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		room.ID, room.Name, room.Genre, room.Description, room.Gradient, room.TextColor, room.Image, room.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create room: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM rooms WHERE id = $1",
		room.ID,
	).Scan(&createdAt, &updatedAt)

	room.CreatedAt = createdAt.Format(time.RFC3339)
	room.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, room)
}

// UpdateRoom updates an existing room
func UpdateRoom(c *gin.Context) {
	id := c.Param("id")

	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	room.ID = id

	_, err := database.DB.Exec(
		"UPDATE rooms SET name = $1, genre = $2, description = $3, gradient = $4, text_color = $5, image = $6, active = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8",
		room.Name, room.Genre, room.Description, room.Gradient, room.TextColor, room.Image, room.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update room: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	err = database.DB.QueryRow(
		"SELECT created_at, updated_at FROM rooms WHERE id = $1",
		id,
	).Scan(&createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Room not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated room"})
		return
	}

	room.CreatedAt = createdAt.Format(time.RFC3339)
	room.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, room)
}

// DeleteRoom deletes a room
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM rooms WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete room: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Room deleted successfully"})
}
