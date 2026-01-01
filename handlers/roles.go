package handlers

import (
	"database/sql"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Role represents a user role
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Active      bool     `json:"active"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// GetRoles returns all roles
func GetRoles(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description, permissions, active, created_at, updated_at FROM roles ORDER BY name")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch roles"})
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		var createdAt, updatedAt time.Time
		var permissions pq.StringArray
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &permissions, &role.Active, &createdAt, &updatedAt)
		if err != nil {
			continue
		}
		role.Permissions = []string(permissions)
		role.CreatedAt = createdAt.Format(time.RFC3339)
		role.UpdatedAt = updatedAt.Format(time.RFC3339)
		roles = append(roles, role)
	}

	c.JSON(200, roles)
}

// GetRoleByID returns a specific role
func GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	var role Role
	var createdAt, updatedAt time.Time
	var permissions pq.StringArray
	err := database.DB.QueryRow(
		"SELECT id, name, description, permissions, active, created_at, updated_at FROM roles WHERE id = $1",
		id,
	).Scan(&role.ID, &role.Name, &role.Description, &permissions, &role.Active, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch role"})
		return
	}

	role.Permissions = []string(permissions)
	role.CreatedAt = createdAt.Format(time.RFC3339)
	role.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, role)
}

// CreateRole creates a new role
func CreateRole(c *gin.Context) {
	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	role.ID = uuid.New().String()
	role.Active = true

	_, err := database.DB.Exec(
		"INSERT INTO roles (id, name, description, permissions, active) VALUES ($1, $2, $3, $4, $5)",
		role.ID, role.Name, role.Description, pq.Array(role.Permissions), role.Active,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create role: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var permissions pq.StringArray
	database.DB.QueryRow(
		"SELECT permissions, created_at, updated_at FROM roles WHERE id = $1",
		role.ID,
	).Scan(&permissions, &createdAt, &updatedAt)

	role.Permissions = []string(permissions)
	role.CreatedAt = createdAt.Format(time.RFC3339)
	role.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(201, role)
}

// UpdateRole updates an existing role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	role.ID = id

	_, err := database.DB.Exec(
		"UPDATE roles SET name = $1, description = $2, permissions = $3, active = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5",
		role.Name, role.Description, pq.Array(role.Permissions), role.Active, id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update role: " + err.Error()})
		return
	}

	var createdAt, updatedAt time.Time
	var permissions pq.StringArray
	err = database.DB.QueryRow(
		"SELECT permissions, created_at, updated_at FROM roles WHERE id = $1",
		id,
	).Scan(&permissions, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch updated role"})
		return
	}

	role.Permissions = []string(permissions)
	role.CreatedAt = createdAt.Format(time.RFC3339)
	role.UpdatedAt = updatedAt.Format(time.RFC3339)

	c.JSON(200, role)
}

// DeleteRole deletes a role
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	// Check if role is being used by any users
	var userCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role_id = $1", id).Scan(&userCount)
	if err == nil && userCount > 0 {
		c.JSON(400, gin.H{"error": "Cannot delete role: it is assigned to users"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete role: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Role deleted successfully"})
}
