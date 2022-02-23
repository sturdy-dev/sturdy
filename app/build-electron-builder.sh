#!/usr/bin/env bash

set -euo pipefail

DO_UPLOAD=0
STURDY_SYNC_VERSION="v0.9.0"
CODESIGN=1
NOTARIZE=1
CHANNEL="Beta"

while [[ $# -gt 0 ]]; do
	case "$1" in
	--sturdy-sync-version)
		STURDY_SYNC_VERSION="$2"
		shift
		shift
		;;
	--upload)
		DO_UPLOAD=1
		shift
		;;
	--no-codesign)
		CODESIGN=0
		shift
		;;
	--no-notarize)
		NOTARIZE=0
		shift
		;;
	--stable)
		CHANNEL=""
		shift
		;;
	esac
done

source build-common.sh

if [ "$STURDY_SYNC_VERSION" == "" ]; then
	echoerr "--sturdy-sync-version is not set!"
	exit 1
fi

if ((NOTARIZE)); then
	setup_darwin_notarize
fi

build darwin amd64
build darwin arm64

build windows amd64 zip

build linux amd64
#build linux arm64

if ((DO_UPLOAD)); then
	invalidate_cloudfront "${CHANNEL}"
fi
