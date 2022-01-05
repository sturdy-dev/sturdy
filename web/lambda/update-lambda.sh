#!/bin/bash

set -x

FUNCTION_NAME="cloudfront-edge-redirect"
DISTRIBUTION_ID="E1S00JE3I1WHG2"

while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --production)
      FUNCTION_NAME="cloudfront-index-redirect"
      DISTRIBUTION_ID="ETTHKSFEKJVNO"
      shift
      ;;
    *)
      # unknown option
      echo "Unknown argument: $1"
      exit 1;
      ;;
  esac
done

# Build
zip cloudfront-index-redirect.zip index.js

# Upload
UPDATE_CODE_RES=$(AWS_REGION=us-east-1 aws lambda update-function-code --function-name "${FUNCTION_NAME}" --zip-file fileb://cloudfront-index-redirect.zip)
CODE_SHA=$(echo "$UPDATE_CODE_RES" | jq --raw-output '.CodeSha256')

sleep 5

# Tag a new version
PUBLISH_VERSION_RES=$(AWS_REGION=us-east-1 aws lambda publish-version --function-name "${FUNCTION_NAME}" --code-sha256 "${CODE_SHA}")
FUNCTION_ARN=$(echo "$PUBLISH_VERSION_RES" | jq --raw-output '.FunctionArn')

sleep 5

# Get the existing configuration
aws cloudfront get-distribution-config --id "${DISTRIBUTION_ID}" > distribution.json
ETAG=$(jq --raw-output '.ETag' distribution.json)

# Update CloudFront
jq ".DistributionConfig.DefaultCacheBehavior.LambdaFunctionAssociations.Items[0].LambdaFunctionARN = \"${FUNCTION_ARN}\" | .DistributionConfig" distribution.json > config.json
aws cloudfront update-distribution --distribution-config file://config.json --id "${DISTRIBUTION_ID}" --if-match "${ETAG}"