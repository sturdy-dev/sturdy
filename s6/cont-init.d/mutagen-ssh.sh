#!/usr/bin/with-contenv sh

MUTAGEN_SSH_HOST="0.0.0.0"
MUTAGEN_SSH_PORT="22"
MUTAGEN_SSH_ADDR="${MUTAGEN_SSH_HOST}:${MUTAGEN_SSH_PORT}"
MUTAGEN_SSH_KEYS_DIR="/var/data/mutagen-ssh/keys"
MUTAGEN_SSH_KEY_NAME="ed25519"
MUTAGEN_SSH_KEY_PATH="${MUTAGEN_SSH_KEYS_DIR}/${MUTAGEN_SSH_KEY_NAME}"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

generate_keys() {
  if [[ ! -f ${MUTAGEN_SSH_KEY_PATH} ]]; then
    echo "Generating Mutagen SSH keys"
    mkdir -p ${MUTAGEN_SSH_KEYS_DIR}
    ssh-keygen -o -a 100 -t ed25519 -f ${MUTAGEN_SSH_KEYS_DIR}/${MUTAGEN_SSH_KEY_NAME} -C 'sturdy-server' -P ""
  fi
}

generate_keys
pass_down_env MUTAGEN_SSH_HOST
pass_down_env MUTAGEN_SSH_PORT
pass_down_env MUTAGEN_SSH_ADDR
pass_down_env MUTAGEN_SSH_KEY_PATH
