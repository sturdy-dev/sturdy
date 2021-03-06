#!/usr/bin/env bash

set -euo pipefail

echoerr() { echo "$@" 1>&2; }

function package() {
	BUILDER_ARGS="$@"

	APP_VERSION=$(jq --raw-output '.version' package.json)
	validate_version "$APP_VERSION"

	CODESIGN_EXTRA_ARGS=""

	if ((CODESIGN)) && [[ "$BUILDER_ARGS" =~ "--mac" ]]; then
		setup_darwin_notarize
	else
		CODESIGN_EXTRA_ARGS="-c.mac.identity=null" # Disable mac code signing
	fi

	if ((CODESIGN)) && [[ "$BUILDER_ARGS" =~ "--windows" ]]; then
		setup_windows_codesign
	fi

	PUBLISH_ARGS=""
	if ((DO_UPLOAD)); then
		PUBLISH_ARGS="--publish=always"
	fi

	yarn electron-builder $PUBLISH_ARGS $CODESIGN_EXTRA_ARGS $@

	if ((DO_UPLOAD)); then
		if [[ ! "$APP_VERSION" =~ "alpha" ]] && [[ ! "$APP_VERSION" =~ "beta" ]]; then
			create_latest "$BUILDER_ARGS"
		fi
		invalidate_cloudfront "$BUILDER_ARGS"
	fi
}

function invalidate_cloudfront() {
	echo "--- Invalidating cloudfront cache..."

	BUILDER_ARGS="$@"
	app_version=$(jq --raw-output '.version' package.json)
	declare -a paths

	if [[ "$BUILDER_ARGS" =~ "--mac" ]]; then
		if [[ ! "$APP_VERSION" =~ "alpha" ]] && [[ ! "$APP_VERSION" =~ "beta" ]]; then
			if [[ "$BUILDER_ARGS" =~ "--x64" ]]; then
				paths+=("/client/Sturdy-${app_version}.dmg")
				paths+=("/client/darwin/amd64/Install*")
			fi

			if [[ "$BUILDER_ARGS" =~ "--arm64" ]]; then
				paths+=("/client/Sturdy-${app_version}-arm64.dmg")
				paths+=("/client/darwin/arm64/Install*")
			fi
		fi

		paths+=("/client/alpha-mac.yml")
		paths+=("/client/beta-mac.yml")
		paths+=("/client/latest-mac.yml")
	fi

	if [[ "$BUILDER_ARGS" =~ "--windows" ]]; then
		if [[ ! "$APP_VERSION" =~ "alpha" ]] && [[ ! "$APP_VERSION" =~ "beta" ]]; then
			paths+=("/client/Sturdy-Installer-${app_version}.exe")
			paths+=("/client/windows/amd64/Sturdy-Installer.exe")
		fi
		paths+=("/client/alpha.yml")
		paths+=("/client/beta.yml")
		paths+=("/client/latest.yml")
	fi

	if [[ "$BUILDER_ARGS" =~ "--linux" ]]; then
		if [[ "$BUILDER_ARGS" =~ "--x64" ]]; then
			if [[ ! "$APP_VERSION" =~ "alpha" ]] && [[ ! "$APP_VERSION" =~ "beta" ]]; then
				paths+=("/client/Sturdy_${app_version}_amd64.deb")
				paths+=("/client/linux/amd64/Sturdy-Latest.deb")
				paths+=("/client/Sturdy-${app_version}.x86_64.rpm")
				paths+=("/client/linux/amd64/Sturdy-Latest.rpm")
				paths+=("/client/Sturdy-${app_version}.AppImage")
				paths+=("/client/linux/amd64/Sturdy.AppImage")
			fi
			paths+=("/client/alpha-linux.yml")
			paths+=("/client/beta-linux.yml")
			paths+=("/client/latest-linux.yml")
		fi

		if [[ "$BUILDER_ARGS" =~ "--arm64" ]]; then
			if [[ ! "$APP_VERSION" =~ "alpha" ]] && [[ ! "$APP_VERSION" =~ "beta" ]]; then
				paths+=("/client/Sturdy_${app_version}_arm64.deb")
				paths+=("/client/linux/arm64/Sturdy-Latest.deb")
				paths+=("/client/Sturdy-${app_version}.aarch64.rpm")
				paths+=("/client/linux/arm64/Sturdy-Latest.rpm")
				paths+=("/client/Sturdy-${app_version}-arm64.AppImage")
				paths+=("/client/linux/arm64/Sturdy.AppImage")
			fi
			paths+=("/client/alpha-linux-arm64.yml")
			paths+=("/client/beta-linux-arm64.yml")
			paths+=("/client/latest-linux-arm64.yml")
		fi
	fi

	aws cloudfront create-invalidation --distribution-id EUQY8O4OTQKLV --paths "${paths[@]}"
}

