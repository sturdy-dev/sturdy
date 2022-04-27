#!/usr/bin/env bash

set -euo pipefail

HAS_ERROR=0

pushd api
go run ./cmd/api --help 1>/dev/null 2>/dev/null && echo "oss: OK" || {
	HAS_ERROR=1
	echo "oss: FAIL"
}

go run -tags enterprise ./cmd/api --help 1>/dev/null 2>/dev/null && echo "enterprise: OK" || {
	HAS_ERROR=1
	echo "enterprise: FAIL"
}
go run -tags cloud ./cmd/api --help 1>/dev/null 2>/dev/null && echo "cloud: OK" || {
	HAS_ERROR=1
	echo "cloud: FAIL"
}
popd

if [[ $HAS_ERROR -eq 1 ]]; then
	exit 1
fi
