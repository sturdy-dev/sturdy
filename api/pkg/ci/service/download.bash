#!/bin/bash

set -euo pipefail

echoerr() { echo "$@" 1>&2; }

function change_id() {
	jq --raw-output '.change_id' < sturdy.json
}

function workspace_id() {
	jq --raw-output '.workspace_id' < sturdy.json
}

function snapshot_id() {
	jq --raw-output '.snapshot_id' < sturdy.json
}

function get_workspace_url() {
	local id
	local snapshot_id
	local res

	id=$1
	snapshot_id=$2

	echoerr "[Sturdy] Downloading workspace ${id} at ${snapshot_id}"

	res=$(
		curl 'https://__PUBLIC_API__HOSTNAME__/graphql' \
			--silent --show-error --fail-with-body \
			-H 'Content-Type: application/json' \
			-H 'Accept: application/json' \
			-H 'Authorization: bearer __JWT__' \
			--data-binary "{\"query\":\"query { workspace(id: \\\"${id}\\\") { id downloadTarGz(input: {snapshotID: \\\"${snapshot_id}\\\" }) { url } } }\"}"
	)

  url=$(echo "$res" | jq --raw-output '.data.workspace.downloadTarGz.url')

  if [ "${url}" != "null" ]; then
          echo "$url"
  else
          echoerr "[Sturdy] Failed to download"
          echoerr "$res"
          return 1
  fi
}

function get_change_url() {
	local id
	local res

	id=$1

	echoerr "[Sturdy] Downloading change ${id}"

	res=$(
		curl 'https://__PUBLIC_API__HOSTNAME__/graphql' \
			--silent --show-error --fail-with-body \
			-H 'Content-Type: application/json' \
			-H 'Accept: application/json' \
			-H 'Authorization: bearer __JWT__' \
			--data-binary "{\"query\":\"query { change(id: \\\"${id}\\\") { id title downloadTarGz { url } } }\"}"
	)

  url=$(echo "$res" | jq --raw-output '.data.change.downloadTarGz.url')

  if [ "${url}" != "null" ]; then
          echo "$url"
  else
          echoerr "[Sturdy] Failed to download"
          echoerr "$res"
          return 1
  fi
}

function download() {
	curl "$1" --silent > archive.tar.gz
}

function extract() {
	echoerr "[Sturdy] Extracting..."

	tar -xzf archive.tar.gz -C tmp_output

	echoerr "[Sturdy] Contents is now available in ./tmp_output"
}

function pre_cleanup() {
	echo "[Sturdy] found existing tmp_output, cleaning it up"
	rm -rf tmp_output
}

function prepare() {
	[ -d "tmp_output" ] && pre_cleanup
	mkdir tmp_output > /dev/null 2>&1 || true
}

prepare

CHANGE_ID="$(change_id)"
WORKSPACE_ID="$(workspace_id)"
SNAPSHOT_ID="$(snapshot_id)"

if [ -n "${WORKSPACE_ID}" ] && [ "${WORKSPACE_ID}" != "null" ] && [ -n "${SNAPSHOT_ID}" ] && [ "${SNAPSHOT_ID}" != "null" ]; then
    download "$(get_workspace_url "$WORKSPACE_ID" "$SNAPSHOT_ID")"
    extract
elif [ -n "${CHANGE_ID}" ] && [ "${CHANGE_ID}" != "null" ]; then
    download "$(get_change_url "$CHANGE_ID")"
    extract
else 
    echoerr "[Sturdy] No workspace or change id found, exiting"
    exit 1
fi
