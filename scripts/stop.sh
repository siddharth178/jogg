#!/bin/bash

echo "Importing the environment variables"
source "scripts/config.sh"

echo "Stopping the Docker Container for Postgres DB"
docker stop "${DOCKER_CONTAINER_NAME}"
