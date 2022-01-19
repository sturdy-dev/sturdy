FROM golang:1.17.6-alpine3.15 as mutagen-ssh-builder
WORKDIR /go/src/mutagen-ssh
# cache mutagen-ssh depencencies
COPY ./mutagen-ssh/go.mod ./go.mod
COPY ./mutagen-ssh/go.sum ./go.sum
RUN go mod download
# build mutagen-ssh
COPY ./mutagen-ssh .
RUN go build -v -o /usr/bin/mutagen-ssh mutagen-ssh/cmd/mutagen-ssh

FROM golang:1.17.6-alpine3.15 as api-builder
# install github.com/libgit2/git2go dependencies
RUN apk update \
    && apk add --no-cache \
        libgit2-dev=1.3.0-r0 \
        pkgconfig \
        gcc \
        libc-dev 
WORKDIR /go/src/backend
# cache api dependencies
COPY ./backend/go.mod ./go.mod
COPY ./backend/go.sum ./go.sum
RUN go mod download -x
# build backend
COPY ./backend ./
RUN go build -tags enterprise,static,system_libgit2 -v -o /usr/bin/api mash/cmd/api

FROM jasonwhite0/rudolfs:0.3.5 as rudolfs-builder

FROM node:17.3.1-alpine3.15 as web-builder
WORKDIR /web
RUN apk update \
    && apk add --no-cache \
        python3 \
        make \
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
ARG REPROXY_SHA256_SUM="35dd1cc3568533a0b6e1109e7ba630d60e2e39716eea28d3961c02f0feafee8e"
ADD "https://github.com/umputun/reproxy/releases/download/${REPROXY_VERSION}/reproxy_${REPROXY_VERSION}_linux_arm64.tar.gz" /tmp/reproxy.tar.gz
RUN SHELL="/bin/ash" \
    set -o pipefail \
    && sha256sum "/tmp/reproxy.tar.gz" \
    && echo "${REPROXY_SHA256_SUM}  /tmp/reproxy.tar.gz" | sha256sum -c \
    && tar -xzf /tmp/reproxy.tar.gz -C /usr/bin \
    && rm /tmp/reproxy.tar.gz

FROM alpine:3.15
# postgresql
RUN apk update \
    && apk add --no-cache \
        postgresql14=14.1-r5
# rudolfs
RUN apk update \
    && apk add --no-cache \
        openssl
COPY --from=rudolfs-builder /rudolfs /usr/bin/rudolfs
# api
RUN apk update \
    && apk add --no-cache \
        git-lfs=3.0.2-r0 \
        libgit2=1.3.0-r0
COPY --from=api-builder /usr/bin/api /usr/bin/api
# mutagen-ssh
RUN apk update \
    && apk add --no-cache \
        openssh-keygen
COPY --from=mutagen-ssh-builder /usr/bin/mutagen-ssh /usr/bin/mutagen-ssh
# web
COPY --from=web-builder /web/dist/oneliner /web/dist
# reproxy
COPY --from=reproxy-builder /usr/bin/reproxy /usr/bin/reproxy
# s6-overlay
ARG S6_OVERLAY_VERSION="v2.2.0.3"
ARG S6_OVERLAY_SHA256_SUM="a24ebad7b9844cf9a8de70a26795f577a2e682f78bee9da72cf4a1a7bfd5977e"
ADD "https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-aarch64-installer" /tmp/s6-overlay-installer
RUN SHELL="/bin/ash" \
    set -o pipefail \
    && sha256sum "/tmp/s6-overlay-installer" \
    && echo "${S6_OVERLAY_SHA256_SUM}  /tmp/s6-overlay-installer" | sha256sum -c \
    && chmod +x /tmp/s6-overlay-installer
RUN "/tmp/s6-overlay-installer" /
COPY s6 /etc
ENV S6_KILL_GRACETIME=0
ENV S6_SERVICES_GRACETIME=0
ENTRYPOINT [ "/init" ]
