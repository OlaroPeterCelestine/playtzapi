#!/bin/bash

# Database Backup Script for Playtz API
# This script backs up the PostgreSQL database every 10 minutes

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${PROJECT_ROOT}/backups"
LOG_FILE="${PROJECT_ROOT}/backup.log"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f "${PROJECT_ROOT}/.env" ]; then
    export $(cat "${PROJECT_ROOT}/.env" | grep -v '^#' | xargs)
fi

# Function to log messages
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# Function to parse DATABASE_URL
parse_database_url() {
    local url="$1"
    
    # Remove postgresql:// or postgres:// prefix
    url="${url#postgresql://}"
    url="${url#postgres://}"
    
    # Extract user:password@host:port/database
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
    # Fallback to individual environment variables
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-5432}"
    DB_USER="${DB_USER:-postgres}"
    DB_PASSWORD="${DB_PASSWORD:-}"
    DB_NAME="${DB_NAME:-railway}"
    
    if [ -z "$DB_PASSWORD" ]; then
        log "${YELLOW}⚠${NC} DB_PASSWORD not set, trying without password"
    fi
fi

# Add PostgreSQL to PATH if installed via Homebrew (prioritize version 17, then 16, then 14)
if [ -d "/opt/homebrew/opt/postgresql@17/bin" ]; then
    export PATH="/opt/homebrew/opt/postgresql@17/bin:$PATH"
elif [ -d "/opt/homebrew/opt/postgresql@16/bin" ]; then
    export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"
elif [ -d "/opt/homebrew/opt/postgresql@14/bin" ]; then
    export PATH="/opt/homebrew/opt/postgresql@14/bin:$PATH"
elif [ -d "/opt/homebrew/opt/postgresql/bin" ]; then
    export PATH="/opt/homebrew/opt/postgresql/bin:$PATH"
fi

# Check if pg_dump is available
if ! command -v pg_dump &> /dev/null; then
    log "${RED}✗${NC} pg_dump not found. Please install PostgreSQL client tools."
    log "   macOS: brew install postgresql"
    log "   Ubuntu: sudo apt-get install postgresql-client"
    exit 1
fi

# Generate backup filename with timestamp
TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql"
BACKUP_FILE_COMPRESSED="${BACKUP_FILE}.gz"

# Set PGPASSWORD environment variable for pg_dump
export PGPASSWORD="$DB_PASSWORD"

# Perform backup
log "${YELLOW}→${NC} Starting database backup..."
log "   Database: $DB_NAME"
log "   Host: $DB_HOST:$DB_PORT"
log "   User: $DB_USER"

# Run pg_dump
if pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    --no-password \
    --verbose \
    --clean \
    --if-exists \
    --create \
    > "$BACKUP_FILE" 2>> "$LOG_FILE"; then
    
    # Compress backup
    if gzip "$BACKUP_FILE"; then
        BACKUP_SIZE=$(du -h "$BACKUP_FILE_COMPRESSED" | cut -f1)
        log "${GREEN}✓${NC} Backup completed successfully: $BACKUP_FILE_COMPRESSED ($BACKUP_SIZE)"
        
        # Clean up old backups (keep last 100 backups)
        cd "$BACKUP_DIR"
        BACKUP_COUNT=$(ls -1 backup_*.sql.gz 2>/dev/null | wc -l)
        
        if [ "$BACKUP_COUNT" -gt 100 ]; then
            log "${YELLOW}→${NC} Cleaning up old backups (keeping last 100)..."
            ls -1t backup_*.sql.gz 2>/dev/null | tail -n +101 | xargs rm -f
            log "${GREEN}✓${NC} Cleanup completed"
        fi
        
        exit 0
    else
        log "${RED}✗${NC} Failed to compress backup"
        exit 1
    fi
else
    log "${RED}✗${NC} Backup failed. Check log file: $LOG_FILE"
    rm -f "$BACKUP_FILE"
    exit 1
fi

