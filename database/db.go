package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	var connStr string
	var dbName string

	// Check if DATABASE_URL is provided (Railway/Heroku style)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// Use DATABASE_URL directly - lib/pq supports PostgreSQL URLs
		connStr = databaseURL
		
		// Extract database name for logging
		parsedURL, err := url.Parse(databaseURL)
		if err == nil {
			dbName = strings.TrimPrefix(parsedURL.Path, "/")
		}
	} else {
		// Fallback to individual environment variables
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}

		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}

		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "celestine" // Default user
		}

		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = ""
		}

		dbName = os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "1029"
		}

		// Build connection string - ensure dbname is explicitly set
		if dbPassword != "" {
			connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				dbHost, dbPort, dbUser, dbPassword, dbName)
		} else {
			connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
				dbHost, dbPort, dbUser, dbName)
		}
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database '%s': %w", dbName, err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
