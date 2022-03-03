#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUN_INSTALL=1

while [[ $# -gt 0 ]]; do
  case "$1" in
  --skip-install)
    RUN_INSTALL=0
    shift
    ;;
  esac
done

pushd "${CWD}/../app"

if ((RUN_INSTALL)); then
  yarn install
  ./local-dev-download-deps.sh
fi

yarn dev
