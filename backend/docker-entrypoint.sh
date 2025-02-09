#!/bin/sh
set -e

echo ""
echo "Running database migrations..."
export GOOSE_DRIVER=postgres
export GOOSE_MIGRATION_DIR=/app/migrations
export GOOSE_DBSTRING=$DATABASE_URL
goose up

echo ""
echo "Starting application..."
exec /app/main
