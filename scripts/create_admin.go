package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✅ Connected to database")

	// Check if admin role exists
	var adminRoleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID)
	
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

		_, err = db.Exec(
			"INSERT INTO roles (id, name, description, permissions, active) VALUES ($1, $2, $3, $4, $5)",
			adminRoleID, "admin", "Administrator with full system access", pq.Array(adminPermissions), true,
		)
		if err != nil {
			log.Fatalf("Failed to create admin role: %v", err)
		}
		fmt.Println("✅ Created admin role")
	} else if err != nil {
		log.Fatalf("Error checking for admin role: %v", err)
	} else {
		fmt.Println("✅ Admin role already exists")
	}

	// Check if admin user exists
	var adminUserID string
	err = db.QueryRow("SELECT id FROM users WHERE username = 'admin' OR email = 'admin@playtz.com'").Scan(&adminUserID)
	
	if err == sql.ErrNoRows {
		// Create admin user
		adminUserID = uuid.New().String()
		password := "admin123" // Simple password for testing
		
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		_, err = db.Exec(
			"INSERT INTO users (id, email, username, password_hash, first_name, last_name, role_id, active, password_change_required) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
			adminUserID, "admin@playtz.com", "admin", string(hashedPassword), "Admin", "User", adminRoleID, true, false,
		)
		if err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		fmt.Println("✅ Created admin user")
		fmt.Println("")
		fmt.Println("📋 Admin Credentials:")
		fmt.Println("   Username: admin")
		fmt.Println("   Email: admin@playtz.com")
		fmt.Println("   Password: admin123")
		fmt.Println("")
	} else if err != nil {
		log.Fatalf("Error checking for admin user: %v", err)
	} else {
		// Admin user exists, reset password to admin123
		password := "admin123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		_, err = db.Exec(
			"UPDATE users SET password_hash = $1, password_change_required = $2 WHERE id = $3",
			string(hashedPassword), false, adminUserID,
		)
		if err != nil {
			log.Fatalf("Failed to update admin password: %v", err)
		}
		fmt.Println("✅ Admin user already exists - password reset to admin123")
		fmt.Println("")
		fmt.Println("📋 Admin Credentials:")
		fmt.Println("   Username: admin")
		fmt.Println("   Email: admin@playtz.com")
		fmt.Println("   Password: admin123")
		fmt.Println("")
	}

	fmt.Println("✅ Setup complete!")
}

