#!/bin/bash
# SIS-INV: Reset Admin Role & Password Wrapper

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔄 Resetting administrator account..."

cd "$PROJECT_ROOT/backend" && go run cmd/reset-admin/main.go
