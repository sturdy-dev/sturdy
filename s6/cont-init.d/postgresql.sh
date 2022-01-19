#!/usr/bin/with-contenv sh

# based on: https://raw.githubusercontent.com/sameersbn/docker-postgresql/master/runtime/functions

set -eou pipefail

POSTGRESQL_DB_NAME="sturdy"
POSTGRESQL_DB_USER="sturdy"
POSTGRESQL_DB_PASS="strudy"
POSTGRESQL_LISTEN_HOST="127.0.0.1"
POSTGRESQL_LISTEN_PORT="5432"
POSTGRESQL_DB_ADDR="${POSTGRESQL_LISTEN_HOST}:${POSTGRESQL_LISTEN_PORT}"

PG_HOME="/var/data/postgresql"
PG_DATADIR="${PG_HOME}/data"
PG_RUNDIR=/run/postgresql
PG_USER="postgres"
PG_PASSWORD="postgres"
PG_HBA_CONF="${PG_DATADIR}/pg_hba.conf"
PG_CONF="${PG_DATADIR}/postgresql.conf"

DB_NAME="${POSTGRESQL_DB_NAME}"
DB_USER="${POSTGRESQL_DB_USER}"
DB_PASS="${POSTGRESQL_DB_PASS}"
DB_TEMPLATE="template1"

pass_down_env() {
  local key=${1}
  local value=$(eval "echo \"\$$key\"")
  echo -n "$value" >"/var/run/s6/container_environment/${key}"
}

## Execute command as PG_USER
exec_as_postgres() {
  s6-setuidgid "${PG_USER}" "$@"
}

set_hba_param() {
  local value=${1}
  echo "${value}" >>${PG_HBA_CONF}
}

initialize_database() {
  if [[ ! -f ${PG_DATADIR} ]]; then
    echo "Initializing database..."

    if [[ -n $PG_PASSWORD ]]; then
      echo "${PG_PASSWORD}" >/tmp/pwfile
    fi

    exec_as_postgres initdb --pgdata=${PG_DATADIR} \
      --username=${PG_USER} --encoding=unicode --auth=trust ${PG_PASSWORD:+--pwfile=/tmp/pwfile} >/dev/null
  fi

  # configure path to data_directory
  set_postgresql_param "data_directory" "${PG_DATADIR}"

  echo "Trusting connections from the local network..."
  set_hba_param "host all all samenet trust"

  # allow remote connections to postgresql database
  set_hba_param "host all all 0.0.0.0/0 md5"
}

set_postgresql_param() {
  local key=${1}
  local value=${2}
  local verbosity=${3:-verbose}

  if [[ -n ${value} ]]; then
    local current=$(exec_as_postgres sed -n -e "s/^\(${key} = '\)\([^ ']*\)\(.*\)$/\2/p" ${PG_CONF})
    if [[ "${current}" != "${value}" ]]; then
      if [[ ${verbosity} == verbose ]]; then
        echo "‣ Setting postgresql.conf parameter: ${key} = '${value}'"
      fi
      value="$(echo "${value}" | sed 's|[&]|\\&|g')"
      exec_as_postgres sed -i "s|^[#]*[ ]*${key} = .*|${key} = '${value}'|" ${PG_CONF}
    fi
  fi
}

create_rundir() {
  echo "Initializing rundir..."
  mkdir -p ${PG_RUNDIR} ${PG_RUNDIR}/main.pg_stat_tmp
  chmod -R 0755 ${PG_RUNDIR}
  chmod g+s ${PG_RUNDIR}
  chown -R ${PG_USER}:${PG_USER} ${PG_RUNDIR}
}

create_user() {
  echo "Creating database user: ${DB_USER}"
  if [[ -z $(psql -U ${PG_USER} -Atc "SELECT 1 FROM pg_catalog.pg_user WHERE usename = '${DB_USER}'") ]]; then
    psql -U ${PG_USER} -c "CREATE ROLE \"${DB_USER}\" with LOGIN CREATEDB PASSWORD '${DB_PASS}';" >/dev/null
  fi
}

create_database() {
  echo "Creating database: ${DB_NAME}"
  if [[ -z $(psql -U ${PG_USER} -Atc "SELECT 1 FROM pg_catalog.pg_database WHERE datname = '${DB_NAME}'") ]]; then
    psql -U ${PG_USER} -c "CREATE DATABASE \"${DB_NAME}\" WITH TEMPLATE = \"${DB_TEMPLATE}\";;" >/dev/null
  fi

  if [[ -n ${DB_USER} ]]; then
    echo "‣ Granting access to ${DB_USER} on ${DB_NAME}"
    psql -U ${PG_USER} -c "GRANT ALL PRIVILEGES ON DATABASE \"${DB_NAME}\" to \"${DB_USER}\";" >/dev/null
  fi
}

create_datadir() {
  echo "Initializing datadir..."
  mkdir -p ${PG_HOME}
  if [[ -d ${PG_DATADIR} ]]; then
    chmod 0600 $(find ${PG_DATADIR} -type f)
    chmod 0700 $(find ${PG_DATADIR} -type d)
  fi
  chown -R ${PG_USER}:${PG_USER} ${PG_HOME}
}

configure_postgresql() {
  create_rundir
  create_datadir
  initialize_database

  # start postgres server internally for the creation of users and databases
  exec_as_postgres pg_ctl -D ${PG_DATADIR} -w start >/dev/null
  create_user
  create_database
  # stop the postgres server
  exec_as_postgres pg_ctl -D ${PG_DATADIR} -w stop >/dev/null
}

configure_postgresql

pass_down_env POSTGRESQL_DB_NAME
pass_down_env POSTGRESQL_DB_USER
pass_down_env POSTGRESQL_DB_PASS
pass_down_env POSTGRESQL_DB_ADDR
pass_down_env PG_DATADIR
