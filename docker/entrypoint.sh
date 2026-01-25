#!/bin/sh
set -eu

mkdir -p /data

mkdir -p /app/db
touch /data/app.db

chown -R app:app /data

rm -f /app/db/app.db
ln -s /data/app.db /app/db/app.db

migrate -source file:///app/db/migrations -database sqlite3:///data/app.db up

exec su-exec app "$@"

