# Railway Backup Setup Guide

This guide explains how to set up automated database backups on Railway.

## Option 1: Railway Cron Job (Recommended)

Railway supports cron jobs through their platform. You can set up a scheduled task:

### Steps:

1. **Add a new service in Railway** for backups:
   - Go to your Railway project
   - Click "New" → "Empty Service"
   - Name it "backup-service"

2. **Configure the service:**
   - Set the build command: `go build -o backup cmd/backup/main.go`
   - Set the start command: `./backup`
   - Add environment variables:
     - `DATABASE_URL` (automatically available)
     - `PORT` (not needed for cron)

3. **Set up Railway Cron:**
   - In Railway dashboard, go to your backup service
   - Add a cron schedule: `*/10 * * * *` (every 10 minutes)
   - Or use Railway's scheduled tasks feature

## Option 2: Separate Backup Service (Always Running)

Create a separate Railway service that runs continuously and performs backups:

### Steps:

1. **Create backup service:**
   ```toml
   # In railway.toml or Railway dashboard
   [services.backup]
   buildCommand = "go build -o backup-scheduler cmd/backup-scheduler/main.go"
   startCommand = "./backup-scheduler"
   ```

2. **Deploy as separate service:**
   - This service will run 24/7
   - Performs backups every 10 minutes
   - Uses minimal resources

## Option 3: Add to Main Application

You can add backup functionality to your main application:

### Implementation:

Add a background goroutine in `main.go`:

```go
// In main.go, after database initialization
go func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        // Run backup
        cmd := exec.Command("bash", "scripts/backup_railway.sh")
        cmd.Env = os.Environ()
        cmd.Run()
    }
}()
```

## Option 4: Railway Scheduled Tasks (New Feature)

Railway now supports scheduled tasks:

1. Go to your project → Settings
2. Add a scheduled task
3. Set schedule: `*/10 * * * *` (every 10 minutes)
4. Command: `go run cmd/backup/main.go`

## Environment Variables on Railway

Railway automatically provides:
- `DATABASE_URL` - Your PostgreSQL connection string
- `PORT` - Server port (not needed for backups)

Optional variables you can add:
- `BACKUP_S3_BUCKET` - If you want to upload to S3
- `AWS_ACCESS_KEY_ID` - For S3 uploads
- `AWS_SECRET_ACCESS_KEY` - For S3 uploads

## Backup Storage on Railway

**Important:** Railway's filesystem is ephemeral. Backups stored locally will be lost when the service restarts.

### Solutions:

1. **Upload to S3/Cloud Storage** (Recommended)
   - Configure `BACKUP_S3_BUCKET`
   - Backups automatically upload to S3
   - Script includes S3 upload code

2. **Use Railway Volumes** (Persistent Storage)
   - Create a Railway volume
   - Mount it to `/backups`
   - Backups persist across restarts

3. **Send to External Service**
   - Use webhook to send backups
   - Store in your own server
   - Use Railway's API to download

## Recommended Setup

For production, I recommend:

1. **Use Railway Scheduled Tasks** (if available)
   - Most reliable
   - No need for always-running service
   - Cost-effective

2. **Upload to S3**
   - Reliable storage
   - Version control
   - Easy to restore

3. **Keep last 10 backups on Railway**
   - For quick access
   - Use volume for persistence

## Testing on Railway

1. Deploy the backup service
2. Check logs: `railway logs --service backup`
3. Verify backups are created
4. Check S3 (if configured)

## Files Created

- `scripts/backup_railway.sh` - Railway-compatible backup script
- `cmd/backup/main.go` - Single backup command
- `cmd/backup-scheduler/main.go` - Continuous backup scheduler

## Quick Start

```bash
# Test backup locally (simulates Railway)
DATABASE_URL="your_railway_url" go run cmd/backup/main.go

# Or use the script directly
DATABASE_URL="your_railway_url" bash scripts/backup_railway.sh
```

