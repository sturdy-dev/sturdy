# syntax=docker/dockerfile:1

ARG XX_VERSION=1.1.0

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM golang:1.18.0-bullseye as ssh-builder
WORKDIR /go/src/ssh
COPY ./ssh/scripts/build-mutagen.sh ./scripts/build-mutagen.sh
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/go-mod \
    export GOMODCACHE=/root/.cache/go-mod && \
    bash ./scripts/build-mutagen.sh
# cache ssh depencencies
COPY ./ssh/go.mod ./go.mod
COPY ./ssh/go.sum ./go.sum
COPY ./ssh .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/go-mod \
    GOMODCACHE=/root/.cache/go-mod \
    go build -v -o /usr/bin/ssh getsturdy.com/ssh/cmd/ssh

FROM debian:11.2-slim as ssh
RUN apt-get update && apt-get install -y --no-install-recommends --allow-downgrades \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=ssh-builder /usr/bin/ssh /usr/bin/ssh
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta2 /usr/bin/mutagen-agent-v0.12.0-beta2
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta6 /usr/bin/mutagen-agent-v0.12.0-beta6
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta7 /usr/bin/mutagen-agent-v0.12.0-beta7
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.13.0-beta2 /usr/bin/mutagen-agent-v0.13.0-beta2

FROM golang:1.18.0-bullseye as libgit-builder

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    build-essential \
    clang \
    libssl-dev \
    libz-dev \
    cmake

COPY --from=xx / /

WORKDIR /static
COPY --chmod=0744 scripts/docker/build_libgit2.sh .
RUN ./build_libgit2.sh build_libgit2

