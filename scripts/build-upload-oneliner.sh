#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

VERSION=$(date +%Y-%m-%d-%H-%M-%S)
IMAGE="getsturdy/server"

docker buildx build \
  --platform linux/arm64,linux/amd64 \
  --target oneliner \
  --build-arg API_BUILD=enterprise \
  --tag "${IMAGE}:latest" \
  --tag "${IMAGE}:${VERSION}" \
  --push \
  "$CWD/.."
