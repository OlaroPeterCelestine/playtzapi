# Playtz API

A comprehensive REST API for Playtz 102.9 radio station, built with Go and Gin.

## 📁 Project Structure

```
playtzapi/
├── docs/              # Documentation files
│   ├── API_ROUTES.md  # Complete API routes documentation
│   ├── DEPLOYMENT.md  # Railway deployment guide
│   └── ...
├── scripts/           # Utility scripts
│   ├── backup_*.sh    # Database backup scripts
│   ├── test_endpoints.sh  # API testing script
│   └── ...
├── handlers/          # API route handlers
├── database/          # Database connection and migrations
├── middleware/        # Authentication middleware
├── auth/              # Session management
├── models/            # Data models
├── static/            # HTML pages (login, dashboard)
└── cmd/               # Command-line tools
    ├── backup/        # Backup command
    └── backup-scheduler/  # Backup scheduler
```

## 📚 Documentation

- [API Routes](./docs/API_ROUTES.md) - Complete API documentation
- [Deployment Guide](./docs/DEPLOYMENT.md) - How to deploy to Railway
- [Backup Setup](./docs/RAILWAY_BACKUP_QUICKSTART.md) - Automated backups

## 🚀 Quick Start

1. **Setup Environment:**
   ```bash
   cp .env.example .env  # Create .env file
   # Edit .env with your credentials
   ```

2. **Run Database Migrations:**
   ```bash
   go run main.go  # Migrations run automatically
   ```

3. **Start Server:**
   ```bash
   go run main.go
   ```

4. **Test Endpoints:**
   ```bash
   ./scripts/test_endpoints.sh
   ```

## 🔐 Admin Access

- Login: `http://localhost:8080/admin/login`
- Dashboard: `http://localhost:8080/admin/dashboard`

## 💾 Backups

Automated backups run every 10 minutes:
- Database backups: `backups/backup_*.sql.gz`
- .env backups: `backups/.env_backup_*.gz`

See [Backup Documentation](./docs/RAILWAY_BACKUP_QUICKSTART.md) for Railway setup.
