#!/usr/bin/env bash

set -euo pipefail

echoerr() { echo "$@" 1>&2; }

function reset() {
  rm -rf dist out
}

function build_app() {
  yarn build
}

function download_sturdy_sync() {
    BIN_DIR=assets/bin

    rm -rf $BIN_DIR
    mkdir $BIN_DIR

    OS=$1
    ARCH=$2

    ARCHIVE_FORMAT="tar.gz"
    if [ "$OS" == "windows" ]; then
      ARCHIVE_FORMAT="zip"
    fi

    ARCHIVE_NAME="sturdy-$STURDY_SYNC_VERSION-$OS-$ARCH.$ARCHIVE_FORMAT"

    curl -s -Lo "$BIN_DIR/$ARCHIVE_NAME" "https://getsturdy.com/client/$ARCHIVE_NAME"

    if [ "$ARCHIVE_FORMAT" == "tar.gz" ]; then
      tar xzf "$BIN_DIR/$ARCHIVE_NAME" -C $BIN_DIR
    elif [ "$ARCHIVE_FORMAT" == "zip" ]; then
      unzip "$BIN_DIR/$ARCHIVE_NAME" -d $BIN_DIR/
    else
      echoerr "Unsupported archive format: $ARCHIVE_FORMAT"
      exit 1
    fi

    rm "$BIN_DIR/$ARCHIVE_NAME"
}

function build() {
  OS=$1
  ARCH=$2

  echo "--- Building for $OS/$ARCH"

  reset
  download_sturdy_sync "$OS" "$ARCH"

  rm -rf out/make

  build_app

  ELECTRON_BUILDER_OS="--mac"
  ELECTRON_BUILDER_ARCH="--x64"

  if [ "$ARCH" == "arm64" ]; then
    ELECTRON_BUILDER_ARCH="--arm64"
  fi

  if [ "$OS" == "windows" ]; then
    ELECTRON_BUILDER_OS="--win"
  fi

  if [ "$OS" == "linux" ]; then
    ELECTRON_BUILDER_OS="--linux"
  fi

  if (( CODESIGN )) && [ "$OS" == "windows" ] ; then
    prep_codesign_windows
  fi

  CODESIGN_EXTRA_ARGS=""
  if [ "$OS" == "darwin" ] ; then
    if ! (( CODESIGN )); then
      CODESIGN_EXTRA_ARGS="-c.mac.identity=null" # Disable code signing
      echo "--- Skipping code signing!"
    fi
  fi

  CHANNEL_PATH_SUFFIX=""

  if [ ! -z "$CHANNEL" ]; then
    CHANNEL_PATH_SUFFIX=$(echo "-$CHANNEL" | tr '[:upper:]' '[:lower:]')
  fi

  BUILDER_CONFIG_YML="electron-builder-${OS}-${ARCH}.yml"

  yq eval ".publish[0].url |= \"https://autoupdate.getsturdy.com/client$CHANNEL_PATH_SUFFIX/$OS/$ARCH\" | .publish[1].path=\"client$CHANNEL_PATH_SUFFIX/$OS/$ARCH\"" \
    electron-builder.yml > $BUILDER_CONFIG_YML

  if (( NOTARIZE )) && [ "$OS" == "darwin" ]; then
    yq -i eval '.afterSign = "./afterSignHook.js"' "$BUILDER_CONFIG_YML"
  fi

  if [ ! -z "$CHANNEL" ]; then
    yq -i eval "
      .productName += \" $CHANNEL\",
      .appId += \"$CHANNEL_PATH_SUFFIX\",
      .extraMetadata.name += \" $CHANNEL\",
      .linux.desktop.Name += \" $CHANNEL\"
    " "$BUILDER_CONFIG_YML"
  fi

  PUBLISH_ARGS=""
  if (( DO_UPLOAD )); then
    PUBLISH_ARGS="--publish=always"
  fi

  yarn electron-builder "$ELECTRON_BUILDER_OS" "$ELECTRON_BUILDER_ARCH" \
    --config "$BUILDER_CONFIG_YML" \
    $PUBLISH_ARGS \
    $CODESIGN_EXTRA_ARGS

  if (( DO_UPLOAD )) && [ -z "$CHANNEL" ]; then
      create_latest "$OS" "$ARCH"
  fi
}

function prep_codesign_windows() {
  TMP_DIR=$(mktemp -d)
  aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.crt | jq --raw-output '.SecretString' > "$TMP_DIR/sign.crt"
  aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.key | jq --raw-output '.SecretString' > "$TMP_DIR/sign.key"

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

function invalidate_cloudfront() {
  CHANNEL="$1"
  CHANNEL_CLIENT_DIR="client"

  if [ ! -z "$CHANNEL" ]; then
    CHANNEL_CLIENT_DIR=$(echo "client-$CHANNEL" | tr '[:upper:]' '[:lower:]')
  fi

  aws cloudfront create-invalidation --distribution-id EUQY8O4OTQKLV \
    --paths "/${CHANNEL_CLIENT_DIR}/darwin/amd64/latest-mac.yml" \
     "/${CHANNEL_CLIENT_DIR}/darwin/arm64/latest-mac.yml" \
     "/${CHANNEL_CLIENT_DIR}/windows/amd64/latest.yml" \
     "/${CHANNEL_CLIENT_DIR}/linux/amd64/latest.yml" \
     "/${CHANNEL_CLIENT_DIR}/darwin/amd64/Install*" \
     "/${CHANNEL_CLIENT_DIR}/darwin/arm64/Install*" \
     "/${CHANNEL_CLIENT_DIR}/linux/amd64/*" \
     "/${CHANNEL_CLIENT_DIR}/windows/amd64/Sturdy-Installer*"
}

function create_latest() {
  OS=$1
  ARCH=$2

  app_version=$(jq --raw-output '.version' package.json)

  if [ "$OS" == "darwin" ] && [ "$ARCH" == "amd64" ]; then
    aws s3 cp \
        "s3://autoupdate.getsturdy.com/client/darwin/amd64/Sturdy-${app_version}.dmg" \
        "s3://autoupdate.getsturdy.com/client/darwin/amd64/Install Sturdy.dmg"
  fi

  if [ "$OS" == "darwin" ] && [ "$ARCH" == "arm64" ]; then
    aws s3 cp \
        "s3://autoupdate.getsturdy.com/client/darwin/arm64/Sturdy-${app_version}-arm64.dmg" \
        "s3://autoupdate.getsturdy.com/client/darwin/arm64/Install Sturdy.dmg"
  fi

  if [ "$OS" == "windows" ] && [ "$ARCH" == "amd64" ]; then
    aws s3 cp \
      "s3://autoupdate.getsturdy.com/client/windows/amd64/Sturdy-Installer-${app_version}.exe" \
      "s3://autoupdate.getsturdy.com/client/windows/amd64/Sturdy-Installer.exe"
  fi
}

function setup_darwin_notarize() {
  printf 'Apple ID (email): '
  read APPLE_ID
  printf 'Password: '
  read -s APPLE_ID_PASSWORD
  export APPLE_ID
  export APPLE_ID_PASSWORD
}