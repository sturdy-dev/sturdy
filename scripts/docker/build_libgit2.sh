#!/usr/bin/env bash

# This file is largely inspired by fluxcd/golang-with-libgit2
# https://github.com/fluxcd/golang-with-libgit2
#
# It has been modified to only build libgit2 from source
#
# License: Apache 2.0

set -euxo pipefail

LIBGIT2_URL="${LIBGIT2_URL:-https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.tar.gz}"

TARGET_DIR="${TARGET_DIR:-/usr/local/$(xx-info triple)}"
BUILD_ROOT_DIR="${BUILD_ROOT_DIR:-/build}"
SRC_DIR="${BUILD_ROOT_DIR}/src"

TARGET_ARCH="${TARGET_ARCH:-$(uname -m)}"
if command -v xx-info; then
    TARGET_ARCH="$(xx-info march)"
fi

C_COMPILER="${CC:-/usr/bin/gcc}"
CMAKE_PARAMS=""
if command -v xx-clang; then
    C_COMPILER="/usr/bin/xx-clang"
    CMAKE_PARAMS="$(xx-clang --print-cmake-defines)"
fi

function download_source(){
    mkdir -p "$2"

    curl --max-time 120 -o "$2/source.tar.gz" -LO "$1"
    tar -C "$2" --strip 1 -xzvf "$2/source.tar.gz"
    rm "$2/source.tar.gz"
}

function build_libgit2(){
    download_source "${LIBGIT2_URL}" "${SRC_DIR}/libgit2"

    pushd "${SRC_DIR}/libgit2"

    mkdir -p build

    pushd build

    SSL_LIBRARY="${TARGET_DIR}/lib/libssl.a"
    CRYPTO_LIBRARY="${TARGET_DIR}/lib/libcrypto.a"
    if [[ ! $OSTYPE == darwin* ]] && [ "${TARGET_ARCH}" = "x86_64" ]; then
        SSL_LIBRARY="${TARGET_DIR}/lib64/libssl.a"
        CRYPTO_LIBRARY="${TARGET_DIR}/lib64/libcrypto.a"
    fi

    # Set osx arch only when cross compiling on darwin
    if [[ $OSTYPE == darwin* ]] && [ ! "${TARGET_ARCH}" = "$(uname -m)" ]; then
        CMAKE_PARAMS=-DCMAKE_OSX_ARCHITECTURES="${TARGET_ARCH}"
    fi

    cmake "${CMAKE_PARAMS}" \
        -DCMAKE_C_COMPILER="${C_COMPILER}" \
        -DCMAKE_INSTALL_PREFIX="${TARGET_DIR}" \
        -DTHREADSAFE:BOOL=ON \
        -DBUILD_CLAR:BOOL=OFF \
        -DBUILD_SHARED_LIBS=OFF \
        -DCMAKE_POSITION_INDEPENDENT_CODE:BOOL=ON \
        -DCMAKE_C_FLAGS=-fPIC \
        -DUSE_SSH:BOOL=ON \
        -DHAVE_LIBSSH2_MEMORY_CREDENTIALS:BOOL=ON \
        -DDEPRECATE_HARD:BOOL=ON \
        -DUSE_BUNDLED_ZLIB:BOOL=ON \
        -DUSE_HTTPS:STRING=OpenSSL \
        -DREGEX_BACKEND:STRING=builtin \
        -DOPENSSL_SSL_LIBRARY="${SSL_LIBRARY}" \
        -DOPENSSL_CRYPTO_LIBRARY="${CRYPTO_LIBRARY}" \
        -DCMAKE_INCLUDE_PATH="${TARGET_DIR}/include" \
        -DCMAKE_LIBRARY_PATH="${TARGET_DIR}/lib" \
        -DCMAKE_PREFIX_PATH="${TARGET_DIR}" \
        -DCMAKE_BUILD_TYPE="RelWithDebInfo" \
        ..

    cmake --build . --target install

    popd
    popd
}

"$@"
