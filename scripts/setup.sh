#!/bin/bash

echo "Importing the environment variables"
source "scripts/config.sh"

echo "Destroying existing database and docker containers"
docker rm -f "${DOCKER_CONTAINER_NAME}" || true
rm -rf "${POSTGRES_LOCAL_DATA_VOLUME}" || true

echo "Creating the local directory to store the postgres db data"
mkdir -p "${POSTGRES_LOCAL_DATA_VOLUME}"

echo "Creating docker network bridge (ignore error, if run the second time)"
docker network create -d bridge "${DOCKER_BRIDGE_NETWORK}" || true

echo "Running the postgres db in a docker container"
docker run --detach \
  --name="${DOCKER_CONTAINER_NAME}" \
  --hostname="${DOCKER_CONTAINER_HOSTNAME}" \
  --network="${DOCKER_BRIDGE_NETWORK}" \
  --publish "${POSTGRES_HOST_PORT}:${POSTGRES_NODE_PORT}" \
  --volume "${POSTGRES_LOCAL_DATA_VOLUME}:${DOCKER_DATA_VOLUME}" \
  --volume "${POSTGRES_LOCAL_SCRIPT_VOLUME}:${DOCKER_SCRIPT_VOLUME}" \
  --env POSTGRES_USER="${POSTGRES_USER}"\
  --env POSTGRES_PASSWORD="${POSTGRES_PASSWORD}"\
  --env POSTGRES_DB="${POSTGRES_DB_NAME}"\
  "${POSTGRES_DOCKER_IMAGE}"

until docker exec -t "${DOCKER_CONTAINER_NAME}"\
 psql -U "${POSTGRES_USER}" -c "select 1" -d "${POSTGRES_DB_NAME}" >/dev/null;\
 do echo "Waiting for database to be setup"; sleep 1; done

echo "Updating database to latest schema: Creating the tables, constraints and indices"
docker run --rm -i -v "${PWD}/scripts/db-migrations":/migrations \
      --network host migrate/migrate:v4.14.1 \
      -path=/migrations \
      -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_HOST_PORT}/${POSTGRES_DB_NAME}?sslmode=disable" \
      up

echo "Executing a command inside the container: Import seed data"
docker exec -ti "${DOCKER_CONTAINER_NAME}" \
  sh -c "psql -U ${POSTGRES_USER} -d ${POSTGRES_DB_NAME} -f ${DOCKER_SCRIPT_VOLUME}/import-data.sql"

echo "Local Jogg Database Setup Completed"
