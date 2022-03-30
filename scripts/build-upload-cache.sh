#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

IMAGE="getsturdy/server"

docker buildx build \
	--platform linux/arm64,linux/amd64 \
	--target libgit-builder \
	--cache-to=getsturdy/server:cache \
	--cache-from=getsturdy/server:cache \
	--tag "${IMAGE}:latest" \
	"$CWD/.."
