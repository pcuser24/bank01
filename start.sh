#!/bin/sh

# Exit immediately if a simple command exits with a non-zero status
set -e

# List environment variables
echo "List environment variables"
env

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "run the app"
# Example: entrypoint.sh server start -> server start
exec "$@"
