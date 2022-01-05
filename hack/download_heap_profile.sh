#!/usr/bin/env bash

set -euo pipefail

RFC3339_DATE="$(date +%Y-%m-%dT%H:%m:%SZ%z)"
OUT_FILE="${RFC3339_DATE}.heap.out"
SECONDS="1"

echo "Recording ${SECONDS}s heap profile..."

curl \
    --output "${OUT_FILE}" \
    "http://localhost:6060/debug/pprof/heap?seconds=${SECONDS}"

echo "${OUT_FILE}"
