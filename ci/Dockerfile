FROM golang:1.18.0-buster AS builder

RUN apt-get update && apt-get install -y git gcc cmake libssl-dev

# Build libssh2 (dep for libgit2 -DGIT_SSH=TRUE)
WORKDIR /libssh2
RUN git clone https://github.com/libssh2/libssh2.git . && git checkout libssh2-1.9.0
RUN mkdir bin && cd bin && cmake .. && cmake --build . --target install

# Build libgit2
WORKDIR /libgit2
RUN git clone https://github.com/libgit2/libgit2.git . && git checkout v1.3.0
RUN mkdir build && cd build && cmake .. -DGIT_SSH=TRUE -DGIT_SSH_MEMORY_CREDENTIALS=TRUE && cmake --build . --target install

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


ENV LD_LIBRARY_PATH=/usr/local/lib
