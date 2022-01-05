#!/usr/bin/env bash

QUERY=$(cat receivers.sql)
PGPASSWORD=$(cat db-pwd) psql -h driva.cqawetpfgboc.eu-north-1.rds.amazonaws.com -U driva -d driva -t -A -F"," -c "${QUERY}" > output.csv
