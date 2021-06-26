#!/bin/bash

export POSTGRES_DOCKER_IMAGE=postgres:12

export POSTGRES_NODE_PORT=5432
export POSTGRES_HOST_PORT=5432
export POSTGRES_LOCAL_DATA_VOLUME=~/local-database/jogg-db
export POSTGRES_LOCAL_SCRIPT_VOLUME="${PWD}/scripts/seed-data"
export POSTGRES_USER="jogg"
export POSTGRES_PASSWORD="jogg"
export POSTGRES_DB_NAME="jogg"

export DOCKER_CONTAINER_NAME=joggdb
export DOCKER_CONTAINER_HOSTNAME=joggdb
export DOCKER_BRIDGE_NETWORK=joggnet
export DOCKER_DATA_VOLUME=/postgres/postgres-data
export DOCKER_SCRIPT_VOLUME=/postgres/seed-data
