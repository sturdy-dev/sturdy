#!/bin/bash

set -euo pipefail

yarn install
yarn spectaql spectaql.yml
aws s3 sync ./public/ s3://schema.getsturdy.com/
# aws s3 cp ../pkg/graphql/schema.graphql s3://schema.getsturdy.com/sturdy.graphql
aws cloudfront create-invalidation --distribution-id "E1GIL34XU5VF0F" --paths "/*"
