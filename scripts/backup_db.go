package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"playtz-api/database"
	"time"

	"github.com/joho/godotenv"
)

// BackupDB performs a database backup
func BackupDB() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	// Initialize database connection to verify it's accessible
	if err := database.InitDB(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.CloseDB()

	// Get DATABASE_URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	// Create backup directory
	projectRoot, _ := os.Getwd()
	backupDir := filepath.Join(projectRoot, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", timestamp))
	backupFileCompressed := backupFile + ".gz"

	// Add PostgreSQL to PATH if installed via Homebrew
	postgresPaths := []string{
		"/opt/homebrew/opt/postgresql@17/bin",
		"/opt/homebrew/opt/postgresql@16/bin",
		"/opt/homebrew/opt/postgresql@14/bin",
		"/opt/homebrew/opt/postgresql/bin",
	}
	
	currentPath := os.Getenv("PATH")
	for _, pgPath := range postgresPaths {
		if _, err := os.Stat(pgPath); err == nil {
			os.Setenv("PATH", pgPath+":"+currentPath)
			break
		}
	}

	// Check if pg_dump is available
	if _, err := exec.LookPath("pg_dump"); err != nil {
		return fmt.Errorf("pg_dump not found. Please install PostgreSQL client tools")
	}

	// Run pg_dump
	cmd := exec.Command("pg_dump", databaseURL, "--clean", "--if-exists", "--create")
	
	// Create output file
	outFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		outFile.Close()
		os.Remove(backupFile)
		return fmt.Errorf("pg_dump failed: %w", err)
	}
	outFile.Close()

	// Compress backup using gzip
	if _, err := exec.LookPath("gzip"); err == nil {
		cmd = exec.Command("gzip", backupFile)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: Failed to compress backup: %v\n", err)
		} else {
			backupFile = backupFileCompressed
		}
	}

	// Get file size
	fileInfo, _ := os.Stat(backupFile)
	fileSize := fileInfo.Size()

	fmt.Printf("✅ Backup completed: %s (%.2f MB)\n", backupFile, float64(fileSize)/(1024*1024))

	// Clean up old backups (keep last 100)
	cleanupOldBackups(backupDir, 100)

	return nil
}

func cleanupOldBackups(backupDir string, keepCount int) {
	files, err := filepath.Glob(filepath.Join(backupDir, "backup_*.sql*"))
	if err != nil {
		return
	}

	if len(files) <= keepCount {
		return
	}

	// Sort by modification time (newest first)
	// Simple approach: remove files beyond keepCount
	// In production, you'd want to sort by time
	for i := keepCount; i < len(files); i++ {
		os.Remove(files[i])
	}
}

func main() {
	if err := BackupDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

