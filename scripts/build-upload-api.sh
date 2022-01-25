#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Authenticate the docker daemon to the registry
aws ecr get-login-password --region eu-north-1 |
  docker login --username AWS --password-stdin 902160009014.dkr.ecr.eu-north-1.amazonaws.com

VERSION=$(date +%Y-%m-%d-%H-%M-%S)
ECR_NAME="902160009014.dkr.ecr.eu-north-1.amazonaws.com/api:${VERSION}"

docker buildx build \
  --platform linux/amd64 \
  --tag "$ECR_NAME" \
  --build-arg API_BUILD=cloud \
  --target api \
  --push \
  "$CWD/.."

echo "Successfully built and pushed $ECR_NAME :-)"
