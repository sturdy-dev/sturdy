FROM golang:1.17.3-buster AS builder

RUN apt-get update && apt-get install -y git gcc cmake libssl-dev

# Build libssh2 (dep for libgit2 -DGIT_SSH=TRUE)
WORKDIR /libssh2
RUN git clone https://github.com/libssh2/libssh2.git . && git checkout libssh2-1.9.0
RUN mkdir bin && cd bin && cmake .. && cmake --build . --target install

# Build libgit2
WORKDIR /libgit2
RUN git clone https://github.com/libgit2/libgit2.git . && git checkout v1.3.0
RUN mkdir build && cd build && cmake .. -DGIT_SSH=TRUE -DGIT_SSH_MEMORY_CREDENTIALS=TRUE && cmake --build . --target install

# Build driva
WORKDIR /src

# Cache dependency downloads
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build (and cache!) common + large dependencies
RUN go build -v go.uber.org/zap && \
    go build -v github.com/gin-gonic/gin && \
    go build -v github.com/google/go-github/v39/github && \
    go build -v github.com/gliderlabs/ssh && \
    go build -v github.com/aws/aws-sdk-go/service/s3/s3manager

COPY . .

RUN go build -tags enterprise,cloud -v -o driva mash/cmd/api

# Copy over artifacts to a new container, without all the bloat form the build step
FROM debian:buster

RUN apt-get update && apt-get install -y git libssl-dev curl

WORKDIR /

# Install git-lfs
ARG GIT_LFS_VERSION="3.0.1"
ARG GIT_LFS_SHA256SUM="29706bf26d26a4e3ddd0cad02a1d05ff4f332a2fab4ecab3bbffbb000d6a5797"
RUN mkdir -p /tmp/lfs \
    && cd /tmp/lfs \
    && curl -sLO "https://github.com/git-lfs/git-lfs/releases/download/v${GIT_LFS_VERSION}/git-lfs-linux-amd64-v${GIT_LFS_VERSION}.tar.gz" \
    && sha256sum "git-lfs-linux-amd64-v${GIT_LFS_VERSION}.tar.gz" \
    && echo "${GIT_LFS_SHA256SUM}  git-lfs-linux-amd64-v${GIT_LFS_VERSION}.tar.gz" | sha256sum -c \
    && tar -xvf "git-lfs-linux-amd64-v${GIT_LFS_VERSION}.tar.gz" \
    && bash ./install.sh \
    && cd / \
    && rm -rf /tmp/lfs

# Copy libgit
COPY --from=builder /usr/local/lib/libgit2.so /usr/local/lib/
COPY --from=builder /usr/local/lib/libgit2.so.1.3 /usr/local/lib/
COPY --from=builder /usr/local/lib/libgit2.so.1.3.0 /usr/local/lib/
ENV LD_LIBRARY_PATH=/usr/local/lib

# Copy driva
COPY --from=builder /src/driva /usr/bin/driva

# Copy runtime migrations
COPY /db/migrations /db/migrations

# Smoke test binaries
RUN git version && git-lfs version
