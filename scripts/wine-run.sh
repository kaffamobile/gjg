#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <exe> [args...]" >&2
  exit 2
fi

DIR=$(cd "$(dirname "$0")" && pwd)
"$DIR/wine-check.sh" >/dev/null

export WINEDEBUG="-all"
exe=$1; shift || true

exec wine "$exe" "$@"

