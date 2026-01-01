# Quick Start: Database Backup

## Start Automated Backups (Every 10 Minutes)

### Method 1: Go Scheduler (Easiest)

```bash
# Run in background
cd /Users/olaro/Documents/GitHub/playtzapi
nohup ./scripts/backup_scheduler > backup_scheduler.log 2>&1 &
```

### Method 2: Cron Job

```bash
# Set up cron job (runs every 10 minutes)
./scripts/setup_backup_cron.sh
```

### Method 3: Manual Test

```bash
# Test backup once
./scripts/backup_db.sh

# Or using Go
go run scripts/backup_db.go
```

## Verify Backups

```bash
# List backups
ls -lh backups/

# View latest backup
ls -t backups/ | head -1
```

## Stop Scheduler

```bash
# Find and kill process
ps aux | grep backup_scheduler
kill <PID>
```

## View Logs

```bash
# Backup logs
tail -f backup.log

# Scheduler logs
tail -f backup_scheduler.log
```

That's it! Your database will be backed up every 10 minutes automatically.

