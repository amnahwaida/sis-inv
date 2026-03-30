#!/bin/bash
# SIS-INV: Automated PostgreSQL Backup Script
# Based on PRD F06.3

# Configuration
PROJECT_ROOT="/home/vannyezha/project/sis_inv"
BACKUP_DIR="$PROJECT_ROOT/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Database credentials (from .env)
DB_USER="sekolah_admin"
DB_NAME="inventaris_db"
DB_PASS="sekolah_secret_2024"
DB_HOST="localhost"
DB_PORT="5432"

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Set password for pg_dump to avoid prompt
export PGPASSWORD=$DB_PASS

echo "[$(date)] Starting SIS-INV database backup..."

# Database backup to SQL file
BACKUP_FILE="$BACKUP_DIR/db_$DATE.sql"
docker exec sisinv_db pg_dump -U $DB_USER $DB_NAME > $BACKUP_FILE

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
