FROM golang:1.17.6-alpine3.15 as api-builder
# install github.com/libgit2/git2go dependencies
RUN apk update \
    && apk add --no-cache \
        libgit2-dev=1.3.0-r0 \
        pkgconfig \
        gcc \
        libc-dev 
WORKDIR /go/src/backend
# cache backend dependencies
COPY ./backend/go.mod ./go.mod
COPY ./backend/go.sum ./go.sum
RUN go mod download -x
# build backend
COPY ./backend ./
RUN go build -tags enterprise,static,system_libgit2 -v -o /usr/bin/backend mash/cmd/api

FROM jasonwhite0/rudolfs:0.3.5 as rudolfs-builder

FROM alpine:3.15
# postgresql
RUN apk update \
    && apk add --no-cache \
        postgresql14=14.1-r5 \
    && mkdir /run/postgresql
# rudolfs
COPY --from=rudolfs-builder /rudolfs /usr/bin/rudolfs
# libgit2 & git-lfs
RUN apk update \
    && apk add --no-cache \
        git-lfs=3.0.2-r0 \
        libgit2=1.3.0-r0
# backend
COPY --from=api-builder /usr/bin/backend /usr/bin/backend
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
