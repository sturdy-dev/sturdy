#!/usr/bin/env bash

set -euo pipefail

HAS_ERROR=0

pushd api
go run ./cmd/api --help 1>/dev/null 2>/dev/null || {
	echo "failed to build oss"
	HAS_ERROR=1
}
go run -tags enterprise ./cmd/api --help 1>/dev/null 2>/dev/null || {
	echo "failed to build enterprise"
	HAS_ERROR=1
}
go run -tags cloud ./cmd/api --help 1>/dev/null 2>/dev/null || {
	echo "failed to build cloud"
	HAS_ERROR=1
}
popd

if [[ $HAS_ERROR -eq 1 ]]; then
	exit 1
else
	echo 'All good!'
fi
