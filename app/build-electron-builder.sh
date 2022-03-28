#!/usr/bin/env bash

set -euo pipefail

DO_UPLOAD=0
CODESIGN=1
NOTARIZE=1
CHANNEL="Beta"
DO_BUILD=1

ELECTRON_BUILDER_ARCHS=""
ELECTRON_BUILDER_PLATFORMS=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	--no-build)
		DO_BUILD=0
		shift
		;;
	--arm64)
		ELECTRON_BUILDER_ARCHS="$ELECTRON_BUILDER_ARCHS --arm64"
		shift
		;;
	--amd64)
		ELECTRON_BUILDER_ARCHS="$ELECTRON_BUILDER_ARCHS --x64"
		shift
		;;
	--windows)
		ELECTRON_BUILDER_PLATFORMS="$ELECTRON_BUILDER_PLATFORMS --windows"
		shift
		;;
	--linux)
		ELECTRON_BUILDER_PLATFORMS="$ELECTRON_BUILDER_PLATFORMS --linux"
		shift
		;;
	--mac)
		ELECTRON_BUILDER_PLATFORMS="$ELECTRON_BUILDER_PLATFORMS --mac"
		shift
		;;
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

if [[ -z "$ELECTRON_BUILDER_ARCHS" ]]; then
	# default amd64 and arm64
	ELECTRON_BUILDER_ARCHS="--x64 --arm64"
fi

if [[ -z "$ELECTRON_BUILDER_PLATFORMS" ]]; then
	# default all platforms
	ELECTRON_BUILDER_PLATFORMS="--linux --windows --mac"
fi

source build-common.sh

if ((DO_BUILD)); then
	yarn build
fi

package "$ELECTRON_BUILDER_ARCHS" "$ELECTRON_BUILDER_PLATFORMS"
