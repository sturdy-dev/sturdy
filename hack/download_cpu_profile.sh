#!/usr/bin/env bash

set -euo pipefail

RFC3339_DATE="$(date +%Y-%m-%dT%H:%m:%SZ%z)"
OUT_FILE="${RFC3339_DATE}.cpu.out"
SECONDS="30"

echo "Recording ${SECONDS}s cpu profile..."

curl \
  --output "${OUT_FILE}" \
  "http://localhost:6060/debug/pprof/profile?seconds=${SECONDS}"

echo "${OUT_FILE}"
