#!/bin/bash
set -e

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-penfeel}
DB_SSLMODE=${DB_SSLMODE:-disable}
OUTPUT_FILE=${OUTPUT_FILE:-"./schema.sql"}

# Build connection string
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Dumping database schema to ${OUTPUT_FILE}..."

# Use pg_dump to create a schema-only dump
pg_dump --schema-only --no-owner --no-acl -d "${DB_URL}" -f "${OUTPUT_FILE}"

# If there are enum tables or other reference data to include
if [ "$INCLUDE_DATA" = "true" ]; then
  echo "Including reference data in dump..."
  
  # Get a list of tables ending with _enum (adjust this query as needed)
  TABLES=$(psql "${DB_URL}" -t -c "SELECT tablename FROM pg_tables WHERE tablename LIKE '%_enum' AND schemaname = 'public';")
  
  for TABLE in $TABLES; do
    echo "Adding data from table: $TABLE"
    # Append the data as INSERT statements
    pg_dump --data-only --no-owner --no-acl -d "${DB_URL}" -t "public.${TABLE}" >> "${OUTPUT_FILE}"
  done
fi

echo "Schema dump completed successfully!" 