#! /usr/bin/env bash

set -x
set -eo pipefail

if ! [ -x "$(command -v psql)" ]; then
    echo >&2 "Error: psql is not installed."
    exit 1
fi

DB_USER=postgres
DB_PASSWORD=password

DB_NAME=factory_ci
DB_PORT=5432

# Allow to skip Docker if a dockerized Postgres database is already running
if [[ -z "${SKIP_DOCKER}" ]]
then
    docker run \
        -e POSTGRES_USER=${DB_USER} \
        -e POSTGRES_PASSWORD=${DB_PASSWORD} \
        -e POSTGRES_DB=${DB_NAME} \
        -p "${DB_PORT}":5432 \
        -d postgres \
        postgres -N 1000
fi

export PGPASSWORD="${DB_PASSWORD}"
until psql -h "localhost" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c "\q" ; do
    >&2 echo "Postgres is still unavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up and running on port ${DB_PORT}"

DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}
export DATABASE_URL

# Run sql scripts within /sql directory
psql -v ON_ERROR_STOP=1 -h "localhost" --username "$DB_USER" --dbname "$DB_NAME" -f ./scripts/sql/create_projects_table.sql

>&2 echo "Postgres has been migrated, ready to go!"
