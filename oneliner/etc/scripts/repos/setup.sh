#!/bin/bash

set -eou pipefail

mkdir -p "/var/data/repos"

# symlink to /repos because that is the default repository location by convention.
ln -s "/var/data/repos" "/repos"
