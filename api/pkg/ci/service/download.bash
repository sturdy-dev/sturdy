#!/bin/bash

set -euo pipefail

echoerr() { echo "$@" 1>&2; }

function change_id() {
	cat sturdy.json | jq --raw-output '.change_id'
}

function workspace_id() {
	cat sturdy.json | jq --raw-output '.workspace_id'
}

function get_workspace_url() {
	local id
	local res

	id=$1

	echoerr "[Sturdy] Downloading workspace ${id}"

	res=$(
		curl 'https://__PUBLIC_API__HOSTNAME__/graphql' \
			--silent --show-error --fail \
			-H 'Content-Type: application/json' \
			-H 'Accept: application/json' \
			-H 'Authorization: bearer __JWT__' \
			--data-binary "{\"query\":\"query { workspace(id: \\\"${id}\\\") { id downloadTarGz { url } } }\"}"
	)

	echo "$res" | jq --raw-output '.data.workspace.downloadTarGz.url'
}

function get_change_url() {
	local id
	local res

	id=$1

	echoerr "[Sturdy] Downloading change ${id}"

	res=$(
		curl 'https://__PUBLIC_API__HOSTNAME__/graphql' \
			--silent --show-error --fail \
			-H 'Content-Type: application/json' \
			-H 'Accept: application/json' \
			-H 'Authorization: bearer __JWT__' \
			--data-binary "{\"query\":\"query { change(id: \\\"${id}\\\") { id title downloadTarGz { url } } }\"}"
	)

	echo "$res" | jq --raw-output '.data.change.downloadTarGz.url'
}

function download() {
	curl $1 --silent >archive.tar.gz
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

	mkdir tmp_output 2 /dev/null &>1 || true
}

prepare

CHANGE_ID="$(change_id)"
WORKSPACE_ID="$(workspace_id)"
if [ -n "${WORKSPACE_ID}" ]; then
    download "$(get_workspace_url "$WORKSPACE_ID")"
    extract
else if [ -n "${CHANGE_ID}" ]; then
    download "$(get_change_url "$CHANGE_ID")"
    extract
else 
    echoerr "[Sturdy] No workspace or change id found, exiting"
    exit 1
fi
