#!/usr/bin/env bash

set -euo pipefail

# Authenticate the docker daemon to the registry
aws ecr get-login-password --region eu-north-1 \
  | docker login --username AWS --password-stdin 902160009014.dkr.ecr.eu-north-1.amazonaws.com

# Build in Docker to link and ship against a specific version of libgit2
# Cross compilation is not possible with CGO

VERSION=$(date +%Y-%m-%d-%H-%M-%S);
ECR_NAME="902160009014.dkr.ecr.eu-north-1.amazonaws.com/api:${VERSION}"

docker buildx build \
    --platform linux/amd64 \
    --tag "$ECR_NAME" \
    --push \
    .

echo "Successfully built and pushed $ECR_NAME :-)"
