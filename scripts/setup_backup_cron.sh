#!/bin/bash

# Setup script to install the backup cron job

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_SCRIPT="${SCRIPT_DIR}/backup_db.sh"

# Make backup script executable
chmod +x "$BACKUP_SCRIPT"

# Get absolute path to backup script
BACKUP_SCRIPT_ABS="$(cd "$(dirname "$BACKUP_SCRIPT")" && pwd)/$(basename "$BACKUP_SCRIPT")"

echo "Setting up database backup cron job..."
echo "Backup script: $BACKUP_SCRIPT_ABS"
echo ""

# Check if cron job already exists
ENV_BACKUP_SCRIPT="${SCRIPT_DIR}/backup_env.sh"
ENV_BACKUP_SCRIPT_ABS="$(cd "$(dirname "$ENV_BACKUP_SCRIPT")" && pwd)/$(basename "$ENV_BACKUP_SCRIPT")"
CRON_CMD="*/10 * * * * $BACKUP_SCRIPT_ABS && $ENV_BACKUP_SCRIPT_ABS >> ${PROJECT_ROOT}/backup_cron.log 2>&1"

if crontab -l 2>/dev/null | grep -q "$BACKUP_SCRIPT_ABS"; then
    echo "⚠️  Cron job already exists. Removing old entry..."
    crontab -l 2>/dev/null | grep -v "$BACKUP_SCRIPT_ABS" | crontab -
fi

# Add new cron job
(crontab -l 2>/dev/null; echo "$CRON_CMD") | crontab -

echo "✅ Cron job installed successfully!"
echo ""
echo "The database will be backed up every 10 minutes."
echo "Backups will be stored in: ${PROJECT_ROOT}/backups"
echo "Logs will be written to: ${PROJECT_ROOT}/backup.log"
echo ""
echo "To view current cron jobs: crontab -l"
echo "To remove this cron job: crontab -e (then remove the line)"

