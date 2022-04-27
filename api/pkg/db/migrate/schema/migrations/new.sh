#!/usr/bin/env bash

set -euo pipefail

# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
migrate create \
    -ext sql \
    -dir "$(dirname $0)" \
    -seq "$1"
