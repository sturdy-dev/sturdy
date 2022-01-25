#!/usr/bin/env bash

set -euxo pipefail

VERSION=""
DO_UPLOAD=0
LOCAL_INSTALL=0
BUILD_DARWIN_ONLY=0
HOMEBREW_NAME="sturdy"
WINDOWS_INSTALLER_NAME="windows-installer.ps1"
MUTAGEN_RELEASE="e42274bc5746bfc4312d7fe026b73a5f5c0b6b34" # sturdy-v0.13.0-beta2
MUTAGEN_PATH="${HOME}/src/mutagen"
HOMEBREW_TAP_REPO_PATH="${HOME}/src/sturdy-dev-homebrew-tap"

while [[ $# -gt 0 ]]; do
  case "$1" in
  --version)
    VERSION="$2"
    shift
    shift
    ;;
  --upload)
    DO_UPLOAD=1
    shift
    ;;
  --install)
    LOCAL_INSTALL=1
    shift
    ;;
  --build-darwin-only)
    BUILD_DARWIN_ONLY=1
    shift
    ;;
  --beta)
    HOMEBREW_NAME="sturdybeta"
    WINDOWS_INSTALLER_NAME="windows-installer-beta.ps1"
    shift
    ;;
  esac
done

VERSION_WITHOUT_V="${VERSION/v/}"

errecho() { echo >&2 $@; }

if [ -z "$VERSION" ]; then
  errecho "--version is not set! aborting!"
  exit 1
fi

# Prepare Mutagen
git -C $MUTAGEN_PATH fetch
git -C $MUTAGEN_PATH checkout $MUTAGEN_RELEASE
git -C $MUTAGEN_PATH status

