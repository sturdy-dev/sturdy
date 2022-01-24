#!/usr/bin/with-contenv sh

API_HOST="127.0.0.1"
API_HTTP_PORT="3000"
API_GIT_PORT="3001"
API_HTTPPPROF_LISTEN_HOST="127.0.0.1"
API_HTTPPPROF_LISTEN_PORT="6060"
API_HTTPPPROF_LISTEN_ADDR="${API_HTTPPPROF_LISTEN_HOST}:${API_HTTPPPROF_LISTEN_PORT}"
API_HTTP_ADDR="${API_HOST}:${API_HTTP_PORT}"
API_GIT_ADDR="${API_HOST}:${API_GIT_PORT}"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

pass_down_env API_HTTP_ADDR
pass_down_env API_GIT_ADDR
pass_down_env API_HTTPPPROF_LISTEN_ADDR
