#!/usr/bin/with-contenv sh

REPOS_PATH="/var/data/repos"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

mkdir -p "${REPOS_PATH}"

pass_down_env REPOS_PATH
