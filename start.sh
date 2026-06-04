#!/bin/sh

set -e

echo "run db migration"
# This command was used migrate binary tool. The code updated to use migrate of golang
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"