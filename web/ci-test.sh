#!/bin/bash

set -euo pipefail

# Symlink node_modules from the cached docker image
rm -rf ./node_modules || true
ln -s /worker/node_modules ./node_modules

echo "--- yarn codegen"
DEBUG=1 yarn codegen

echo "--- yarn lint"
yarn lint

echo "--- yarn test"
yarn test
