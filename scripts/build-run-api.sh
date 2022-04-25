#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pushd "${CWD}/../api"

TAGS=""
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --enterprise)
      TAGS="--tags enterprise"
      shift
      ;;
    *)
      # unknown option
      echo "Unknown argument: $1"
      exit 1;
      ;;
  esac
done

go run -v ${TAGS} ./cmd/api --http.addr 127.0.0.1:3000 --analytics.disable
