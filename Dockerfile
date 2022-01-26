FROM golang:1.17.6-alpine3.15 as ssh-builder
WORKDIR /go/src/ssh
RUN apk update \
    && apk add --no-cache \
        bash \
        git
COPY ./ssh/scripts/build-mutagen.sh ./scripts/build-mutagen.sh
RUN bash ./scripts/build-mutagen.sh
# cache ssh depencencies
COPY ./ssh/go.mod ./go.mod
COPY ./ssh/go.sum ./go.sum
RUN go mod download
# build ssh
COPY ./ssh .
RUN go build -v -o /usr/bin/ssh getsturdy.com/ssh/cmd/ssh

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
RUN go mod download -x
# build api
ARG API_BUILD=notset
COPY ./api ./
RUN go build -tags "${API_BUILD},static,system_libgit2" -v -o /usr/bin/api getsturdy.com/api/cmd/api

FROM alpine:3.15 as api
RUN apk update \
    && apk add --no-cache \
        git \
        git-lfs=3.0.2-r0 \
        libgit2=1.3.0-r0
COPY --from=api-builder /usr/bin/api /usr/bin/api
ENTRYPOINT [ "/usr/bin/api" ]

FROM jasonwhite0/rudolfs:0.3.5 as rudolfs-builder

FROM node:17.3.1-alpine3.15 as web-builder
WORKDIR /web
RUN apk update \
    && apk add --no-cache \
        python3=3.9.7-r4 \
        make=4.3-r0 \
        g++
# cache web dependencies
COPY ./web/package.json ./package.json
COPY ./web/yarn.lock ./yarn.lock
RUN yarn install --frozen-lockfile
# build web
COPY ./web .
RUN yarn build:oneliner

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
        postgresql14=14.1-r5 \
        openssl=1.1.1l-r8 \
        git \
        git-lfs=3.0.2-r0 \
        libgit2=1.3.0-r0 \
        openssh-keygen=8.8_p1-r1 \
        bash \
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
# s6-overlay
ARG S6_OVERLAY_VERSION="v2.2.0.3"
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN if [[ "$(uname -m)" == 'x86_64' ]]; then \
        ARCH='x86'; \
        S6_OVERLAY_SHA256_SUM="2f38adf08dc3aba06324c73f155e5f5b97b2d15d1c1bf7d95e7a09716247c4ca"; \
    elif [[ "$(uname -m)" == 'aarch64' ]]; then \
        ARCH='aarch64'; \
        S6_OVERLAY_SHA256_SUM="a24ebad7b9844cf9a8de70a26795f577a2e682f78bee9da72cf4a1a7bfd5977e"; \
    fi \
    && wget --quiet --output-document "/tmp/s6-overlay-installer" "https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-${ARCH}-installer" \
    && sha256sum "/tmp/s6-overlay-installer" \
    && echo "${S6_OVERLAY_SHA256_SUM}  /tmp/s6-overlay-installer" | sha256sum -c \
    && chmod +x /tmp/s6-overlay-installer \
    && /tmp/s6-overlay-installer / \
    && rm /tmp/s6-overlay-installer
COPY s6 /etc
ENV S6_KILL_GRACETIME=0 \
    S6_SERVICES_GRACETIME=0
# 80 is a port for web + api
# 22 is a port for ssh
EXPOSE 80 22
VOLUME [ "/var/data" ]
ENTRYPOINT [ "/init" ]
