package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
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

	fmt.Println("🔍 Connected to database")
	fmt.Println("")

	// Get admin role ID
	var adminRoleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID)
	if err == sql.ErrNoRows {
		log.Fatal("Admin role does not exist. Please run the server first to create it.")
	} else if err != nil {
		log.Fatalf("Error getting admin role: %v", err)
	}

	fmt.Printf("✅ Found admin role: %s\n", adminRoleID)
	fmt.Println("")

	// Get all users except admin
	rows, err := db.Query(`
		SELECT id, username, email 
		FROM users 
		WHERE username != 'admin' AND email != 'admin@playtz.com'
	`)
	if err != nil {
		log.Fatalf("Error querying users: %v", err)
	}
	defer rows.Close()

	var usersToDelete []struct {
		ID       string
		Username string
		Email    string
	}

	for rows.Next() {
		var user struct {
			ID       string
			Username string
			Email    string
		}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		usersToDelete = append(usersToDelete, user)
	}

	fmt.Printf("📋 Found %d users to delete (excluding admin)\n", len(usersToDelete))
	if len(usersToDelete) > 0 {
		fmt.Println("Users to be deleted:")
		for _, user := range usersToDelete {
			fmt.Printf("  - %s (%s) - ID: %s\n", user.Username, user.Email, user.ID)
		}
		fmt.Println("")
	}

	// Delete all users except admin
	result, err := db.Exec(`
		DELETE FROM users 
		WHERE username != 'admin' AND email != 'admin@playtz.com'
	`)
	if err != nil {
		log.Fatalf("Error deleting users: %v", err)
	}

	deletedCount, _ := result.RowsAffected()
	fmt.Printf("✅ Deleted %d users\n", deletedCount)
	fmt.Println("")

	// Ensure admin user exists and has admin role
	var adminUserID string
	err = db.QueryRow("SELECT id FROM users WHERE username = 'admin' OR email = 'admin@playtz.com'").Scan(&adminUserID)

	if err == sql.ErrNoRows {
		// Create admin user
		fmt.Println("📝 Creating admin user...")
		
		// Generate UUID for admin user
		adminUserID = "822181e3-f800-4139-ba8a-cfb2f61d9ee6" // Fixed ID for consistency
		
		// Hash password
		password := "admin123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}
		
		_, err = db.Exec(`
			INSERT INTO users (id, email, username, password_hash, first_name, last_name, role_id, active, password_change_required) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, adminUserID, "admin@playtz.com", "admin", string(hashedPassword), "Admin", "User", adminRoleID, true, false)
		
		if err != nil {
			log.Fatalf("Error creating admin user: %v", err)
		}
		fmt.Println("✅ Created admin user")
	} else if err != nil {
		log.Fatalf("Error checking for admin user: %v", err)
	} else {
		// Update admin user to ensure it has admin role
		fmt.Println("📝 Updating admin user...")
		_, err = db.Exec(`
			UPDATE users 
			SET role_id = $1, active = true, password_change_required = false
			WHERE id = $2
		`, adminRoleID, adminUserID)
		if err != nil {
			log.Fatalf("Error updating admin user: %v", err)
		}
		fmt.Println("✅ Updated admin user")
	}

	// Verify admin user
	var adminUsername, adminEmail, adminRoleName string
	err = db.QueryRow(`
		SELECT u.username, u.email, r.name 
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, adminUserID).Scan(&adminUsername, &adminEmail, &adminRoleName)

	if err != nil {
		log.Fatalf("Error verifying admin user: %v", err)
	}

	fmt.Println("")
	fmt.Println("✅ CLEANUP COMPLETE")
	fmt.Println("===================")
	fmt.Printf("Admin Username: %s\n", adminUsername)
	fmt.Printf("Admin Email:    %s\n", adminEmail)
	fmt.Printf("Admin Role:     %s\n", adminRoleName)
	fmt.Printf("Admin Password: admin123\n")
	fmt.Println("")
	fmt.Println("✅ Database now contains only the admin user")
}

