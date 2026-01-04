#!/usr/bin/env bash

LIMIT="*"
if [[ -n "$1" ]]
then
  LIMIT="$1"
fi

set -euo pipefail

cd "$(dirname "$0")/.."

export MODE_TEST=1

BASE_DIR="$(pwd)"

# cd "$BASE_DIR"
go run gotest.tools/gotestsum@latest --format pkgname ./...
