#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker buildx build \
  --platform linux/amd64 \
  --target oneliner \
  --build-arg API_BUILD=enterprise \
  "$CWD/.."

# --tag "$ECR_NAME" \
# --push \
