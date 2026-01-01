package handlers

import (
	"database/sql"
	"playtz-api/auth"
	"playtz-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	User      *User  `json:"user,omitempty"`
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Username and password are required"})
		return
	}

	// Get user from database
	var user User
	var passwordHash string
	var createdAt, updatedAt time.Time
	var roleName *string

	err := database.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.password_hash, u.first_name, u.last_name, 
		       u.role_id, u.active, u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 OR u.email = $1
	`, req.Username).Scan(
		&user.ID, &user.Email, &user.Username, &passwordHash,
		&user.FirstName, &user.LastName, &user.RoleID, &user.Active,
		&createdAt, &updatedAt, &roleName,
	)

	if err == sql.ErrNoRows {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to authenticate user"})
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(403, gin.H{"error": "Account is inactive"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	// Set role name
	if roleName != nil {
		user.RoleName = *roleName
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Create session
	sessionStore := auth.GetSessionStore()
	sessionID := sessionStore.CreateSession(
		user.ID,
		user.Username,
		user.Email,
		user.RoleID,
		user.RoleName,
	)

	// Set session cookie
	c.SetCookie("session_id", sessionID, 600, "/", "", false, true) // 10 minutes, httpOnly

	c.JSON(200, LoginResponse{
		Success:   true,
		Message:   "Login successful",
		SessionID: sessionID,
		User:      &user,
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err == nil && sessionID != "" {
		sessionStore := auth.GetSessionStore()
		sessionStore.DeleteSession(sessionID)
	}

	// Clear cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	sess := session.(*auth.Session)
	
	// Get fresh user data
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

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword handles password change for authenticated users
func ChangePassword(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	sess := session.(*auth.Session)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Current password and new password (min 6 characters) are required"})
		return
	}

	// Get current password hash from database
	var passwordHash string
	err := database.DB.QueryRow(
		"SELECT password_hash FROM users WHERE id = $1",
		sess.UserID,
	).Scan(&passwordHash)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.CurrentPassword))
	if err != nil {
		c.JSON(401, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password and clear password_change_required flag
	_, err = database.DB.Exec(
		"UPDATE users SET password_hash = $1, password_change_required = false, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		string(hashedPassword), sess.UserID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

