package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

// RailwayBackup runs a database backup on Railway
func main() {
	// Load environment variables (Railway provides these automatically)
	if err := godotenv.Load(); err != nil {
		// On Railway, .env might not exist, that's okay
		fmt.Printf("Note: .env file not found (this is normal on Railway)\n")
	}

	// Check if DATABASE_URL is set
	if os.Getenv("DATABASE_URL") == "" {
		fmt.Fprintf(os.Stderr, "Error: DATABASE_URL not set\n")
		os.Exit(1)
	}

	fmt.Printf("[%s] Starting Railway database backup...\n", time.Now().Format("2006-01-02 15:04:05"))

	// Run the backup script
	scriptPath := "scripts/backup_railway.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Try alternative path
		scriptPath = "./scripts/backup_railway.sh"
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[%s] Backup completed successfully\n", time.Now().Format("2006-01-02 15:04:05"))
}

