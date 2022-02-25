#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker buildx build --load \
	--target oneliner \
	--build-arg API_BUILD_TAGS=enterprise \
	--build-arg VERSION="development" \
	--tag "sturdy-oneliner:oss" \
	"$CWD/.."

docker run --interactive \
    --pull never \
    --publish 30080:80 \
    --volume "$HOME/.sturdydata-development:/var/data" \
    "sturdy-oneliner:oss"