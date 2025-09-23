#!/usr/bin/env bash
set -euo pipefail

if command -v wine >/dev/null 2>&1; then
  echo "wine: $(wine --version)"
  exit 0
fi

echo "wine not found. Install it (e.g., sudo apt-get install wine)." >&2
exit 1

