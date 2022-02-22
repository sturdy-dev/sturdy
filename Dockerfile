# syntax=docker/dockerfile:1
FROM golang:1.17.6-alpine3.15 as ssh-builder
WORKDIR /go/src/ssh
RUN apk update \
    && apk add --no-cache \
    bash \
    git
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

FROM alpine:3.15 as ssh
RUN apk update \
    apk add --no-cache \
    ca-certificates=20211220-r0 
COPY --from=ssh-builder /usr/bin/ssh /usr/bin/ssh
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta2 /usr/bin/mutagen-agent-v0.12.0-beta2
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta6 /usr/bin/mutagen-agent-v0.12.0-beta6
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta7 /usr/bin/mutagen-agent-v0.12.0-beta7
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.13.0-beta2 /usr/bin/mutagen-agent-v0.13.0-beta2

FROM golang:1.17.6-alpine3.15 as api-builder
# github.com/libgit2/git2go dependencies
RUN apk update \
    && apk add --no-cache \
    libgit2-dev=1.3.0-r0 \
    pkgconfig \
    gcc \
    libc-dev=0.7.2-r3
WORKDIR /go/src/api
# cache api dependencies
COPY ./api/go.mod ./go.mod
COPY ./api/go.sum ./go.sum

# build api
ARG API_BUILD_TAGS
ARG VERSION
COPY ./api ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/go-mod \
    GOMODCACHE=/root/.cache/go-mod \
    go build \
    -tags "${API_BUILD_TAGS},static,system_libgit2" \
    -ldflags "-X getsturdy.com/api/pkg/version.Version=${VERSION}" \
    -v -o /usr/bin/api getsturdy.com/api/cmd/api

FROM alpine:3.15 as api
RUN apk update \
    && apk add --no-cache \
    git \
    git-lfs=3.0.2-r0 \
    libgit2=1.3.0-r0
COPY --from=api-builder /usr/bin/api /usr/bin/api
ENTRYPOINT [ "/usr/bin/api" ]

FROM jasonwhite0/rudolfs:0.3.5 as rudolfs-builder

FROM --platform=$BUILDPLATFORM node:17.3.1-alpine3.15 as web-builder
# The website is the same for linux/amd64 and linux/arm64 (output is html), setting --platform to run all builds on the
# native host platform. (Skips emulation!)
WORKDIR /web
RUN apk update \
    && apk add --no-cache \
    python3=3.9.7-r4 \
    make=4.3-r0 \
    g++
# cache web dependencies
COPY ./web/package.json ./package.json
COPY ./web/yarn.lock ./yarn.lock
# The --network-timeout is here to prevent network issues when building linux/amd64 images on linux/arm64 hosts
RUN --mount=type=cache,target=/root/.yarn YARN_CACHE_FOLDER=/root/.yarn \
    yarn install --frozen-lockfile \
    --network-timeout 1000000000
# build web
COPY ./web .
RUN yarn build:oneliner

FROM golang:1.17.6-alpine3.15 as sslmux-builder
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/go-mod \
    GOMODCACHE=/root/.cache/go-mod \
    go install -v github.com/JamesDunne/sslmux@v0.0.0-20180531161153-81a78ca8247d

FROM alpine:3.15 as reproxy-builder
ARG REPROXY_VERSION="v0.11.0"
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN if [[ "$(uname -m)" == 'aarch64' ]]; then \
    ARCH='arm64'; \
    REPROXY_SHA256_SUM='35dd1cc3568533a0b6e1109e7ba630d60e2e39716eea28d3961c02f0feafee8e'; \
    elif [[ "$(uname -m)" == 'x86_64' ]]; then \
    ARCH='x86_64'; \
    REPROXY_SHA256_SUM='100a1389882b8ab68ae94f37e9222f5f928ece299d8cfdf5b26c9f12f902c23a'; \
    fi \
    && wget --quiet --output-document "/tmp/reproxy.tar.gz" "https://github.com/umputun/reproxy/releases/download/${REPROXY_VERSION}/reproxy_${REPROXY_VERSION}_linux_${ARCH}.tar.gz" \
    && sha256sum "/tmp/reproxy.tar.gz" \
    && echo "${REPROXY_SHA256_SUM}  /tmp/reproxy.tar.gz" | sha256sum -c \
    && tar -xzf /tmp/reproxy.tar.gz -C /usr/bin \
    && rm /tmp/reproxy.tar.gz

FROM alpine:3.15 as oneliner
# postgresql
# openssl is needed by rudolfs to generate secret
# git, git-lfs and libgit2 are needed by api
# openssh-keygen is needed by ssh to generate ssh keys
# ca-cerificates is needed by ssh to connect to tls hosts
RUN apk update \
    && apk add --no-cache \
    postgresql14=14.2-r0 \
    openssl=1.1.1l-r8 \
    git \
    git-lfs=3.0.2-r0 \
    libgit2=1.3.0-r0 \
    openssh-keygen=8.8_p1-r1 \
    ca-certificates=20211220-r0 
COPY --from=rudolfs-builder /rudolfs /usr/bin/rudolfs
COPY --from=api-builder /usr/bin/api /usr/bin/api
COPY --from=ssh-builder /usr/bin/ssh /usr/bin/ssh
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta2 /usr/bin/mutagen-agent-v0.12.0-beta2
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta6 /usr/bin/mutagen-agent-v0.12.0-beta6
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.12.0-beta7 /usr/bin/mutagen-agent-v0.12.0-beta7
COPY --from=ssh-builder /go/src/ssh/mutagen-agent-v0.13.0-beta2 /usr/bin/mutagen-agent-v0.13.0-beta2
COPY --from=web-builder /web/dist/oneliner /web/dist
COPY --from=reproxy-builder /usr/bin/reproxy /usr/bin/reproxy
COPY --from=sslmux-builder /go/bin/sslmux /usr/bin/sslmux
# s6-overlay
ARG S6_OVERLAY_VERSION="3.0.0.2" \
    S6_OVERLAY_NOARCH_SHA256_SUM="17880e4bfaf6499cd1804ac3a6e245fd62bc2234deadf8ff4262f4e01e3ee521" \
    S6_OVERLAY_SYMLINKS_ARCH_SHA256_SUM="6ee2b8580b23c0993b1e8c66b58777f32f6ff031ba0192cccd53a31e62942c70" \
    S6_OVERLAY_SYMLINKS_NOARCH_SHA256_SUM="d67c9b436ef59ffefd4f083f07b2869662af40b2ea79a069b147dd0c926db2d3"
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
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
ENV LANG="en_US.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    LC_ALL="C" \
    S6_KILL_GRACETIME=0 \
    S6_SERVICES_GRACETIME=0 \
    S6_CMD_WAIT_FOR_SERVICES_MAXTIME=30000 \
    STURDY_GITHUB_APP_ID=0 \
    STURDY_GITHUB_APP_CLIENT_ID= \
    STURDY_GITHUB_APP_SECRET= \
    STURDY_GITHUB_APP_PRIVATE_KEY_PATH= \
    STURDY_API_ALLOW_CORS_ORIGINS=http://127.0.0.1:80
# 80 is a port for web, api and ssh
# 22 is a port for ssh
EXPOSE 80 22
VOLUME [ "/var/data" ]
ENTRYPOINT [ "/init" ]
