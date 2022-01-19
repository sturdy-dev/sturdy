#!/usr/bin/with-contenv sh

RUDOLFS_LISTEN_HOST="127.0.0.1"
RUDOLFS_LISTEN_PORT="8888"
RUDOLFS_ADDR="${RUDOLFS_LISTEN_HOST}:${RUDOLFS_LISTEN_PORT}"
RUDOLFS_KEY="$(openssl rand -hex 32)"
RUDOLFS_PATH="/var/data/rudolfs"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

pass_down_env RUDOLFS_ADDR
pass_down_env RUDOLFS_KEY
pass_down_env RUDOLFS_PATH
