#!/bin/bash
# SIS-INV: Persistent Full Database Backup (Manual Archive)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="$PROJECT_ROOT/backups/archives"
DATE=$(date +%Y%m%d_%H%M%S)

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

echo "[$(date)] 🗄️  Starting PERSISTENT full backup..."
echo "⚠️  Note: This backup will NOT be automatically deleted."

# Check if docker is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Error: Docker container '$CONTAINER_NAME' is not running."
    exit 1
fi

BACKUP_FILE="$BACKUP_DIR/full_archive_$DATE.sql"

# Run backup inside docker and stream to host file
docker exec -e PGPASSWORD="$DB_PASS" "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ] && [ -s "$BACKUP_FILE" ]; then
    gzip "$BACKUP_FILE"
    echo "✅ Persistent backup completed!"
    echo "📍 Location: backups/archives/full_archive_$DATE.sql.gz"
else
    echo "❌ FAILED to create persistent backup."
    rm -f "$BACKUP_FILE"
    exit 1
fi
