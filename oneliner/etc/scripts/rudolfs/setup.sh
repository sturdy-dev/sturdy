#!/bin/bash

set -euo pipefail

log() {
	echo "rudolfs-prepare $@"
}

DIR="/var/data/rudolfs"
KEY="$DIR/key"

if [[ ! -f "$KEY" ]]; then
	log "Generating key..."
	mkdir -p "$DIR"
	openssl rand -hex 32 >"$DIR/key"
fi
