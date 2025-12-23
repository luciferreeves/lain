#!/usr/bin/env bash
set -e

DEFAULT_DB="lain"
DEFAULT_OWNER="$(whoami)"
DEFAULT_HOST="localhost"
DEFAULT_PORT="5432"

read -p "Enter database name [$DEFAULT_DB]: " DB_NAME
DB_NAME="${DB_NAME:-$DEFAULT_DB}"

read -p "Enter database owner username [$DEFAULT_OWNER]: " DB_OWNER
DB_OWNER="${DB_OWNER:-$DEFAULT_OWNER}"

read -s -p "Enter password for $DB_OWNER (optional): " DB_PASSWORD
echo

# ---------- connection strategy ----------

PSQL_BASE="psql -d postgres -U $DB_OWNER"

if [ -n "$DB_PASSWORD" ]; then
    export PGPASSWORD="$DB_PASSWORD"
fi

# Prefer TCP (Docker-safe)
PSQL_TCP="$PSQL_BASE -h $DEFAULT_HOST -p $DEFAULT_PORT"
# Fallback: local socket / peer-auth
PSQL_LOCAL="$PSQL_BASE"

echo
echo "Detecting PostgreSQL connection mode..."

if $PSQL_TCP -c "SELECT 1" >/dev/null 2>&1; then
    PSQL_CMD="$PSQL_TCP"
    CONNECTION_INFO="TCP ($DEFAULT_HOST:$DEFAULT_PORT)"
elif $PSQL_LOCAL -c "SELECT 1" >/dev/null 2>&1; then
    PSQL_CMD="$PSQL_LOCAL"
    CONNECTION_INFO="Local socket (peer / trust)"
else
    echo "❌ Unable to connect to PostgreSQL."
    echo "Tried:"
    echo "  - TCP:      $DEFAULT_HOST:$DEFAULT_PORT"
    echo "  - Local:    UNIX socket"
    echo
    echo "Fix one of the following:"
    echo "  • Docker Postgres not running"
    echo "  • Wrong user ($DB_OWNER)"
    echo "  • Password required but not provided"
    echo "  • pg_hba.conf blocks access"
    exit 1
fi

echo "Connected via: $CONNECTION_INFO"

# ---------- existence check ----------

DB_EXISTS=$($PSQL_CMD -Atc \
  "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';" \
  || true)

if [ "$DB_EXISTS" = "1" ]; then
    ACTION="reset"
else
    ACTION="create"
fi

echo
echo "=== Database $ACTION ==="
echo "Database: $DB_NAME"
echo "Owner:    $DB_OWNER"
echo "Auth:     $CONNECTION_INFO"
echo

read -p "Proceed? [Y/n]: " CONFIRM
if [[ "$CONFIRM" =~ ^[Nn]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# ---------- drop + create ----------

echo
echo "Applying changes..."

$PSQL_CMD -c "DROP DATABASE IF EXISTS \"$DB_NAME\";"
$PSQL_CMD -c "CREATE DATABASE \"$DB_NAME\" OWNER \"$DB_OWNER\";"

echo
echo "✅ Database '$DB_NAME' ready."

# ---------- DSN ----------

if [ -n "$DB_PASSWORD" ]; then
    DSN="postgresql://$DB_OWNER:$DB_PASSWORD@$DEFAULT_HOST:$DEFAULT_PORT/$DB_NAME?sslmode=disable"
else
    DSN="postgresql://$DB_OWNER@$DEFAULT_HOST:$DEFAULT_PORT/$DB_NAME?sslmode=disable"
fi

echo
echo "DSN:"
echo "$DSN"

unset PGPASSWORD
