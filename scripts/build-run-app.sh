#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pushd "${CWD}/../app"

yarn install
./local-dev-download-deps.sh
yarn dev