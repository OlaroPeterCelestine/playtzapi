package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Event represents an event
type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Date        string `json:"date"`
	Time        string `json:"time,omitempty"`
	Location    string `json:"location,omitempty"`
	Image       string `json:"image,omitempty"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// GetEvents returns all events
func GetEvents(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, description, date, time, location, image, active, created_at, updated_at FROM events ORDER BY date DESC, created_at DESC")
	if err != nil {
		c.JSON(500, []Event{})
		return
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var createdAt, updatedAt time.Time
		var date time.Time
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &date, &event.Time, &event.Location, &event.Image, &event.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		event.Date = date.Format("2006-01-02")
		event.CreatedAt = createdAt.Format(time.RFC3339)
		event.UpdatedAt = updatedAt.Format(time.RFC3339)
		events = append(events, event)
	}

	c.JSON(200, events)
}

// GetEventByID returns a specific event
func GetEventByID(c *gin.Context) {
	id := c.Param("id")

	var event Event
	var createdAt, updatedAt time.Time
	var date time.Time
	err := database.DB.QueryRow(
		"SELECT id, title, description, date, time, location, image, active, created_at, updated_at FROM events WHERE id = $1",
		id,
	).Scan(&event.ID, &event.Title, &event.Description, &date, &event.Time, &event.Location, &event.Image, &event.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch event"})
		return
	}

	event.Date = date.Format("2006-01-02")
	event.CreatedAt = createdAt.Format(time.RFC3339)
	event.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, event)
}

// CreateEvent creates a new event
func CreateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	event.ID = uuid.New().String()
	event.Active = true

	_, err := database.DB.Exec(
		"INSERT INTO events (id, title, description, date, time, location, image, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		event.ID, event.Title, event.Description, event.Date, event.Time, event.Location, event.Image, event.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create event: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var date time.Time
	database.DB.QueryRow(
		"SELECT date, created_at, updated_at FROM events WHERE id = $1",
		event.ID,
	).Scan(&date, &createdAt, &updatedAt)

	event.Date = date.Format("2006-01-02")
	event.CreatedAt = createdAt.Format(time.RFC3339)
	event.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, event)
}

// UpdateEvent updates an existing event
func UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	event.ID = id

	_, err := database.DB.Exec(
		"UPDATE events SET title = $1, description = $2, date = $3, time = $4, location = $5, image = $6, active = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $8",
		event.Title, event.Description, event.Date, event.Time, event.Location, event.Image, event.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update event: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var date time.Time
	err = database.DB.QueryRow(
		"SELECT date, created_at, updated_at FROM events WHERE id = $1",
		id,
	).Scan(&date, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated event"})
		return
	}

	event.Date = date.Format("2006-01-02")
	event.CreatedAt = createdAt.Format(time.RFC3339)
	event.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, event)
}

// DeleteEvent deletes an event
func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete event: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Event deleted successfully"})
}
