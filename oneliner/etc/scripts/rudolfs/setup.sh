#!/usr/bin/env sh

set -euo pipefail

DIR="/var/data/rudolfs"
KEY="$DIR/key"

if [[ ! -f "$KEY" ]]; then
  echo "Generating key..."
  mkdir -p "$DIR"
  openssl rand -hex 32 >"$DIR/key"
fi
