#!/usr/bin/env bash

QUERY=$(cat receivers.sql)
PGPASSWORD=$(cat ../db-pwd-new) psql -h database-1.cluster-cqawetpfgboc.eu-north-1.rds.amazonaws.com -U postgres -d sturdy -t -A -F"," -c "${QUERY}" > output.csv
