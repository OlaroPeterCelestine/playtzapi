#!/bin/bash

# Railway-compatible Database Backup Script
# This script is designed to run on Railway and upload backups to cloud storage

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${PROJECT_ROOT}/backups"
LOG_FILE="${PROJECT_ROOT}/backup.log"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Function to log messages
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Parse DATABASE_URL
parse_database_url() {
    local url="$1"
    url="${url#postgresql://}"
    url="${url#postgres://}"
    
    if [[ $url =~ ^([^:]+):([^@]+)@([^:]+):([^/]+)/(.+)$ ]]; then
        DB_USER="${BASH_REMATCH[1]}"
        DB_PASSWORD="${BASH_REMATCH[2]}"
        DB_HOST="${BASH_REMATCH[3]}"
        DB_PORT="${BASH_REMATCH[4]}"
        DB_NAME="${BASH_REMATCH[5]}"
        return 0
    fi
    return 1
}

# Get database connection details
if [ -n "$DATABASE_URL" ]; then
    if parse_database_url "$DATABASE_URL"; then
        log "${GREEN}✓${NC} Parsed DATABASE_URL"
    else
        log "${RED}✗${NC} Failed to parse DATABASE_URL"
        exit 1
    fi
else
    log "${RED}✗${NC} DATABASE_URL not set"
    exit 1
fi

# Check if pg_dump is available (Railway should have it)
if ! command -v pg_dump &> /dev/null; then
    log "${RED}✗${NC} pg_dump not found"
    exit 1
fi

# Generate backup filename
TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql"
BACKUP_FILE_COMPRESSED="${BACKUP_FILE}.gz"

# Set PGPASSWORD
export PGPASSWORD="$DB_PASSWORD"

log "${YELLOW}→${NC} Starting database backup on Railway..."
log "   Database: $DB_NAME"
log "   Host: $DB_HOST:$DB_PORT"

# Run pg_dump
if pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    --no-password \
    --clean \
    --if-exists \
    --create \
    > "$BACKUP_FILE" 2>> "$LOG_FILE"; then
    
    # Compress backup
    if gzip "$BACKUP_FILE"; then
        BACKUP_SIZE=$(du -h "$BACKUP_FILE_COMPRESSED" | cut -f1)
        log "${GREEN}✓${NC} Backup completed: $BACKUP_FILE_COMPRESSED ($BACKUP_SIZE)"
        
        # Upload to cloud storage if configured (optional)
        if [ -n "$BACKUP_S3_BUCKET" ] && command -v aws &> /dev/null; then
            log "${YELLOW}→${NC} Uploading to S3..."
            aws s3 cp "$BACKUP_FILE_COMPRESSED" "s3://${BACKUP_S3_BUCKET}/backups/" && \
                log "${GREEN}✓${NC} Uploaded to S3" || \
                log "${RED}✗${NC} S3 upload failed"
        fi
        
        # Clean up old backups (keep last 10 on Railway to save space)
        cd "$BACKUP_DIR"
        BACKUP_COUNT=$(ls -1 backup_*.sql.gz 2>/dev/null | wc -l)
        
        if [ "$BACKUP_COUNT" -gt 10 ]; then
            log "${YELLOW}→${NC} Cleaning up old backups (keeping last 10)..."
            ls -1t backup_*.sql.gz 2>/dev/null | tail -n +11 | xargs rm -f
        fi
        
        exit 0
    else
        log "${RED}✗${NC} Failed to compress backup"
        exit 1
    fi
else
    log "${RED}✗${NC} Backup failed. Check log: $LOG_FILE"
    rm -f "$BACKUP_FILE"
    exit 1
fi

