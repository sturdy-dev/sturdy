#!/bin/bash

# This requires coreutils (brew install coreutils)

set -euo pipefail
set -x

cd ~/src

# Reset
sturdy-sync sync terminate --all
sturdy-sync daemon stop
rm ~/.sturdy

test_libgit2() {
  # Init libgit2 codebase
  rm -rf libgit2
  sturdy init 55c23cd2-a0ea-4e65-8655-0c0703e7aef1 libgit2
  sturdy status
  sturdy stop
  sturdy-sync daemon stop
  sturdy status
  sturdy start
  sturdy status
  # sturdy-sync sync terminate --all
  # sturdy stop
}

move_binary() {
  # move the sturdy-sync binary
  SYMLINK=$(which sturdy-sync)
  PRE_PATH=$(greadlink -f $(which sturdy-sync))
  NEW_PATH="${PRE_PATH}-moved"
  mv $PRE_PATH "${NEW_PATH}"

  # move the symlink
  rm $SYMLINK
  ln -s "$NEW_PATH" "$SYMLINK"
}

test_libgit2

move_binary



# test again!
# test_libgit2

# rm -rf libgit2
# sturdy init 55c23cd2-a0ea-4e65-8655-0c0703e7aef1 libgit2
sturdy status
sturdy start
sturdy init 55c23cd2-a0ea-4e65-8655-0c0703e7aef1 libgit3