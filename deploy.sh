#!/usr/bin/env bash

# Run this script on a server :-)
# ./deploy.sh [VERSION]

set -euo pipefail

VERSION=$1

# Authenticate Docker
aws ecr get-login-password --region eu-north-1 | \
  sudo docker login --username AWS --password-stdin 902160009014.dkr.ecr.eu-north-1.amazonaws.com

# Make sure that /repos is mounted
df -h | grep "/repos"

# Stop the application if running
sudo docker stop driva || true
sudo docker rm driva || true

# Start!
sudo docker run \
  --name driva \
  --network sturdy \
  -e DD_AGENT_HOST=dd-agent \
  -d \
  --entrypoint driva \
  -p 3000:3000 \
  -p 3002:3002 \
  -p 6060:6060 \
  --ulimit nofile=90000:90000 \
  --volume /repos:/repos \
  --volume /secrets:/secrets \
  --log-driver=awslogs \
  --log-opt awslogs-region=eu-north-1 \
  --log-opt awslogs-group=driva \
  --log-opt awslogs-create-group=true \
  --restart always \
  "902160009014.dkr.ecr.eu-north-1.amazonaws.com/api:${VERSION}" \
  --hostname="$HOSTNAME" \
  --http-listen-addr="0.0.0.0:3000" \
  --git-listen-addr="0.0.0.0:3002" \
  --http-pprof-listen-addr="0.0.0.0:6060" \
  --repos-base-path="/repos" \
  --github-app-id=98002 \
  --github-app-client-id="Iv1.7166a6dda97db2e0" \
  --github-app-secret="something-secret-d92e7452" \
  --github-app-private-key-path="/secrets/sturdy-devtools.2021-07-01.private-key.pem" \
  --production-logger \
  --send-posthog-events \
  --db="postgres://driva:$(cat /home/ec2-user/db-pwd)@driva.cqawetpfgboc.eu-north-1.rds.amazonaws.com:5432/driva?sslmode=disable" \
  --send-invites-worker \
  --gmail-token-json-path /secrets/gmail-gustav-westling-token.json \
  --gmail-credentials-json-path /secrets/gmail-gustav-westling-credentials.json
