# Quick Start: Railway Backups

## ✅ YES - Backups can run on Railway!

I've created Railway-compatible backup scripts. Here's how to set it up:

## Option 1: Separate Backup Service (Recommended)

### Steps:

1. **In Railway Dashboard:**
   - Go to your project
   - Click "New" → "Empty Service"
   - Name it: `backup-service`

2. **Configure the service:**
   - **Build Command:**
     ```
     go mod download && go mod tidy && go build -o backup-scheduler cmd/backup-scheduler/main.go
     ```
   - **Start Command:**
     ```
     ./backup-scheduler
     ```
   - **Environment Variables:**
     - `DATABASE_URL` - Copy from your main service (Railway will auto-suggest)
     - All other env vars are optional

3. **Deploy:**
   - Railway will build and deploy
   - The service runs 24/7
   - Backups happen every 10 minutes automatically

## Option 2: Add to Main Service (Simpler)

Add backup functionality directly to your main application:

1. **Update `main.go`** - Add this after database initialization:

```go
// Start backup scheduler in background
go func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        cmd := exec.Command("bash", "scripts/backup_railway.sh")
        cmd.Env = os.Environ()
        cmd.Run()
    }
}()
```

2. **Deploy** - Your main service will now backup every 10 minutes

## Option 3: Railway Scheduled Tasks

If Railway supports scheduled tasks:

1. Go to your project → Settings → Scheduled Tasks
2. Add new task:
   - **Schedule:** `*/10 * * * *` (every 10 minutes)
   - **Command:** `go run cmd/backup/main.go`
3. Save

## Important Notes

⚠️ **Railway's filesystem is ephemeral** - backups stored locally will be lost on restart!

### Solutions:

1. **Upload to S3** (Recommended)
   - Add environment variables:
     - `BACKUP_S3_BUCKET=your-bucket-name`
     - `AWS_ACCESS_KEY_ID=your-key`
     - `AWS_SECRET_ACCESS_KEY=your-secret`
   - Backups automatically upload to S3

2. **Use Railway Volumes** (Persistent Storage)
   - Create a volume in Railway
   - Mount to `/backups`
   - Backups persist across restarts

3. **Send to External Service**
   - Use webhook/API to send backups
   - Store on your own server

## Testing

Test locally first:

```bash
# Simulate Railway environment
export DATABASE_URL="your_railway_database_url"
go run cmd/backup/main.go
```

## Files Created

- ✅ `scripts/backup_railway.sh` - Railway-compatible backup script
- ✅ `cmd/backup/main.go` - Single backup command
- ✅ `cmd/backup-scheduler/main.go` - Continuous scheduler
- ✅ `RAILWAY_BACKUP.md` - Full documentation

## Quick Deploy

**Easiest method:** Create a separate backup service in Railway dashboard with:
- Build: `go build -o backup-scheduler cmd/backup-scheduler/main.go`
- Start: `./backup-scheduler`
- Env: `DATABASE_URL` (from main service)

That's it! Backups will run automatically on Railway every 10 minutes! 🚀

