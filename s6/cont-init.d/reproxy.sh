#!/usr/bin/with-contenv sh

REPROXY_LISTEN_HOST="0.0.0.0"
REPROXY_LISTEN_PORT="80"
REPROXY_LISTEN_ADDR="${REPROXY_LISTEN_HOST}:${REPROXY_LISTEN_PORT}"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

pass_down_env REPROXY_LISTEN_ADDR
