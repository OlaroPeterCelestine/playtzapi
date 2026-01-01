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

// BackupScheduler runs database backups every 10 minutes
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
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

	// Run initial backup
	fmt.Println("🚀 Starting database backup scheduler...")
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
	fmt.Printf("[%s] Starting backup...\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// Run the backup script
	cmd := exec.Command("go", "run", "scripts/backup_db.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Backup failed: %v\n", err)
	} else {
		fmt.Printf("✅ Backup completed successfully\n")
	}
	
	fmt.Println()
}

func runEnvBackup() {
	fmt.Printf("[%s] Backing up .env file...\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// Run the .env backup script
	cmd := exec.Command("bash", "scripts/backup_env.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  .env backup failed: %v\n", err)
	} else {
		fmt.Printf("✅ .env backup completed\n")
	}
}

