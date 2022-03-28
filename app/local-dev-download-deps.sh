#!/usr/bin/env bash

set -euo pipefail
set -x

function download_sturdy_sync() {
	BIN_DIR=assets/bin
	rm -rf $BIN_DIR
	mkdir $BIN_DIR
	OS=$1
	ARCH=$2
	ARCHIVE_FORMAT="tar.gz"
	if [ "$OS" == "windows" ]; then
		ARCHIVE_FORMAT="zip"
	fi
	ARCHIVE_NAME="sturdy-$STURDY_SYNC_VERSION-$OS-$ARCH.$ARCHIVE_FORMAT"
	curl -s -Lo "$BIN_DIR/$ARCHIVE_NAME" "https://getsturdy.com/client/$ARCHIVE_NAME"
	if [ "$ARCHIVE_FORMAT" == "tar.gz" ]; then
		tar xzf "$BIN_DIR/$ARCHIVE_NAME" -C $BIN_DIR
	elif [ "$ARCHIVE_FORMAT" == "zip" ]; then
		unzip "$BIN_DIR/$ARCHIVE_NAME" -d $BIN_DIR/
	else
		echoerr "Unsupported archive format: $ARCHIVE_FORMAT"
		exit 1
	fi
	rm "$BIN_DIR/$ARCHIVE_NAME"

}

STURDY_SYNC_VERSION="v0.9.0"

while [[ $# -gt 0 ]]; do
	case "$1" in
	--sturdy-sync-version)
		STURDY_SYNC_VERSION="$2"
		shift
		shift
		;;
	esac
done

if [ "$STURDY_SYNC_VERSION" == "" ]; then
	echoerr "--sturdy-sync-version is not set!"
	exit 1
fi

ARCH=amd64
OS="darwin"

if [[ $(uname -m) == "arm64" ]]; then
	ARCH=arm64
fi

if [[ $(uname) == "Linux" ]]; then
	OS="linux"
fi

download_sturdy_sync "$OS" "$ARCH"