function setup_windows_codesign() {
	echo "--- Setting up Windows code signing..."

	TMP_DIR=$(mktemp -d)
	aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.crt | jq --raw-output '.SecretString' >"$TMP_DIR/sign.crt"
	aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.key | jq --raw-output '.SecretString' >"$TMP_DIR/sign.key"

	WIN_CSC_LINK="$TMP_DIR/sign.pfx"
	WIN_CSC_KEY_PASSWORD="$(head -c 64 /dev/urandom | base64)"

	export WIN_CSC_LINK
	export WIN_CSC_KEY_PASSWORD

	openssl pkcs12 \
		-export \
		-out "$WIN_CSC_LINK" \
		-in "$TMP_DIR/sign.crt" \
		-inkey "$TMP_DIR/sign.key" \
		-password pass:"$WIN_CSC_KEY_PASSWORD"
}

function validate_version() {
	local version="$1"
	# https://semver.org/
	SEMVER_REGEX="^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$"
	(echo "$version" | grep -Eq "$SEMVER_REGEX") || (echo "$version: invalid semver, see https://semver.org/" && exit 1)
}

function create_latest() {
	BUILDER_ARGS="$@"
	app_version=$(jq --raw-output '.version' package.json)

	echo "--- Creating latest version files..."

	if [[ "$BUILDER_ARGS" =~ "--mac" ]] && [[ "$BUILDER_ARGS" =~ "--x64" ]]; then
		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}.dmg" \
			"s3://autoupdate.getsturdy.com/client/darwin/amd64/Install Sturdy.dmg"
	fi

	if [[ "$BUILDER_ARGS" =~ "--mac" ]] && [[ "$BUILDER_ARGS" =~ "--arm64" ]]; then
		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}-arm64.dmg" \
			"s3://autoupdate.getsturdy.com/client/darwin/arm64/Install Sturdy.dmg"
	fi

	if [[ "$BUILDER_ARGS" =~ "--windows" ]] && [[ "$BUILDER_ARGS" =~ "--x64" ]]; then
		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-Installer-${app_version}.exe" \
			"s3://autoupdate.getsturdy.com/client/windows/amd64/Sturdy-Installer.exe"
	fi

	if [[ "$BUILDER_ARGS" =~ "--linux" ]] && [[ "$BUILDER_ARGS" =~ "--x64" ]]; then
		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy_${app_version}_amd64.deb" \
			"s3://autoupdate.getsturdy.com/client/linux/amd64/Sturdy-Latest.deb"

		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}.x86_64.rpm" \
			"s3://autoupdate.getsturdy.com/client/linux/amd64/Sturdy-Latest.rpm"

		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}.AppImage" \
			"s3://autoupdate.getsturdy.com/client/linux/amd64/Sturdy.AppImage"
	fi

	if [[ "$BUILDER_ARGS" =~ "--linux" ]] && [[ "$BUILDER_ARGS" =~ "--arm64" ]]; then
		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy_${app_version}_arm64.deb" \
			"s3://autoupdate.getsturdy.com/client/linux/arm64/Sturdy-Latest.deb"

		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}.aarch64.rpm" \
			"s3://autoupdate.getsturdy.com/client/linux/arm64/Sturdy-Latest.rpm"

		aws s3 cp \
			"s3://autoupdate.getsturdy.com/client/Sturdy-${app_version}-arm64.AppImage" \
			"s3://autoupdate.getsturdy.com/client/linux/arm64/Sturdy.AppImage"
	fi
}

function setup_darwin_notarize() {
	echo "--- Setting up Apple code signing..."
	printf 'Apple ID (email): '
	read APPLE_ID
	printf 'Password: '
	read -s APPLE_ID_PASSWORD
	echo
	export APPLE_ID
	export APPLE_ID_PASSWORD
}
