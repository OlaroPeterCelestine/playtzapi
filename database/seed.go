package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdmin creates an admin user if it doesn't exist
// This runs automatically on startup to ensure admin exists
func SeedAdmin() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Check if admin role exists
	var adminRoleID string
	err := DB.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID)

	if err == sql.ErrNoRows {
		// Create admin role
		adminRoleID = uuid.New().String()
		adminPermissions := []string{
			"users.read", "users.write", "users.delete",
			"roles.read", "roles.write", "roles.delete",
			"news.read", "news.write", "news.delete",
			"events.read", "events.write", "events.delete",
			"mixes.read", "mixes.write", "mixes.delete",
			"merchandise.read", "merchandise.write", "merchandise.delete",
			"orders.read", "orders.write", "orders.delete",
			"careers.read", "careers.write", "careers.delete",
			"rooms.read", "rooms.write", "rooms.delete",
			"admin.dashboard", "admin.settings",
		}

		_, err = DB.Exec(
			"INSERT INTO roles (id, name, description, permissions, active) VALUES ($1, $2, $3, $4, $5)",
			adminRoleID, "admin", "Administrator with full system access", pq.Array(adminPermissions), true,
		)
		if err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}
		log.Println("✅ Created admin role")
	} else if err != nil {
		return fmt.Errorf("error checking for admin role: %w", err)
	} else {
		log.Println("✅ Admin role exists")
	}

	// Check if admin user exists
	var adminUserID string
	err = DB.QueryRow("SELECT id FROM users WHERE username = 'admin' OR email = 'admin@playtz.com'").Scan(&adminUserID)

	if err == sql.ErrNoRows {
		// Create admin user
		adminUserID = uuid.New().String()
		password := "admin123" // Default password

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		_, err = DB.Exec(
			"INSERT INTO users (id, email, username, password_hash, first_name, last_name, role_id, active, password_change_required) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
			adminUserID, "admin@playtz.com", "admin", string(hashedPassword), "Admin", "User", adminRoleID, true, false,
		)
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		log.Println("✅ Created admin user (username: admin, password: admin123)")
	} else if err != nil {
		return fmt.Errorf("error checking for admin user: %w", err)
	} else {
		log.Println("✅ Admin user exists")
	}

	return nil
}

