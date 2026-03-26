#!/usr/bin/env bash
set -euo pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

cd "$ROOT"

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is required"
  exit 1
fi

go run -mod=mod entgo.io/ent/cmd/ent generate ./pkg/ent/schema

go run -mod=mod entgo.io/ent/cmd/ent migrate --dev-url "$DATABASE_URL" --dir ./pkg/ent/migrate
