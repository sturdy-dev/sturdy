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
LATEST_VERSION="$(curl -s https://registry.hub.docker.com/v1/repositories/getsturdy/server/tags | jq -r '.[].name' | grep -v cache | tail -1)"
VERSION=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	--major)
		MAJOR="$(($(echo "$LATEST_VERSION" | cut -d. -f1) + 1))"
		MINOR="0"
		PATCH="0"
		VERSION="$MAJOR.$MINOR.$PATCH"
		DOCKER_VERSION_TAG_ARG="--tag ${IMAGE}:${VERSION}"
		shift
		;;
	--minor)
		MAJOR="$(echo "$LATEST_VERSION" | cut -d. -f1)"
		MINOR="$(($(echo "$LATEST_VERSION" | cut -d. -f2) + 1))"
		PATCH="0"
		VERSION="$MAJOR.$MINOR.$PATCH"
		DOCKER_VERSION_TAG_ARG="--tag ${IMAGE}:${VERSION}"
		shift
		;;
	--patch)
		MAJOR="$(echo "$LATEST_VERSION" | cut -d. -f1)"
		MINOR="$(echo "$LATEST_VERSION" | cut -d. -f2)"
		PATCH="$(($(echo "$LATEST_VERSION" | cut -d. -f3) + 1))"
		VERSION="$MAJOR.$MINOR.$PATCH"
		DOCKER_VERSION_TAG_ARG="--tag ${IMAGE}:${VERSION}"
		shift
		;;
	--version)
		VERSION="$2"
		DOCKER_VERSION_TAG_ARG="--tag ${IMAGE}:${VERSION}"
		shift
		shift
		;;
	--image)
		IMAGE="$2"
		shift
		shift
		;;
	--push)
		PUSH_ARG="--push"
		shift
		;;
	esac
done

if [[ -z "$VERSION" ]]; then
	echo "version number is not set"
	echo "you can set it via --version or --major/--minor/--patch"
	exit 1
fi

echo "image: ${IMAGE}"
echo "version: ${VERSION}"
sleep 1
echo

docker buildx build \
	--platform linux/arm64,linux/amd64 \
	--target oneliner \
	--cache-to=getsturdy/server:cache \
	--cache-from=getsturdy/server:cache \
	--build-arg API_BUILD_TAGS=enterprise \
	--build-arg VERSION="${VERSION}" \
	--tag "${IMAGE}:latest" \
	${DOCKER_VERSION_TAG_ARG} \
	${PUSH_ARG} \
	"$CWD/.."
