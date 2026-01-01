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
	// For cross-origin (Vercel), we need SameSite=None and Secure=true
	// Using Header directly to ensure SameSite=None is set correctly
	c.Header("Set-Cookie", "session_id="+sessionID+"; Path=/; Max-Age=600; HttpOnly; Secure; SameSite=None")

	c.JSON(200, LoginResponse{
		Success:   true,
		Message:   "Login successful",
		SessionID: sessionID,
		User:      &user,
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	// Try to get session from cookie or header
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
	}

	// Delete session from store if it exists
	if sessionID != "" {
		sessionStore := auth.GetSessionStore()
		sessionStore.DeleteSession(sessionID)
	}

	// Clear cookie - use Header directly for cross-origin support
	// Set multiple cookie clearing headers to ensure it works across browsers
	// Also try to clear with different domain settings
	c.Header("Set-Cookie", "session_id=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=None")
	c.Header("Set-Cookie", "session_id=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure; SameSite=None")
	
	// Additional cookie clearing for better compatibility
	c.SetCookie("session_id", "", -1, "/", "", true, true)

	c.JSON(200, gin.H{
		"message": "Logged out successfully",
		"success": true,
	})
}

// GetCurrentUserOptional returns the current authenticated user if available, or null if not authenticated
// This prevents 401 errors on login page when checking auth status
func GetCurrentUserOptional(c *gin.Context) {
	// Try to get session from cookie or header
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
	}

	if sessionID == "" {
		// Not authenticated - return 200 with null (not an error)
		c.JSON(200, gin.H{"authenticated": false, "user": nil})
		return
	}

	// Check if session exists
	sessionStore := auth.GetSessionStore()
	session, exists := sessionStore.GetSession(sessionID)
	if !exists {
		// Session invalid or expired - return 200 with null (not an error)
		c.JSON(200, gin.H{"authenticated": false, "user": nil})
		return
	}

	// Get fresh user data
	var user User
	var createdAt, updatedAt time.Time
	var roleName *string

	err = database.DB.QueryRow(`
		SELECT u.id, u.email, u.username, u.first_name, u.last_name, 
		       u.role_id, u.active, u.created_at, u.updated_at, r.name as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, session.UserID).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.FirstName, &user.LastName, &user.RoleID, &user.Active,
		&createdAt, &updatedAt, &roleName,
	)

	if err != nil {
		// Database error - return 200 with null (don't expose internal errors)
		c.JSON(200, gin.H{"authenticated": false, "user": nil})
		return
	}

	if roleName != nil {
		user.RoleName = *roleName
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Authenticated - return user
	c.JSON(200, gin.H{"authenticated": true, "user": user})
}

// GetCurrentUser returns the current authenticated user (requires auth)
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

