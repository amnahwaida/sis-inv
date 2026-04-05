#!/bin/bash
# SIS-INV: Reset Admin Role & Password Wrapper

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONTAINER_NAME="sisinv_backend"

echo "🔄 Resetting administrator account..."

# 1. Try running inside Docker container first (Recommended)
if docker ps | grep -q "$CONTAINER_NAME"; then
    echo "🐳 Running inside Docker container: $CONTAINER_NAME"
    docker exec -it "$CONTAINER_NAME" ./reset-admin
    exit $?
fi

# 2. Fallback to local Go execution if Docker is not running or container missing
echo "⚠️  Docker container '$CONTAINER_NAME' not found. Falling back to local Go execution..."
if command -v go >/dev/null 2>&1; then
    cd "$PROJECT_ROOT/backend" && go run cmd/reset-admin/main.go
else
    echo "❌ Error: Neither Docker container nor local Go installation found."
    exit 1
fi