# TODO(gustav): figure out why running this step above (before building libgit2) breaks things
RUN apt-get update && apt-get install -y libssh2-1-dev && rm -rf /var/lib/apt/lists/*

FROM libgit-builder as api-builder

WORKDIR /go/src/api

# build api
ARG API_BUILD_TAGS
ARG VERSION
COPY ./api ./

ENV LD_LIBRARY_PATH=/usr/local/lib
ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/root/.cache/go-build,id=go-build \
    --mount=type=cache,target=/root/.cache/go-mod,id=go-cache \
    export LIBRARY_PATH="/usr/local/$(xx-info triple):/usr/local/$(xx-info triple)/lib64:${LIBRARY_PATH}" && \
    export PKG_CONFIG_PATH="/usr/local/$(xx-info triple)/lib/pkgconfig:/usr/local/$(xx-info triple)/lib64/pkgconfig:${PKG_CONFIG_PATH}" && \
    export FLAGS="$(pkg-config --static --libs --cflags libssh2 openssl libgit2)" && \
    export CGO_LDFLAGS="${FLAGS} -static" && \
    export GOMODCACHE=/root/.cache/go-mod && \
    go build \
    -tags "${API_BUILD_TAGS},netgo,osusergo,static_build" \
    -ldflags "-X getsturdy.com/api/pkg/version.Version=${VERSION}" \
    -v -o /usr/bin/api  getsturdy.com/api/cmd/api

FROM debian:11.2-slim as api
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends --allow-downgrades \
    git-lfs=2.13.2-1+b5 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=api-builder /usr/bin/api /usr/bin/api
ENTRYPOINT [ "/usr/bin/api" ]

# for amd64, use a prebuild rudolfs image
FROM jasonwhite0/rudolfs:0.3.5 as rudolfs-builder-amd64
# for arm64, compile from source, since there is no prebuild image
FROM --platform=$BUILDPLATFORM rust:1.55 as rudolfs-builder-arm64
ENV PKG_CONFIG_ALLOW_CROSS="1" \
    DEBIAN_FRONTEND="noninteractive" \
    CARGO_BUILD_TARGET="aarch64-unknown-linux-gnu"
RUN apt-get update \
    && apt-get -y --no-install-recommends --allow-downgrades install \
    musl-tools \
    ca-certificates \
    git \
    && git clone https://github.com/jasonwhite/rudolfs
WORKDIR /rudolfs
SHELL ["/bin/bash", "-c"]
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/rudolfs/target \
    git checkout 0.3.5 \
    && rustup target add "${CARGO_BUILD_TARGET}" \
    && cargo build --target "${CARGO_BUILD_TARGET}" --release  \
    && mkdir -p /build \
    && cp "target/${CARGO_BUILD_TARGET}/release/rudolfs" /build/ \
    && strip /build/rudolfs

FROM debian:11.2-slim as rudolfs
VOLUME ["/data"]
RUN apt-get update \
    && apt-get install -y --no-install-recommends --allow-downgrades \
    ca-certificates
# use the correct binary depending on the architecture. we do this to avoid building amd64 version ourselves, 
# as it requires us to run qemu emulation which is very slow.
ARG TARGETARCH
COPY --from=rudolfs-builder-arm64 /build/rudolfs /arm64/rudolfs
COPY --from=rudolfs-builder-amd64 /rudolfs /amd64/rudolfs
RUN cp "/${TARGETARCH}/rudolfs" /usr/bin/rudolfs
EXPOSE 8080
ENTRYPOINT ["/usr/bin/rudolfs"]
CMD ["--cache-dir", "/data"]

FROM --platform=$BUILDPLATFORM node:17.3.1-alpine3.15 as web-builder
# The website is the same for linux/amd64 and linux/arm64 (output is html), setting --platform to run all builds on the
# native host platform. (Skips emulation!)
WORKDIR /web
RUN apk update \
    && apk add --no-cache \
    python3=3.9.7-r4 \
    make=4.3-r0 \
    g++ \
    git
# cache web dependencies
COPY ./web/package.json ./package.json
COPY ./web/yarn.lock ./yarn.lock
# The --network-timeout is here to prevent network issues when building linux/amd64 images on linux/arm64 hosts
RUN --mount=type=cache,target=/root/.yarn,id=web-installer-3 YARN_CACHE_FOLDER=/root/.yarn \
    yarn install --frozen-lockfile \
    --network-timeout 1000000000
# build web
COPY ./web .
RUN yarn build:oneliner

FROM golang:1.18.0-bullseye as sslmux-builder
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/go-mod \
    GOMODCACHE=/root/.cache/go-mod \
    go install -v github.com/JamesDunne/sslmux@v0.0.0-20180531161153-81a78ca8247d

FROM debian:11.2-slim  as reproxy-builder
ARG REPROXY_VERSION="v0.11.0"
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN if [[ "$(uname -m)" == 'aarch64' ]]; then \
    ARCH='arm64'; \
    REPROXY_SHA256_SUM='35dd1cc3568533a0b6e1109e7ba630d60e2e39716eea28d3961c02f0feafee8e'; \
    elif [[ "$(uname -m)" == 'x86_64' ]]; then \
    ARCH='x86_64'; \
    REPROXY_SHA256_SUM='100a1389882b8ab68ae94f37e9222f5f928ece299d8cfdf5b26c9f12f902c23a'; \
    fi \
    && apt-get update && apt-get install -y wget \
    && wget --quiet --output-document "/tmp/reproxy.tar.gz" "https://github.com/umputun/reproxy/releases/download/${REPROXY_VERSION}/reproxy_${REPROXY_VERSION}_linux_${ARCH}.tar.gz" \
    && sha256sum "/tmp/reproxy.tar.gz" \
    && echo "${REPROXY_SHA256_SUM}  /tmp/reproxy.tar.gz" | sha256sum -c \
    && tar -xzf /tmp/reproxy.tar.gz -C /usr/bin \
    && rm /tmp/reproxy.tar.gz

FROM debian:11.2-slim as oneliner
# postgresql
# openssl is needed by rudolfs to generate secret
# git, git-lfs and libgit2 are needed by api
# openssh-keygen is needed by ssh to generate ssh keys
# ca-cerificates is needed by ssh to connect to tls hosts
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends --allow-downgrades \
    curl=7.74.0-1.3+deb11u1 \
    ca-certificates=20210119 \
    gnupg \
    && curl https://www.postgresql.org/media/keys/ACCC4CF8.asc \
    | gpg --dearmor \
    | tee /etc/apt/trusted.gpg.d/apt.postgresql.org.gpg >/dev/null \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ bullseye-pgdg main" > /etc/apt/sources.list.d/postgresql.list \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends --allow-downgrades \
    postgresql-14=14.2-1.pgdg110+1 \
    openssl=1.1.1k-1+deb11u2 \
    git=1:2.30.2-1 \
    git-lfs=2.13.2-1+b5 \
    keychain=2.8.5-2 \
    wget \
    bash \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# s6-overlay
ARG S6_OVERLAY_VERSION="3.0.0.2" \
    S6_OVERLAY_NOARCH_SHA256_SUM="17880e4bfaf6499cd1804ac3a6e245fd62bc2234deadf8ff4262f4e01e3ee521" \
    S6_OVERLAY_SYMLINKS_ARCH_SHA256_SUM="6ee2b8580b23c0993b1e8c66b58777f32f6ff031ba0192cccd53a31e62942c70" \
    S6_OVERLAY_SYMLINKS_NOARCH_SHA256_SUM="d67c9b436ef59ffefd4f083f07b2869662af40b2ea79a069b147dd0c926db2d3"
RUN ARCH="$(uname -m)" \
    && if [[ "$ARCH" == 'x86_64' ]]; then \
    S6_OVERLAY_ARCH_SHA256_SUM="a4c039d1515812ac266c24fe3fe3c00c48e3401563f7f11d09ac8e8b4c2d0b0c"; \
    elif [[ "$ARCH" == 'aarch64' ]]; then \
    S6_OVERLAY_ARCH_SHA256_SUM="e6c15e22dde00af4912d1f237392ac43a1777633b9639e003ba3b78f2d30eb33"; \
    fi \
    && wget --quiet --output-document "/tmp/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && sha256sum "/tmp/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && echo "${S6_OVERLAY_NOARCH_SHA256_SUM}  /tmp/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" | sha256sum -c \
    && tar -C / -Jxpf "/tmp/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && rm "/tmp/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    \
    && wget --quiet --output-document "/tmp/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" \
    && sha256sum "/tmp/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" \
    && echo "${S6_OVERLAY_ARCH_SHA256_SUM}  /tmp/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" | sha256sum -c \
    && tar -C / -Jxpf "/tmp/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" \
    && rm "/tmp/s6-overlay-${ARCH}-${S6_OVERLAY_VERSION}.tar.xz" \
    \
    && wget --quiet --output-document "/tmp/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && sha256sum "/tmp/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && echo "${S6_OVERLAY_SYMLINKS_NOARCH_SHA256_SUM}  /tmp/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" | sha256sum -c \
    && tar -C / -Jxpf "/tmp/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && rm "/tmp/s6-overlay-symlinks-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    \
    && wget --quiet --output-document "/tmp/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz" "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz" \
    && sha256sum "/tmp/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz" \
    && echo "${S6_OVERLAY_SYMLINKS_ARCH_SHA256_SUM}  /tmp/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz" | sha256sum -c \
    && tar -C / -Jxpf "/tmp/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz" \
    && rm "/tmp/s6-overlay-symlinks-arch-${S6_OVERLAY_VERSION}.tar.xz"
COPY oneliner/etc /etc

# use the correct binary depending on the architecture. we do this to avoid building amd64 version ourselves,
# as it requires us to run qemu emulation which is very slow.
ARG TARGETARCH
COPY --from=rudolfs-builder-arm64 /build/rudolfs /arm64/rudolfs
COPY --from=rudolfs-builder-amd64 /rudolfs /amd64/rudolfs
RUN cp "/${TARGETARCH}/rudolfs" /usr/bin/rudolfs
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta2 /usr/bin/mutagen-agent-v0.12.0-beta2
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta6 /usr/bin/mutagen-agent-v0.12.0-beta6
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta7 /usr/bin/mutagen-agent-v0.12.0-beta7
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.13.0-beta2 /usr/bin/mutagen-agent-v0.13.0-beta2
COPY --from=reproxy-builder /usr/bin/reproxy /usr/bin/reproxy
COPY --from=sslmux-builder /go/bin/sslmux /usr/bin/sslmux
COPY --from=ssh-builder /usr/bin/ssh /usr/bin/ssh
COPY --from=web-builder /web/dist/oneliner /web/dist
COPY --from=api-builder /usr/bin/api /usr/bin/api


ENV LANG="en_US.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    LC_ALL="C" \
    S6_KILL_GRACETIME=0 \
    S6_SERVICES_GRACETIME=0 \
    S6_CMD_WAIT_FOR_SERVICES_MAXTIME=300000 \
    STURDY_GITHUB_APP_ID=0 \
    STURDY_GITHUB_APP_CLIENT_ID= \
    STURDY_GITHUB_APP_SECRET= \
    STURDY_GITHUB_APP_PRIVATE_KEY_PATH= \
    STURDY_API_ALLOW_CORS_ORIGINS=http://127.0.0.1:80 \
    STURDY_ANALYTICS_DISABLE=false

# Expose sslmux on port 7000, acting as the entrypoint to the application
EXPOSE 7000

VOLUME [ "/var/data" ]
ENTRYPOINT [ "/init" ]
