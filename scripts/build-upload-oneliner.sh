#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# build-upload-oneliner.sh makes a full release of Sturdy Enterprise for self-hosting.
# By default, the images is tagged as "latest" (but not pushed).
#
# Set --version to a version number (such as 1.14.0) to push an image tagged with the version number, as well as "latest".
# Set --push to push the container to Docker Hub ( https://hub.docker.com/r/getsturdy/server ).

IMAGE="getsturdy/server"
VERSION=$(date +%Y-%m-%d-%H-%M-%S)
DOCKER_VERSION_TAG_ARG=""
PUSH_ARG=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	--version)
		VERSION="$2"
		DOCKER_VERSION_TAG_ARG="--tag ${IMAGE}:${VERSION}"
		shift
		shift
		;;
	--push)
		PUSH_ARG="--push"
		shift
		;;
	esac
done

VERSION="1.0.0"

docker buildx build \
  --platform linux/arm64,linux/amd64 \
  --target oneliner \
  --build-arg API_BUILD_TAGS=enterprise \
  --build-arg VERSION="${VERSION}" \
  --tag "${IMAGE}:latest" \
  ${DOCKER_VERSION_TAG_ARG} \
  ${PUSH_ARG} \
  "$CWD/.."
