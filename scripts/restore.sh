#!/bin/bash
# SIS-INV: Manual PostgreSQL Restore Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_FILE=$1

# Load .env from project root
if [ -f "$PROJECT_ROOT/.env" ]; then
    export $(grep -v '^#' "$PROJECT_ROOT/.env" | xargs)
fi

# Database credentials
DB_USER=${DB_USER:-"sekolah_admin"}
DB_NAME=${DB_NAME:-"inventaris_db"}
DB_PASS=${DB_PASSWORD:-"sekolah_secret_2024"}
CONTAINER_NAME="sisinv_db"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: ./restore.sh <path_to_backup_file.sql[.gz]>"
    exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Error: File '$BACKUP_FILE' not found."
    exit 1
fi

# Check if docker is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Error: Docker container '$CONTAINER_NAME' is not running."
    exit 1
fi

echo "⚠️  WARNING: This will overwrite current data in '$DB_NAME'."
read -p "Are you sure you want to proceed? (y/N) " confirm
if [[ $confirm != [yY] ]]; then
    echo "Restore cancelled."
    exit 0
fi

echo "[$(date)] 🔄 Starting database restore..."

# Decompress if needed
if [[ "$BACKUP_FILE" == *.gz ]]; then
    gunzip -c "$BACKUP_FILE" | docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME"
else
    docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" "$DB_NAME" < "$BACKUP_FILE"
fi

if [ $? -eq 0 ]; then
    echo "✅ Restore completed successfully."
else
    echo "❌ FAILED to restore database."
    exit 1
fi
