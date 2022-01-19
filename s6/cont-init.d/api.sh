#!/usr/bin/with-contenv sh

API_HOST="127.0.0.1"
API_HTTP_PORT="3000"
API_GIT_PORT="3001"
API_HTTP_ADDR="${API_HOST}:${API_HTTP_PORT}"
API_GIT_ADDR="${API_HOST}:${API_GIT_PORT}"
API_REPOS_PATH="/var/data/api/repos"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

mkdir -p "${API_REPOS_PATH}"

pass_down_env API_HTTP_ADDR
pass_down_env API_GIT_ADDR
pass_down_env API_REPOS_PATH
