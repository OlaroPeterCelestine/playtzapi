package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NewsArticle represents a news article
type NewsArticle struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Excerpt   string `json:"excerpt,omitempty"`
	Author    string `json:"author,omitempty"`
	Image     string `json:"image,omitempty"`
	Published bool   `json:"published"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// GetNews returns all news articles
func GetNews(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, content, author, image, published, created_at, updated_at FROM news ORDER BY created_at DESC")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch news"})
		return
	}
	defer rows.Close()

	var articles []NewsArticle
	for rows.Next() {
		var article NewsArticle
		var createdAt, updatedAt time.Time
		err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Author, &article.Image, &article.Published, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		article.CreatedAt = createdAt.Format(time.RFC3339)
		article.UpdatedAt = updatedAt.Format(time.RFC3339)
		if len(article.Content) > 200 {
			article.Excerpt = article.Content[:200] + "..."
		} else {
			article.Excerpt = article.Content
		}
		articles = append(articles, article)
	}

	c.JSON(200, articles)
}

// GetNewsByID returns a specific news article
func GetNewsByID(c *gin.Context) {
	id := c.Param("id")

	var article NewsArticle
	var createdAt, updatedAt time.Time
	err := database.DB.QueryRow(
		"SELECT id, title, content, author, image, published, created_at, updated_at FROM news WHERE id = $1",
		id,
	).Scan(&article.ID, &article.Title, &article.Content, &article.Author, &article.Image, &article.Published, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "News article not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch news article"})
		return
	}

	article.CreatedAt = createdAt.Format(time.RFC3339)
	article.UpdatedAt = updatedAt.Format(time.RFC3339)
	if len(article.Content) > 200 {
		article.Excerpt = article.Content[:200] + "..."
	} else {
		article.Excerpt = article.Content
	}

	c.JSON(200, article)
}

// CreateNews creates a new news article
func CreateNews(c *gin.Context) {
	var article NewsArticle
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate ID
	article.ID = uuid.New().String()
	article.Published = false

	_, err := database.DB.Exec(
		"INSERT INTO news (id, title, content, author, image, published) VALUES ($1, $2, $3, $4, $5, $6)",
		article.ID, article.Title, article.Content, article.Author, article.Image, article.Published,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create news: " + err.Error()})
		return
	}

	// Fetch created article
	var createdAt, updatedAt time.Time
	database.DB.QueryRow(
		"SELECT created_at, updated_at FROM news WHERE id = $1",
		article.ID,
	).Scan(&createdAt, &updatedAt)

	article.CreatedAt = createdAt.Format(time.RFC3339)
	article.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, article)
}

// UpdateNews updates an existing news article
func UpdateNews(c *gin.Context) {
	id := c.Param("id")

	var article NewsArticle
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	article.ID = id

	_, err := database.DB.Exec(
		"UPDATE news SET title = $1, content = $2, author = $3, image = $4, published = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6",
		article.Title, article.Content, article.Author, article.Image, article.Published, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update news: " + err.Error()})
		return
	}

	// Fetch updated article
	var createdAt, updatedAt time.Time
	err = database.DB.QueryRow(
		"SELECT created_at, updated_at FROM news WHERE id = $1",
		id,
	).Scan(&createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "News article not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated article"})
		return
	}

	article.CreatedAt = createdAt.Format(time.RFC3339)
	article.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, article)
}

// DeleteNews deletes a news article
func DeleteNews(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM news WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete news: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "News article not found"})
		return
	}

	c.JSON(200, gin.H{"message": "News article deleted successfully"})
}
