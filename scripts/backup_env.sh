#!/bin/bash

# Backup .env file script
# This backs up the .env file to the backups directory

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${PROJECT_ROOT}/backups"
ENV_FILE="${PROJECT_ROOT}/.env"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "⚠️  .env file not found at $ENV_FILE"
    exit 1
fi

# Generate backup filename with timestamp
TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
BACKUP_FILE="${BACKUP_DIR}/.env_backup_${TIMESTAMP}"

# Copy .env file to backup location
cp "$ENV_FILE" "$BACKUP_FILE"

# Compress the backup
if command -v gzip &> /dev/null; then
    gzip "$BACKUP_FILE"
    BACKUP_FILE="${BACKUP_FILE}.gz"
fi

# Get file size
BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)

echo "✅ .env file backed up: $BACKUP_FILE ($BACKUP_SIZE)"

# Clean up old .env backups (keep last 50)
cd "$BACKUP_DIR"
ENV_BACKUP_COUNT=$(ls -1 .env_backup_*.gz 2>/dev/null | wc -l)

if [ "$ENV_BACKUP_COUNT" -gt 50 ]; then
    echo "🧹 Cleaning up old .env backups (keeping last 50)..."
    ls -1t .env_backup_*.gz 2>/dev/null | tail -n +51 | xargs rm -f
    echo "✅ Cleanup completed"
fi

exit 0

