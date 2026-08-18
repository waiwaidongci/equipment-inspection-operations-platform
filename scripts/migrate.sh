#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${POSTGRES_DSN:-}" ]]; then
  echo "POSTGRES_DSN is required" >&2
  exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
for file in "$ROOT_DIR"/migrations/*.sql; do
  echo "Applying $file"
  psql "$POSTGRES_DSN" -v ON_ERROR_STOP=1 -f "$file"
done
