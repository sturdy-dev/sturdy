#!/usr/bin/env bash

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

pushd "$CWD/../api"

EXIT_CODE=0

OSS_ENTERPISE_DEPENDENCIES=$(go list -f '{{ join .Deps "\n"}}' ./cmd/api | grep getsturdy | grep enterprise)
if [ -n "${OSS_ENTERPISE_DEPENDENCIES}" ]; then
	echo "ERROR: getsturdy.com/api/cmd/api links Enterprise modules:"
	echo "${OSS_ENTERPISE_DEPENDENCIES}"
	EXIT_CODE=1
fi

OSS_CLOUD_DEPENDENCIES=$(go list -f '{{ join .Deps "\n"}}' getsturdy.com/api/cmd/api | grep getsturdy | grep cloud)
if [ -n "${OSS_CLOUD_DEPENDENCIES}" ]; then
	echo "ERROR: getsturdy.com/api/cmd/api links Cloud modules:"
	echo "${OSS_CLOUD_DEPENDENCIES}"
	EXIT_CODE=1
fi

ENTERPRISE_CLOUD_DEPENDENCIES=$(go list -f '{{ join .Deps "\n"}}' -tags enterprise getsturdy.com/api/cmd/api | grep getsturdy | grep cloud)
if [ -n "${ENTERPRISE_CLOUD_DEPENDENCIES}" ]; then
	echo "ERROR: getsturdy.com/api/cmd/api -tags enterprise links Cloud modules:"
	echo "${ENTERPRISE_CLOUD_DEPENDENCIES}"
	EXIT_CODE=1
fi

if [ $EXIT_CODE -eq 0 ]; then
	echo "All good!"
fi

exit $EXIT_CODE
