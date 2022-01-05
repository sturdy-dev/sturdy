#!/usr/bin/env bash

set -euo pipefail

# Run this script to install everything we need on a new server :-)

# Install Docker and a psql client
sudo yum install docker postgresql amazon-cloudwatch-agent

# Self-configure amazon-cloudwatch-agent
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s

# Authenticate the docker daemon to the registry
aws ecr get-login-password --region eu-north-1 | \
  sudo docker login --username AWS --password-stdin 902160009014.dkr.ecr.eu-north-1.amazonaws.com

# Mount external disk
# sudo mount /dev/nvme1n1 /repos
