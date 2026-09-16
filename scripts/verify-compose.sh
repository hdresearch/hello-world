#!/bin/sh
set -eu

command -v curl >/dev/null 2>&1 || { echo "curl is required" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is required" >&2; exit 1; }

base_url="${HELLO_WORLD_URL:-http://127.0.0.1:8080}"
attempt=0
until curl --fail --silent --show-error "$base_url/healthz" >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge 60 ]; then
    echo "Frontend did not become healthy at $base_url" >&2
    exit 1
  fi
  sleep 1
done

before=$(curl --fail --silent --show-error "$base_url/api/visits" | jq -er '.count')
after=$(curl --fail --silent --show-error --request POST "$base_url/api/visits" | jq -er '.count')
final=$(curl --fail --silent --show-error "$base_url/api/visits" | jq -er '.count')

if [ "$after" -ne $((before + 1)) ] || [ "$final" -ne "$after" ]; then
  echo "Visit counter verification failed: before=$before after=$after final=$final" >&2
  exit 1
fi

echo "Hello World stack verified: $before -> $after"