build_upload() {
  GOOS=$1
  GOARCH=$2
  GOARM=""

  ARCHIVE_FORMAT="tar.gz"
  # The third argument can be "zip" (for Windows builds)
  if [ $# -ge 3 ] && [ "$3" == "zip" ]; then
    ARCHIVE_FORMAT="zip"
  fi

  # Third argument can also be a GOARM value
  if [ $# -ge 3 ] && [ "$3" != "zip" ]; then
    GOARM=$3
  fi

  STURDY_BIN_NAME="sturdy"
  STURDY_SYNC_BIN_NAME="sturdy-sync"
  if [ "$GOOS" == "windows" ]; then
    STURDY_BIN_NAME="${STURDY_BIN_NAME}.exe"
    STURDY_SYNC_BIN_NAME="${STURDY_SYNC_BIN_NAME}.exe"
  fi

  BUILD_TARGET_NAME="$VERSION-$GOOS-$GOARCH"
  # Add GOARM to name if set
  if [ "$GOARM" != "" ]; then
    BUILD_TARGET_NAME+="${GOARM}"
  fi
  BUILD_TARGET_NAME+=".$ARCHIVE_FORMAT"

  ARCHIVE_NAME="sturdy-$BUILD_TARGET_NAME"

  OUTPUT_DIR=$(mktemp -d)

  errecho "Building ${ARCHIVE_NAME}"

  # Build binary
  GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build -o "${OUTPUT_DIR}/${STURDY_BIN_NAME}" \
    -ldflags "-X getsturdy.com/client/cmd/sturdy/version.Version=$VERSION" \
    getsturdy.com/client/cmd/sturdy

  # Build mutagen
  cd $MUTAGEN_PATH
  GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build -o "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}" \
    -ldflags "-X github.com/mutagen-io/mutagen/pkg/synchronization.SturdyVersion=$VERSION" \
    -ldflags "-X github.com/mutagen-io/mutagen/pkg/sturdy/api.clientVersion=sturdy-sync/$VERSION" \
    github.com/mutagen-io/mutagen/cmd/mutagen
  cd - >/dev/null 2>&1

  # Cleanup (if exist)
  rm "$ARCHIVE_NAME" || true

  if [ "$GOOS" == "windows" ]; then
    sign_windows_binaries
  fi

  # Build archived version (to use with Homebrew)
  if [ "$ARCHIVE_FORMAT" == "tar.gz" ]; then
    tar -s "/-$BUILD_TARGET_NAME//" \
      -s "#${OUTPUT_DIR}##" \
      -zcvf "$ARCHIVE_NAME" \
      "${OUTPUT_DIR}/${STURDY_BIN_NAME}" "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}"
  elif [ "$ARCHIVE_FORMAT" == "zip" ]; then
    zip "$ARCHIVE_NAME" --junk-paths "${OUTPUT_DIR}/${STURDY_BIN_NAME}" "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}"
  else
    errecho "Unexpected archive format ${ARCHIVE_FORMAT}"
    exit 1
  fi

  if ((DO_UPLOAD)); then
    aws s3 cp --quiet ${ARCHIVE_NAME} "s3://getsturdy.com/client/${ARCHIVE_NAME}"
  else
    errecho "In dry-run, not uploading!"
  fi

  # Windows direct downloads
  if [ "$GOOS" == "windows" ]; then
    if ((DO_UPLOAD)); then
      aws s3 cp --quiet "${OUTPUT_DIR}/${STURDY_BIN_NAME}" "s3://getsturdy.com/client/sturdy-$VERSION-$GOOS-$GOARCH.exe"
      aws s3 cp --quiet "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}" "s3://getsturdy.com/client/sturdy-sync-$VERSION-$GOOS-$GOARCH.exe"
    else
      errecho "In dry-run, not uploading windows binaries!"
    fi
  fi

  echo $ARCHIVE_NAME
}

windows_codesign() {
  local input
  local output
  input=$1
  output=$2

  TMP_DIR=$(mktemp -d)
  aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.crt | jq --raw-output '.SecretString' >"${TMP_DIR}/sign.crt"
  aws secretsmanager get-secret-value --secret-id sturdy/codesign/sturdy_sweden_ab.key | jq --raw-output '.SecretString' >"${TMP_DIR}/sign.key"

  osslsigncode sign \
    -certs "${TMP_DIR}/sign.crt" \
    -key "${TMP_DIR}/sign.key" \
    -n "Sturdy" \
    -i "https://getsturdy.com/" \
    -t http://timestamp.digicert.com \
    -in "${input}" \
    -out "${output}"

  # Cleanup
  rm -f "${TMP_DIR}/sign.crt" "${TMP_DIR}/sign.key"
}

sign_windows_binaries() {
  windows_codesign "${OUTPUT_DIR}/${STURDY_BIN_NAME}" "${OUTPUT_DIR}/signed-${STURDY_BIN_NAME}"
  windows_codesign "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}" "${OUTPUT_DIR}/signed-${STURDY_SYNC_BIN_NAME}"

  # replace the non-signed versions
  mv "${OUTPUT_DIR}/signed-${STURDY_BIN_NAME}" "${OUTPUT_DIR}/${STURDY_BIN_NAME}"
  mv "${OUTPUT_DIR}/signed-${STURDY_SYNC_BIN_NAME}" "${OUTPUT_DIR}/${STURDY_SYNC_BIN_NAME}"
}

# Update InstallationInstructions.vue version number
update_web_install() {
  t=$(mktemp)
  cat ../../../web/src/components/install/InstallationInstructions.vue |
    sed "s/const latestVersion = '.*'/const latestVersion = '${VERSION}'/" >$t
  mv $t ../../../web/src/components/install/InstallationInstructions.vue
}

create_formula() {
  local NAME_TITLECASE="$(tr '[:lower:]' '[:upper:]' <<<${HOMEBREW_NAME:0:1})${HOMEBREW_NAME:1}"
  cat >"${HOMEBREW_TAP_REPO_PATH}/Formula/${HOMEBREW_NAME}.rb" <<EOF
class ${NAME_TITLECASE} < Formula
    desc "Sturdy Client"
    homepage "https://getsturdy.com/"
    version "${VERSION_WITHOUT_V}"

    if OS.mac? && Hardware::CPU.intel?
        url "https://getsturdy.com/client/${DARWIN_AMD64_TAR_NAME}"
        sha256 "$(sha256sum ${DARWIN_AMD64_TAR_NAME} | awk '{print $1}')"
    elsif OS.mac? && Hardware::CPU.arm?
        url "https://getsturdy.com/client/${DARWIN_ARM64_TAR_NAME}"
        sha256 "$(sha256sum ${DARWIN_ARM64_TAR_NAME} | awk '{print $1}')"
    elsif OS.linux? && Hardware::CPU.intel?
        url "https://getsturdy.com/client/${LINUX_AMD64_TAR_NAME}"
        sha256 "$(sha256sum ${LINUX_AMD64_TAR_NAME} | awk '{print $1}')"
    elsif OS.linux? && Hardware::CPU.arm?
        url "https://getsturdy.com/client/${LINUX_ARM64_TAR_NAME}"
        sha256 "$(sha256sum ${LINUX_ARM64_TAR_NAME} | awk '{print $1}')"
    end

    def install
        bin.install "sturdy"
        bin.install "sturdy-sync"
    end
end
EOF
}

commit_and_push_formula() {
  git -C "$HOMEBREW_TAP_REPO_PATH" diff
  git -C "$HOMEBREW_TAP_REPO_PATH" add "Formula/${HOMEBREW_NAME}.rb"
  git -C "$HOMEBREW_TAP_REPO_PATH" commit -m "sturdy ${VERSION}"
  git -C "$HOMEBREW_TAP_REPO_PATH" push
}

update_windows() {
  t=$(mktemp)
  cat windows-installer.ps1 |
    sed "s/\$VERSION=\".*\"/\$VERSION=\"${VERSION}\"/" >$t
  mv $t windows-installer.ps1
}

upload_windows() {
  aws s3 cp windows-installer.ps1 "s3://getsturdy.com/client/${WINDOWS_INSTALLER_NAME}"
}

# Build and upload files to S3
DARWIN_AMD64_TAR_NAME=$(build_upload darwin amd64)
if ((BUILD_DARWIN_ONLY == 0)); then
  DARWIN_ARM64_TAR_NAME=$(build_upload darwin arm64)
  LINUX_AMD64_TAR_NAME=$(build_upload linux amd64)
  LINUX_ARM64_TAR_NAME=$(build_upload linux arm64)
  LINUX_ARMv5_TAR_NAME=$(build_upload linux arm 5)
  LINUX_ARMv6_TAR_NAME=$(build_upload linux arm 6)
  LINUX_ARMv7_TAR_NAME=$(build_upload linux arm 7)
  WINDOWS_AMD64_ZIP_NAME=$(build_upload windows amd64 zip)
fi

update_windows
update_web_install

# Push update to Homebrew Tap
if ((DO_UPLOAD)); then
  create_formula
  commit_and_push_formula
  upload_windows
else
  errecho "In dry-run, not creating formula!"
fi

if ((LOCAL_INSTALL)); then
  rm -rf "/usr/local/Cellar/sturdy" || true

  mkdir -p "/usr/local/Cellar/sturdy/${VERSION_WITHOUT_V}/bin"
  tar xzvf ${DARWIN_AMD64_TAR_NAME} -C "/usr/local/Cellar/sturdy/${VERSION_WITHOUT_V}/bin"

  rm "/usr/local/bin/sturdy" || true
  rm "/usr/local/bin/sturdy-sync" || true

  ln -s "/usr/local/Cellar/sturdy/${VERSION_WITHOUT_V}/bin/sturdy" "/usr/local/bin/sturdy"
  ln -s "/usr/local/Cellar/sturdy/${VERSION_WITHOUT_V}/bin/sturdy-sync" "/usr/local/bin/sturdy-sync"
fi
