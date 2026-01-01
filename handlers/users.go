package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user account
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    string `json:"role_id"`
	RoleName  string `json:"role_name,omitempty"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"` // Optional - will generate default if not provided
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    string `json:"role_id"`
}

// CreateUserResponse represents user creation response
type CreateUserResponse struct {
	User            User   `json:"user"`
	DefaultPassword string `json:"default_password,omitempty"` // Only returned when default password is generated
	Message         string `json:"message,omitempty"`
}

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT u.id, u.email, u.username, u.first_name, u.last_name, u.role_id, u.active, u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var createdAt, updatedAt time.Time
		var roleName *string
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.RoleID, &user.Active, &createdAt, &updatedAt, &roleName)
		if err != nil {
			continue
		}
		if roleName != nil {
			user.RoleName = *roleName
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
		users = append(users, user)
	}

	c.JSON(200, users)
}

// GetUserByID returns a specific user
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user User
	var createdAt, updatedAt time.Time
	var roleName *string
	err := database.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.first_name, u.last_name, u.role_id, u.active, u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, id).Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.RoleID, &user.Active, &createdAt, &updatedAt, &roleName)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	if roleName != nil {
		user.RoleName = *roleName
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, user)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	userID := uuid.New().String()
	_, err = database.DB.Exec(
		"INSERT INTO users (id, email, username, password_hash, first_name, last_name, role_id, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		userID, req.Email, req.Username, string(hashedPassword), req.FirstName, req.LastName, req.RoleID, true,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var roleName *string
	database.DB.QueryRow(`
		SELECT u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, userID).Scan(&createdAt, &updatedAt, &roleName)

	user := User{
		ID:        userID,
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    req.RoleID,
		Active:    true,
		CreatedAt: createdAt.Format(time.RFC3339),
		UpdatedAt: updatedAt.Format(time.RFC3339),
	}
	if roleName != nil {
		user.RoleName = *roleName
	}

	c.JSON(201, user)
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	user.ID = id

	_, err := database.DB.Exec(
		"UPDATE users SET email = $1, username = $2, first_name = $3, last_name = $4, role_id = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6",
		user.Email, user.Username, user.FirstName, user.LastName, user.RoleID, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var roleName *string
	err = database.DB.QueryRow(`
		SELECT u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, id).Scan(&createdAt, &updatedAt, &roleName)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)
	if roleName != nil {
		user.RoleName = *roleName
	}

	c.JSON(200, user)
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// UpdateUserRole updates a user's role
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	var update struct {
		RoleID string `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET role_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", update.RoleID, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user role: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "User role updated",
		"role_id": update.RoleID,
	})
}
