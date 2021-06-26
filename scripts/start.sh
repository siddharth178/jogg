#!/bin/bash

echo "Importing the environment variables"
source "scripts/config.sh"

echo "Running the postgres db in a docker container"
docker start "${DOCKER_CONTAINER_NAME}"
