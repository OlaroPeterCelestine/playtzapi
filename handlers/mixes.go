package handlers

import (
	"database/sql"
	"playtz-api/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Mix represents a music mix
type Mix struct {
	ID          string  `json:"id"`
	RoomID      string  `json:"room_id"`
	Title       string  `json:"title"`
	Artist      string  `json:"artist"`
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	Tracks      int     `json:"tracks"`
	Color       string  `json:"color"`
	TextColor   string  `json:"text_color"`
	BorderColor string  `json:"border_color"`
	Image       string  `json:"image,omitempty"`
	AudioURL    string  `json:"audio_url,omitempty"`
	TrackList   []Track `json:"track_list,omitempty"`
	Active      bool    `json:"active"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

// Track represents a track in a mix
type Track struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
	Link     string `json:"link,omitempty"` // Streaming URL (audio or video)
	Type     string `json:"type,omitempty"` // "audio" or "video"
}

// GetMixes returns all mixes, optionally filtered by room
func GetMixes(c *gin.Context) {
	roomID := c.Query("room_id")

	var rows *sql.Rows
	var err error

	if roomID != "" {
		rows, err = database.DB.Query(
			"SELECT id, room_id, title, artist, description, duration, tracks, color, text_color, border_color, image, audio_url, active, created_at, updated_at FROM mixes WHERE room_id = $1 ORDER BY created_at DESC",
			roomID,
		)
	} else {
		rows, err = database.DB.Query(
			"SELECT id, room_id, title, artist, description, duration, tracks, color, text_color, border_color, image, audio_url, active, created_at, updated_at FROM mixes ORDER BY created_at DESC",
		)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch mixes"})
		return
	}
	defer rows.Close()

	var mixes []Mix
	for rows.Next() {
		var mix Mix
		var createdAt, updatedAt time.Time
		err := rows.Scan(&mix.ID, &mix.RoomID, &mix.Title, &mix.Artist, &mix.Description, &mix.Duration, &mix.Tracks, &mix.Color, &mix.TextColor, &mix.BorderColor, &mix.Image, &mix.AudioURL, &mix.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		mix.CreatedAt = createdAt.Format(time.RFC3339)
		mix.UpdatedAt = updatedAt.Format(time.RFC3339)

		// Get track count from tracks table
		var trackCount int
		database.DB.QueryRow("SELECT COUNT(*) FROM tracks WHERE mix_id = $1", mix.ID).Scan(&trackCount)
		mix.Tracks = trackCount

		mixes = append(mixes, mix)
	}

	c.JSON(200, mixes)
}

// GetMixByID returns a specific mix
func GetMixByID(c *gin.Context) {
	id := c.Param("id")

	var mix Mix
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, room_id, title, artist, description, duration, tracks, color, text_color, border_color, image, audio_url, active, created_at, updated_at FROM mixes WHERE id = $1",
		id,
	).Scan(&mix.ID, &mix.RoomID, &mix.Title, &mix.Artist, &mix.Description, &mix.Duration, &mix.Tracks, &mix.Color, &mix.TextColor, &mix.BorderColor, &mix.Image, &mix.AudioURL, &mix.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Mix not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch mix"})
		return
	}

	mix.CreatedAt = createdAt.Format(time.RFC3339)
	mix.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Get tracks for this mix
	rows, err := database.DB.Query(
		"SELECT number, title, artist, duration, link, type FROM tracks WHERE mix_id = $1 ORDER BY number",
		id,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var track Track
			rows.Scan(&track.Number, &track.Title, &track.Artist, &track.Duration, &track.Link, &track.Type)
			mix.TrackList = append(mix.TrackList, track)
		}
		mix.Tracks = len(mix.TrackList)
	}

	c.JSON(200, mix)
}

// CreateMix creates a new mix
func CreateMix(c *gin.Context) {
	var mix Mix
	if err := c.ShouldBindJSON(&mix); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	mix.ID = uuid.New().String()
	mix.Active = true
	mix.Tracks = 0

	_, err := database.DB.Exec(
		"INSERT INTO mixes (id, room_id, title, artist, description, duration, tracks, color, text_color, border_color, image, audio_url, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		mix.ID, mix.RoomID, mix.Title, mix.Artist, mix.Description, mix.Duration, mix.Tracks, mix.Color, mix.TextColor, mix.BorderColor, mix.Image, mix.AudioURL, mix.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create mix: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM mixes WHERE id = $1",
		mix.ID,
	).Scan(&createdAt, &updatedAt)

	mix.CreatedAt = createdAt.Format(time.RFC3339)
	mix.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, mix)
}

// UpdateMix updates an existing mix
func UpdateMix(c *gin.Context) {
	id := c.Param("id")

	var mix Mix
	if err := c.ShouldBindJSON(&mix); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	mix.ID = id

	_, err := database.DB.Exec(
		"UPDATE mixes SET room_id = $1, title = $2, artist = $3, description = $4, duration = $5, color = $6, text_color = $7, border_color = $8, image = $9, audio_url = $10, active = $11, updated_at = CURRENT_TIMESTAMP WHERE id = $12",
		mix.RoomID, mix.Title, mix.Artist, mix.Description, mix.Duration, mix.Color, mix.TextColor, mix.BorderColor, mix.Image, mix.AudioURL, mix.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update mix: " + err.Error()})
		return
	}

	// Update tracks count
	var trackCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM tracks WHERE mix_id = $1", id).Scan(&trackCount)
	database.DB.Exec("UPDATE mixes SET tracks = $1 WHERE id = $2", trackCount, id)

	var createdAt, updatedAt time.Time
	err = database.DB.QueryRow(
		"SELECT created_at, updated_at FROM mixes WHERE id = $1",
		id,
	).Scan(&createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Mix not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated mix"})
		return
	}

	mix.Tracks = trackCount
	mix.CreatedAt = createdAt.Format(time.RFC3339)
	mix.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, mix)
}

// DeleteMix deletes a mix
func DeleteMix(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM mixes WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete mix: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Mix not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Mix deleted successfully"})
}

// AddTrackToMix adds a track to a mix
func AddTrackToMix(c *gin.Context) {
	mixID := c.Param("id")

	var track Track
	if err := c.ShouldBindJSON(&track); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Get next track number
	var maxNumber int
	database.DB.QueryRow("SELECT COALESCE(MAX(number), 0) FROM tracks WHERE mix_id = $1", mixID).Scan(&maxNumber)
	track.Number = maxNumber + 1

	// Determine type from link if not provided
	if track.Type == "" {
		if len(track.Link) > 0 {
			if contains([]string{".mp4", ".webm", ".m3u8"}, track.Link) {
				track.Type = "video"
			} else {
				track.Type = "audio"
			}
		} else {
			track.Type = "audio"
		}
	}

	trackID := uuid.New().String()
	_, err := database.DB.Exec(
		"INSERT INTO tracks (id, mix_id, number, title, artist, duration, link, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		trackID, mixID, track.Number, track.Title, track.Artist, track.Duration, track.Link, track.Type,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add track: " + err.Error()})
		return
	}

	// Update mix track count
	var trackCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM tracks WHERE mix_id = $1", mixID).Scan(&trackCount)
	database.DB.Exec("UPDATE mixes SET tracks = $1 WHERE id = $2", trackCount, mixID)

	c.JSON(200, gin.H{
		"message": "Track added to mix",
		"track":   track,
	})
}

// AddTracksToMix adds multiple tracks to a mix via links
func AddTracksToMix(c *gin.Context) {
	mixID := c.Param("id")

	var request struct {
		Links []string `json:"links"` // Array of track URLs/links
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if len(request.Links) == 0 {
		c.JSON(400, gin.H{"error": "No track links provided"})
		return
	}

	// Get starting track number
	var maxNumber int
	database.DB.QueryRow("SELECT COALESCE(MAX(number), 0) FROM tracks WHERE mix_id = $1", mixID).Scan(&maxNumber)

	tracks := make([]Track, 0, len(request.Links))
	for i, link := range request.Links {
		if link == "" {
			continue
		}

		trackType := "audio"
		if contains([]string{".mp4", ".webm", ".m3u8"}, link) {
			trackType = "video"
		}

		trackID := uuid.New().String()
		trackNumber := maxNumber + i + 1

		_, err := database.DB.Exec(
			"INSERT INTO tracks (id, mix_id, number, link, type) VALUES ($1, $2, $3, $4, $5)",
			trackID, mixID, trackNumber, link, trackType,
		)

		if err == nil {
			tracks = append(tracks, Track{
				Number: trackNumber,
				Link:   link,
				Type:   trackType,
			})
		}
	}

	// Update mix track count
	var trackCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM tracks WHERE mix_id = $1", mixID).Scan(&trackCount)
	database.DB.Exec("UPDATE mixes SET tracks = $1 WHERE id = $2", trackCount, mixID)

	c.JSON(200, gin.H{
		"message": "Tracks added to mix",
		"count":   len(tracks),
		"tracks":  tracks,
	})
}

// RemoveTrackFromMix removes a track from a mix
func RemoveTrackFromMix(c *gin.Context) {
	mixID := c.Param("id")

	trackNumber := c.Query("track_number")
	if trackNumber == "" {
		c.JSON(400, gin.H{"error": "Track number required"})
		return
	}

	num, err := strconv.Atoi(trackNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid track number"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM tracks WHERE mix_id = $1 AND number = $2", mixID, num)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to remove track: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Track not found"})
		return
	}

	// Update mix track count
	var trackCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM tracks WHERE mix_id = $1", mixID).Scan(&trackCount)
	database.DB.Exec("UPDATE mixes SET tracks = $1 WHERE id = $2", trackCount, mixID)

	c.JSON(200, gin.H{
		"message":      "Track removed from mix",
		"track_number": trackNumber,
	})
}

// Helper function to check if string contains any of the substrings
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if len(str) >= len(s) && str[len(str)-len(s):] == s {
			return true
		}
	}
	return false
}
