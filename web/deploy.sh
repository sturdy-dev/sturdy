#!/usr/bin/env bash

set -euo pipefail

BUCKET="gustav-staging.driva.dev"
DISTRIBUTION_ID="E1S00JE3I1WHG2"
BUILD=1
TEST=1

POSITIONAL=()
while [[ $# -gt 0 ]]; do
	key="$1"
	case $key in
	--production)
		BUCKET="getsturdy.com"
		DISTRIBUTION_ID="ETTHKSFEKJVNO"
		shift
		;;
	--no-build)
		BUILD=0
		shift
		;;
	--no-test)
		TEST=0
		shift
		;;
	*)
		# unknown option
		echo "Unknown argument: $1"
		exit 1
		;;
	esac
done

if ((BUILD)); then
	yarn install
	yarn codegen
fi

if ((TEST)); then
	yarn run lint
	yarn run test
fi

if ((BUILD)); then
	yarn build:client
	yarn build:prerender
fi

aws s3 sync ./dist/static "s3://${BUCKET}" --cache-control max-age=600
aws s3 sync ./dist/client/assets "s3://${BUCKET}/assets" --cache-control max-age=600
aws s3 cp ./dist/client/client-side-render.html "s3://${BUCKET}/client-side-render.html" --cache-control max-age=600
aws s3 cp ./robots.txt "s3://${BUCKET}/robots.txt" --cache-control max-age=600
aws cloudfront create-invalidation --distribution-id "${DISTRIBUTION_ID}" --paths "/*"

# Mark deployment in Sentry
DATE=$(date '+%Y-%m-%d %H:%M:%S')
VERSION="${DATE} by ${USER}"
curl https://sentry.io/api/hooks/release/builtin/5901793/509f0445006ff8bffb976f46ea4b61c0ac618a1a6ed64ae4a6a833e520c4638b/ \
	-X POST \
	-H 'Content-Type: application/json' \
	-d "{\"version\":\"${VERSION}\",\"shortVersion\":\"${DATE}\"}"

echo
echo "Done!"
echo "Release Name: ${VERSION}"
echo
