package handlers

import (
	"playtz-api/auth"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminDashboardData represents data for the admin dashboard
type AdminDashboardData struct {
	User         User   `json:"user"`
	Stats        Stats  `json:"stats"`
	RecentNews   []NewsArticle `json:"recent_news,omitempty"`
	RecentEvents []Event `json:"recent_events,omitempty"`
	RecentOrders []Order `json:"recent_orders,omitempty"`
}

// Stats represents dashboard statistics
type Stats struct {
	TotalUsers      int `json:"total_users"`
	TotalNews       int `json:"total_news"`
	TotalEvents     int `json:"total_events"`
	TotalMerchandise int `json:"total_merchandise"`
	TotalOrders     int `json:"total_orders"`
	PendingOrders   int `json:"pending_orders"`
}

// GetAdminDashboard returns dashboard data based on user role
func GetAdminDashboard(c *gin.Context) {
	session, _ := c.Get("session")
	sess := session.(*auth.Session)

	// Get user details
	var user User
	var createdAt, updatedAt time.Time
	var roleName *string

	err := database.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.first_name, u.last_name, 
		       u.role_id, u.active, u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, sess.UserID).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.FirstName, &user.LastName, &user.RoleID, &user.Active,
		&createdAt, &updatedAt, &roleName,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user data"})
		return
	}

	if roleName != nil {
		user.RoleName = *roleName
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Get statistics based on role
	stats := Stats{}
	
	// All roles can see basic counts
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	database.DB.QueryRow("SELECT COUNT(*) FROM news").Scan(&stats.TotalNews)
	database.DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&stats.TotalEvents)
	database.DB.QueryRow("SELECT COUNT(*) FROM merchandise").Scan(&stats.TotalMerchandise)
	database.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	database.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = 'pending'").Scan(&stats.PendingOrders)

	dashboardData := AdminDashboardData{
		User:  user,
		Stats: stats,
	}

	// Role-based data access
	switch user.RoleName {
	case "Admin", "Super Admin":
		// Admins see everything
		dashboardData.RecentNews = getRecentNews(5)
		dashboardData.RecentEvents = getRecentEvents(5)
		dashboardData.RecentOrders = getRecentOrders(5)
	case "Editor", "Content Manager":
		// Editors see news and events
		dashboardData.RecentNews = getRecentNews(5)
		dashboardData.RecentEvents = getRecentEvents(5)
	case "Manager":
		// Managers see orders
		dashboardData.RecentOrders = getRecentOrders(5)
	}

	c.JSON(200, dashboardData)
}

// Helper functions
func getRecentNews(limit int) []NewsArticle {
	rows, err := database.DB.Query(`
		SELECT id, title, content, author, image, published, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return []NewsArticle{}
	}
	defer rows.Close()

	var articles []NewsArticle
	for rows.Next() {
		var article NewsArticle
		var createdAt, updatedAt time.Time
		rows.Scan(&article.ID, &article.Title, &article.Content, &article.Author,
			&article.Image, &article.Published, &createdAt, &updatedAt)
		article.CreatedAt = createdAt.Format(time.RFC3339)
		article.UpdatedAt = updatedAt.Format(time.RFC3339)
		articles = append(articles, article)
	}
	return articles
}

func getRecentEvents(limit int) []Event {
	rows, err := database.DB.Query(`
		SELECT id, title, description, date, time, location, image, active, created_at, updated_at
		FROM events
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return []Event{}
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var createdAt, updatedAt time.Time
		rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date,
			&event.Time, &event.Location, &event.Image, &event.Active,
			&createdAt, &updatedAt)
		event.CreatedAt = createdAt.Format(time.RFC3339)
		event.UpdatedAt = updatedAt.Format(time.RFC3339)
		events = append(events, event)
	}
	return events
}

func getRecentOrders(limit int) []Order {
	rows, err := database.DB.Query(`
		SELECT id, user_id, total, status, shipping_address, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return []Order{}
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var shippingAddrJSON string
		var createdAt, updatedAt time.Time
		rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status,
			&shippingAddrJSON, &createdAt, &updatedAt)
		order.CreatedAt = createdAt.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}
	return orders
}

