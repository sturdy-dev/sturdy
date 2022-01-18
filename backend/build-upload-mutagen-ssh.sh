#!/usr/bin/env bash

set -euo pipefail

mkdir -p tmp

# clone mutagen
[ ! -d "tmp/mutagen/.git" ] && git clone git@github.com:sturdy-dev/mutagen.git tmp/mutagen
git -C "tmp/mutagen" fetch

# Authenticate the docker daemon to the registry
aws ecr get-login-password --region eu-north-1 |
  docker login --username AWS --password-stdin 902160009014.dkr.ecr.eu-north-1.amazonaws.com

VERSION=$(date +%Y-%m-%d-%H-%M-%S)
ECR_NAME="902160009014.dkr.ecr.eu-north-1.amazonaws.com/mutagen-ssh:${VERSION}"

docker buildx build \
  -f Dockerfile.mutagen-ssh \
  --platform linux/amd64 \
  --tag "$ECR_NAME" \
  --push \
  .

echo "Successfully built and pushed $ECR_NAME :-)"
