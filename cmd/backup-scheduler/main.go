package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// RailwayBackupScheduler runs backups every 10 minutes on Railway
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Note: .env file not found (normal on Railway)\n")
	}

	// Create backup directory
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create backup directory: %v\n", err)
		os.Exit(1)
	}

	// Channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("🚀 Railway Backup Scheduler Started")
	fmt.Println("   Backups will run every 10 minutes")
	fmt.Println("   Press Ctrl+C to stop")
	fmt.Println()

	// Run first backup immediately
	runBackup()
	runEnvBackup()

	// Create ticker for 10-minute intervals
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Backup loop
	for {
		select {
		case <-ticker.C:
			runBackup()
			runEnvBackup()
		case <-sigChan:
			fmt.Println("\n🛑 Shutting down backup scheduler...")
			return
		}
	}
}

func runBackup() {
	fmt.Printf("[%s] Starting database backup...\n", time.Now().Format("2006-01-02 15:04:05"))

	// Use the Railway backup script
	scriptPath := "scripts/backup_railway.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "./scripts/backup_railway.sh"
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Backup failed: %v\n", err)
	} else {
		fmt.Printf("✅ Backup completed successfully\n")
	}
	fmt.Println()
}

func runEnvBackup() {
	fmt.Printf("[%s] Backing up .env file...\n", time.Now().Format("2006-01-02 15:04:05"))

	scriptPath := "scripts/backup_env.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "./scripts/backup_env.sh"
	}

	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  .env backup failed: %v\n", err)
	} else {
		fmt.Printf("✅ .env backup completed\n")
	}
}

