#!/bin/bash
# SIS-INV: Manual PostgreSQL Backup Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="$PROJECT_ROOT/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Load .env from project root
if [ -f "$PROJECT_ROOT/.env" ]; then
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# Database credentials
DB_USER=${DB_USER:-"sekolah_admin"}
DB_NAME=${DB_NAME:-"inventaris_db"}
DB_PASS=${DB_PASSWORD:-"sekolah_secret_2024"}
CONTAINER_NAME="sisinv_db"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] 📦 Starting database backup..."

# Check if docker is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Error: Docker container '$CONTAINER_NAME' is not running."
    exit 1
fi

BACKUP_FILE="$BACKUP_DIR/db_$DATE.sql"

# Run backup inside docker and stream to host file
docker exec -e PGPASSWORD="$DB_PASS" "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ] && [ -s "$BACKUP_FILE" ]; then
    gzip "$BACKUP_FILE"
    echo "✅ Backup completed: backups/db_$DATE.sql.gz"
    
    # Retention
    find "$BACKUP_DIR" -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
else
    echo "❌ FAILED to create database backup."
    rm -f "$BACKUP_FILE"
    exit 1
fi

