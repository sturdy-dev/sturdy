#!/bin/bash

set -euo pipefail

log() {
	echo "[ssh-prepare] $@"
}

SSH_KEYS_DIR="/var/data/ssh/keys"
SSH_KEY_NAME="ed25519"
SSH_KEY_PATH="${SSH_KEYS_DIR}/${SSH_KEY_NAME}"

generate_keys() {
	mkdir -p "${SSH_KEYS_DIR}"
	if [[ ! -f ${SSH_KEY_PATH} ]]; then
		log "Generating Mutagen SSH keys"
		mkdir -p ${SSH_KEYS_DIR}
		ssh-keygen -o -a 100 -t ed25519 -f ${SSH_KEYS_DIR}/${SSH_KEY_NAME} -C 'sturdy-server' -P ""
	fi
}

generate_keys
