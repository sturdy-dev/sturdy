#!/bin/bash

set -eou pipefail

banner() {
  RED='\033[0;31m'
  BLUE='\033[0;34m'
  CYAN='\033[0;36m'
  NC='\033[0m' # No Color
  printf "${CYAN}Starting Sturdy\nIf this is the first time starting the server, it might take a few minutes...${NC}\n"
  printf "\n${RED}While you're waiting, why don't you... \n${NC}"
  printf "${CYAN}‣${NC} Check out the documentation\t https://getsturdy.com/docs \n"
  printf "${CYAN}‣${NC} or join the community 🐣\t https://discord.gg/fQcH9QAVpX \n"
  printf "\n${RED}Made with <3, in Stockholm, Sweden${NC}\n\n"
}

banner

mkdir -p "/var/data/repos"

# symlink to /repos because that is the default repository location by convention.
ln -s "/var/data/repos" "/repos"
