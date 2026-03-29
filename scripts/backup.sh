#!/bin/bash

# Configuration
BACKUP_DIR="./backups"
DB_CONTAINER="sis-inv-db-1"
DB_USER="sekolah_admin"
DB_NAME="inventaris_db"
RETENTION_DAYS=7

# Date format for the backup file
DATE=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/db_backup_$DATE.sql.gz"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Starting database backup at $DATE"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Enable pipefail so if pg_dump fails, the whole pipe fails
set -o pipefail

# Try to find the exact container name
ACTUAL_CONTAINER=$(docker ps -a --format '{{.Names}}' | grep "db" | grep "sis" | head -n 1)
if [ -z "$ACTUAL_CONTAINER" ]; then
    echo -e "${RED}✗ Error: Database container not found!${NC}"
    exit 1
fi

echo "Found database container: $ACTUAL_CONTAINER"

# Run pg_dump inside the container and pipe it through gzip on the host
if docker exec "$ACTUAL_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$BACKUP_FILE"; then
    echo -e "${GREEN}✓ Backup successful! Saved to: $BACKUP_FILE${NC}"
else
    echo -e "${RED}✗ Backup failed!${NC}"
    rm -f "$BACKUP_FILE" # Clean up empty file
    exit 1
fi

# Clean up old backups strictly older than $RETENTION_DAYS days
echo "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -type f -name "db_backup_*.sql.gz" -mtime +$RETENTION_DAYS -exec rm {} \;

echo -e "${GREEN}✓ Cleanup complete.${NC}"
exit 0
