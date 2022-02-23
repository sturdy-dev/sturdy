#!/usr/bin/env bash

set -euo pipefail
set -x

STURDY_SYNC_VERSION=""

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

source build-common.sh

ARCH=amd64
OS="darwin"

if [[ $(uname -m) == "arm64" ]]; then
  ARCH=arm64
fi

if [[ $(uname) == "Linux" ]]; then
  OS="linux"
fi

download_sturdy_sync "$OS" "$ARCH"
