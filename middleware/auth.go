package middleware

import (
	"playtz-api/auth"

	"github.com/gin-gonic/gin"
)

// RequireAuth middleware checks if user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// Try to get from header as fallback
			sessionID = c.GetHeader("X-Session-ID")
			if sessionID == "" {
				c.JSON(401, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}
		}

		sessionStore := auth.GetSessionStore()
		session, exists := sessionStore.GetSession(sessionID)
		if !exists {
			c.JSON(401, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		// Store session in context
		c.Set("session", session)
		c.Set("user_id", session.UserID)
		c.Set("role_id", session.RoleID)
		c.Set("role_name", session.RoleName)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, exists := c.Get("role_name")
		if !exists {
			c.JSON(403, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		userRole := roleName.(string)
		hasRole := false
		for _, requiredRole := range roleNames {
			if userRole == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission middleware checks if user has required permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would require fetching role permissions from database
		// For now, we'll use role-based access
		c.Next()
	}
}

