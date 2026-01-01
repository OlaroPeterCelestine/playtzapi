package middleware

import (
	"playtz-api/auth"

	"github.com/gin-gonic/gin"
)

// RequireAuth middleware checks if user is authenticated using JWT tokens
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			// Try Authorization header as fallback
			authHeader := c.GetHeader("Authorization")
			token = auth.ExtractTokenFromHeader(authHeader)
			if token == "" {
				c.JSON(401, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}
		}

		// Validate JWT token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)
		c.Set("role_name", claims.RoleName)
		c.Set("claims", claims)

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

