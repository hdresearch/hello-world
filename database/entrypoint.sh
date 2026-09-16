#!/bin/sh
set -eu

: "${BACKEND_TOKEN:?BACKEND_TOKEN is required}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD is required}"
export POSTGRES_DB="${POSTGRES_DB:-webstack}"
export POSTGRES_USER="${POSTGRES_USER:-postgres}"

postgres_pid=""
gateway_pid=""

shutdown() {
  trap - INT TERM
  [ -z "$gateway_pid" ] || kill -TERM "$gateway_pid" 2>/dev/null || true
  [ -z "$postgres_pid" ] || kill -TERM "$postgres_pid" 2>/dev/null || true
  [ -z "$gateway_pid" ] || wait "$gateway_pid" 2>/dev/null || true
  [ -z "$postgres_pid" ] || wait "$postgres_pid" 2>/dev/null || true
  exit 0
}
trap shutdown INT TERM

docker-entrypoint.sh postgres -c listen_addresses=127.0.0.1 &
postgres_pid=$!

attempt=0
until pg_isready --host 127.0.0.1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge 60 ] || ! kill -0 "$postgres_pid" 2>/dev/null; then
    echo "PostgreSQL did not become ready" >&2
    wait "$postgres_pid" || true
    exit 1
  fi
  sleep 1
done

PGPASSWORD="$POSTGRES_PASSWORD" psql \
  --host 127.0.0.1 \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  --set ON_ERROR_STOP=1 \
  --file /opt/hello-world/migrations/001_visits.sql

gosu postgres /usr/local/bin/hello-world-gateway &
gateway_pid=$!

while kill -0 "$postgres_pid" 2>/dev/null && kill -0 "$gateway_pid" 2>/dev/null; do
  sleep 1
done

echo "A required database service exited" >&2
kill -TERM "$gateway_pid" "$postgres_pid" 2>/dev/null || true
wait "$gateway_pid" 2>/dev/null || true
wait "$postgres_pid" 2>/dev/null || true
exit 1
