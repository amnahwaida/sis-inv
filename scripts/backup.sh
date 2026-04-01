#!/bin/bash
# SIS-INV: Automated PostgreSQL Backup Script
# Based on PRD F06.3

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="$PROJECT_ROOT/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Load .env from project root if it exists
if [ -f "$PROJECT_ROOT/.env" ]; then
    # Simple export for common .env formats
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# Database credentials (use env vars, with defaults)
DB_USER=${DB_USER:-"sekolah_admin"}
DB_NAME=${DB_NAME:-"inventaris_db"}
DB_PASS=${DB_PASSWORD:-"sekolah_secret_2024"}
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Set password for pg_dump to avoid prompt
export PGPASSWORD=$DB_PASS

echo "[$(date)] Starting SIS-INV database backup..."

# Database backup to SQL file
BACKUP_FILE="$BACKUP_DIR/db_$DATE.sql"
docker exec -e PGPASSWORD=$DB_PASS sisinv_db pg_dump -U $DB_USER $DB_NAME > $BACKUP_FILE

if [ $? -eq 0 ]; then
    # Compress backup
    gzip $BACKUP_FILE
    echo "[$(date)] Database backup completed: db_$DATE.sql.gz"
else
    echo "[$(date)] ❌ FAILED to create database backup"
    exit 1
fi

# Clean up old backups (older than 30 days)
echo "[$(date)] Cleaning up old backups (>${RETENTION_DAYS} days)..."
find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Unset password
unset PGPASSWORD

echo "[$(date)] SIS-INV Backup Process Finished Successfully."
echo "[$(date)] ---" >> $BACKUP_DIR/backup.log
