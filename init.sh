#!/bin/sh

echo "Waiting for DB to be ready at $DB_HOST:$DB_PORT..."
until nc -z "$DB_HOST" "$DB_PORT"; do
  sleep 1
done

echo "DB is ready, impoter data..."
/app/seeder -d true
/app/seeder -t pharmacy -p /app/data/pharmacies.json
/app/seeder -t user -p /app/data/users.json

echo "Starting application..."
exec /app/phantomApp
