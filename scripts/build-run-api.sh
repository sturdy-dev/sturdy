#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pushd "${CWD}/../api"

go run -v ./cmd/api --http.addr 127.0.0.1:3000 --analytics.disable
