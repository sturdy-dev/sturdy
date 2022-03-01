#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pushd "${CWD}/../web"

yarn install
yarn codegen
yarn dev